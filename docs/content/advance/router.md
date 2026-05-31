---
title: "Domestic Direct Connect and Ad Blocking"
draft: false
weight: 3
---

### Note: the original Trojan does not support this feature

Trojan-Go's built-in router module enables direct connections to domestic websites — traffic destined for Chinese mainland IPs and domains bypasses the proxy entirely.

The router can be configured with three policies (`bypass`, `proxy`, `block`) on the client side. On the server side only the `block` policy is available.

Example configuration:

```json
{
  "run_type": "client",
  "local_addr": "127.0.0.1",
  "local_port": 1080,
  "remote_addr": "your_server",
  "remote_port": 443,
  "password": ["your_password"],
  "ssl": {
    "sni": "your-domain-name.com"
  },
  "mux": {
    "enabled": true
  },
  "router": {
    "enabled": true,
    "bypass": [
      "geoip:cn",
      "geoip:private",
      "geosite:cn",
      "geosite:geolocation-cn"
    ],
    "block": ["geosite:category-ads"],
    "proxy": ["geosite:geolocation-!cn"]
  }
}
```

This activates the router in whitelist mode: connections to mainland China or private network IPs/domains connect directly, ad network domains are blocked, and everything else is proxied.

The required databases `geoip.dat` and `geosite.dat` are included in the release archives. They come from V2Ray's [domain-list-community](https://github.com/v2fly/domain-list-community) and [geoip](https://github.com/v2fly/geoip) projects.

**GeoSite tags** use the form `geosite:<tag>`, for example `geosite:cn`, `geosite:geolocation-!cn`, `geosite:category-ads-all`, `geosite:bilibili`. All available tags are in the [`data`](https://github.com/v2fly/domain-list-community/tree/master/data) directory of the domain-list-community repository. For detailed usage see [V2Ray Routing — Predefined Domain Lists](https://www.v2fly.org/config/routing.html).

**GeoIP tags** use the form `geoip:<country_code>`, for example `geoip:cn`, `geoip:hk`, `geoip:us`, `geoip:private`. `geoip:private` is a special entry covering private/reserved IP ranges. Other entries cover the IP address blocks of individual countries and regions; see [Wikipedia country codes](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) for codes.

You can also write custom routing rules. For example, to block all traffic to `example.com` and its subdomains, and the subnet `192.168.1.0/24`:

```json
"block": [
    "domain:example.com",
    "cidr:192.168.1.0/24"
]
```

Supported rule prefixes:

| Prefix    | Match type               |
| --------- | ------------------------ |
| `domain:` | Subdomain match          |
| `full:`   | Exact domain match       |
| `regexp:` | Regular expression match |
| `cidr:`   | CIDR match               |

See the "Complete Configuration Reference" section for full details.
