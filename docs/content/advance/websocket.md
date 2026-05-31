---
title: "Using WebSocket for CDN Relay and Man-in-the-Middle Mitigation"
draft: false
weight: 2
---

### Note: the original Trojan does not support this feature

Trojan-Go supports TLS + WebSocket transport, making it possible to relay traffic through a CDN.

The Trojan protocol itself carries no encryption — its security depends on the outer TLS layer. However, once traffic passes through a CDN, TLS is transparent to the CDN provider, meaning they can inspect the plaintext. **If you are using an untrusted CDN (any CDN operated by a company registered or licensed in mainland China should be treated as untrusted), you must enable Shadowsocks AEAD encryption for WebSocket traffic to prevent identification and inspection.**

To enable WebSocket support, add a `websocket` block to both the server and client configuration files, set `enabled` to `true`, and fill in the `path` and `host` fields. Example:

```json
"websocket": {
    "enabled": true,
    "path": "/your-websocket-path",
    "host": "example.com"
}
```

`host` is the hostname, typically your domain. It is optional on the client; if left empty, `remote_addr` is used.

`path` is the WebSocket URL path. It must begin with `/`. There are no other restrictions beyond a valid URL format, but server and client `path` values must match. Choose a long, non-obvious path to reduce the risk of GFW direct active probing.

The client's `host` is sent in the WebSocket HTTP upgrade request to the CDN and must be valid. The `path` must match on both sides for the WebSocket handshake to succeed.

A complete client configuration example:

```json
{
  "run_type": "client",
  "local_addr": "127.0.0.1",
  "local_port": 1080,
  "remote_addr": "www.your_awesome_domain_name.com",
  "remote_port": 443,
  "password": ["your_password"],
  "websocket": {
    "enabled": true,
    "path": "/your-websocket-path",
    "host": "example.com"
  },
  "shadowsocks": {
    "enabled": true,
    "password": "12345678"
  }
}
```
