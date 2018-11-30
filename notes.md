INFO[0036] name=Flower care addr=C4:7C:8D:66:D5:27 rssi=-43

list-attributes C4:7C:8D:66:D5:27
connect C4:7C:8D:66:D5:27
select-attribute /org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0002/char0003
select-attribute 00002a00-0000-1000-8000-00805f9b34fb
read 0003

08:38:52::$ gatttool --device=C4:7C:8D:66:D5:27 --char-read -a 0x03
Characteristic value/descriptor: 46 6c 6f 77 65 72 20 63 61 72 65


08:53:32::$ sudo hcitool leinfo C4:7C:8D:66:D5:27
Requesting information ...
	Handle: 3585 (0x0e01)
	LMP Version: 4.0 (0x6) LMP Subversion: 0x706
	Manufacturer: RivieraWaves S.A.S (96)
	Features: 0x01 0x00 0x00 0x00 0x00 0x00 0x00 0x00


Attempting to connect to C4:7C:8D:66:D5:27
[CHG] Device C4:7C:8D:66:D5:27 Connected: yes
Connection successful
[NEW] Primary Service
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service000c
	00001801-0000-1000-8000-00805f9b34fb
	Generic Attribute Profile
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service000c/char000d
	00002a05-0000-1000-8000-00805f9b34fb
	Service Changed
[NEW] Descriptor
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service000c/char000d/desc000f
	00002902-0000-1000-8000-00805f9b34fb
	Client Characteristic Configuration
[NEW] Primary Service
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0010
	0000fe95-0000-1000-8000-00805f9b34fb
	Xiaomi Inc.
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0010/char0011
	00000001-0000-1000-8000-00805f9b34fb
	SDP
[NEW] Descriptor
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0010/char0011/desc0013
	00002902-0000-1000-8000-00805f9b34fb
	Client Characteristic Configuration
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0010/char0014
	00000002-0000-1000-8000-00805f9b34fb
	Unknown
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0010/char0016
	00000004-0000-1000-8000-00805f9b34fb
	Unknown
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0010/char0018
	00000007-0000-1000-8000-00805f9b34fb
	ATT
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0010/char001a
	00000010-0000-1000-8000-00805f9b34fb
	UPNP
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0010/char001c
	00000013-0000-1000-8000-00805f9b34fb
	Unknown
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0010/char001e
	00000014-0000-1000-8000-00805f9b34fb
	Hardcopy Data Channel
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0010/char0020
	00001001-0000-1000-8000-00805f9b34fb
	Browse Group Descriptor Service Class
[NEW] Descriptor
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0010/char0020/desc0022
	00002902-0000-1000-8000-00805f9b34fb
	Client Characteristic Configuration
[NEW] Primary Service
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0023
	0000fef5-0000-1000-8000-00805f9b34fb
	Dialog Semiconductor GmbH
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0023/char0024
	8082caa8-41a6-4021-91c6-56f9b954cc34
	Vendor specific
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0023/char0026
	724249f0-5ec3-4b5f-8804-42345af08651
	Vendor specific
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0023/char0028
	6c53db25-47a1-45fe-a022-7c92fb334fd4
	Vendor specific
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0023/char002a
	9d84b9a3-000c-49d8-9183-855b673fda31
	Vendor specific
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0023/char002c
	457871e8-d516-4ca1-9116-57d0b17b9cb2
	Vendor specific
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0023/char002e
	5f78df94-798c-46f5-990a-b3eb6a065c88
	Vendor specific
[NEW] Descriptor
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0023/char002e/desc0030
	00002902-0000-1000-8000-00805f9b34fb
	Client Characteristic Configuration
[NEW] Primary Service
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0031
	00001204-0000-1000-8000-00805f9b34fb
	Generic Telephony
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0031/char0032
	00001a00-0000-1000-8000-00805f9b34fb
	Unknown
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0031/char0034
	00001a01-0000-1000-8000-00805f9b34fb
	Unknown
[NEW] Descriptor
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0031/char0034/desc0036
	00002902-0000-1000-8000-00805f9b34fb
	Client Characteristic Configuration
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0031/char0037
	00001a02-0000-1000-8000-00805f9b34fb
	Unknown
[NEW] Descriptor
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service0031/char0037/desc0039
	00001a02-0000-1000-8000-00805f9b34fb
	Unknown
[NEW] Primary Service
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service003a
	00001206-0000-1000-8000-00805f9b34fb
	UPNP IP Service
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service003a/char003b
	00001a11-0000-1000-8000-00805f9b34fb
	Unknown
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service003a/char003d
	00001a10-0000-1000-8000-00805f9b34fb
	Unknown
[NEW] Descriptor
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service003a/char003d/desc003f
	00002902-0000-1000-8000-00805f9b34fb
	Client Characteristic Configuration
[NEW] Characteristic
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service003a/char0040
	00001a12-0000-1000-8000-00805f9b34fb
	Unknown
[NEW] Descriptor
	/org/bluez/hci0/dev_C4_7C_8D_66_D5_27/service003a/char0040/desc0042
	00001a12-0000-1000-8000-00805f9b34fb
	Unknown
[CHG] Device C4:7C:8D:66:D5:27 ServicesResolved: yes
[CHG] Device C4:7C:8D:66:D5:27 ServicesResolved: no
[CHG] Device C4:7C:8D:66:D5:27 Connected: no

https://github.com/open-homeautomation/miflora/blob/master/miflora/miflora_poller.py#L11
https://github.com/golang/go/wiki/Modules
https://github.com/currantlabs/gatt/blob/master/examples/explorer.go

$ gatttool -b C4:7C:8D:66:D5:27 --characteristics
handle = 0x0002, char properties = 0x02, char value handle = 0x0003, uuid = 00002a00-0000-1000-8000-00805f9b34fb
handle = 0x0004, char properties = 0x02, char value handle = 0x0005, uuid = 00002a01-0000-1000-8000-00805f9b34fb
handle = 0x0006, char properties = 0x0a, char value handle = 0x0007, uuid = 00002a02-0000-1000-8000-00805f9b34fb
handle = 0x0008, char properties = 0x02, char value handle = 0x0009, uuid = 00002a04-0000-1000-8000-00805f9b34fb
handle = 0x000d, char properties = 0x22, char value handle = 0x000e, uuid = 00002a05-0000-1000-8000-00805f9b34fb
handle = 0x0011, char properties = 0x1a, char value handle = 0x0012, uuid = 00000001-0000-1000-8000-00805f9b34fb
handle = 0x0014, char properties = 0x02, char value handle = 0x0015, uuid = 00000002-0000-1000-8000-00805f9b34fb
handle = 0x0016, char properties = 0x12, char value handle = 0x0017, uuid = 00000004-0000-1000-8000-00805f9b34fb
handle = 0x0018, char properties = 0x08, char value handle = 0x0019, uuid = 00000007-0000-1000-8000-00805f9b34fb
handle = 0x001a, char properties = 0x08, char value handle = 0x001b, uuid = 00000010-0000-1000-8000-00805f9b34fb
handle = 0x001c, char properties = 0x0a, char value handle = 0x001d, uuid = 00000013-0000-1000-8000-00805f9b34fb
handle = 0x001e, char properties = 0x02, char value handle = 0x001f, uuid = 00000014-0000-1000-8000-00805f9b34fb
handle = 0x0020, char properties = 0x10, char value handle = 0x0021, uuid = 00001001-0000-1000-8000-00805f9b34fb
handle = 0x0024, char properties = 0x0a, char value handle = 0x0025, uuid = 8082caa8-41a6-4021-91c6-56f9b954cc34
handle = 0x0026, char properties = 0x0a, char value handle = 0x0027, uuid = 724249f0-5ec3-4b5f-8804-42345af08651
handle = 0x0028, char properties = 0x02, char value handle = 0x0029, uuid = 6c53db25-47a1-45fe-a022-7c92fb334fd4
handle = 0x002a, char properties = 0x0a, char value handle = 0x002b, uuid = 9d84b9a3-000c-49d8-9183-855b673fda31
handle = 0x002c, char properties = 0x0e, char value handle = 0x002d, uuid = 457871e8-d516-4ca1-9116-57d0b17b9cb2
handle = 0x002e, char properties = 0x12, char value handle = 0x002f, uuid = 5f78df94-798c-46f5-990a-b3eb6a065c88
handle = 0x0032, char properties = 0x0a, char value handle = 0x0033, uuid = 00001a00-0000-1000-8000-00805f9b34fb
handle = 0x0034, char properties = 0x1a, char value handle = 0x0035, uuid = 00001a01-0000-1000-8000-00805f9b34fb
handle = 0x0037, char properties = 0x02, char value handle = 0x0038, uuid = 00001a02-0000-1000-8000-00805f9b34fb
handle = 0x003b, char properties = 0x02, char value handle = 0x003c, uuid = 00001a11-0000-1000-8000-00805f9b34fb
handle = 0x003d, char properties = 0x1a, char value handle = 0x003e, uuid = 00001a10-0000-1000-8000-00805f9b34fb
handle = 0x0040, char properties = 0x02, char value handle = 0x0041, uuid = 00001a12-0000-1000-8000-00805f9b34fb

Service: 0000120400001000800000805f9b34fb
  Characteristic  00001a0000001000800000805f9b34fb
    properties    read write
    value         0000 | "\x00\x00"
  Characteristic  00001a0100001000800000805f9b34fb
    properties    read write notify
    value         aabbccddeeff99887766000000000000 | "\xaa\xbb\xcc\xdd\xee\xff\x99\x88wf\x00\x00\x00\x00\x00\x00"
  Descriptor      2902 (Client Characteristic Configuration)
    value         0000 | "\x00\x00"
  Characteristic  00001a0200001000800000805f9b34fb
    properties    read
    value         6415322e372e30 | "d\x152.7.0"
  Descriptor      00001a0200001000800000805f9b34fb
    value         6415322e372e30 | "d\x152.7.0"
