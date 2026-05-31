---
title: "Overview"
draft: false
weight: 1
---

The core components of Trojan-Go are:

- `tunnel` — concrete implementations of each protocol
- `proxy` — the proxy core
- `config` — configuration registration and parsing
- `redirector` — active-probe deception module
- `statistic` — user authentication and statistics

Source code for each is found in the corresponding directory.

## `tunnel.Tunnel` — the Tunnel abstraction

Trojan-Go abstracts every protocol (including routing) as a _tunnel_ (`tunnel.Tunnel` interface). Each tunnel may start a server (`tunnel.Server`) and a client (`tunnel.Client`). A server can peel off and accept _streams_ (`tunnel.Conn`) and _packets_ (`tunnel.PacketConn`) from the tunnel below it. A client can create streams and packets through the tunnel below it.

Each tunnel is unaware of what lies beneath it, but it does know what lies above it.

All tunnels require the layer below to provide stream or packet transport (or both). All tunnels must provide stream transport to the layer above, though not necessarily packet transport.

A tunnel may have only a server, only a client, or both. Tunnels that have both can serve as the transport tunnel between a Trojan-Go client and server.

Note the distinction between a _Trojan-Go_ client/server and a _tunnel_ client/server. The following diagram illustrates the relationship:

```text
  Inbound                          GFW                            Outbound
-------> Tunnel A server → Tunnel B client -----> Tunnel B server → Tunnel C client ------->
            (Trojan-Go client side)                  (Trojan-Go server side)
```

The lowest tunnel is the transport layer — a tunnel that neither pulls from nor creates a stream/packet through another tunnel, acting as Tunnel A or C above.

- `transport` — pluggable transport layer
- `socks` — SOCKS5 proxy, server side only
- `tproxy` — transparent proxy, server side only
- `dokodemo` — reverse proxy, server side only
- `freedom` — free egress, client side only

These tunnels create streams and packets directly from TCP/UDP sockets and accept no underlying tunnel.

All other tunnels — as long as the layer below satisfies their stream/packet requirements — can be combined and stacked in any number and order. These act as Tunnel B in the diagram:

- `trojan`
- `websocket`
- `mux`
- `simplesocks`
- `tls`
- `router` (routing, client side only)

None of these care about the implementation below. They can however dispatch incoming streams and packets to different tunnels above based on content.

For a typical Trojan-Go client and server with WebSocket and Mux, the tunnel stacks are:

**Client**

- Inbound (tree):
  - `transport` (root)
    - `adapter` — detects HTTP vs SOCKS5 and dispatches
      - `http` (leaf)
      - `socks` (leaf)
- Outbound (chain):
  - `transport → tls → websocket → trojan → mux → simplesocks`

**Server**

- Inbound (tree):
  - `transport` (root)
    - `tls` — detects HTTP vs non-HTTP and dispatches
      - `websocket`
        - `trojan` (leaf)
          - `mux`
            - `simplesocks` (leaf)
      - `trojan` — detects mux vs plain Trojan and dispatches (leaf)
        - `mux`
          - `simplesocks` (leaf)
- Outbound (chain):
  - `freedom`

## `proxy.Proxy` — the proxy core

The proxy core listens on the leaf nodes of the inbound tunnel tree, extracts streams and packets along with their metadata, and forwards them to the single outbound tunnel client.

There can be multiple inbound protocol stacks (e.g. the client accepting both SOCKS5 and HTTP at the same time; the server accepting both plain Trojan and WebSocket-carried Trojan). There is exactly one outbound.

The inbound stacks are described as a multi-branch tree; the outbound is a simple list (chain). Each multi-child tree node has the ability to accurately detect and dispatch streams/packets to the correct child — which is consistent with the assumption that each protocol knows which protocols it can carry above it.
