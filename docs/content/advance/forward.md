---
title: "Tunnels and Reverse Proxies"
draft: false
weight: 5
---

You can use Trojan-Go to build tunnels. A typical use case is creating a local, pollution-free DNS resolver:

```json
{
  "run_type": "forward",
  "local_addr": "127.0.0.1",
  "local_port": 53,
  "remote_addr": "your_awesome_server",
  "remote_port": 443,
  "target_addr": "8.8.8.8",
  "target_port": 53,
  "password": ["your_awesome_password"]
}
```

A `forward` mode instance is essentially a client, but requires `target_addr` and `target_port` to specify the tunnel's destination.

With this configuration, local TCP and UDP port 53 are monitored. All traffic sent to local port 53 is forwarded through the TLS tunnel to the remote server, which then connects to `8.8.8.8:53`, and the response travels back through the tunnel. You can point `127.0.0.1` at your DNS resolver and get unpolluted results identical to those from the remote server.

Using the same principle, you can build a local Google mirror:

```json
{
  "run_type": "forward",
  "local_addr": "127.0.0.1",
  "local_port": 443,
  "remote_addr": "your_awesome_server",
  "remote_port": 443,
  "target_addr": "www.google.com",
  "target_port": 443,
  "password": ["your_awesome_password"]
}
```

Visiting `https://127.0.0.1` will show the Google homepage. Note that the browser will display a certificate warning because Google's certificate is issued for `google.com`, not `127.0.0.1`.

Similarly, you can use `forward` to tunnel other proxy protocols. For example, to carry Shadowsocks traffic over Trojan-Go when the remote host is running a Shadowsocks server on `127.0.0.1:12345` alongside a Trojan-Go server on port 443:

```json
{
  "run_type": "forward",
  "local_addr": "0.0.0.0",
  "local_port": 54321,
  "remote_addr": "your_awesome_server",
  "remote_port": 443,
  "target_addr": "127.0.0.1",
  "target_port": 12345,
  "password": ["your_awesome_password"]
}
```

Any TCP/UDP connection to local port 54321 is equivalent to connecting to remote port 12345. A Shadowsocks client can connect to `localhost:54321`, and its traffic will travel through the Trojan tunnel to the Shadowsocks server at `remote:12345`.
