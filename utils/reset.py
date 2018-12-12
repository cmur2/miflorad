#!/usr/bin/python

# This script allows hard-resetting the Intel Wireless Bluetooth 8265 chip
# (ID 8087:0a2b) built into newer Thinkpads which tends to get stuck as other
# people noticed before: https://bbs.archlinux.org/viewtopic.php?id=193813

# Setup: pip install pyusb
# Usage: sudo python reset.py

from usb.core import find as finddev
dev = finddev(idVendor=0x8087, idProduct=0x0a2b)
dev.reset()
