---
title: "Pluggable Transport Plugin Development"
draft: false
weight: 150
---

Trojan-Go encourages the development of transport plugins to expand the variety of protocols and deepen the strategic toolkit against the GFW.

A transport plugin replaces the TLS layer of the `transport` tunnel, handling encryption and obfuscation in its place.

Plugins communicate with Trojan-Go over TCP sockets and are completely decoupled from Trojan-Go itself, so you can write them in any language and any design pattern. The [SIP003](https://shadowsocks.org/en/spec/Plugin.html) standard is the recommended baseline — plugins conforming to it work with both Trojan-Go and Shadowsocks.

After a plugin is enabled, Trojan-Go passes **plaintext TCP traffic** to the plugin. The plugin only needs to handle inbound TCP and may re-encode it into any format it wishes — QUIC, HTTP, even ICMP.

Trojan-Go plugin design principles differ slightly from Shadowsocks:

1. **The plugin itself handles encryption, obfuscation, and integrity verification, and must resist replay attacks.**

2. **The plugin should impersonate an existing, widely-used service X and embed encrypted content within that service's traffic.**

3. **When the server-side plugin detects tampered or replayed content, it must hand the connection to Trojan-Go rather than dropping it.** Specifically, it forwards all already-read and unread bytes to Trojan-Go and creates a bidirectional bridge. Trojan-Go then connects to a real X server and lets the attacker interact with it directly.

Rationale:

- Principle 1: the Trojan protocol has no encryption of its own. Replacing TLS with a plugin means **fully trusting the plugin's security**.
- Principle 2: inherits Trojan's philosophy — the best place to hide a tree is a forest.
- Principle 3: fully leverages Trojan-Go's active-probing resistance. Even if the GFW probes your server, it will behave identically to service X.

### Example

1. Your plugin impersonates MySQL traffic. The GFW notices an abnormally high-volume "MySQL" connection and actively connects to probe it.
2. The GFW sends a probe payload. The server-side plugin validates it, finds it is not proxy traffic, and hands the connection to Trojan-Go.
3. Trojan-Go detects the anomaly and redirects it to a real MySQL server. The GFW then interacts with a genuine MySQL server and cannot distinguish your host from a real one.

Even if your protocol does not fully satisfy principles 2 and 3, or even principle 1, development is still encouraged. The GFW only audits and blocks popular, well-known protocols. As long as a custom protocol is not publicly published, it tends to have excellent longevity.
