---
title: "URL Scheme (Draft)"
draft: false
weight: 200
---

## Changelog

- `encryption` format changed to `ss;method:password`

## Overview

Thanks to @DuckSoft, @StudentMain, and @phlinhng for their discussion and contributions to the Trojan-Go URL scheme. **This is a draft — it requires more real-world usage and discussion.**

The Trojan-Go **client** accepts URLs to locate server resources. Design principles:

- Conforms to URL format standards
- Human-readable and machine-friendly
- A URL locates a Trojan-Go node resource for easy sharing

Base64-encoded data must not be embedded in URLs. Base64 does not guarantee transmission security; if you need secure URL sharing, encrypt the plaintext URL rather than modifying its format.

## Format

`$()` denotes a value that must be `encodeURIComponent`-encoded.

```text
trojan-go://
    $(trojan-password)
    @
    trojan-host
    :
    port
/?
    sni=$(tls-sni.com)&
    type=$(original|ws|h2|h2+ws)&
        host=$(websocket-host.com)&
        path=$(/websocket/path)&
    encryption=$(ss;aes-256-gcm;ss-password)&
    plugin=$(...)
#$(descriptive-text)
```

Example:

```text
trojan-go://password1234@google.com/?sni=microsoft.com&type=ws&host=youtube.com&path=%2Fgo&encryption=ss%3Baes-256-gcm%3Afuckgfw
```

Since Trojan-Go is compatible with Trojan, the Trojan URL scheme:

```text
trojan://password@remote_host:remote_port
```

is accepted and treated as equivalent to:

```text
trojan-go://password@remote_host:remote_port
```

Once the server uses a feature incompatible with Trojan, `trojan-go://` must be used. This prevents Trojan-Go URLs from being accidentally accepted by a plain Trojan client.

## Field reference

All parameter names and constant strings are case-sensitive.

### `trojan-password`

The Trojan password. Required, must not be empty, should contain only printable ASCII characters. Must be `encodeURIComponent`-encoded.

### `trojan-host`

Node IP or domain name. Required, must not be empty. IPv6 addresses must be enclosed in square brackets. IDN domains (e.g. `中文.cn`) must use `xn--...` format.

### `port`

Node port. Defaults to `443` if omitted. Must be an integer in `[1, 65535]`.

### `tls` / `allowInsecure`

This field does not exist. TLS is always enabled unless a transport plugin disables it. TLS verification must be enabled. Nodes whose server identity cannot be verified against a root CA are not suitable for sharing.

### `sni`

Custom TLS SNI. Defaults to `trojan-host` if omitted; must not be empty. Must be `encodeURIComponent`-encoded.

### `type`

Transport type. Defaults to `original` if omitted; must not be empty.

| Value      | Behavior                                        |
| ---------- | ----------------------------------------------- |
| `original` | Standard Trojan transport, no CDN-friendly path |
| `ws`       | WebSocket over TLS                              |

Future values may include `h2` and `h2+ws`.

### `host`

Custom HTTP `Host` header. Optional; defaults to `trojan-host`. Can be empty, though this may produce unexpected behavior.

> **Note:** For non-standard ports (not 80 or 443), RFC standards require the port to be appended to the hostname in the `Host` header, e.g. `example.com:44333`.

Must be `encodeURIComponent`-encoded.

### `path`

Valid when `type` is `ws`, `h2`, or `h2+ws`. Required, must not be empty, must begin with `/`. May include `&`, `#`, and `?` characters as long as the overall path is a valid URL path. Must be `encodeURIComponent`-encoded.

### `mux`

This field does not exist. Servers support `mux` by default. Whether to enable mux is a client preference and should not be encoded in a server-locating URL.

### `encryption`

Cryptographic layer for securing Trojan traffic. Optional; defaults to `none`. Must not be empty if specified.

For Shadowsocks encryption:

```text
ss;method:password
```

`method` must be one of:

- `aes-128-gcm`
- `aes-256-gcm`
- `chacha20-ietf-poly1305`

`password` must not be empty and should be printable ASCII. Semicolons in `password` do not need to be escaped. Must be `encodeURIComponent`-encoded.

### `plugin`

Reserved for future use. Optional; must not be empty if specified.

### URL fragment (`#...`)

Node description. Should not be omitted or left empty. Must be `encodeURIComponent`-encoded.
