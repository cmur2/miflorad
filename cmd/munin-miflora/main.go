package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	impl "miflorad/common/ble"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
)

const discoveryTimeout = 4 * time.Second
const connectionTimeout = 4 * time.Second

func readData(client ble.Client, profile *ble.Profile) {
	prefix := flag.Args()[0]

	regexNonAlphaNumeric, err4 := regexp.Compile("[^a-z0-9]+")
	if err4 != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile regex, err: %s\n", err4)
	}
	id := regexNonAlphaNumeric.ReplaceAllString(strings.ToLower(flag.Args()[1]), "")

	metaData, err := impl.RequestVersionBattery(client, profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to request version battery, err: %s\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "Firmware version: %s\n", metaData.FirmwareVersion)

	fmt.Fprintf(os.Stdout, "%s.miflora.%s.battery_level %d %d\n", prefix, id, metaData.BatteryLevel, time.Now().Unix())
	fmt.Fprintf(os.Stdout, "%s.miflora.%s.firmware_version %d %d\n", prefix, id, metaData.NumericFirmwareVersion(), time.Now().Unix())

	// for the newer models a magic number must be written before we can read the current data
	if metaData.FirmwareVersion >= "2.6.6" {
		err2 := impl.RequestModeChange(client, profile)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Failed to request mode change, err: %s\n", err2)
			return
		}
	}

	sensorData, err3 := impl.RequestSensorData(client, profile)
	if err3 != nil {
		fmt.Fprintf(os.Stderr, "Failed to request sensor data, err: %s\n", err3)
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

	{
		device, err := dev.NewDevice("default")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open device, err: %s\n", err)
			os.Exit(1)
		}
		ble.SetDefaultDevice(device)
	}

	filter := func(adv ble.Advertisement) bool {
		return strings.ToUpper(adv.Addr().String()) == strings.ToUpper(flag.Args()[1])
	}

	fmt.Fprintln(os.Stderr, "Scanning...")
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), discoveryTimeout))
	client, err := ble.Connect(ctx, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to %s, err: %s\n", flag.Args()[1], err)
		os.Exit(1)
	}

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

	profile, err := client.DiscoverProfile(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to discover profile, err: %s\n", err)
		os.Exit(1)
	}

	readData(client, profile)

	fmt.Fprintln(os.Stderr, "Connection done")

	client.CancelConnection()

	<-done
}
