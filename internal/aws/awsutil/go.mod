module github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil

go 1.19

require (
	github.com/aws/aws-sdk-go v1.44.301
	github.com/stretchr/testify v1.12.1
	go.uber.org/zap v1.24.0
	golang.org/x/net v0.12.0
)

require (
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/text v0.11.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

retract (
	v0.76.2
	v0.76.1
	v0.65.0
)
