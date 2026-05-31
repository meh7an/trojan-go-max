---
title: "SimpleSocks Protocol"
draft: false
weight: 50
---

SimpleSocks is a lightweight proxy protocol with no authentication mechanism — essentially Trojan without the SHA-224 password hash. The purpose of stripping authentication is to reduce per-connection overhead inside a multiplexed tunnel.

SimpleSocks is only used on connections that have already been demultiplexed by smux; it is always carried inside smux.

SimpleSocks is even simpler than SOCKS5. Its header format:

```text
+-----+------+----------+----------+-----------+
| CMD | ATYP | DST.ADDR | DST.PORT |  Payload  |
+-----+------+----------+----------+-----------+
|  1  |  1   | Variable |    2     |  Variable |
+-----+------+----------+----------+-----------+
```

Field definitions are identical to those in the Trojan protocol.
