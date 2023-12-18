module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver

go 1.19

require (
	github.com/google/go-cmp v0.6.0
	github.com/hashicorp/go-version v1.6.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.81.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest v0.81.0
	github.com/stretchr/testify v1.8.4
	github.com/testcontainers/testcontainers-go v0.21.0
	go.opentelemetry.io/collector/component v0.91.0
	go.opentelemetry.io/collector/config/confighttp v0.81.0
	go.opentelemetry.io/collector/config/configopaque v0.81.0
	go.opentelemetry.io/collector/config/configtls v0.81.0
	go.opentelemetry.io/collector/confmap v0.91.0
	go.opentelemetry.io/collector/consumer v0.91.0
	go.opentelemetry.io/collector/featuregate v1.0.0
	go.opentelemetry.io/collector/pdata v1.0.0
	go.opentelemetry.io/collector/receiver v0.91.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.26.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/containerd/containerd v1.6.19 // indirect
	github.com/cpuguy83/dockercfg v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v24.0.4+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/imdario/mergo v0.3.15 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/knadh/koanf/v2 v2.0.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20220423185008-bf980b35cac4 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.81.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/opencontainers/runc v1.1.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/cors v1.9.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/collector v0.91.0 // indirect
	go.opentelemetry.io/collector/config/configauth v0.81.0 // indirect
	go.opentelemetry.io/collector/config/configcompression v0.81.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.91.0 // indirect
	go.opentelemetry.io/collector/config/internal v0.81.0 // indirect
	go.opentelemetry.io/collector/extension v0.91.0 // indirect
	go.opentelemetry.io/collector/extension/auth v0.81.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.42.0 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	golang.org/x/exp v0.0.0-20230510235704-dd950f8aeaea // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest => ../../pkg/pdatatest

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil => ../../pkg/pdatautil

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal => ../../internal/coreinternal

// see https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/24240
replace github.com/docker/docker v24.0.4+incompatible => github.com/docker/docker v24.0.5-0.20230719162248-f022632503d1+incompatible

retract (
	v0.76.2
	v0.76.1
	v0.65.0
)
