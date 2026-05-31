# Trojan-Go [![Go Report Card](https://goreportcard.com/badge/github.com/p4gefau1t/trojan-go)](https://goreportcard.com/report/github.com/p4gefau1t/trojan-go) [![Downloads](https://img.shields.io/github/downloads/p4gefau1t/trojan-go/total?label=downloads&logo=github&style=flat-square)](https://img.shields.io/github/downloads/p4gefau1t/trojan-go/total?label=downloads&logo=github&style=flat-square)

A full Trojan proxy implementation in Go, compatible with the original Trojan protocol and configuration file format. Secure, efficient, lightweight, and easy to use.

Trojan-Go supports [multiplexing](#multiplexing) for improved concurrency; a [router module](#router) for domestic/overseas traffic splitting; [CDN relay](#websocket) via WebSocket over TLS; [secondary encryption](#aead-encryption) of Trojan traffic via Shadowsocks AEAD; and a [pluggable transport layer](#transport-plugins) that allows replacing TLS with other encrypted tunnels.

Pre-built binaries are available on the [Releases page](https://github.com/p4gefau1t/trojan-go/releases). Extract and run directly — no additional dependencies required.

For questions, bug reports, or ideas, join the [Telegram discussion group](https://t.me/trojan_go_chat).

## Overview

**For the full documentation and configuration guide, see the [Trojan-Go Docs](https://p4gefau1t.github.io/trojan-go).**

Trojan-Go is compatible with the vast majority of original Trojan features, including:

- TLS tunnel transport
- UDP proxying
- Transparent proxy (NAT mode; iptables setup reference [here](https://github.com/shadowsocks/shadowsocks-libev/tree/v3.3.1#transparent-proxy))
- Resistance against GFW passive and active detection
- MySQL / **PostgreSQL** data persistence
- MySQL / **PostgreSQL** user authentication
- Per-user traffic statistics and quota enforcement

In addition, Trojan-Go extends the original with:

- Quick-deploy "easy mode"
- Automatic SOCKS5 / HTTP proxy detection
- TProxy-based transparent proxy (TCP / UDP)
- Cross-platform with no special dependencies
- Multiplexing via [smux](https://github.com/xtaci/smux) to reduce latency and improve concurrency
- Custom router module for domestic direct connect, ad blocking, and other routing rules
- WebSocket transport for CDN relay (WebSocket over TLS) and GFW man-in-the-middle mitigation
- TLS fingerprint impersonation to resist GFW Client Hello fingerprinting
- gRPC-based API for user management and speed limiting
- Pluggable transport layer — replace TLS with other protocols or plaintext, with full Shadowsocks SIP003 plugin support
- YAML configuration file support

## GUI Clients

Trojan-Go's server is compatible with all original Trojan clients (Igniter, ShadowRocket, etc.). The following clients support Trojan-Go extended features (WebSocket, Mux, etc.):

- [Qv2ray](https://github.com/Qv2ray/Qv2ray) — cross-platform (Windows / macOS / Linux), uses the Trojan-Go core, supports all extended features.
- [Igniter-Go](https://github.com/p4gefau1t/trojan-go-android) — Android client, forked from Igniter with the core replaced by Trojan-Go, supports all extended features.

## Usage

1. **Quick start — server and client (easy mode)**

    Server:
    ```shell
    sudo ./trojan-go -server -remote 127.0.0.1:80 -local 0.0.0.0:443 -key ./your_key.key -cert ./your_cert.crt -password your_password
    ```

    Client:
    ```shell
    ./trojan-go -client -remote example.com:443 -local 127.0.0.1:1080 -password your_password
    ```

2. **Start with a config file — client / server / transparent proxy / relay (normal mode)**

    ```shell
    ./trojan-go -config config.json
    ```

3. **Start client from a URL**

    ```shell
    ./trojan-go -url 'trojan-go://password@cloudflare.com/?type=ws&path=%2Fpath&host=your-site.com'
    ```

4. **Deploy with Docker**

    ```shell
    docker run \
        --name trojan-go \
        -d \
        -v /etc/trojan-go/:/etc/trojan-go \
        --network host \
        p4gefau1t/trojan-go
    ```

    Or with a custom config path:

    ```shell
    docker run \
        --name trojan-go \
        -d \
        -v /path/to/host/config:/path/in/container \
        --network host \
        p4gefau1t/trojan-go \
        /path/in/container/config.json
    ```

## Features

Trojan-Go and Trojan are generally interoperable, but using any extended feature (multiplexing, WebSocket, etc.) breaks compatibility with the original Trojan client.

### Portability

The compiled Trojan-Go binary is a single self-contained executable with no external dependencies. You can easily compile (or cross-compile) it for your server, PC, Raspberry Pi, or router. Build tags let you strip unused modules to reduce binary size.

For example, to cross-compile a client-only binary for MIPS Linux:

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -tags "client" -trimpath -ldflags "-s -w -buildid="
```

See the [Trojan-Go Docs](https://p4gefau1t.github.io/trojan-go) for the full build tag reference.

### Ease of Use

The configuration format is compatible with the original Trojan but significantly simplified — unspecified fields are filled with sensible defaults.

Server config `server.json`:

```json
{
  "run_type": "server",
  "local_addr": "0.0.0.0",
  "local_port": 443,
  "remote_addr": "127.0.0.1",
  "remote_port": 80,
  "password": ["your_awesome_password"],
  "ssl": {
    "cert": "your_cert.crt",
    "key": "your_key.key",
    "sni": "www.your-awesome-domain-name.com"
  }
}
```

Client config `client.json`:

```json
{
  "run_type": "client",
  "local_addr": "127.0.0.1",
  "local_port": 1080,
  "remote_addr": "www.your-awesome-domain-name.com",
  "remote_port": 443,
  "password": ["your_awesome_password"]
}
```

YAML is also supported. The following `client.yaml` is equivalent to `client.json` above:

```yaml
run-type: client
local-addr: 127.0.0.1
local-port: 1080
remote-addr: www.your-awesome-domain-name.com
remote-port: 443
password:
  - your_awesome_password
```

### WebSocket

Trojan-Go supports TLS + WebSocket transport, enabling CDN relay.

Add a `websocket` block to both the server and client configs to enable it:

```json
"websocket": {
    "enabled": true,
    "path": "/your-websocket-path",
    "hostname": "www.your-awesome-domain-name.com"
}
```

`hostname` is optional. Server and client `path` values must match. When WebSocket is enabled on the server, it simultaneously accepts both WebSocket and standard Trojan connections, so clients without WebSocket support can still connect.

### Multiplexing

On poor network connections, a single TLS handshake can take a long time. Trojan-Go supports multiplexing via [smux](https://github.com/xtaci/smux), carrying multiple TCP connections over a single TLS tunnel to reduce handshake latency under high concurrency.

> Multiplexing does not increase raw link speed, but reduces latency and improves the experience under many concurrent requests, such as loading image-heavy pages.

Enable it on the client with:

```json
"mux": {
    "enabled": true
}
```

Only the client needs this setting; the server detects and adapts automatically.

### Router

The built-in router module enables domestic direct connect, overseas proxy, and custom routing rules.

Three policies are available:

- `proxy` — route through the TLS tunnel to the Trojan server.
- `bypass` — connect directly from the local device.
- `block` — drop the connection.

```json
"router": {
    "enabled": true,
    "bypass": [
        "geoip:cn",
        "geoip:private",
        "full:localhost"
    ],
    "block": [
        "cidr:192.168.1.1/24"
    ],
    "proxy": [
        "domain:google.com"
    ],
    "default_policy": "proxy"
}
```

### AEAD Encryption

Trojan-Go supports Shadowsocks AEAD as a secondary encryption layer over the Trojan protocol, ensuring that WebSocket traffic cannot be identified or inspected by an untrusted CDN:

```json
"shadowsocks": {
    "enabled": true,
    "password": "my-password"
}
```

Both server and client must enable it with a matching password.

### Transport Plugins

Trojan-Go supports pluggable transport layers and is compatible with the Shadowsocks [SIP003](https://shadowsocks.org/en/wiki/Plugin.html) plugin standard. Example using `v2ray-plugin`:

> **This configuration is not secure — for demonstration only.**

Server:
```json
"transport_plugin": {
    "enabled": true,
    "type": "shadowsocks",
    "command": "./v2ray-plugin",
    "arg": ["-server", "-host", "www.example.com"]
}
```

Client:
```json
"transport_plugin": {
    "enabled": true,
    "type": "shadowsocks",
    "command": "./v2ray-plugin",
    "arg": ["-host", "www.example.com"]
}
```

### MySQL / PostgreSQL User Management

Trojan-Go supports database-backed user management compatible with the original Trojan's MySQL schema. PostgreSQL is also supported as a drop-in alternative.

MySQL config:

```json
"mysql": {
    "enabled": true,
    "server_addr": "localhost",
    "server_port": 3306,
    "database": "trojan",
    "username": "trojan",
    "password": "your_db_password",
    "check_rate": 60
}
```

PostgreSQL config:

```json
"postgresql": {
    "enabled": true,
    "server_addr": "localhost",
    "server_port": 5432,
    "database": "trojan",
    "username": "trojan",
    "password": "your_db_password",
    "check_rate": 60,
    "ssl_mode": "disable"
}
```

Both backends use the same `users` table schema. `password` is the SHA-224 hex hash of the plaintext password. Traffic values (`download`, `upload`, `quota`) are in bytes. If `download + upload > quota`, the server rejects that user's connections.

```sql
-- MySQL
CREATE TABLE users (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL,
    password CHAR(56) NOT NULL,
    quota BIGINT NOT NULL DEFAULT 0,
    download BIGINT UNSIGNED NOT NULL DEFAULT 0,
    upload BIGINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    INDEX (password)
);

-- PostgreSQL
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL,
    password CHAR(56) NOT NULL,
    quota BIGINT NOT NULL DEFAULT 0,
    download BIGINT NOT NULL DEFAULT 0,
    upload BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX ON users (password);
```

## Building

> Go >= 1.14 is required.

Using Make:

```shell
git clone https://github.com/p4gefau1t/trojan-go.git
cd trojan-go
make
make install   # optional: installs systemd service files
```

Using Go directly:

```shell
go build -tags "full"
```

Cross-compilation examples:

```shell
# 64-bit Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags "full"

# Apple Silicon
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -tags "full"

# 64-bit Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags "full"
```

Commonly used build tags:

| Tag | Description |
|-----|------------|
| `full` | All modules (default for releases) |
| `mini` | client + server + forward + nat + mysql + postgresql |
| `client` | Client mode only |
| `server` | Server mode only |
| `mysql` | MySQL user authentication backend |
| `postgresql` | PostgreSQL user authentication backend |

> **PostgreSQL note:** before building with the `postgresql` or `full` tag for the first time, run `go get github.com/lib/pq@v1.10.7 && go mod tidy` to pin the driver hash in `go.sum`.

## Acknowledgements

- [Trojan](https://github.com/trojan-gfw/trojan)
- [V2Fly](https://github.com/v2fly)
- [utls](https://github.com/refraction-networking/utls)
- [smux](https://github.com/xtaci/smux)
- [go-tproxy](https://github.com/LiamHaworth/go-tproxy)

## Stargazers over time

[![Stargazers over time](https://starchart.cc/p4gefau1t/trojan-go.svg)](https://starchart.cc/p4gefau1t/trojan-go)
