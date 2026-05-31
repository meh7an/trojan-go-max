---
title: "Transparent Proxy"
draft: false
weight: 11
---

### Note: the original Trojan does not fully support this feature (UDP)

Trojan-Go supports transparent TCP/UDP proxying via `tproxy`.

To enable transparent proxy mode, take a valid client configuration (see the Basic Configuration section) and change `run_type` to `nat`, then adjust the local listening port as needed.

Next, add the iptables rules below. This example assumes the gateway has two network interfaces: one facing the LAN and one facing the internet. LAN inbound packets are handed to Trojan-Go, which forwards them through the TLS tunnel to the remote Trojan-Go server via the internet interface. Replace `$SERVER_IP`, `$TROJAN_GO_PORT`, and `$INTERFACE` with your own values.

```shell
# Create the TROJAN_GO chain
iptables -t mangle -N TROJAN_GO

# Bypass the Trojan-Go server address
iptables -t mangle -A TROJAN_GO -d $SERVER_IP -j RETURN

# Bypass private and reserved addresses
iptables -t mangle -A TROJAN_GO -d 0.0.0.0/8 -j RETURN
iptables -t mangle -A TROJAN_GO -d 10.0.0.0/8 -j RETURN
iptables -t mangle -A TROJAN_GO -d 127.0.0.0/8 -j RETURN
iptables -t mangle -A TROJAN_GO -d 169.254.0.0/16 -j RETURN
iptables -t mangle -A TROJAN_GO -d 172.16.0.0/12 -j RETURN
iptables -t mangle -A TROJAN_GO -d 192.168.0.0/16 -j RETURN
iptables -t mangle -A TROJAN_GO -d 224.0.0.0/4 -j RETURN
iptables -t mangle -A TROJAN_GO -d 240.0.0.0/4 -j RETURN

# Mark packets that did not match the rules above
iptables -t mangle -A TROJAN_GO -j TPROXY -p tcp --on-port $TROJAN_GO_PORT --tproxy-mark 0x01/0x01
iptables -t mangle -A TROJAN_GO -j TPROXY -p udp --on-port $TROJAN_GO_PORT --tproxy-mark 0x01/0x01

# Redirect all TCP/UDP packets arriving on $INTERFACE to the TROJAN_GO chain
iptables -t mangle -A PREROUTING -p tcp -i $INTERFACE -j TROJAN_GO
iptables -t mangle -A PREROUTING -p udp -i $INTERFACE -j TROJAN_GO

# Route marked packets back through the loopback interface
ip route add local default dev lo table 100
ip rule add fwmark 1 lookup 100
```

After the rules are in place, **start Trojan-Go with root privileges**:

```shell
sudo trojan-go
```
