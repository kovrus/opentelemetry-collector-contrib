// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metricsgenerationprocessor

import (
	"fmt"
	"path/filepath"
	"testing"

	"go.opentelemetry.io/collector/confmap/confmaptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		id           config.ComponentID
		expected     config.Processor
		errorMessage string
	}{
		{
			id: config.NewComponentIDWithName(typeStr, ""),
			expected: &Config{
				ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
				Rules: []Rule{
					{
						Name:      "new_metric",
						Unit:      "percent",
						Type:      "calculate",
						Metric1:   "metric1",
						Metric2:   "metric2",
						Operation: "percent",
					},
					{
						Name:      "new_metric",
						Unit:      "unit",
						Type:      "scale",
						Metric1:   "metric1",
						ScaleBy:   1000,
						Operation: "multiply",
					},
				},
			},
		},
		{
			id:           config.NewComponentIDWithName(typeStr, "missing_new_metric"),
			errorMessage: fmt.Sprintf("missing required field %q", nameFieldName),
		},
		{
			id:           config.NewComponentIDWithName(typeStr, "missing_type"),
			errorMessage: fmt.Sprintf("missing required field %q", typeFieldName),
		},
		{
			id:           config.NewComponentIDWithName(typeStr, "invalid_generation_type"),
			errorMessage: fmt.Sprintf("%q must be in %q", typeFieldName, generationTypeKeys()),
		},
		{
			id:           config.NewComponentIDWithName(typeStr, "missing_operand1"),
			errorMessage: fmt.Sprintf("missing required field %q", metric1FieldName),
		},
		{
			id:           config.NewComponentIDWithName(typeStr, "missing_operand2"),
			errorMessage: fmt.Sprintf("missing required field %q for generation type %q", metric2FieldName, calculate),
		},
		{
			id:           config.NewComponentIDWithName(typeStr, "missing_scale_by"),
			errorMessage: fmt.Sprintf("field %q required to be greater than 0 for generation type %q", scaleByFieldName, scale),
		},
		{
			id:           config.NewComponentIDWithName(typeStr, "invalid_operation"),
			errorMessage: fmt.Sprintf("%q must be in %q", operationFieldName, operationTypeKeys()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
			require.NoError(t, err)

			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()
			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, config.UnmarshalProcessor(sub, cfg))

			if tt.expected != nil {
				assert.NoError(t, cfg.Validate())
				assert.Equal(t, tt.expected, cfg)
			} else {
				assert.EqualError(t, cfg.Validate(), tt.errorMessage)
			}
		})
	}
}
