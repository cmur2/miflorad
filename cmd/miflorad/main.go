package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	common "miflorad/common"
	impl "miflorad/common/ble"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"

	"github.com/pkg/errors"
)

// program version, will be populated on build
var version string

var (
	scanTimeout = flag.Duration("scantimeout", 10*time.Second, "timeout after that a scan per peripheral will be aborted")
	interval    = flag.Duration("interval", 25*time.Second, "metrics collection interval")
	readRetries = flag.Int("readretries", 2, "number of times reading will be attempted per peripheral")
)

type peripheral struct {
	id                string
	lastMetaDataFetch time.Time
	metaData          common.VersionBatteryResponse
}

var allPeripherals []*peripheral

var (
	countSuccess = 0
	countFailure = 0
)

func checkTooShortInterval() error {
	numPeripherals := int64(len(flag.Args()))
	numReadRetries := int64(*readRetries)
	if (*scanTimeout).Nanoseconds()*numReadRetries*numPeripherals >= (*interval).Nanoseconds() {
		return errors.New(fmt.Sprintf(
			"The interval of %s is too short given the scan timeout of %s for %d peripheral(s) with %d retries each! Exiting...\n",
			*interval, *scanTimeout, numPeripherals, *readRetries))
	}
	return nil
}

func getVersion() string {
	if version == "" {
		return "dev"
	} else {
		return version
	}
}

func readData(peripheral *peripheral, client ble.Client) error {
	// re-request meta data (for battery level) if last check more than 24 hours ago
	// Source: https://github.com/open-homeautomation/miflora/blob/ffd95c3e616df8843cc8bff99c9b60765b124092/miflora/miflora_poller.py#L92
	if time.Since(peripheral.lastMetaDataFetch) >= 24*time.Hour {
		metaData, err := impl.RequestVersionBattery(client)
		if err != nil {
			return errors.Wrap(err, "can't request version battery")
		}
		peripheral.metaData = metaData
		peripheral.lastMetaDataFetch = time.Now()
	}

	if peripheral.metaData.RequiresModeChangeBeforeRead() {
		err2 := impl.RequestModeChange(client)
		if err2 != nil {
			return errors.Wrap(err2, "can't request mode change")
		}
	}

	sensorData, err3 := impl.RequestSensorData(client)
	if err3 != nil {
		return errors.Wrap(err3, "can't request sensor data")
	}

	fmt.Println(sensorData.Temperature, sensorData.Brightness)

	return nil
}

func connectPeripheral(peripheral *peripheral) error {
	fmt.Fprintf(os.Stderr, "Scanning for %s...\n", peripheral.id)

	// only way to get back the found advertisement, must be buffered!
	foundAdvertisementChannel := make(chan ble.Advertisement, 1)

	filter := func(adv ble.Advertisement) bool {
		if strings.ToUpper(adv.Addr().String()) == strings.ToUpper(peripheral.id) {
			foundAdvertisementChannel <- adv
			return true
		}
		return false
	}

	timeConnectStart := time.Now()

	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *scanTimeout))
	client, err := ble.Connect(ctx, filter)
	if err != nil {
		return errors.Wrapf(err, "can't connect to %s", peripheral.id)
	}

	timeConnectTook := time.Since(timeConnectStart).Seconds()
	fmt.Println(timeConnectTook)
	// fmt.Fprintf(os.Stdout, "%s.miflora.%s.connect_time %.2f %d\n", prefix, id, timeConnectTook, time.Now().Unix())

	// foundAdvertisement := <-foundAdvertisementChannel
	// fmt.Fprintf(os.Stdout, "%s.miflora.%s.rssi %d %d\n", prefix, id, foundAdvertisement.RSSI(), time.Now().Unix())

	// Source: https://github.com/go-ble/ble/blob/master/examples/basic/explorer/main.go#L53
	// Make sure we had the chance to print out the message.
	done := make(chan struct{})
	// Normally, the connection is disconnected by us after our exploration.
	// However, it can be asynchronously disconnected by the remote peripheral.
	// So we wait(detect) the disconnection in the go routine.
	go func() {
		<-client.Disconnected()
		close(done)
	}()

	timeReadoutStart := time.Now()

	if _, err := client.DiscoverProfile(true); err != nil {
		return errors.Wrap(err, "can't descover profile")
	}

	err2 := readData(peripheral, client)

	timeReadoutTook := time.Since(timeReadoutStart).Seconds()
	fmt.Println(timeReadoutTook)
	// fmt.Fprintf(os.Stdout, "%s.miflora.%s.readout_time %.2f %d\n", prefix, id, timeReadoutTook, time.Now().Unix())

	client.CancelConnection()

	<-done

	return err2
}

func readPeripheral(peripheral *peripheral) error {
	var err error
	for retry := 0; retry < *readRetries; retry++ {
		err = connectPeripheral(peripheral)
		// stop retrying once we have a success, last err will be returned (or nil)
		if err == nil {
			break
		}
	}
	return err
}

func readAllPeripherals(quit chan struct{}) {
	for _, peripheral := range allPeripherals {
		// check for quit signal (non-blocking) and terminate
		select {
		case <-quit:
			return
		default:
		}

		err := readPeripheral(peripheral)
		if err != nil {
			countFailure++
			fmt.Fprintf(os.Stderr, "Failed to read peripheral %s, err: %s\n", peripheral.id, err)
			continue
		}
		countSuccess++
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] peripheral-id [peripheral-ids...] \n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := checkTooShortInterval(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "miflorad version %s\n", getVersion())

	device, err := dev.NewDevice("default")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open device, err: %s\n", err)
		os.Exit(1)
	}
	ble.SetDefaultDevice(device)

	intervalTicker := time.NewTicker(*interval)
	quit := make(chan struct{})

	go func() {
		fmt.Fprintf(os.Stderr, "Starting miflorad loop with %s interval...\n", *interval)

		// populate all peripherals data structure
		allPeripherals = make([]*peripheral, len(flag.Args()))
		for i, peripheralID := range flag.Args() {
			allPeripherals[i] = &peripheral{
				id:                peripheralID,
				lastMetaDataFetch: time.Unix(0, 0), // force immediate 1st request
			}
		}

		// main loop
		readAllPeripherals(quit)
		for range intervalTicker.C {
			readAllPeripherals(quit)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	signal := <-signals
	fmt.Fprintf(os.Stderr, "Received %s! Stopping...\n", signal)
	intervalTicker.Stop()
	close(quit)
	// wait for last readPeripheral to finish (worst case)
	time.Sleep(*scanTimeout * time.Duration(*readRetries))

	fmt.Fprintf(os.Stderr, "Failures:  %d\n", countFailure)
	fmt.Fprintf(os.Stderr, "Successes: %d\n", countSuccess)

	if err := device.Stop(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close device, err: %s\n", err)
		os.Exit(1)
	}
}
