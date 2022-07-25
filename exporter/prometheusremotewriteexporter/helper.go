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

package prometheusremotewriteexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"

import (
	"context"
	"errors"
	"fmt"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter/tenant"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"sort"

	"github.com/prometheus/prometheus/prompb"
)

// batchTimeSeries splits series into multiple batch write requests.
func batchTimeSeries(
	ctx context.Context,
	source tenant.Source,
	tsMap map[string]*prompb.TimeSeries,
	maxBatchByteSize int,
) (map[string][]*prompb.WriteRequest, error) {
	if len(tsMap) == 0 {
		return nil, errors.New("invalid tsMap: cannot be empty map")
	}

	var requests = make(map[string][]*prompb.WriteRequest)
	var tsArray = make(map[string][]prompb.TimeSeries)
	var sizeOfCurrentBatch = make(map[string]int)

	for _, v := range tsMap {
		sizeOfSeries := v.Size()

		tenant, err := source.GetTenant(ctx, v)
		if err != nil {
			return nil, consumererror.NewPermanent(fmt.Errorf("failed to determine the tenant: %w", err))
		}

		tenantSizeOfCurrentBatch, ok := sizeOfCurrentBatch[tenant]
		if !ok {
			tenantSizeOfCurrentBatch = 0
		}

		if tenantSizeOfCurrentBatch+sizeOfSeries >= maxBatchByteSize {
			wrapped := convertTimeseriesToRequest(tsArray[tenant])
			requests[tenant] = append(requests[tenant], wrapped)

			tsArray[tenant] = make([]prompb.TimeSeries, 0)
			tenantSizeOfCurrentBatch = 0
		}

		tsArray[tenant] = append(tsArray[tenant], *v)
		sizeOfCurrentBatch[tenant] = tenantSizeOfCurrentBatch + sizeOfSeries
	}

	for tenant, v := range tsArray {
		if len(v) != 0 {
			wrapped := convertTimeseriesToRequest(v)
			requests[tenant] = append(requests[tenant], wrapped)
		}
	}

	return requests, nil
}

func convertTimeseriesToRequest(tsArray []prompb.TimeSeries) *prompb.WriteRequest {
	// the remote_write endpoint only requires the timeseries.
	// otlp defines it's own way to handle metric metadata
	return &prompb.WriteRequest{
		// Prometheus requires time series to be sorted by Timestamp to avoid out of order problems.
		// See:
		// * https://github.com/open-telemetry/wg-prometheus/issues/10
		// * https://github.com/open-telemetry/opentelemetry-collector/issues/2315
		Timeseries: orderBySampleTimestamp(tsArray),
	}
}

func orderBySampleTimestamp(tsArray []prompb.TimeSeries) []prompb.TimeSeries {
	for i := range tsArray {
		sL := tsArray[i].Samples
		sort.Slice(sL, func(i, j int) bool {
			return sL[i].Timestamp < sL[j].Timestamp
		})
	}
	return tsArray
}
