// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package routingprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor"

import (
	"context"
	"errors"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/telemetryquerylanguage/contexts/tqllogs"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/telemetryquerylanguage/contexts/tqlmetrics"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/telemetryquerylanguage/contexts/tqltraces"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

// router routes logs, metrics and traces using the configured attributes and
// attribute sources.
// Upon routing it also groups the logs, metrics and spans into a joint upper level
// structure (plog.Logs, pmetric.Metrics and ptrace.Traces respectively) in order
// to not cause higher CPU usage in the exporters when exproting data (it's always
// better to batch before exporting).
type router[E component.Exporter] struct {
	logger    *zap.Logger
	extractor *extractor

	defaultExportersIDs []string
	defaultExporters    []E

	routes    map[Condition][]string
	exporters map[Condition][]E
}

// newRouter creates a new router instance with its type parameter constrained
// to component.Exporter. router can be instantiated with component.TracesExporter,
// component.MetricsExporter, and component.LogsExporter type arguments.
func newRouter[E component.Exporter](
	defaultExporters []string,
	routes map[Condition][]string,
	context *extractor,
	logger *zap.Logger,
) router[E] {
	return router[E]{
		logger:              logger,
		routes:              routes,
		defaultExportersIDs: defaultExporters,
		extractor:           context,
		exporters:           make(map[Condition][]E),
	}
}

// signal is a type constraint that permits plog.Logs, pmetric.Metrics,
// and ptrace.Traces types.
type signal interface {
	plog.Logs | pmetric.Metrics | ptrace.Traces
}

type routedSignal[E component.Exporter, S signal] struct {
	signal    S
	exporters []E
}

func (r *router[E]) RouteMetrics(ctx context.Context, tm pmetric.Metrics) []routedSignal[E, pmetric.Metrics] {
	if r.extractor == nil {
		return r.routeMetrics(ctx, tm)
	}
	return []routedSignal[E, pmetric.Metrics]{r.routeMetricsForContext(ctx, tm)}
}

func (r *router[E]) routeMetrics(_ context.Context, tm pmetric.Metrics) []routedSignal[E, pmetric.Metrics] {
	// This way we're not ending up with all the metrics split up which would cause
	// higher CPU usage.
	routingMap := map[Condition]struct {
		exporters  []E
		resMetrics pmetric.ResourceMetricsSlice
	}{}

	ctx := tqlmetrics.MetricTransformContext{}
	resMetricsSlice := tm.ResourceMetrics()
	for i := 0; i < resMetricsSlice.Len(); i++ {
		resMetrics := resMetricsSlice.At(i)
		ctx.Resource = resMetrics.Resource()

		for c, routes := range r.exporters {
			stmt, ok := c.(*StmtCondition)
			if !ok {
				continue
			}
			groupKey := &StmtCondition{}
			exp := r.defaultExporters

			if stmt.Condition(ctx) {
				stmt.Function(ctx)
				exp = routes
				groupKey = stmt
			}

			if rEntry, ok := routingMap[groupKey]; ok {
				resMetrics.MoveTo(rEntry.resMetrics.AppendEmpty())
			} else {
				new := pmetric.NewResourceMetricsSlice()
				resMetrics.MoveTo(new.AppendEmpty())
				routingMap[groupKey] = struct {
					exporters  []E
					resMetrics pmetric.ResourceMetricsSlice
				}{
					exporters:  exp,
					resMetrics: new,
				}
			}
		}
	}

	// Now that we have all the ResourceMetrics grouped, let's create pmetric.Metrics
	// for each group and add it to the returned routedMetrics slice.
	ret := make([]routedSignal[E, pmetric.Metrics], 0, len(routingMap))
	for _, rEntry := range routingMap {
		metrics := pmetric.NewMetrics()
		metrics.ResourceMetrics().EnsureCapacity(rEntry.resMetrics.Len())
		rEntry.resMetrics.MoveAndAppendTo(metrics.ResourceMetrics())

		ret = append(ret, routedSignal[E, pmetric.Metrics]{
			signal:    metrics,
			exporters: rEntry.exporters,
		})
	}

	return ret
}

func (r *router[E]) routeMetricsForContext(ctx context.Context, tm pmetric.Metrics) routedSignal[E, pmetric.Metrics] {
	value := r.extractor.extractFromContext(ctx)
	exp, ok := r.exporters[ValueCondition{value}]
	if !ok {
		return routedSignal[E, pmetric.Metrics]{
			signal:    tm,
			exporters: r.defaultExporters,
		}
	}

	return routedSignal[E, pmetric.Metrics]{
		signal:    tm,
		exporters: exp,
	}
}

func (r *router[E]) RouteTraces(ctx context.Context, tr ptrace.Traces) []routedSignal[E, ptrace.Traces] {
	if r.extractor == nil {
		return r.routeTracesForResource(ctx, tr)
	}
	return []routedSignal[E, ptrace.Traces]{r.routeTracesForContext(ctx, tr)}
}

func (r *router[E]) routeTracesForResource(_ context.Context, tr ptrace.Traces) []routedSignal[E, ptrace.Traces] {
	// routingEntry is used to group ptrace.ResourceSpans that are routed to
	// the same set of exporters.
	// This way we're not ending up with all the logs split up which would cause
	// higher CPU usage.
	routingMap := map[Condition]struct {
		exporters []E
		resSpans  ptrace.ResourceSpansSlice
	}{}
	ctx := tqltraces.SpanTransformContext{}
	resSpansSlice := tr.ResourceSpans()
	for i := 0; i < resSpansSlice.Len(); i++ {
		resSpans := resSpansSlice.At(i)
		ctx.Resource = resSpans.Resource()

		for c, routes := range r.exporters {
			stmt, ok := c.(*StmtCondition)
			if !ok {
				continue
			}

			groupKey := &StmtCondition{}
			exp := r.defaultExporters
			if stmt.Condition(ctx) {
				stmt.Function(ctx)
				exp = routes
				groupKey = stmt
			}

			if rEntry, ok := routingMap[groupKey]; ok {
				resSpans.MoveTo(rEntry.resSpans.AppendEmpty())
			} else {
				new := ptrace.NewResourceSpansSlice()
				resSpans.MoveTo(new.AppendEmpty())
				routingMap[groupKey] = struct {
					exporters []E
					resSpans  ptrace.ResourceSpansSlice
				}{
					exporters: exp,
					resSpans:  new,
				}
			}
		}
	}
	// Now that we have all the ResourceSpans grouped, let's create ptrace.Traces
	// for each group and add it to the returned routedTraces slice.
	ret := make([]routedSignal[E, ptrace.Traces], 0, len(routingMap))
	for _, rEntry := range routingMap {
		traces := ptrace.NewTraces()
		traces.ResourceSpans().EnsureCapacity(rEntry.resSpans.Len())
		rEntry.resSpans.MoveAndAppendTo(traces.ResourceSpans())

		ret = append(ret, routedSignal[E, ptrace.Traces]{
			signal:    traces,
			exporters: rEntry.exporters,
		})
	}

	return ret
}

func (r *router[E]) routeTracesForContext(ctx context.Context, tr ptrace.Traces) routedSignal[E, ptrace.Traces] {
	value := r.extractor.extractFromContext(ctx)
	exp, ok := r.exporters[ValueCondition{value}]
	if !ok {
		return routedSignal[E, ptrace.Traces]{
			signal:    tr,
			exporters: r.defaultExporters,
		}
	}

	return routedSignal[E, ptrace.Traces]{
		signal:    tr,
		exporters: exp,
	}
}

func (r *router[E]) RouteLogs(ctx context.Context, tl plog.Logs) []routedSignal[E, plog.Logs] {
	if r.extractor == nil {
		return r.routeLogsForResource(ctx, tl)
	}
	return []routedSignal[E, plog.Logs]{r.routeLogsForContext(ctx, tl)}
}

func (r *router[E]) routeLogsForResource(_ context.Context, tl plog.Logs) []routedSignal[E, plog.Logs] {
	// routingEntry is used to group plog.ResourceLogs that are routed to
	// the same set of exporters.
	// This way we're not ending up with all the logs split up which would cause
	// higher CPU usage.
	routingMap := map[Condition]struct {
		exporters []E
		resLogs   plog.ResourceLogsSlice
	}{}

	ctx := tqllogs.LogTransformContext{}
	for i := 0; i < tl.ResourceLogs().Len(); i++ {
		resLogs := tl.ResourceLogs().At(i)
		ctx.Resource = resLogs.Resource()

		for c, routes := range r.exporters {
			stmt, ok := c.(*StmtCondition)
			if !ok {
				continue
			}

			groupKey := &StmtCondition{}
			exp := r.defaultExporters
			if stmt.Condition(ctx) {
				stmt.Function(ctx)
				exp = routes
				groupKey = stmt
			}

			if rEntry, ok := routingMap[groupKey]; ok {
				resLogs.MoveTo(rEntry.resLogs.AppendEmpty())
			} else {
				new := plog.NewResourceLogsSlice()
				resLogs.MoveTo(new.AppendEmpty())

				routingMap[groupKey] = struct {
					exporters []E
					resLogs   plog.ResourceLogsSlice
				}{
					exporters: exp,
					resLogs:   new,
				}
			}
		}
	}
	// Now that we have all the ResourceLogs grouped, let's create plog.Logs
	// for each group and add it to the returned routedLogs slice.
	ret := make([]routedSignal[E, plog.Logs], 0, len(routingMap))
	for _, rEntry := range routingMap {
		logs := plog.NewLogs()
		logs.ResourceLogs().EnsureCapacity(rEntry.resLogs.Len())
		rEntry.resLogs.MoveAndAppendTo(logs.ResourceLogs())

		ret = append(ret, routedSignal[E, plog.Logs]{
			signal:    logs,
			exporters: rEntry.exporters,
		})
	}

	return ret
}

func (r *router[E]) routeLogsForContext(ctx context.Context, tl plog.Logs) routedSignal[E, plog.Logs] {
	value := r.extractor.extractFromContext(ctx)
	exp, ok := r.exporters[ValueCondition{value}]
	if !ok {
		return routedSignal[E, plog.Logs]{
			signal:    tl,
			exporters: r.defaultExporters,
		}
	}
	return routedSignal[E, plog.Logs]{
		signal:    tl,
		exporters: exp,
	}
}

func (r *router[E]) registerExportersForType(
	exporters map[config.DataType]map[config.ComponentID]component.Exporter,
	typ config.DataType,
) error {
	err := r.registerExporters(exporters[typ])
	if err != nil {
		if errors.Is(err, errDefaultExporterNotFound) || errors.Is(err, errExporterNotFound) {
			r.logger.Warn(
				"can't find the exporter for the routing processor for this pipeline type."+
					" This is OK if you did not specify this processor for that pipeline type",
				zap.Any("pipeline_type", typ),
				zap.Error(err),
			)
		} else {
			return err
		}
	}

	return nil
}

func (r *router[E]) registerExporters(exporters map[config.ComponentID]component.Exporter) error {
	available := make(map[string]component.Exporter)
	for id, exp := range exporters {
		exporter, ok := exp.(E)
		if !ok {
			return fmt.Errorf("the exporter %q isn't a %T exporter", id.String(), new(E))
		}
		available[id.String()] = exporter
	}

	// default exporters
	if err := r.registerExportersForDefaultRoute(available); err != nil {
		return err
	}

	// for each routing table entry
	for stmt, exp := range r.routes {
		if err := r.registerExportersForRoute(stmt, exp, available); err != nil {
			return err
		}
	}

	return nil
}

// registerExportersForDefaultRoute registers the configured default exporters
// using the provided available exporters map.
func (r *router[E]) registerExportersForDefaultRoute(available map[string]component.Exporter) error {
	for _, exp := range r.defaultExportersIDs {
		v, ok := available[exp]
		if !ok {
			return fmt.Errorf("error registering default exporter %q: %w",
				exp, errDefaultExporterNotFound,
			)
		}
		r.defaultExporters = append(r.defaultExporters, v.(E))
	}
	return nil
}

// registerExportersForRoute registers the requested exporters using the provided
// available exporters map to check if they were available.
func (r *router[E]) registerExportersForRoute(
	route Condition,
	requested []string,
	available map[string]component.Exporter) error {
	r.logger.Debug("Registering exporter for route",
		//zap.String("route", route),
		zap.Any("requested", requested),
	)

	for _, exp := range requested {
		v, ok := available[exp]
		if !ok {
			return fmt.Errorf("error registering route %q for exporter %q: %w",
				route, exp, errExporterNotFound,
			)
		}
		r.exporters[route] = append(r.exporters[route], v.(E))
	}
	return nil
}

func (r *router[E]) hasExporters() bool {
	if len(r.exporters) == 0 && len(r.defaultExporters) == 0 {
		return false
	}
	return true
}
