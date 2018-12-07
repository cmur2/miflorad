package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"miflorad/common"

	"github.com/currantlabs/gatt"
	"github.com/currantlabs/gatt/examples/option"
)

const discoveryTimeout = 4 * time.Second
const connectionTimeout = 4 * time.Second

type DiscoveryResult struct {
	p    gatt.Peripheral
	a    *gatt.Advertisement
	rssi int
}

var discoveryDone = make(chan DiscoveryResult)
var connectionDone = make(chan struct{})

func onStateChanged(device gatt.Device, state gatt.State) {
	fmt.Fprintln(os.Stderr, "State:", state)
	switch state {
	case gatt.StatePoweredOn:
		fmt.Fprintln(os.Stderr, "Scanning...")
		device.Scan([]gatt.UUID{}, false)
		return
	default:
		device.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	id := strings.ToUpper(flag.Args()[0])
	// fmt.Fprintln(os.Stderr, p.ID())
	if strings.ToUpper(p.ID()) != id {
		return
	}

	// Stop scanning once we've got the peripheral we're looking for.
	p.Device().StopScanning()

	discoveryDone <- DiscoveryResult{p, a, rssi}
}

func onPeriphConnected(p gatt.Peripheral, err error) {
	fmt.Fprintln(os.Stderr, "Connected")
	// defer p.Device().CancelConnection(p)

	if err := p.SetMTU(500); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set MTU, err: %s\n", err)
	}

	// Discover services and characteristics
	{
		_, err := p.DiscoverServices([]gatt.UUID{common.MifloraServiceUUID})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to discover services, err: %s\n", err)
			return
		}
	}
	for _, service := range p.Services() {
		_, err := p.DiscoverCharacteristics(nil, service)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to discover characteristics, err: %s\n", err)
			return
		}
	}

	metaData, err := common.MifloraRequestVersionBattery(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to request version battery, err: %s\n", err)
		return
	}
	fmt.Fprintf(os.Stdout, "Battery level:    %d%%\n", metaData.BatteryLevel)
	fmt.Fprintf(os.Stdout, "Firmware version: %s\n", metaData.FirmwareVersion)

	// for the newer models a magic number must be written before we can read the current data
	if metaData.FirmwareVersion >= "2.6.6" {
		err2 := common.MifloraRequestModeChange(p)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Failed to request mode change, err: %s\n", err2)
			return
		}
	}

	sensorData, err3 := common.MifloraRequstSensorData(p)
	if err3 != nil {
		fmt.Fprintf(os.Stderr, "Failed to request sensor data, err: %s\n", err3)
		return
	}
	fmt.Fprintf(os.Stdout, "Temparature:      %.1f °C\n", sensorData.Temperature)
	fmt.Fprintf(os.Stdout, "Brightness:       %d lux\n", sensorData.Brightness)
	fmt.Fprintf(os.Stdout, "Moisture:         %d %%\n", sensorData.Moisture)
	fmt.Fprintf(os.Stdout, "Conductivity:     %d µS/cm\n", sensorData.Conductivity)
}

func onPeriphDisconnected(p gatt.Peripheral, err error) {
	fmt.Fprintln(os.Stderr, "Disconnected")
	close(connectionDone)
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] peripheral-id\n", os.Args[0])
		os.Exit(1)
	}

	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open device, err: %s\n", err)
		os.Exit(1)
	}

	// Register discovery handler
	device.Handle(gatt.PeripheralDiscovered(onPeriphDiscovered))

	device.Init(onStateChanged)

	var discoveryResult DiscoveryResult

	select {
	case discoveryResult = <-discoveryDone:
		fmt.Fprintln(os.Stderr, "Discovery done")
	case <-time.After(discoveryTimeout):
		fmt.Fprintln(os.Stderr, "Discovery timed out")
		device.StopScanning()
		device.Stop()
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Discovered peripheral ID:%s, NAME:(%s), RSSI:%d\n", discoveryResult.p.ID(), discoveryResult.p.Name(), discoveryResult.rssi)

	// Register connection handlers
	device.Handle(
		gatt.PeripheralConnected(onPeriphConnected),
		gatt.PeripheralDisconnected(onPeriphDisconnected),
	)

	device.Connect(discoveryResult.p)

	select {
	case <-connectionDone:
		fmt.Fprintln(os.Stderr, "Connection done")
	case <-time.After(connectionTimeout):
		fmt.Fprintln(os.Stderr, "Connection timed out")
		fmt.Fprintln(os.Stderr, "A")
		// device.CancelConnection(discoveryResult.p)
		fmt.Fprintln(os.Stderr, "B")
	}

	fmt.Fprintln(os.Stderr, "C")
	// device.Stop()
	fmt.Fprintln(os.Stderr, "D")
}
