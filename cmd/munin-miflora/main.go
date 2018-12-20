package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"miflorad/common"
	impl "miflorad/common/ble"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
)

const discoveryTimeout = 10 * time.Second

func readData(client ble.Client, profile *ble.Profile) {
	prefix := flag.Args()[0]
	id := common.MifloraGetAlphaNumericID(flag.Args()[1])

	metaData, err := impl.RequestVersionBattery(client, profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to request version battery, err: %s\n", err)
		fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
		return
	}

	fmt.Fprintf(os.Stderr, "Firmware version: %s\n", metaData.FirmwareVersion)

	fmt.Fprintf(os.Stdout, "%s.miflora.%s.battery_level %d %d\n", prefix, id, metaData.BatteryLevel, time.Now().Unix())
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.firmware_version %d %d\n", prefix, id, metaData.NumericFirmwareVersion(), time.Now().Unix())

	if metaData.RequiresModeChangeBeforeRead() {
		err2 := impl.RequestModeChange(client, profile)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Failed to request mode change, err: %s\n", err2)
			fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
			return
		}
	}

	sensorData, err3 := impl.RequestSensorData(client, profile)
	if err3 != nil {
		fmt.Fprintf(os.Stderr, "Failed to request sensor data, err: %s\n", err3)
		fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
		return
	}
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.temperature %.1f %d\n", prefix, id, sensorData.Temperature, time.Now().Unix())
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.brightness %d %d\n", prefix, id, sensorData.Brightness, time.Now().Unix())
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.moisture %d %d\n", prefix, id, sensorData.Moisture, time.Now().Unix())
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.conductivity %d %d\n", prefix, id, sensorData.Conductivity, time.Now().Unix())
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] prefix peripheral-id\n", os.Args[0])
		os.Exit(1)
	}

	prefix := flag.Args()[0]
	id := common.MifloraGetAlphaNumericID(flag.Args()[1])

	{
		device, err := dev.NewDevice("default")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open device, err: %s\n", err)
			fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
			os.Exit(1)
		}
		ble.SetDefaultDevice(device)
	}

	// only way to get back the found advertisement, must be buffered!
	foundAdvertisementChannel := make(chan ble.Advertisement, 1)

	filter := func(adv ble.Advertisement) bool {
		if strings.ToUpper(adv.Addr().String()) == strings.ToUpper(flag.Args()[1]) {
			foundAdvertisementChannel <- adv
			return true
		}
		return false
	}

	timeConnectStart := time.Now()

	fmt.Fprintln(os.Stderr, "Scanning...")
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), discoveryTimeout))
	client, err := ble.Connect(ctx, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to %s, err: %s\n", flag.Args()[1], err)
		fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
		os.Exit(1)
	}

	timeConnectTook := time.Since(timeConnectStart).Seconds()
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.connect_time %.2f %d\n", prefix, id, timeConnectTook, time.Now().Unix())

	foundAdvertisement := <-foundAdvertisementChannel
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.rssi %d %d\n", prefix, id, foundAdvertisement.RSSI(), time.Now().Unix())

	// Source: https://github.com/go-ble/ble/blob/master/examples/basic/explorer/main.go#L53
	// Make sure we had the chance to print out the message.
	done := make(chan struct{})
	// Normally, the connection is disconnected by us after our exploration.
	// However, it can be asynchronously disconnected by the remote peripheral.
	// So we wait(detect) the disconnection in the go routine.
	go func() {
		<-client.Disconnected()
		fmt.Fprintln(os.Stderr, "Disconnected")
		close(done)
	}()

	fmt.Fprintln(os.Stderr, "Connected")

	timeReadoutStart := time.Now()

	profile, err := client.DiscoverProfile(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to discover profile, err: %s\n", err)
		fmt.Fprintf(os.Stdout, "%s.miflora.%s.failed 1 %d\n", prefix, id, time.Now().Unix())
		os.Exit(1)
	}

	readData(client, profile)

	timeReadoutTook := time.Since(timeReadoutStart).Seconds()
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.readout_time %.2f %d\n", prefix, id, timeReadoutTook, time.Now().Unix())

	fmt.Fprintln(os.Stderr, "Connection done")

	client.CancelConnection()

	<-done
}
