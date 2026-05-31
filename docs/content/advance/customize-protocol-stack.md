---
title: "Custom Protocol Stack"
draft: false
weight: 8
---

### Note: the original Trojan does not support this feature

Trojan-Go allows advanced users to define a fully custom protocol stack. In custom mode, Trojan-Go surrenders control of the stack and lets you compose the underlying protocol layers yourself. Examples of what becomes possible:

- Multiple nested TLS layers
- TLS → WebSocket → TLS → Shadowsocks AEAD
- Trojan protocol over Shadowsocks AEAD directly on TCP
- Unwrapping an inbound Trojan TLS stream and re-wrapping it in a new outbound TLS Trojan stream

**Do not use this feature unless you understand networking. Incorrect configuration can break Trojan-Go or cause performance and security problems.**

Trojan-Go abstracts every protocol as a _tunnel_. Each tunnel may expose a client (sender), a server (receiver), or both. Customizing the stack means defining how tunnels are composed.

### Read the "Overview" page in the Developer Guide before continuing — make sure you understand how Trojan-Go works.

Supported tunnels and their properties:

| Tunnel      | Needs stream from below | Needs packet from below | Provides stream above | Provides packet above | Can be inbound | Can be outbound |
| ----------- | ----------------------- | ----------------------- | --------------------- | --------------------- | -------------- | --------------- |
| transport   | n                       | n                       | y                     | y                     | y              | y               |
| dokodemo    | n                       | n                       | y                     | y                     | y              | n               |
| tproxy      | n                       | n                       | y                     | y                     | y              | n               |
| tls         | y                       | n                       | y                     | n                     | y              | y               |
| trojan      | y                       | n                       | y                     | y                     | y              | y               |
| mux         | y                       | n                       | y                     | n                     | y              | y               |
| simplesocks | y                       | n                       | y                     | y                     | y              | y               |
| shadowsocks | y                       | n                       | y                     | n                     | y              | y               |
| websocket   | y                       | n                       | y                     | n                     | y              | y               |
| freedom     | n                       | n                       | y                     | y                     | n              | y               |
| socks       | y                       | y                       | y                     | y                     | y              | n               |
| http        | y                       | n                       | y                     | n                     | y              | n               |
| router      | y                       | y                       | y                     | y                     | n              | y               |
| adapter     | n                       | n                       | y                     | y                     | y              | n               |

Custom stacks are defined by naming nodes (tags), providing their configuration, and then describing directed paths composed of those tags.

For a typical Trojan-Go server the paths are:

- **Inbound** (two paths; `tls` auto-detects and dispatches Trojan and WebSocket traffic):
  - `transport → tls → trojan`
  - `transport → tls → websocket → trojan`
- **Outbound** (single chain):
  - `router → freedom`

Inbound paths form a **multi-branch tree** rooted at the first node; graphs that are not trees produce undefined behavior. The outbound must be a single **chain**.

Every path must satisfy:

1. Start with a tunnel that does **not** require a stream or packet from below (`transport`, `adapter`, `tproxy`, `dokodemo`, etc.).
2. End with a tunnel that provides **both** stream and packet to the layer above (`trojan`, `simplesocks`, `freedom`, etc.).
3. All tunnels on an outbound chain must be usable as outbound; all tunnels on all inbound paths must be usable as inbound.

To activate a custom stack, set `run_type` to `custom`. All configuration fields other than `inbound` and `outbound` are ignored.

### Client example (`client.yaml`)

```yaml
run-type: custom

inbound:
  node:
    - protocol: adapter
      tag: adapter
      config:
        local-addr: 127.0.0.1
        local-port: 1080
    - protocol: socks
      tag: socks
      config:
        local-addr: 127.0.0.1
        local-port: 1080
  path:
    - - adapter
      - socks

outbound:
  node:
    - protocol: transport
      tag: transport
      config:
        remote-addr: your_server
        remote-port: 443

    - protocol: tls
      tag: tls
      config:
        ssl:
          sni: localhost
          key: server.key
          cert: server.crt

    - protocol: trojan
      tag: trojan
      config:
        password:
          - 12345678

  path:
    - - transport
      - tls
      - trojan
```

### Server example (`server.yaml`)

```yaml
run-type: custom

inbound:
  node:
    - protocol: websocket
      tag: websocket
      config:
        websocket:
          enabled: true
          hostname: example.com
          path: /ws

    - protocol: transport
      tag: transport
      config:
        local-addr: 0.0.0.0
        local-port: 443
        remote-addr: 127.0.0.1
        remote-port: 80

    - protocol: tls
      tag: tls
      config:
        remote-addr: 127.0.0.1
        remote-port: 80
        ssl:
          sni: localhost
          key: server.key
          cert: server.crt

    - protocol: trojan
      tag: trojan1
      config:
        remote-addr: 127.0.0.1
        remote-port: 80
        password:
          - 12345678

    - protocol: trojan
      tag: trojan2
      config:
        remote-addr: 127.0.0.1
        remote-port: 80
        password:
          - 87654321

  path:
    - - transport
      - tls
      - trojan1
    - - transport
      - tls
      - websocket
      - trojan2

outbound:
  node:
    - protocol: freedom
      tag: freedom

  path:
    - - freedom
```
