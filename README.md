# Go-Stand-Up-Pls

Golang cli to control TiMotion Desks via Bluetooth

Works with [TWD1 series Bluetooth adapter](https://www.timotion.com/en/products/accessories/twd1-series).

## Installation

```bash
go install github.com/2easy/go-Stand-Up-Pls@latest
```

## Usage

You can scan Bluetooth devices to discover their addresses by running:

```bash
$ go-Stand-Up-Pls scan
2024/07/12 20:52:58 INFO Scanning BLE devices
2024/07/12 20:52:58 INFO found device address=24eb8b22-88d6-4c1d-8bdb-2829eed16868 rssi=-95 local_name=""
2024/07/12 20:52:59 INFO found device address=3daf1e2e-8b2e-4285-81df-5fd7da72c5c1 rssi=-82 local_name="stand UP- 0829"
2024/07/12 20:52:59 INFO found device address=3daf1e2e-8b2e-4285-81df-5fd7da72c5c1 rssi=-82 local_name="stand UP- 0829"
2024/07/12 20:52:59 INFO found device address=a51d418d-88cb-4b81-b419-21d3b4d42890 rssi=-94 local_name=""
2024/07/12 20:52:59 INFO found device address=0fcfb86f-58cf-bc01-1e8b-95e02d2da82a rssi=-95 local_name=""
2024/07/12 20:52:59 INFO found device address=7c8ed629-9aba-4775-a53e-a802d1910f68 rssi=-90 local_name=""
2024/07/12 20:52:59 INFO found device address=12e9528f-0370-4d59-af4f-76f06bc07e10 rssi=-84 local_name="STANMORE II"
2024/07/12 20:52:59 INFO found device address=ef203c08-a1fe-8db3-e01c-30dbd1ceb468 rssi=-90 local_name=""
2024/07/12 20:52:59 INFO found device address=7d506273-3054-4d78-db4b-abbb9c388dbb rssi=-84 local_name=""
2024/07/12 20:52:59 INFO found device address=4e06a0f2-0c04-9f51-f4f0-66f68d7be805 rssi=-82 local_name=""
2024/07/12 20:52:59 INFO found device address=4e06a0f2-0c04-9f51-f4f0-66f68d7be805 rssi=-83 local_name=""
2024/07/12 20:52:59 INFO found device address=3daf1e2e-8b2e-4285-81df-5fd7da72c5c1 rssi=-65 local_name="stand UP- 0829"
2024/07/12 20:52:59 INFO found device address=3daf1e2e-8b2e-4285-81df-5fd7da72c5c1 rssi=-66 local_name="stand UP- 0829"
2024/07/12 20:52:59 INFO found device address=6625dd97-8bc6-ad6f-29f2-9597d2910b0e rssi=-66 local_name=""
2024/07/12 20:52:59 INFO found device address=7352fd28-aaa5-4480-a25e-6031e0d6e07f rssi=-89 local_name="STANMORE II"
```

Your desk should be named similar to `stand UP- XXXX`.

You can set your schedule using the address found out above to set a list of positions to cycle
through:

```bash
$ go-Stand-Up-Pls cycle --desk-address 3daf1e2e-8b2e-4285-81df-5fd7da72c5c1 --position 105:15 --position 65:25 --repeat 2
time=2024-07-12T21:03:06.996+02:00 level=INFO msg="Scanning BLE devices"
time=2024-07-12T21:03:07.034+02:00 level=INFO msg="found device" address=ef203c08-a1fe-8db3-e01c-30dbd1ceb468 rssi=-90 local_name=""
time=2024-07-12T21:03:07.034+02:00 level=INFO msg="found device" address=743319bf-a765-f082-a0c6-411020af9ec8 rssi=-90 local_name=""
time=2024-07-12T21:03:07.048+02:00 level=INFO msg="found device" address=841d56b8-bf96-9724-3be7-3fc070752665 rssi=-85 local_name=""
time=2024-07-12T21:03:07.338+02:00 level=INFO msg="found device" address=3daf1e2e-8b2e-4285-81df-5fd7da72c5c1 rssi=-60 local_name="stand UP- 0829"
time=2024-07-12T21:03:07.506+02:00 level=INFO msg="Connected to the device" address=3daf1e2e-8b2e-4285-81df-5fd7da72c5c1 local_name="stand UP- 0829"
time=2024-07-12T21:03:08.684+02:00 level=INFO msg="Device successfully initialised" height=65
time=2024-07-12T21:03:11.839+02:00 level=INFO msg="Moving to position for specified duration" position=105 duration=15m0s
time=2024-07-12T21:03:11.840+02:00 level=INFO msg="Staring to move" direction=up height=74 target_height=105 speed=64
time=2024-07-12T21:03:11.840+02:00 level=INFO msg="Finished moving" direction=up height=105 target_height=105 speed=64
```

In this case the program will cycle twice through the two specified positions:

1. `105 centimeters` for `15 minutes`
1. `65 centimeters` for `25 minutes`

Enjoy!

## TiMOTION Trademark Notice

The TiMOTION name and logo are trademarks of TiMOTION Technology Co. Ltd. All Rights Reserved. Is not part of the licensing for this project.
