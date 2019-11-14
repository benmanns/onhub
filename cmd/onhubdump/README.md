# onhubdump

onhubdump will download, parse, and print a JSON version of the OnHub
diagnostic report. It's a command from the
[onhub](https://github.com/benmanns/onhub) Go package.

## Installation

```sh
go get github.com/benmanns/onhub/cmd/onhubdump
```

## Running

```sh
onhubdump
```

This returns a JSON dump of the data from `http://192.168.86.1/api/v1/diagnostic-report`.

If you use a different subnet for your router:

```sh
onhubdump http://192.168.85.1/api/v1/diagnostic-report
```

If you want to run against a local, already downloaded dump:

```sh
onhubdump ./path/to/diagnostic-report
```

You will likely want to either save a copy of the diagnostic report or save a
copy of the parsed output from onhubdump rather than re-running the command as
generating the report takes a few seconds each time.

```sh
curl -o diagnostic-report http://192.168.86.1/api/v1/diagnostic-report
onhubdump ./diagnostic-report
```

Alternatively:

```sh
onhubdump > diagnostic-report.json
```

Use the [jq](https://stedolan.github.io/jq/) tool to format and manipulate
output on the command line.

Unbuffered:

```sh
onhubdump | jq
```

Buffered report:

```sh
curl -o diagnostic-report http://192.168.86.1/api/v1/diagnostic-report
onhubdump ./diagnostic-report | jq
```

Buffered output:

```sh
onhubdump > diagnostic-report.json
jq < diagnostic-report.json
```

## Example outputs

Note: these examples are unbuffered, but I recommend you use
`onhubdump ./your-report | ...` or `jq ... < your-report.json` to avoid taxing
your router.

### Get version

```sh
onhubdump | jq '.version'
```

```
"8350.53.0 (Official Build) stable-channel whirlwind"
```

### List file paths

```sh
onhubdump | jq '.files | .[] | .path'
```

```
"/etc/resolv.conf"
"/sys/firmware/log"
"/tmp/debug-log"
"/var/log/boot.log"
"/var/log/net.log"
"/var/log/net.1.log"
"/var/log/update_engine/update_engine.XXXXXXXX-XXXXXX"
"/var/log/webservd/YYYY-MM-DD.log"
"/var/log/webservd/YYYY-MM-DD.log"
"/var/lib/ap/monitor/wan_idle_usage"
"/var/log/messages"
"/var/log/messages.1"
```

### Read raw file contents

```sh
onhubdump | jq -r '.files | .[] | select(.path == "/var/log/messages") | .content'
```

```
2019-11-11T16:28:53.788464+00:00 INFO ap-hal[1184]: [INFO:qca_wifi_phy.cc(268)] Unable to read qos metrics file
2019-11-11T16:28:53.789888+00:00 ERR ap-monitor[1487]: [ERROR:qos_metrics_plugin.cc(247)] Unable to get qos_metrics5
2019-11-11T16:28:53.790436+00:00 INFO ap-monitor[1487]: [INFO:qos_metrics_plugin.cc(202)] Unable to read qos metrics from kernel
2019-11-11T16:28:53.790926+00:00 WARNING ap-monitor[1487]: [WARNING:qos_metrics_plugin.cc(117)] Unable to collect QoS Metrics
...
```

### Output network config without escaping

```sh
onhubdump | jq -r '.networkConfig'
```

```
local_network {
  ip_address: "192.168.86.1"
  netmask: "255.255.255.0"
}
wireless_network {
...
```

### List commands

```sh
onhubdump | jq '.commandOutputs | .[] | .command'
"/bin/ifconfig"
"/usr/sbin/iw dev wlan-2400mhz station dump"
"/usr/sbin/iw dev wlan-5000mhz station dump"
"/usr/sbin/iw dev"
"/sbin/brctl showstp br-lan"
"/usr/sbin/ethtool -S wan0"
"/usr/sbin/ethtool -S lan0"
"/bin/route -n"
"/bin/ps auxwwf"
"/sbin/tc -s qdisc show dev wan0"
"/bin/cat /dev/ecm/ecm_db"
"/usr/sbin/iw dev mesh-5000mhz mpath dump"
"/usr/sbin/iw dev mesh-5000mhz mpp dump"
"/usr/sbin/iw dev mesh-5000mhz station dump"
```

### Get Upload/Download speed

```sh
onhubdump | jq '.infoJSON._apCloudStorage._wanSpeedTestResults | { up: ._uploadSpeedBytesPerSecond, down: ._downloadSpeedBytesPerSecond }'
```

```
{
  "up": 649257,
  "down": 10526829
}
```

### Show clients

```sh
onhubdump | jq '.infoJSON._apState._stations | .[] | { name: ._dhcpHostname, connected: ._connected }'
```

```
{
  "name": "My-Computer",
  "connected": true
}
{
  "name": "My-Phone",
  "connected": true
}
{
  "name": "My-Tablet",
  "connected": false
}
...
```

### Show names of connected clients

```sh
onhubdump | jq -r '.infoJSON._apState._stations | .[] | select(._connected) | ._dhcpHostname'
```

```
My-Computer
My-Phone
...
```
