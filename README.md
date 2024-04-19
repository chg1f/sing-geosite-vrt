# sing-box-mixed

## Intro

- Fork from [@SagerNet/sing-geosite](https://github.com/SagerNet/sing-geosite) and replace to [@Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)
- Add [@Loyalsoldier/clash-rules](https://github.com/Loyalsoldier/clash-rules)

## QuickStart

> ALL GENERATES: https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/.rule_set.txt

```json
{
  ...
  "dns": {
    "servers": [
      {
        "tag": "reject-dns",
        "address": "rcode://refused"
      },
      {
        "tag": "cloudflare-doh",
        "address": "https://1.1.1.1/dns-query",
        "address_resolver": "aliyun-doh",
        "detour": "PROXY"
      },
      {
        "tag": "aliyun-doh",
        "address": "https://223.5.5.5/dns-query",
        "detour": "direct-out"
      }
    ],
    "rules": [
      { "server": "fakeip-dns", "clash_mode": "Global" },
      { "server": "aliyun-doh", "clash_mode": "Direct" },
      {
        "server": "reject-dns",
        "rule_set": ["reject"]
      },
      {
        "server": "cloudflare-doh",
        "rule_set": ["telegramcidr", "google", "proxy"]
      },
      {
        "server": "aliyun-doh",
        "rule_set": [
          "geoip-cn",
          "applications",
          "icloud",
          "apple",
          "direct",
          "lancidr",
          "cncidr"
        ]
      }
    ],
    "final": "cloudflare-doh"
  },
  "route": {
    "rule_set": [
      {
        "tag": "geoip-cn",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/SagerNet/sing-geoip/rule-set/geoip-cn.srs"
      },
      {
        "tag": "reject",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/reject.srs"
      },
      {
        "tag": "icloud",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/icloud.srs"
      },
      {
        "tag": "apple",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/apple.srs"
      },
      {
        "tag": "google",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/google.srs"
      },
      {
        "tag": "proxy",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/proxy.srs"
      },
      {
        "tag": "direct",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/direct.srs"
      },
      {
        "tag": "gfw",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/gfw.srs"
      },
      {
        "tag": "tld-not-cn",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/tld-not-cn.srs"
      },
      {
        "tag": "telegramcidr",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/telegramcidr.srs"
      },
      {
        "tag": "cncidr",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/cncidr.srs"
      },
      {
        "tag": "lancidr",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/lancidr.srs"
      },
      {
        "tag": "applications",
        "type": "remote",
        "download_detour": "PROXY",
        "update_interval": "1d",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/chg1f/sing-geosite-mixed/rule-set/applications.srs"
      }
    ],
    "rules": [
      { "outbound": "direct-out", "ip_is_private": true },
      {
        "outbound": "dns-out",
        "type": "logical",
        "mode": "or",
        "rules": [
          { "port": 53 },
          { "protocol": "dns" },
          { "inbound": ["dns-in"] }
        ]
      },
      {
        "outbound": "reject-out",
        "type": "logical",
        "mode": "or",
        "rules": [{ "port": 853 }, { "protocol": "stun" }]
      },
      { "outbound": "PROXY", "clash_mode": "Global" },
      { "outbound": "direct-out", "clash_mode": "Direct" },
      {
        "outbound": "reject-out",
        "rule_set": ["reject"]
      },
      {
        "outbound": "direct-out",
        "type": "logical",
        "mode": "and",
        "rules": [
          {
            "invert": true,
            "rule_set": ["telegramcidr", "google", "proxy"]
          },
          {
            "rule_set": [
              "geoip-cn",
              "applications",
              "icloud",
              "apple",
              "direct",
              "lancidr",
              "cncidr"
            ]
          }
        ]
      }
    ],
    "final": "PROXY"
  }
}
```
