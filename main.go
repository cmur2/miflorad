package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/currantlabs/gatt"
	"github.com/currantlabs/gatt/examples/option"
)

const discoveryTimeout = 4 * time.Second
const connectionTimeout = 4 * time.Second

var mifloraModeChangeData = []byte{0xa0, 0x1f}

var mifloraServiceUUID = gatt.MustParseUUID("00001204-0000-1000-8000-00805f9b34fb")
var mifloraCharModeChangeUUID = gatt.MustParseUUID("00001a00-0000-1000-8000-00805f9b34fb")
var mifloraCharReadSensorDataUUID = gatt.MustParseUUID("00001a01-0000-1000-8000-00805f9b34fb")
var mifloraCharVersionBatteryUUID = gatt.MustParseUUID("00001a02-0000-1000-8000-00805f9b34fb")

var discoveryDone = make(chan gatt.Peripheral)
var connectionDone = make(chan struct{})

func onStateChanged(device gatt.Device, state gatt.State) {
	fmt.Println("State:", state)
	switch state {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning...")
		device.Scan([]gatt.UUID{}, false)
		return
	default:
		device.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	id := strings.ToUpper(flag.Args()[0])
	fmt.Println(p.ID())
	if strings.ToUpper(p.ID()) != id {
		return
	}

	// Stop scanning once we've got the peripheral we're looking for.
	p.Device().StopScanning()

	discoveryDone <- p
}

func findServiceByUUID(services []*gatt.Service, u gatt.UUID) *gatt.Service {
	for _, service := range services {
		if service.UUID().Equal(u) {
			return service
		}
	}
	return nil
}

func findCharacteristicByUUID(characteristics []*gatt.Characteristic, u gatt.UUID) *gatt.Characteristic {
	for _, characteristic := range characteristics {
		if characteristic.UUID().Equal(u) {
			return characteristic
		}
	}
	return nil
}

func onPeriphConnected(p gatt.Peripheral, err error) {
	fmt.Println("Connected")
	defer p.Device().CancelConnection(p)

	if err := p.SetMTU(500); err != nil {
		fmt.Printf("Failed to set MTU, err: %s\n", err)
	}

	// Discovery services
	services, err := p.DiscoverServices([]gatt.UUID{mifloraServiceUUID})
	if err != nil {
		fmt.Printf("Failed to discover services, err: %s\n", err)
		return
	}

	if len(services) == 1 {
		mifloraService := findServiceByUUID(services, mifloraServiceUUID)

		chars, err := p.DiscoverCharacteristics(nil, mifloraService)
		if err != nil {
			fmt.Printf("Failed to discover characteristics, err: %s\n", err)
			return
		}

		mifloraVersionBatteryChar := findCharacteristicByUUID(chars, mifloraCharVersionBatteryUUID)
		bytes, err := p.ReadCharacteristic(mifloraVersionBatteryChar)
		if err != nil {
			fmt.Printf("Failed to read characteristic, err: %s\n", err)
			return
		}
		fmt.Printf("Battery level:    %d%%\n", uint8(bytes[0]))
		fmt.Printf("Firmware version: %s\n", string(bytes[2:]))

		// for the newer models a magic number must be written before we can read the current data
		if string(bytes[2:]) >= "2.6.6" {
			mifloraModeChangeChar := findCharacteristicByUUID(chars, mifloraCharModeChangeUUID)
			err2 := p.WriteCharacteristic(mifloraModeChangeChar, mifloraModeChangeData, false)
			if err2 != nil {
				fmt.Printf("Failed to write characteristic, err: %s\n", err2)
				return
			}
		}

		mifloraSensorDataChar := findCharacteristicByUUID(chars, mifloraCharReadSensorDataUUID)
		bytes2, err3 := p.ReadCharacteristic(mifloraSensorDataChar)
		if err3 != nil {
			fmt.Printf("Failed to read characteristic, err: %s\n", err3)
			return
		}
		fmt.Printf("Temparature:      %f °C\n", float64(binary.LittleEndian.Uint16(bytes2[0:2]))/10.0)
		fmt.Printf("Brightness:       %d lux\n", binary.LittleEndian.Uint32(bytes2[3:7]))
		fmt.Printf("Moisture:         %d %%\n", uint8(bytes2[7]))
		fmt.Printf("Conductivity:     %d µS/cm\n", binary.LittleEndian.Uint16(bytes2[8:10]))
	}
}

func onPeriphDisconnected(p gatt.Peripheral, err error) {
	fmt.Println("Disconnected")
	close(connectionDone)
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalf("usage: %s [options] peripheral-id\n", os.Args[0])
	}

	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	// Register discovery handler
	device.Handle(gatt.PeripheralDiscovered(onPeriphDiscovered))

	device.Init(onStateChanged)

	var peripheral gatt.Peripheral

	select {
	case peripheral = <-discoveryDone:
		fmt.Println("Discovery done")
	case <-time.After(discoveryTimeout):
		// fmt.Println("Discovery timed out")
		log.Fatalf("Discovery timed out\n")
		device.StopScanning()
		device.Stop()
	}

	fmt.Printf("Discovered peripheral ID:%s, NAME:(%s)\n", peripheral.ID(), peripheral.Name())

	// Register connection handlers
	device.Handle(
		gatt.PeripheralConnected(onPeriphConnected),
		gatt.PeripheralDisconnected(onPeriphDisconnected),
	)

	device.Connect(peripheral)

	select {
	case <-connectionDone:
		fmt.Println("Connection done")
	case <-time.After(connectionTimeout):
		fmt.Println("Connection timed out")
	}

	device.Stop()
}
