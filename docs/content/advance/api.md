---
title: "Dynamic User Management via API"
draft: false
weight: 10
---

### Note: the original Trojan does not support this feature

Trojan-Go exposes a gRPC-based API that supports:

- User add / remove / query / modify
- Traffic statistics
- Speed statistics
- IP connection count statistics

Trojan-Go also ships with a built-in API client, so one Trojan-Go instance can control another.

### Setup

Add an `api` section to the server configuration:

```json
{
  "api": {
    "enabled": true,
    "api_addr": "127.0.0.1",
    "api_port": 10000
  }
}
```

Start the Trojan-Go server:

```shell
./trojan-go -config ./server.json
```

You can then manage it from another Trojan-Go instance:

```shell
./trojan-go -api-addr SERVER_API_ADDRESS -api COMMAND
```

`SERVER_API_ADDRESS` is the API address and port (e.g. `127.0.0.1:10000`).

Available `COMMAND` values:

| Command | Description                       |
| ------- | --------------------------------- |
| `list`  | List all users                    |
| `get`   | Get a specific user's information |
| `set`   | Add, remove, or modify a user     |

### Examples

**1. List all users**

```shell
./trojan-go -api-addr 127.0.0.1:10000 -api list
```

All user records are output as JSON. Example response:

```json
[
  {
    "user": {
      "hash": "d63dc919e201d7bc4c825630d2cf25fdc93d4b2f0d46706d29038d01"
    },
    "status": {
      "traffic_total": { "upload_traffic": 36393, "download_traffic": 186478 },
      "speed_current": { "upload_speed": 25210, "download_speed": 72384 },
      "speed_limit": { "upload_speed": 5242880, "download_speed": 5242880 },
      "ip_limit": 50
    }
  }
]
```

All traffic values are in bytes.

**2. Get a single user**

You can identify the target with `-target-password` (plaintext) or `-target-hash` (SHA-224 hex):

```shell
./trojan-go -api-addr 127.0.0.1:10000 -api get -target-password password
```

```shell
./trojan-go -api-addr 127.0.0.1:10000 -api get -target-hash d63dc919e201d7bc4c825630d2cf25fdc93d4b2f0d46706d29038d01
```

**3. Add a user**

```shell
./trojan-go -api-addr 127.0.0.1:10000 -api set -add-profile -target-password password
```

**4. Remove a user**

```shell
./trojan-go -api-addr 127.0.0.1:10000 -api set -delete-profile -target-password password
```

**5. Modify a user**

```shell
./trojan-go -api-addr 127.0.0.1:10000 -api set -modify-profile -target-password password \
    -ip-limit 3 \
    -upload-speed-limit 5242880 \
    -download-speed-limit 5242880
```

This limits the user with password `password` to 5 MiB/s upload and download, and a maximum of 3 simultaneous IP connections. Values of 0 or negative remove the corresponding limit.
