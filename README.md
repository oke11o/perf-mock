# Perf Mock

PerfMock is an open-source mock service for performance and load testing. It offers HTTP and gRPC servers to simulate
load scenarios, helping developers evaluate and enhance application performance. PerfMock is user-friendly and
configurable, ensuring robust and reliable software.

## Environment variables

Default values:
```shell
GRPC_PORT=8091
HTTP_PORT=8092
```


## Changelog

Install https://github.com/miniscruff/changie

You can add changie completion to you favorite shell https://changie.dev/cli/changie_completion/

### Using

See https://changie.dev/guide/quick-start/

Show current version `changie latest`

Show next minor version `changie next minor`

Add new comments - `changie new` - and follow interface

Create changelog release file - `changie batch v0.5.21`

Same for next version - `changie batch $(changie next patch)`

Merge to main CHANGELOG.md file - `changie merge`