---
title: "Using Shadowsocks Plugins / Pluggable Transports"
draft: false
weight: 7
---

### Note: the original Trojan does not support this feature

Trojan-Go supports a pluggable transport layer. In principle, any software with TCP tunneling capability can be used — V2Ray, Shadowsocks, KCP, and so on. Trojan-Go is also compatible with the Shadowsocks SIP003 plugin standard (GoQuiet, v2ray-plugin, etc.) and with Tor pluggable transports such as obfs4 and meek.

When a pluggable transport is enabled, the Trojan-Go **client sends plaintext traffic** to the local plugin, which is responsible for encryption, obfuscation, and delivery to the server-side plugin. The server-side plugin decrypts and decapsulates the traffic, then delivers the **plaintext** to the local Trojan-Go server process.

To use any plugin, add a `transport_plugin` section, specify the plugin executable path, and configure it accordingly.

**Writing your own plugin and protocol is encouraged**, because all existing plugins lack full integration with Trojan-Go's active-probing resistance, and some have no encryption capability at all. See the "Pluggable Transport Plugin Development" page in the Developer Guide if you are interested.

### Example: v2ray-plugin (SIP003)

> **Security warning:** The configuration below transmits unencrypted Trojan protocol over WebSocket. It is for demonstration purposes only. Never use this to bypass the GFW.

**Server:**

```json
"transport_plugin": {
    "enabled": true,
    "type": "shadowsocks",
    "command": "./v2ray-plugin",
    "arg": ["-server", "-host", "www.example.com"]
}
```

**Client:**

```json
"transport_plugin": {
    "enabled": true,
    "type": "shadowsocks",
    "command": "./v2ray-plugin",
    "arg": ["-host", "www.example.com"]
}
```

Note that v2ray-plugin requires a `-server` flag to distinguish server mode from client mode. Refer to v2ray-plugin's documentation for its full options.

After Trojan-Go starts, you will see v2ray-plugin's startup output. The plugin wraps traffic as WebSocket and forwards it.

For non-SIP003 plugins, set `type` to `"other"` and configure `command`, `arg`, and `env` manually.
