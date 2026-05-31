---
title: "Correctly Configuring Trojan-Go"
draft: false
weight: 22
---

The following explains how to configure Trojan-Go correctly so that your proxy node leaves no detectable fingerprint.

Before you begin, you will need:

- A server that has not been blocked by the GFW
- A domain name (free services such as `.tk` are fine)
- Trojan-Go, available from the releases page
- A TLS certificate and private key (free certificates are available from Let's Encrypt and similar CAs)

### Server configuration

The goal is to make your server behave identically to a normal HTTPS website.

First, you need an HTTP server. You can set up a local HTTP service with nginx, Apache, Caddy, or similar, or point to someone else's HTTP server. Its role is to serve a completely normal web page when the GFW actively probes your server.

**You must specify this HTTP server's address in `remote_addr` and `remote_port`. `remote_addr` may be an IP or a domain name. Trojan-Go will test whether the HTTP server is working correctly and will refuse to start if it is not.**

Below is a reasonably secure server configuration `server.json`. It assumes you have an HTTP service running locally on port 80 (required — you can also use an external web server such as `"remote_addr": "example.com"`), and optionally an HTTPS service or a static page returning `400 Bad Request` on port 1234 (optional — you may omit the `fallback_port` field).

```json
{
  "run_type": "server",
  "local_addr": "0.0.0.0",
  "local_port": 443,
  "remote_addr": "127.0.0.1",
  "remote_port": 80,
  "password": ["your_awesome_password"],
  "ssl": {
    "cert": "server.crt",
    "key": "server.key",
    "fallback_port": 1234
  }
}
```

This configuration makes Trojan-Go listen on port 443 on all IP addresses (`0.0.0.0`), using `server.crt` and `server.key` for TLS. Use as complex a password as possible, and make sure the `password` fields match on both the client and server. Note that **Trojan-Go will check whether the HTTP server at `http://remote_addr:remote_port` is functioning correctly and will refuse to start if it is not**.

When a client attempts to connect to the Trojan-Go listening port, the following happens:

- If the TLS handshake succeeds but the content is not a valid Trojan protocol header (e.g. it is an HTTP request or a GFW active probe), Trojan-Go proxies the TLS connection to the local HTTP service at `127.0.0.1:80`. From the outside, the Trojan-Go service looks like an HTTPS website.
- If the TLS handshake succeeds, a valid Trojan protocol header is detected, and the password matches, the server parses the client's request and proxies it. Otherwise it falls back to the behavior above.
- If the TLS handshake fails (the remote is not speaking TLS), Trojan-Go redirects the TCP connection to the service at `127.0.0.1:1234` (if `fallback_port` is set), which returns an HTTP `400 Bad Request` page. If `fallback_port` is not set, the connection is dropped immediately. While optional, setting this field is strongly recommended.

You can verify correct behavior by visiting `https://your-domain-name.com` in a browser — it should display a normal HTTPS-protected web page matching the content on port 80. You can also verify `fallback_port` by visiting `http://your-domain-name.com:443`.

You can even use Trojan-Go as your HTTPS server to serve your website. Visitors will browse it normally through Trojan-Go, and it will not interfere with proxy traffic. However, avoid hosting latency-sensitive services on `remote_port` or `fallback_port`, because Trojan-Go intentionally introduces a small delay when it recognizes non-Trojan traffic in order to resist GFW timing-based detection.

Start the server with:

```shell
./trojan-go -config ./server.json
```

### Client configuration

The corresponding client configuration `client.json`:

```json
{
  "run_type": "client",
  "local_addr": "127.0.0.1",
  "local_port": 1080,
  "remote_addr": "your_awesome_server",
  "remote_port": 443,
  "password": ["your_awesome_password"],
  "ssl": {
    "sni": "your-domain-name.com"
  }
}
```

This opens a SOCKS5/HTTP proxy (auto-detected) on local port 1080, connecting to `your_awesome_server:443`. `your_awesome_server` can be either an IP address or a domain name.

If you filled in a domain name in `remote_addr`, the `sni` field can be omitted. If you used an IP address, `sni` must be the domain name for which you requested the certificate (or the Common Name of a self-signed certificate), and the two must match. Note that the `sni` field is transmitted **in plaintext** in the TLS protocol (its purpose is to let the server select the appropriate certificate). The GFW has proven SNI inspection and blocking capabilities, so do not set `sni` to an already-blocked domain (e.g. `google.com`), as this could cause your server to be blocked as well.

Start the client with:

```shell
./trojan-go -config ./client.json
```

For more about configuration options, see the relevant sections in the left navigation bar.
