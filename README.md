# whymodwhy

`whymodwhy` discovers what packages from the root `go.mod` file should be upgraded in order to upgrade the passed
package name.

It uses `go mod graph` to get the dependencies, so it should work for "ghost" packages that vulnerability tests
tend to find.

```shell
$ whymodwhy github.com/moby/sys/mountinfo
to upgrade 'github.com/moby/sys/mountinfo' these packages must be upgraded:
- github.com/testcontainers/testcontainers-go
- github.com/golang-migrate/migrate/v4
```

```shell
$ whymodwhy -p go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
===== go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc (v0.49.0) =====
	Version: v0.49.0 (last)
		----- Parents -----
		root.com/root_package (indirect)
		cloud.google.com/go/iam (v1.1.8)
		cloud.google.com/go/kms (v1.16.0)
		cloud.google.com/go/storage (v1.40.0)
		github.com/golang-migrate/migrate/v4 (v4.18.1)
		go.step.sm/crypto (v0.45.0)
		google.golang.org/api (v0.180.0)
		cloud.google.com/go (v0.113.0)
		cloud.google.com/go/bigquery (v1.60.0)
		cloud.google.com/go/longrunning (v0.5.7)
		cloud.google.com/go/secretmanager (v1.12.0)
		github.com/smallstep/certificates (v0.26.1)
		github.com/smallstep/cli (v0.26.1)
		----- Deps -----
		go.opentelemetry.io/otel/metric (v1.24.0)
		github.com/davecgh/go-spew (v1.1.1)
		golang.org/x/text (v0.14.0)
		google.golang.org/genproto/googleapis/rpc (v0.0.0-20231106174013-bbf56f31fb17)
		go.opentelemetry.io/otel (v1.24.0)
		go.opentelemetry.io/otel/trace (v1.24.0)
		google.golang.org/grpc (v1.61.0)
		golang.org/x/sys (v0.17.0)
		gopkg.in/yaml.v3 (v3.0.1)
		github.com/stretchr/testify (v1.8.4)
		google.golang.org/protobuf (v1.32.0)
		github.com/pmezard/go-difflib (v1.0.0)
		golang.org/x/net (v0.21.0)
		github.com/go-logr/logr (v1.4.1)
		github.com/go-logr/stdr (v1.2.2)
		github.com/golang/protobuf (v1.5.3)
	Version: v0.48.0
		----- Parents -----
		cloud.google.com/go/firestore (v1.15.0)
		cloud.google.com/go/pubsub (v1.37.0)
	Version: v0.45.0
		----- Parents -----
		github.com/google/certificate-transparency-go (v1.1.7)
```

## Install

```shell
$ go install github.com/rrgmc/whymodwhy@latest 
```

## Author

Rangel Reale (rangelreale@gmail.com)
