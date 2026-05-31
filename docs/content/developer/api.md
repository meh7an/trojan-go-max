---
title: "API Development"
draft: false
weight: 100
---

Trojan-Go implements its API over gRPC, exchanging data via Protocol Buffers. The client API provides traffic and speed information; the server API provides per-user traffic, speed, and online status, with the ability to dynamically add/remove users and set speed limits.

Enable the API module by adding an `api` block to the configuration file:

```json
"api": {
    "enabled": true,
    "api_addr": "0.0.0.0",
    "api_port": 10000,
    "ssl": {
        "enabled": true,
        "cert": "api_cert.crt",
        "key": "api_key.key",
        "verify_client": true,
        "client_cert": [
            "api_client_cert1.crt",
            "api_client_cert2.crt"
        ]
    }
}
```

See the "Complete Configuration Reference" for field descriptions.

To implement an API client, refer to `api/service/api.proto`.
