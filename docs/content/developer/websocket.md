---
title: "WebSocket"
draft: false
weight: 40
---

When traffic is relayed through a CDN, HTTPS is transparent to the CDN provider, which can inspect the WebSocket payload. Because the Trojan protocol itself is unencrypted, adding a Shadowsocks AEAD encryption layer protects both the payload and traffic patterns.

**If you are using a CDN operated by a mainland China ISP, enabling AEAD encryption is mandatory.**

With AEAD enabled, all content carried by WebSocket is encrypted by Shadowsocks AEAD. See the Shadowsocks whitepaper for the exact header format.

With WebSocket support enabled, the full protocol stack is:

| Layer           | Note                          |
| --------------- | ----------------------------- |
| Payload         |                               |
| SimpleSocks     | if multiplexing is enabled    |
| smux            | if multiplexing is enabled    |
| Trojan          |                               |
| Shadowsocks     | if AEAD encryption is enabled |
| WebSocket       |                               |
| Transport layer |                               |
