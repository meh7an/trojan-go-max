---
title: "Multiplexing"
draft: false
weight: 30
---

Trojan-Go implements multiplexing using [smux](https://github.com/xtaci/smux) and defines the SimpleSocks protocol for proxy transport over multiplexed connections.

When multiplexing is enabled, the client first establishes a TLS connection using the standard Trojan protocol format, but sets the Command field to `0x7f` (`protocol.Mux`), signaling that this connection is a multiplexed tunnel (analogous to HTTP's `Upgrade` mechanism). The connection is then handed to the smux client. The server, upon receiving the header, hands the connection to the smux server, which demultiplexes all streams. Each individual smux stream uses the SimpleSocks protocol (Trojan with the authentication portion removed) to specify the proxy destination.

Protocol stack from top to bottom:

| Layer                   |
| ----------------------- |
| Payload                 |
| SimpleSocks             |
| smux                    |
| Trojan (authentication) |
| Underlying transport    |
