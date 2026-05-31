---
title: "Complete Configuration Reference"
draft: false
weight: 30
---

Below is a complete configuration file. Required fields are:

- `run_type`
- `local_addr`
- `local_port`
- `remote_addr`
- `remote_port`

For a server (`server`), `key` and `cert` are also required.

For a client (`client`), reverse-proxy tunnel (`forward`), and transparent proxy (`nat`), `password` is required.

All other fields default to the values shown below.

_Trojan-Go supports YAML syntax as a more human-friendly alternative. The structure is identical to JSON and the behavior is equivalent. To follow YAML naming conventions, replace underscores (`_`) with hyphens (`-`): for example, `remote_addr` becomes `remote-addr` in a YAML file._

```json
{
  "run_type": "<required>",
  "local_addr": "<required>",
  "local_port": "<required>",
  "remote_addr": "<required>",
  "remote_port": "<required>",
  "log_level": 1,
  "log_file": "",
  "password": [],
  "disable_http_check": false,
  "udp_timeout": 60,
  "ssl": {
    "verify": true,
    "verify_hostname": true,
    "cert": "<required for server>",
    "key": "<required for server>",
    "key_password": "",
    "cipher": "",
    "curves": "",
    "prefer_server_cipher": false,
    "sni": "",
    "alpn": ["http/1.1"],
    "session_ticket": true,
    "reuse_session": true,
    "plain_http_response": "",
    "fallback_addr": "",
    "fallback_port": 0,
    "fingerprint": ""
  },
  "tcp": {
    "no_delay": true,
    "keep_alive": true,
    "prefer_ipv4": false
  },
  "mux": {
    "enabled": false,
    "concurrency": 8,
    "idle_timeout": 60
  },
  "router": {
    "enabled": false,
    "bypass": [],
    "proxy": [],
    "block": [],
    "default_policy": "proxy",
    "domain_strategy": "as_is",
    "geoip": "$PROGRAM_DIR$/geoip.dat",
    "geosite": "$PROGRAM_DIR$/geosite.dat"
  },
  "websocket": {
    "enabled": false,
    "path": "",
    "host": ""
  },
  "shadowsocks": {
    "enabled": false,
    "method": "AES-128-GCM",
    "password": ""
  },
  "transport_plugin": {
    "enabled": false,
    "type": "",
    "command": "",
    "option": "",
    "arg": [],
    "env": []
  },
  "forward_proxy": {
    "enabled": false,
    "proxy_addr": "",
    "proxy_port": 0,
    "username": "",
    "password": ""
  },
  "mysql": {
    "enabled": false,
    "server_addr": "localhost",
    "server_port": 3306,
    "database": "",
    "username": "",
    "password": "",
    "check_rate": 60
  },
  "postgresql": {
    "enabled": false,
    "server_addr": "localhost",
    "server_port": 5432,
    "database": "",
    "username": "",
    "password": "",
    "check_rate": 60,
    "ssl_mode": "disable"
  },
  "api": {
    "enabled": false,
    "api_addr": "",
    "api_port": 0,
    "ssl": {
      "enabled": false,
      "key": "",
      "cert": "",
      "verify_client": false,
      "client_cert": []
    }
  }
}
```

---

## Reference

### General options

For `client`, `nat`, and `forward`, `remote_xxxx` should be your Trojan server address and port, and `local_xxxx` is the local SOCKS5/HTTP proxy address (auto-detected).

For `server`, `local_xxxx` is the address on which the Trojan server listens (port 443 is strongly recommended), and `remote_xxxx` is the address to which non-Trojan traffic is proxied — typically the local port 80.

`log_level` sets the logging verbosity. Higher values produce less output.

| Value | Meaning                    |
| ----- | -------------------------- |
| 0     | Debug and above (all logs) |
| 1     | Info and above             |
| 2     | Warning and above          |
| 3     | Error and above            |
| 4     | Fatal and above            |
| 5     | Silent (no output)         |

`log_file` — path of the log output file. If unset, logs go to standard output.

`password` — accepts multiple passwords. In addition to file-based passwords, Trojan-Go supports MySQL and PostgreSQL for password management (see below). A client's password must match one of the passwords in the server's configuration file or in the database before the server will allow the connection.

`disable_http_check` — disables the startup check that verifies the camouflage HTTP server is reachable.

`udp_timeout` — UDP session timeout in seconds.

---

### `ssl` options

`verify` — whether the client (`client`/`nat`/`forward`) verifies the server's certificate. Enabled by default. This should never be set to `false` in production, as it opens the door to man-in-the-middle attacks. If you are using a self-signed certificate, keep `verify` enabled and add the server certificate path to the client's `cert` field instead.

`verify_hostname` — whether the server checks that the client's SNI matches the server's configured SNI. If the SNI field is left empty on the server side, this check is forcibly disabled.

Server-side `cert` and `key` are required and must point to a valid, non-expired certificate and private key. Clients using certificates from a trusted CA do not need to specify `cert`. Clients connecting to a self-signed server must specify the server's certificate in their own `cert` field.

`sni` — the Server Name Indication field in the TLS Client Hello, generally matching the certificate's Common Name. For Let's Encrypt certificates, use your domain name. On the client side, if omitted, `remote_addr` is used. On the server side, if omitted, the Common Name from the certificate is used (wildcard domains such as `*.example.com` are supported).

`fingerprint` — specifies the TLS Client Hello fingerprint to impersonate, using [utls](https://github.com/refraction-networking/utls). Valid values:

| Value       | Behavior                               |
| ----------- | -------------------------------------- |
| `""`        | No fingerprint impersonation (default) |
| `"firefox"` | Impersonate Firefox                    |
| `"chrome"`  | Impersonate Chrome                     |
| `"ios"`     | Impersonate iOS                        |

When a fingerprint is set, the fields `cipher`, `curves`, `alpn`, and `session_ticket` are overridden by that fingerprint's specific values.

`alpn` — application-layer protocol negotiation list. Transmitted in the TLS Client/Server Hello. **If you are using a CDN, an incorrect `alpn` value may cause the CDN to negotiate the wrong application-layer protocol.**

`prefer_server_cipher` — whether the client prefers the cipher suite offered by the server during negotiation.

`cipher` — TLS cipher suites to use. Only modify this if you know exactly what you are doing. Under normal circumstances leave it empty; Trojan-Go will automatically select the best algorithm for the platform and remote end. If specified, separate suite names with colons (`:`), in priority order.

`curves` — elliptic curves preferred for TLS ECDHE. Only modify if you know what you are doing. Separate curve names with colons (`:`), in priority order.

`plain_http_response` — path to a file whose raw bytes are sent as a plaintext TCP response when the TLS handshake fails. Using `fallback_port` instead of this field is recommended.

`fallback_addr` / `fallback_port` — when the TLS handshake fails, Trojan-Go redirects the connection to this address. This is a Trojan-Go feature that helps conceal the server and resist GFW active probing, making port 443 behave exactly like a normal server when probed with non-TLS traffic. If `fallback_addr` is empty, `remote_addr` is used.

`key_log` — path for a TLS key log file. **Recording keys breaks TLS security and must never be used for any purpose other than debugging.**

---

### `mux` — multiplexing options

Multiplexing is a Trojan-Go feature. When both server and client are Trojan-Go, enabling mux reduces latency under high-concurrency conditions (only the client needs to enable it; the server adapts automatically).

Note: multiplexing reduces handshake latency, not link speed. It may actually decrease throughput while increasing CPU and memory usage on both ends.

`enabled` — whether to enable multiplexing.

`concurrency` — maximum number of connections a single TLS tunnel can carry. Default 8. A higher value means lower handshake latency under concurrent load but may reduce throughput. Setting it to 0 or a negative number forces all connections through a single TLS tunnel.

`idle_timeout` — how long (in seconds) an idle TLS tunnel waits before closing. A value of 0 or negative closes idle tunnels immediately.

---

### `router` — routing options

Routing is a Trojan-Go feature. Three policies are available:

- **Proxy** — route the request through the TLS tunnel to the Trojan server.
- **Bypass** — connect directly from the local machine.
- **Block** — drop the connection.

The `proxy`, `bypass`, and `block` lists accept GeoIP/GeoSite tags or custom rules. Clients can configure all three policies; servers can only use `block`.

`enabled` — whether to enable the router module.

`default_policy` — policy applied when no list matches. Default is `"proxy"`. Valid values: `"proxy"`, `"bypass"`, `"block"`.

`domain_strategy` — domain resolution strategy. Default `"as_is"`.

| Value               | Behavior                                                                                                          |
| ------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `"as_is"`           | Match against domain rules only.                                                                                  |
| `"ip_if_non_match"` | Match domain rules first; if no match, resolve to IP and match IP rules. May cause DNS leaks or poisoning.        |
| `"ip_on_demand"`    | Resolve to IP first and match IP rules; if no match, fall back to domain rules. May cause DNS leaks or poisoning. |

`geoip` / `geosite` — paths to the GeoIP and GeoSite database files. Default: `geoip.dat` and `geosite.dat` in the program directory. You can also set the `TROJAN_GO_LOCATION_ASSET` environment variable to override the working directory.

---

### `websocket` options

WebSocket transport is a Trojan-Go feature. For direct connections to a proxy node, enabling WebSocket does **not** improve link speed (it may reduce it) or security. Use it only when routing traffic through a CDN or splitting by path via nginx.

`enabled` — whether to enable WebSocket transport. When enabled on the server, both standard Trojan and WebSocket-over-Trojan connections are accepted. When enabled on the client, all traffic uses WebSocket.

`path` — the URL path for WebSocket. Must begin with `/`. Server and client must match.

`host` — the hostname in the WebSocket HTTP upgrade request. If left empty on the client, `remote_addr` is used. When using a CDN, set this to your domain name.

---

### `shadowsocks` — AEAD encryption options

This option adds a Shadowsocks AEAD encryption layer below the Trojan protocol layer, providing additional encryption inside the already-encrypted TLS tunnel. Both server and client must enable it with matching passwords and methods.

Enable this only when you cannot trust the security of the underlying transport, for example:

- Traffic is being relayed through an untrusted CDN (any CDN operated by a company registered in mainland China should be considered untrusted).
- Your TLS connection is being subjected to a man-in-the-middle attack by the GFW.
- Your certificate is invalid or cannot be verified.
- You are using a pluggable transport layer that does not guarantee cryptographic security.

`enabled` — whether to enable Shadowsocks AEAD encryption of the Trojan protocol layer.

`method` — encryption method. Valid values:

- `"CHACHA20-IETF-POLY1305"`
- `"AES-128-GCM"` (default)
- `"AES-256-GCM"`

`password` — password used to derive the master key. Must be identical on client and server when AEAD is enabled.

---

### `transport_plugin` — pluggable transport options

`enabled` — whether to replace TLS with a pluggable transport. When enabled, Trojan-Go sends **unencrypted Trojan protocol traffic in plaintext** to the plugin, allowing the plugin to apply custom obfuscation and encryption.

`type` — plugin type:

| Value           | Behavior                                                                                                                                                                                                                    |
| --------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `"shadowsocks"` | SIP003-compatible Shadowsocks obfuscation plugin. Trojan-Go rewrites `remote_addr/remote_port/local_addr/local_port` per SIP003, so the plugin communicates directly with the remote end.                                   |
| `"plaintext"`   | Raw TCP without TLS. Trojan-Go does not modify address configuration and does not start a plugin command. For use with nginx TLS offloading or advanced debugging only. **Never use in production to bypass the firewall.** |
| `"other"`       | Other plugins. Address configuration is not modified, but the plugin command is started with the supplied arguments and environment variables.                                                                              |

`command` — path to the plugin executable.

`arg` — plugin startup arguments. A list, e.g. `["-config", "test.json"]`.

`env` — plugin environment variables. A list, e.g. `["VAR1=foo", "VAR2=bar"]`.

`option` — plugin configuration string (SIP003), e.g. `"obfs=http;obfs-host=www.example.com"`.

---

### `tcp` options

`no_delay` — send TCP segments immediately without waiting for the buffer to fill.

`keep_alive` — enable TCP keep-alive probes.

`prefer_ipv4` — prefer IPv4 addresses when resolving hostnames.

---

### `mysql` — MySQL database options

Trojan-Go is compatible with Trojan's MySQL-based user management, though using the API is generally preferred.

`enabled` — whether to use MySQL for user authentication.

`check_rate` — interval in seconds at which Trojan-Go polls MySQL for user data and refreshes its in-memory cache.

All other fields are self-explanatory.

The `users` table schema is identical to the original Trojan definition. Below is an example `CREATE TABLE` statement. Note that `password` stores the SHA-224 hash of the plaintext password (as a hex string), and traffic values (`download`, `upload`, `quota`) are in bytes. You can add and remove users, or adjust quotas, by modifying the database records directly. Trojan-Go automatically updates its active user list based on quotas: if `download + upload > quota`, the server rejects that user's connections.

```sql
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
```

---

### `postgresql` — PostgreSQL database options

Trojan-Go also supports PostgreSQL for user management. The behavior is identical to the MySQL backend; the same quota-based logic and sync interval apply.

> **Note:** This feature is not available in the original Trojan. Build Trojan-Go with the `postgresql` or `full` build tag to include it.

`enabled` — whether to use PostgreSQL for user authentication.

`server_addr` — PostgreSQL host address. Default: `localhost`.

`server_port` — PostgreSQL port. Default: `5432`.

`database` — database name.

`username` — database user.

`password` — database user password.

`check_rate` — interval in seconds at which Trojan-Go polls the database and refreshes its in-memory cache. Default: `60`.

`ssl_mode` — PostgreSQL SSL mode. Default: `"disable"`. Common values: `"disable"`, `"require"`, `"verify-full"`. For production deployments, `"require"` or `"verify-full"` is strongly recommended.

The `users` table schema is compatible with the MySQL definition above. Use the following to create it in PostgreSQL:

```sql
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

---

### `forward_proxy` options

Allows routing Trojan-Go's own traffic through an upstream proxy.

`enabled` — whether to enable the upstream SOCKS5 proxy.

`proxy_addr` — upstream proxy host address.

`proxy_port` — upstream proxy port.

`username` / `password` — proxy credentials. Leave empty to connect without authentication.

---

### `api` options

Trojan-Go provides a gRPC-based API for managing and monitoring the server and client. Capabilities include per-user traffic and speed statistics, dynamic user add/remove, and rate limiting.

`enabled` — whether to enable the API.

`api_addr` — gRPC listen address.

`api_port` — gRPC listen port.

`ssl` — TLS settings for the gRPC endpoint.

- `enabled` — whether to use TLS for gRPC traffic.
- `key`, `cert` — server private key and certificate.
- `verify_client` — whether to authenticate clients with mTLS.
- `client_cert` — list of trusted client certificates when `verify_client` is true.

**Warning: never expose an API endpoint without mutual TLS authentication directly to the internet. Doing so may introduce serious security vulnerabilities.**
