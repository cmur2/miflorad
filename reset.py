#!/usr/bin/python

# Note: pip install pyusb

from usb.core import find as finddev
dev = finddev(idVendor=0x8087, idProduct=0x0a2b)
dev.reset()
