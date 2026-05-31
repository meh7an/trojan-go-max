---
title: "How Trojan Works"
draft: false
weight: 21
---

This page gives a brief overview of how the Trojan protocol works. If you are not interested in how the GFW and Trojan operate, you may skip this section — but reading it is recommended so you can better protect the security and undetectability of your traffic.

## Why Shadowsocks (with stream ciphers) is easy to block

Early firewalls only intercepted and inspected outbound traffic passively. Shadowsocks's encryption was designed so that the transmitted packets had almost no identifiable features — they looked like a completely random bit stream — which was effective against the early GFW.

The current GFW has moved to **active probing**. When it detects a suspicious, unrecognized connection (high traffic volume, random byte stream, high port number, etc.), it will **actively connect** to that server port and replay previously captured traffic (sometimes with deliberate modifications). When a Shadowsocks server detects an abnormal connection it closes it. This abnormal traffic pattern and the abrupt disconnection are treated as fingerprints of a Shadowsocks server, and the server is added to a suspect list. That list may not take effect immediately, but during sensitive periods the servers on it may be temporarily or permanently blocked. Whether a block is applied may also involve human decision-making.

For more detail, see [this report](https://gfw.report/blog/gfw_shadowsocks/).

## How Trojan bypasses the GFW

Unlike Shadowsocks, Trojan does not use a custom encryption protocol to hide itself. Instead it uses the conspicuously standard TLS protocol, making its traffic look identical to normal HTTPS traffic. TLS is a mature encryption system — HTTPS is simply HTTP carried over TLS. A **correctly configured** TLS tunnel guarantees:

- **Confidentiality** — the GFW cannot learn what is being transmitted.
- **Integrity** — any attempt by the GFW to tamper with the encrypted payload is detected by both ends.
- **Non-repudiation** — the GFW cannot forge either party's identity.
- **Forward secrecy** — even if a key is later leaked, the GFW cannot decrypt previously recorded traffic.

Against passive detection, Trojan traffic is indistinguishable from HTTPS. HTTPS accounts for more than half of all internet traffic today, and after a TLS handshake every byte is ciphertext, so there is no practical way to identify Trojan traffic inside it.

Against active probing, when the GFW actively connects to a Trojan server, Trojan correctly recognizes non-Trojan traffic. Unlike Shadowsocks, it does **not** close the connection — it instead proxies it to a normal web server. From the GFW's perspective, the server behaves exactly like an ordinary HTTPS site, with no way to determine that it is a Trojan proxy. This is why Trojan recommends using a legitimate domain name with a certificate signed by a trusted CA: it makes your server completely indistinguishable from a normal HTTPS server under active probing.

The only remaining ways to identify and block Trojan connections are blanket bans (blocking entire IP ranges, certificate types, domain classes, or all outbound HTTPS), or large-scale man-in-the-middle attacks (hijacking all TLS traffic and inspecting content). Double TLS over WebSocket can mitigate man-in-the-middle attacks; see the Advanced Configuration section for details.
