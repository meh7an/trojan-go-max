---
title: "Secondary Encryption with Shadowsocks AEAD"
draft: false
weight: 8
---

### Note: the original Trojan does not support this feature

The Trojan protocol itself carries no encryption — its security depends on the underlying TLS layer. In most situations TLS is sufficient, and re-encrypting Trojan traffic is unnecessary. However, in certain scenarios you may not be able to trust the TLS tunnel's security:

- Traffic is relayed through an untrusted CDN (any CDN operated by a company registered in mainland China should be treated as untrusted).
- Your TLS connection is subject to a GFW man-in-the-middle attack.
- Your certificate is invalid or cannot be verified.
- You are using a pluggable transport layer that does not guarantee cryptographic security.

Trojan-Go supports Shadowsocks AEAD encryption as an additional layer beneath the Trojan protocol layer. Both server and client must enable it simultaneously with identical passwords and methods, or communication will fail.

To enable AEAD encryption, add a `shadowsocks` option:

```json
"shadowsocks": {
    "enabled": true,
    "method": "AES-128-GCM",
    "password": "1234567890"
}
```

If `method` is omitted, `AES-128-GCM` is used by default. For more details, see the "Complete Configuration Reference" section.
