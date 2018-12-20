package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"miflorad/common"
	impl "miflorad/common/gatt"

	"github.com/currantlabs/gatt"
	"github.com/currantlabs/gatt/examples/option"
)

const (
	discoveryTimeout  = 4 * time.Second
	connectionTimeout = 4 * time.Second
)

type DiscoveryResult struct {
	p    gatt.Peripheral
	a    *gatt.Advertisement
	rssi int
}

var (
	discoveryDone  = make(chan DiscoveryResult)
	connectionDone = make(chan struct{})
)

var timeConnectStart time.Time

func onStateChanged(device gatt.Device, state gatt.State) {
	fmt.Fprintln(os.Stderr, "State:", state)
	switch state {
	case gatt.StatePoweredOn:
		timeConnectStart = time.Now()
		fmt.Fprintln(os.Stderr, "Scanning...")
		device.Scan([]gatt.UUID{}, false)
		return
	default:
		device.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	id := strings.ToUpper(flag.Args()[1])
	if strings.ToUpper(p.ID()) != id {
		return
	}

	// Stop scanning once we've got the peripheral we're looking for.
	p.Device().StopScanning()

	discoveryDone <- DiscoveryResult{p, a, rssi}
}

func onPeriphConnected(p gatt.Peripheral, err error) {
	fmt.Fprintln(os.Stderr, "Connected")

	prefix := flag.Args()[0]
	id := common.MifloraGetAlphaNumericID(flag.Args()[1])

	timeConnectTook := time.Since(timeConnectStart).Seconds()
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.connect_time %.2f %d\n", prefix, id, timeConnectTook, time.Now().Unix())

	timeReadoutStart := time.Now()

	// Note: can hang due when device has terminated the connection on it's own already
	// defer p.Device().CancelConnection(p)

	if err := p.SetMTU(500); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set MTU, err: %s\n", err)
	}

	// Discover services and characteristics
	{
		_, err := p.DiscoverServices([]gatt.UUID{gatt.MustParseUUID(common.MifloraServiceUUID)})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to discover services, err: %s\n", err)
			fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
			return
		}
	}
	for _, service := range p.Services() {
		_, err := p.DiscoverCharacteristics(nil, service)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to discover characteristics, err: %s\n", err)
			fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
			return
		}
	}

	metaData, err := impl.MifloraRequestVersionBattery(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to request version battery, err: %s\n", err)
		fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
		return
	}

	fmt.Fprintf(os.Stderr, "Firmware version: %s\n", metaData.FirmwareVersion)

	fmt.Fprintf(os.Stdout, "%s.miflora.%s.battery_level %d %d\n", prefix, id, metaData.BatteryLevel, time.Now().Unix())
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.firmware_version %d %d\n", prefix, id, metaData.NumericFirmwareVersion(), time.Now().Unix())

	if metaData.RequiresModeChangeBeforeRead() {
		err2 := impl.MifloraRequestModeChange(p)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Failed to request mode change, err: %s\n", err2)
			fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
			return
		}
	}

	sensorData, err3 := impl.MifloraRequstSensorData(p)
	if err3 != nil {
		fmt.Fprintf(os.Stderr, "Failed to request sensor data, err: %s\n", err3)
		fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
		return
	}
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.temperature %.1f %d\n", prefix, id, sensorData.Temperature, time.Now().Unix())
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.brightness %d %d\n", prefix, id, sensorData.Brightness, time.Now().Unix())
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.moisture %d %d\n", prefix, id, sensorData.Moisture, time.Now().Unix())
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.conductivity %d %d\n", prefix, id, sensorData.Conductivity, time.Now().Unix())

	timeReadoutTook := time.Since(timeReadoutStart).Seconds()
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.readout_time %.2f %d\n", prefix, id, timeReadoutTook, time.Now().Unix())

	// TODO: report that we are done without closing connection, since it could hang
	close(connectionDone)
}

func onPeriphDisconnected(p gatt.Peripheral, err error) {
	fmt.Fprintln(os.Stderr, "Disconnected")
	close(connectionDone)
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] prefix peripheral-id\n", os.Args[0])
		os.Exit(1)
	}

	prefix := flag.Args()[0]
	id := common.MifloraGetAlphaNumericID(flag.Args()[1])

	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open device, err: %s\n", err)
		fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
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
		fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
		device.StopScanning()
		device.Stop()
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Discovered peripheral ID:%s, NAME:(%s), RSSI:%d\n", discoveryResult.p.ID(), discoveryResult.p.Name(), discoveryResult.rssi)

	fmt.Fprintf(os.Stdout, "%s.miflora.%s.rssi %d %d\n", prefix, id, discoveryResult.rssi, time.Now().Unix())

	// Register connection handlers
	device.Handle(
		gatt.PeripheralConnected(onPeriphConnected),
		// gatt.PeripheralDisconnected(onPeriphDisconnected),
	)

	device.Connect(discoveryResult.p)

	select {
	case <-connectionDone:
		fmt.Fprintln(os.Stderr, "Connection done")
	case <-time.After(connectionTimeout):
		fmt.Fprintln(os.Stderr, "Connection timed out")
		// TODO: can hang due when device has terminated the connection on it's own already
		// device.CancelConnection(discoveryResult.p)
		os.Exit(1)
	}

	// Note: calls CancelConnection() and thus suffers the same problem, kernel will cleanup after our process finishes
	// device.Stop()
}
