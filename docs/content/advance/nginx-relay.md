---
title: "SNI-Based Multi-Path Relay with nginx"
draft: false
weight: 6
---

## Overview

Trojan encapsulates traffic inside TLS. By exploiting the SNI field in TLS, it is possible to use an SNI proxy to route different connection paths through a single port on a relay host.

## Requirements

- **Relay host:** nginx 1.11.5 or later
- **Endpoint host:** any Trojan server (no version requirement)

## Setup

The example below uses two relay hosts and two endpoint hosts. Their domains are `a.example.com`, `b.example.com`, `c.example.com`, and `d.example.com`. There are four connection paths: a→c, a→d, b→c, and b→d.

```text
                    +-----------------+           +--------------------+
                    |                 +---------->+                    |
                    |   VPS RELAY A   |           |   VPS ENDPOINT C   |
              +---->+                 |   +------>+                    |
              |     |  a.example.com  |   |       |   c.example.com    |
              |     |                 +------+    |                    |
+----------+  |     +-----------------+   |  |    +--------------------+
|          |  |                           |  |
|  client  +--+                           |  |
|          |  |                           |  |
+----------+  |     +-----------------+   |  |    +--------------------+
              |     |                 |   |  |    |                    |
              |     |   VPS RELAY B   |   |  +--->+   VPS ENDPOINT D   |
              +---->+                 +---+       |                    |
                    |  b.example.com  |           |   d.example.com    |
                    |                 +---------->+                    |
                    +-----------------+           +--------------------+
```

### Step 1 — Assign path domains and certificates

Assign a domain to each path and resolve them to the appropriate relay host:

```text
a-c.example.com  CNAME  a.example.com
a-d.example.com  CNAME  a.example.com
b-c.example.com  CNAME  b.example.com
b-d.example.com  CNAME  b.example.com
```

Deploy certificates for all target-path domains on the endpoint hosts. Because the DNS records do not point to the endpoint hosts directly, HTTP-based ACME validation will fail — use DNS validation instead. The example below uses AWS Route 53:

```shell
# On host C
certbot certonly --dns-route53 -d a-c.example.com -d b-c.example.com

# On host D
certbot certonly --dns-route53 -d a-d.example.com -d b-d.example.com
```

### Step 2 — Configure the SNI proxy

Use nginx's `ssl_preread` module for SNI proxying. Add the following to `nginx.conf` (outside the `http` block — this is a raw TCP stream service):

Configuration for relay host A (relay host B follows the same pattern):

```nginx
stream {
  map $ssl_preread_server_name $name {
    a-c.example.com   c.example.com;   # route a→c traffic to host C
    a-d.example.com   d.example.com;   # route a→d traffic to host D

    # If this host also runs other services on port 443 (e.g. a web server
    # or Trojan service), have them listen on a different local port (4000
    # here). All TLS requests that do not match an SNI rule above are
    # forwarded to this port. Remove this line if not needed.
    default           localhost:4000;
  }

  server {
    listen      443;
    proxy_pass  $name;
    ssl_preread on;
  }
}
```

### Step 3 — Configure the Trojan service on endpoint hosts

Since a single certificate was issued for all target-path domains, one Trojan server instance can handle all paths. Configuration follows the standard method; irrelevant fields are omitted:

```json
{
  "run_type": "server",
  "local_addr": "0.0.0.0",
  "local_port": 443,
  "ssl": {
    "cert": "/path/to/certificate.crt",
    "key": "/path/to/private.key"
  }
}
```

> **Tip:** If you need separate Trojan server instances for different paths on the same endpoint host (for example, to connect to separate billing services), configure an additional SNI proxy on the endpoint host and route to different local Trojan ports. The process is the same as described above.

## Summary

With this setup, multi-entry, multi-exit, multi-hop Trojan traffic forwarding is possible on a single port. For additional relay hops, apply the same SNI proxy configuration on each intermediate node.
