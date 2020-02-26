# HDHomeRun discover relay

[![Go Report Card](https://goreportcard.com/badge/github.com/a2geek/hdhomerun-discover-relay)](https://goreportcard.com/report/github.com/a2geek/hdhomerun-discover-relay)
[![GitHub release](https://img.shields.io/github/v/release/a2geek/hdhomerun-discover-relay)](https://github.com/a2geek/hdhomerun-discover-relay/releases/latest)

## Purpose

Enable products that use HDHomeRun to run within a virtual machine on a different subnet than the physical device resides.

The discovery mechanism broadcasts on UDP port 65001 to 255.255.255.255. This tool catches that packet, and sends it to the HDHomeRun.

Note that this does not allow all packets back into the VM, but [this advice](https://www.mythtv.org/wiki/Silicondust_HDHomeRun_Dual#Can.27t_Connect_to_HDHR.3F) seems to help with that.

## Compile

```
$ go build -o hdr *.go
```

## Sample run

(it's a bit messy right now)

```
$ sudo ./hdr relay 192.168.123.0/24
Source CIDR = 192.168.123.0/24
HDHomeRun IP(s): [192.168.5.209 192.168.5.117]
Starting...
packet #0, header=ver=4 hdrlen=20 tos=0x0 totallen=48 id=0x49f9 flags=0x2 fragoff=0x0 ttl=64 proto=17 cksum=0xb510 src=192.168.123.11 dst=255.255.255.255, control=ttl=0 src=192.168.123.11 dst=255.255.255.255 ifindex=3, udp=[146 154 253 233 0 28 59 225 0 2 0 12 1 4 255 255 255 255 2 4 255 255 255 255 115 204 125 143]
Payload:
00000000  92 9a fd e9 00 1c 3b e1  00 02 00 0c 01 04 ff ff  |......;.........|
00000010  ff ff 02 04 ff ff ff ff  73 cc 7d 8f              |........s.}.|
Data:
00000000  45 00 00 30 49 f9 40 00  40 11 b5 10 c0 a8 7b 0b  |E..0I.@.@.....{.|
00000010  ff ff ff ff 92 9a fd e9  00 1c 3b e1 00 02 00 0c  |..........;.....|
00000020  01 04 ff ff ff ff 02 04  ff ff ff ff 73 cc 7d 8f  |............s.}.|

UDP Payload:
00000000  92 9a fd e9 00 1c 3b e1  00 02 00 0c 01 04 ff ff  |......;.........|
00000010  ff ff 02 04 ff ff ff ff  73 cc 7d 8f              |........s.}.|
Redirecting to 192.168.5.117
New Packet:
00000000  45 00 00 30 49 f9 40 00  40 11 b5 10 c0 a8 7b 0b  |E..0I.@.@.....{.|
00000010  c0 a8 05 75 92 9a fd e9  00 1c 3b e1 00 02 00 0c  |...u......;.....|
00000020  01 04 ff ff ff ff 02 04  ff ff ff ff 73 cc 7d 8f  |............s.}.|
```

## Help

```
$ ./hdr --help
Usage:
  main [OPTIONS] <discover | relay>

Help Options:
  -h, --help  Show this help message

Available commands:
  discover  Test HDHomeRun discovery mechanism (aliases: d)
  relay     Relay HDHomeRun discovery packets (aliases: r)
```

```
$ ./hdr discover --help
Usage:
  main [OPTIONS] discover

Help Options:
  -h, --help      Show this help message
```

```
$ ./hdr relay --help
Usage:
  hdr [OPTIONS] relay cidr

Help Options:
  -h, --help      Show this help message

[relay command arguments]
  cidr:           Source CIDR for application looking for HDHomeRun
```

## Resources

SiliconDust:
* https://www.silicondust.com/support/linux/
* https://github.com/Silicondust/libhdhomerun (20110518)
* https://info.hdhomerun.com/info/hdhomerun_config

Related information/discussions/applications:
* [Howto: HDHomerun discovery on different LAN segment](https://community.ui.com/questions/Howto-HDHomerun-discovery-on-different-LAN-segment/97db52c6-4add-4ba1-ab0d-27ee6f43db8f)
* [hdhomerun cli tools](https://github.com/patrickshuff/hdhomerun)
* [node hdhomerun](https://github.com/mharsch/node-hdhomerun)
* [how to specify HDHomerun IP](https://tvheadend.org/boards/5/topics/33352?r=33363)
* [udp-broadcast-relay](https://github.com/nomeata/udp-broadcast-relay)
