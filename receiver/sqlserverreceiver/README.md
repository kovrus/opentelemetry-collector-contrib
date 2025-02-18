# Microsoft SQL Server Receiver

<!-- status autogenerated section -->
| Status        |           |
| ------------- |-----------|
| Stability     | [beta]: metrics   |
| Distributions | [contrib], [observiq], [sumo] |
| Issues        | ![Open issues](https://img.shields.io/github/issues-search/open-telemetry/opentelemetry-collector-contrib?query=is%3Aissue%20is%3Aopen%20label%3Areceiver%2Fsqlserver%20&label=open&color=orange&logo=opentelemetry) ![Closed issues](https://img.shields.io/github/issues-search/open-telemetry/opentelemetry-collector-contrib?query=is%3Aissue%20is%3Aclosed%20label%3Areceiver%2Fsqlserver%20&label=closed&color=blue&logo=opentelemetry) |
| [Code Owners](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/CONTRIBUTING.md#becoming-a-code-owner)    | [@djaglowski](https://www.github.com/djaglowski), [@StefanKurek](https://www.github.com/StefanKurek) |

[beta]: https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
[observiq]: https://github.com/observIQ/observiq-otel-collector
[sumo]: https://github.com/SumoLogic/sumologic-otel-collector
<!-- end autogenerated section -->

The `sqlserver` receiver grabs metrics about a Microsoft SQL Server instance using the Windows Performance Counters.
Because of this, it is a Windows only receiver.

## Configuration

The following settings are optional:
- `collection_interval` (default = `10s`): The internal at which metrics should be emitted by this receiver.

To collect from a SQL Server with a named instance, both `computer_name` and `instance_name` are required. For a default SQL Server setup, these settings are optional.
- `computer_name` (optional): The computer name identifies the SQL Server name or IP address of the computer being monitored.
- `instance_name` (optional): The instance name identifies the specific SQL Server instance being monitored.

Example:

```yaml
    receivers:
      sqlserver:
        collection_interval: 10s
```

When a named instance is used, a computer name and a instance name must be specified.
Example with named instance:

```yaml
    receivers:
      sqlserver:
        collection_interval: 10s
        computer_name: CustomServer
        instance_name: CustomInstance
        resource_attributes:
          sqlserver.computer.name:
            enabled: true
          sqlserver.instance.name:
            enabled: true
```

The full list of settings exposed for this receiver are documented [here](./config.go) with detailed sample configurations [here](./testdata/config.yaml).

## Metrics

Details about the metrics produced by this receiver can be found in [documentation.md](./documentation.md)

