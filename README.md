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

## Install

```shell
$ go install github.com/rrgmc/whymodwhy@latest 
```

## Author

Rangel Reale (rangelreale@gmail.com)
