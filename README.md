Runner
---

[![Go Reference](https://pkg.go.dev/badge/github.com/irfansharif/runner.svg)](https://godocs.io/github.com/irfansharif/runner)

Runner is a thin library that's able to precisely measure on-CPU time for
goroutines. It relies on a slightly modified Go runtime for additional
instrumentation in order to function. It's a prototype for what finer-grained
CPU attribution could look like in Go.

### Contributing

The repo includes the necessary Go runtime as a submodule. To get up and
running (assumes you already have `go` installed and in your `PATH`):

```sh
$ git clone git@github.com:irfansharif/runner.git
$ cd runner
$ make go # optional; set up submodules and build the modified go runtime
```

We can now use the modified Go to run tests:
```sh
$ modules/go/bin/go test -v .
```

Alternatively, you can develop using bazel (which already points to mirrored Go
SDK's containing the modified runtime changes for darwin, linux and freebsd).
For unsupported OS/architectures, build the submoduled go as described above
and point bazel to it using the following:

```python
go_local_sdk(
    name = "go_sdk",
    path = "<path to checkout>/modules/go",
)
```

Finally, run the package tests/benchmarks or update BUILD files:

```
$ make test
$ make bench
$ make generate
```
