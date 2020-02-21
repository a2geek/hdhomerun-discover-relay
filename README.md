# HDHomeRun discover relay

## Purpose

Enable products that use HDHomeRun to run within a virtual machine on a different subnet than the physical device resides.

The discovery mechanism broadcasts on UDP port 65001 to 255.255.255.255. This tool catches that packet, and sends it to the HDHomeRun.

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
