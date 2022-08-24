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
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/telemetryquerylanguage/tql"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

var (
	errEmptyRoute                   = errors.New("empty routing attribute provided")
	errNoExporters                  = errors.New("no exporters defined for the route")
	errNoTableItems                 = errors.New("the routing table is empty")
	errNoMissingFromAttribute       = errors.New("the FromAttribute property is empty")
	errDefaultExporterNotFound      = errors.New("default exporter not found")
	errExporterNotFound             = errors.New("exporter not found")
	errNoExportersAfterRegistration = errors.New("provided configuration resulted in no exporter available to accept data")
)

var (
	_ component.TracesProcessor  = (*processorImp)(nil)
	_ component.MetricsProcessor = (*processorImp)(nil)
	_ component.LogsProcessor    = (*processorImp)(nil)
)

type processorImp struct {
	logger *zap.Logger

	routes        map[Condition][]string
	metricsRouter router[component.MetricsExporter]
	logsRouter    router[component.LogsExporter]
	tracesRouter  router[component.TracesExporter]
}

type RouteEntry struct {
	exporters []string
	statement StmtCondition
}

type Condition interface {
	value() string
}

func NewRouteCondition(stmt string) (StmtCondition, error) {
	queries, err := tql.ParseQueries([]string{stmt}, functions(), tqllogs.ParsePath, tqllogs.ParseEnum)
	if err != nil {
		return StmtCondition{}, err
	}
	if len(queries) != 1 {
		return StmtCondition{}, errors.New("")
	}
	return StmtCondition{queries[0]}, nil

}

type StmtCondition struct {
	tql.Query
}

func (r StmtCondition) value() string {
	return ""
}

type ValueCondition struct {
	string
}

func (r ValueCondition) value() string {
	return r.string
}

// newProcessor creates new processor
func newProcessor(logger *zap.Logger, cfg config.Processor) *processorImp {
	logger.Info("building processor")

	oCfg := cfg.(*Config)

	routes := make(map[Condition][]string)
	for _, table := range oCfg.Table {
		if len(table.Statement) != 0 {
			c, err := NewRouteCondition(table.Statement)
			if err != nil {
				print(err)
			}
			routes[&c] = table.Exporters
		} else {
			if oCfg.AttributeSource == resourceAttributeSource {
				stmt := "route() "
				if oCfg.DropRoutingResourceAttribute {
					stmt = fmt.Sprintf("delete_key(resource.attributes, \"%s\") ", oCfg.FromAttribute)
				}
				stmt += " where resource.attributes[\"%s\"] == \"%s\""
				c, err := NewRouteCondition(fmt.Sprintf(stmt, oCfg.FromAttribute, table.Value))
				if err != nil {
					print(err)
				}
				routes[&c] = table.Exporters
			} else {
				routes[ValueCondition{table.Value}] = table.Exporters
			}
		}
	}

	var e *extractor
	if oCfg.AttributeSource == contextAttributeSource {
		e = &extractor{fromAttr: oCfg.FromAttribute, logger: logger}
	}
	return &processorImp{
		logger: logger,

		routes: routes,
		metricsRouter: newRouter[component.MetricsExporter](
			oCfg.DefaultExporters,
			routes,
			e,
			logger,
		),
		logsRouter: newRouter[component.LogsExporter](
			oCfg.DefaultExporters,
			routes,
			e,
			logger,
		),
		tracesRouter: newRouter[component.TracesExporter](
			oCfg.DefaultExporters,
			routes,
			e,
			logger,
		),
	}
}

func (e *processorImp) Start(_ context.Context, host component.Host) error {
	exporters := host.GetExporters()

	err := e.metricsRouter.registerExportersForType(exporters, config.MetricsDataType)
	if err != nil {
		return err
	}

	err = e.logsRouter.registerExportersForType(exporters, config.LogsDataType)
	if err != nil {
		return err
	}
	err = e.tracesRouter.registerExportersForType(exporters, config.TracesDataType)
	if err != nil {
		return err
	}

	if !e.metricsRouter.hasExporters() &&
		!e.logsRouter.hasExporters() &&
		!e.tracesRouter.hasExporters() {
		return errNoExportersAfterRegistration
	}
	return nil
}

func (e *processorImp) Shutdown(context.Context) error {
	return nil
}

func (e *processorImp) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	routedTraces := e.tracesRouter.RouteTraces(ctx, td)
	for _, rt := range routedTraces {
		for _, exp := range rt.exporters {
			// TODO: determine the proper action when errors happen
			if err := exp.ConsumeTraces(ctx, rt.signal); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *processorImp) ConsumeMetrics(ctx context.Context, tm pmetric.Metrics) error {
	routedMetrics := e.metricsRouter.RouteMetrics(ctx, tm)
	for _, rm := range routedMetrics {
		for _, exp := range rm.exporters {
			// TODO: determine the proper action when errors happen
			if err := exp.ConsumeMetrics(ctx, rm.signal); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *processorImp) ConsumeLogs(ctx context.Context, tl plog.Logs) error {
	routedLogs := e.logsRouter.RouteLogs(ctx, tl)
	for _, rl := range routedLogs {
		for _, exp := range rl.exporters {
			// TODO: determine the proper action when errors happen
			if err := exp.ConsumeLogs(ctx, rl.signal); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *processorImp) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}
