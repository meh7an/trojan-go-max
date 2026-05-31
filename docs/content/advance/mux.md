---
title: "Enabling Multiplexing to Improve Concurrent Performance"
draft: false
weight: 1
---

### Note: the original Trojan does not support this feature

Trojan-Go supports multiplexing to improve concurrent network performance.

Trojan is built on TLS. Before a secure TLS connection is established, both ends must perform a key negotiation and exchange — the TLS handshake. Due to GFW interference with TLS handshakes and general congestion on outbound links, a handshake typically takes close to one second or more on ordinary routes. This can noticeably increase latency when browsing or streaming.

Trojan-Go solves this with multiplexing. Each TLS connection carries multiple TCP connections. When a new proxy request arrives, it reuses an existing TLS connection rather than initiating a new TLS handshake from scratch, reducing the latency penalty of repeated handshakes.

Enabling multiplexing does **not** increase link speed (it may reduce it slightly) and may increase CPU and memory load on both client and server. In rough terms, multiplexing trades some throughput and CPU for lower latency. It is most beneficial under high-concurrency workloads, such as browsing pages with many images or firing many UDP requests in parallel.

To enable multiplexing, set the `enabled` field inside `mux` to `true`. Client-side example:

```json
"mux": {
    "enabled": true
}
```

Only the client needs to be configured; the server adapts automatically.

Complete `mux` configuration:

```json
"mux": {
    "enabled": false,
    "concurrency": 8,
    "idle_timeout": 60
}
```

`concurrency` is the maximum number of TCP connections a single TLS tunnel may carry. A higher value means the TLS tunnel is reused more aggressively and handshake-induced latency is lower, but server and client CPU load increases, which may reduce throughput. If TLS handshakes on your route are extremely slow, you can set this to `-1` to force all connections through a single TLS tunnel indefinitely.

`idle_timeout` is how many seconds an idle TLS tunnel waits before being closed. Setting it to `-1` closes idle tunnels immediately, which **may** help reduce unnecessary keep-alive traffic that could trigger GFW probing.
