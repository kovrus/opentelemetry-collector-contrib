module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/parquetexporter

go 1.25.0

require (
	go.opentelemetry.io/collector/component v1.64.0
	go.opentelemetry.io/collector/exporter v0.81.0
	go.opentelemetry.io/collector/pdata v1.64.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/go-version v1.9.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/collector v0.81.0 // indirect
	go.opentelemetry.io/collector/component/componenttest v0.158.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.81.0 // indirect
	go.opentelemetry.io/collector/consumer v0.81.0 // indirect
	go.opentelemetry.io/collector/extension v0.81.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.64.0 // indirect
	go.opentelemetry.io/collector/processor v0.81.0 // indirect
	go.opentelemetry.io/collector/receiver v0.81.0 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/sdk v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
)

retract (
	v0.76.2
	v0.76.1
	v0.65.0
)
