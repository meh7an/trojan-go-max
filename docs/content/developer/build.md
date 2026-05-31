---
title: "Building and Customizing Trojan-Go"
draft: false
weight: 10
---

Go 1.14 or later is required. Verify your compiler version before building. Using `snap` to install and update Go is recommended.

### Building with Make

```shell
make
make install   # optional: installs systemd service files
```

### Building directly with Go

```shell
go build -tags "full"   # build the complete version
```

Cross-compilation is done via the `GOOS` and `GOARCH` environment variables:

```shell
GOOS=windows GOARCH=386 go build -tags "full"     # Windows x86
GOOS=linux   GOARCH=arm64 go build -tags "full"   # Linux ARM64
```

### Custom builds with build tags

Most Trojan-Go modules are pluggable. If you do not need certain features, or want to reduce the binary size, use build tags to include only what you need:

```shell
go build -tags "full"                                        # all modules
go build -tags "client" -trimpath -ldflags="-s -w -buildid="  # client only, symbols stripped
go build -tags "server mysql"                                # server + MySQL support
go build -tags "server postgresql"                           # server + PostgreSQL support
go build -tags "server sqlite"                               # server + SQLite support
go build -tags "server mysql postgresql sqlite"              # server + all database backends
```

The `full` tag is equivalent to:

```shell
go build -tags "api client server forward nat other mysql postgresql sqlite"
```

### Available build tags

| Tag          | Description                                                                   |
| ------------ | ----------------------------------------------------------------------------- |
| `client`     | Client mode (`run_type: client`)                                              |
| `server`     | Server mode (`run_type: server`)                                              |
| `forward`    | Forward/tunnel mode (`run_type: forward`)                                     |
| `nat`        | Transparent proxy mode (`run_type: nat`)                                      |
| `api`        | gRPC API server and client                                                    |
| `mysql`      | MySQL user authentication backend                                             |
| `postgresql` | PostgreSQL user authentication backend                                        |
| `sqlite`     | SQLite user authentication backend (file-based, no server required)           |
| `other`      | Miscellaneous utilities                                                       |
| `mini`       | Minimal build (client + server + forward + nat + mysql + postgresql + sqlite) |
| `full`       | All modules (see expansion above)                                             |

> **Note:** Each database backend requires its driver to be present in `go.sum` before first build:
>
> - MySQL: included in the original `go.sum`
> - PostgreSQL: `go get github.com/lib/pq@v1.10.7 && go mod tidy`
> - SQLite: `go get modernc.org/sqlite@v1.18.2 && go mod tidy`
