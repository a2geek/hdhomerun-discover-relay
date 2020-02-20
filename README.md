# HDHomeRun discover relay

## Purpose

Enable products that use HDHomeRun to run within a virtual machine on a different subnet than the physical device resides.

The discovery mechanism broadcasts on UDP port 65001 to 255.255.255.255. This tool catches that packet, and sends it to the HDHomeRun.
