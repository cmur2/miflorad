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

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/pkg/errors"
)

const mqttConnectTimeout = 10 * time.Second

// program version, will be populated on build
var version string

var (
	scanTimeout       = flag.Duration("scantimeout", 10*time.Second, "timeout after that a scan per peripheral will be aborted")
	readRetries       = flag.Int("readretries", 2, "number of times reading will be attempted per peripheral")
	interval          = flag.Duration("interval", 25*time.Second, "metrics collection interval")
	prefix            = flag.String("prefix", "", "metrics name prefix")
	brokerHost        = flag.String("brokerhost", "localhost", "MQTT broker host to send metrics to")
	brokerUser        = flag.String("brokeruser", "", "MQTT broker user used for authentication")
	brokerPassword    = flag.String("brokerpassword", "", "MQTT broker password used for authentication")
	brokerUseTLS      = flag.Bool("brokerusetls", true, "whether TLS should be used for MQTT broker")
	brokerTopicPrefix = flag.String("brokertopicprefix", "", "MQTT topic prefix for messages")
)

type peripheral struct {
	id                string
	lastMetaDataFetch time.Time
	metaData          common.VersionBatteryResponse
}

var allPeripherals []*peripheral

type mqttLogger struct {
	level string
}

func (logger mqttLogger) Println(a ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("mqtt %s:", logger.level), a)
}

func (logger mqttLogger) Printf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "mqtt %s: "+format, logger.level, a)
}

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

func getMQTTOptions() *mqtt.ClientOptions {
	if *brokerUseTLS {
		return mqtt.NewClientOptions().
			AddBroker(fmt.Sprintf("ssl://%s:8883", *brokerHost)).
			SetUsername(*brokerUser).
			SetPassword(*brokerPassword)
	} else {
		return mqtt.NewClientOptions().
			AddBroker(fmt.Sprintf("tcp://%s:1883", *brokerHost)).
			SetUsername(*brokerUser).
			SetPassword(*brokerPassword)
	}
}

func readData(peripheral *peripheral, client ble.Client, send chan string) error {
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

	id := common.MifloraGetAlphaNumericID(peripheral.id)
	timestamp := time.Now().Unix()

	send <- fmt.Sprintf("%s.miflora.%s.battery_level %d %d", *prefix, id, peripheral.metaData.BatteryLevel, timestamp)
	send <- fmt.Sprintf("%s.miflora.%s.firmware_version %d %d", *prefix, id, peripheral.metaData.NumericFirmwareVersion(), timestamp)

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

	send <- fmt.Sprintf("%s.miflora.%s.temperature %.1f %d", *prefix, id, sensorData.Temperature, timestamp)
	send <- fmt.Sprintf("%s.miflora.%s.brightness %d %d", *prefix, id, sensorData.Brightness, timestamp)
	send <- fmt.Sprintf("%s.miflora.%s.moisture %d %d", *prefix, id, sensorData.Moisture, timestamp)
	send <- fmt.Sprintf("%s.miflora.%s.conductivity %d %d", *prefix, id, sensorData.Conductivity, timestamp)

	return nil
}

func connectPeripheral(peripheral *peripheral, send chan string) error {
	fmt.Fprintf(os.Stderr, "Scanning for %s...\n", peripheral.id)

	id := common.MifloraGetAlphaNumericID(peripheral.id)
	timestamp := time.Now().Unix()

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
	send <- fmt.Sprintf("%s.miflora.%s.connect_time %.2f %d", *prefix, id, timeConnectTook, timestamp)

	foundAdvertisement := <-foundAdvertisementChannel
	send <- fmt.Sprintf("%s.miflora.%s.rssi %d %d", *prefix, id, foundAdvertisement.RSSI(), timestamp)

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

	err2 := readData(peripheral, client, send)

	timeReadoutTook := time.Since(timeReadoutStart).Seconds()
	send <- fmt.Sprintf("%s.miflora.%s.readout_time %.2f %d", *prefix, id, timeReadoutTook, timestamp)

	client.CancelConnection()

	<-done

	return err2
}

func readPeripheral(peripheral *peripheral, send chan string) error {
	var err error
	for retry := 0; retry < *readRetries; retry++ {
		err = connectPeripheral(peripheral, send)
		// stop retrying once we have a success, last err will be returned (or nil)
		if err == nil {
			break
		}
	}
	return err
}

func readAllPeripherals(quit chan struct{}, send chan string) {
	for _, peripheral := range allPeripherals {
		// check for quit signal (non-blocking) and terminate
		select {
		case <-quit:
			return
		default:
		}

		err := readPeripheral(peripheral, send)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read peripheral %s, err: %s\n", peripheral.id, err)
			id := common.MifloraGetAlphaNumericID(peripheral.id)
			send <- fmt.Sprintf("%s.miflora.%s.failed 1 %d", *prefix, id, time.Now().Unix())
			continue
		}
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

	mqtt.WARN = mqttLogger{level: "warning"}
	mqtt.ERROR = mqttLogger{level: "error"}
	mqtt.CRITICAL = mqttLogger{level: "critical"}

	mqttClient := mqtt.NewClient(getMQTTOptions())

	if token := mqttClient.Connect(); token.WaitTimeout(mqttConnectTimeout) && token.Error() != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect MQTT, err: %s\n", token.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Connected to MQTT broker %s\n", *brokerHost)

	device, err := dev.NewDevice("default")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open device, err: %s\n", err)
		os.Exit(1)
	}
	ble.SetDefaultDevice(device)

	intervalTicker := time.NewTicker(*interval)
	quit := make(chan struct{})
	send := make(chan string, 10)

	go func() {
		fmt.Fprintf(os.Stderr, "Starting loop with %s interval...\n", *interval)

		// populate all peripherals data structure
		allPeripherals = make([]*peripheral, len(flag.Args()))
		for i, peripheralID := range flag.Args() {
			allPeripherals[i] = &peripheral{
				id:                peripheralID,
				lastMetaDataFetch: time.Unix(0, 0), // force immediate 1st request
			}
		}

		// main loop
		readAllPeripherals(quit, send)
		for range intervalTicker.C {
			readAllPeripherals(quit, send)
		}
	}()

	go func() {
		for metric := range send {
			// fmt.Fprintln(os.Stdout, metric)
			token := mqttClient.Publish(*brokerTopicPrefix+*prefix, 1, false, metric)
			if token.WaitTimeout(1*time.Second) && token.Error() != nil {
				fmt.Fprintf(os.Stderr, "Failed to publish MQTT, err: %s\n", token.Error())
				continue
			}
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

	mqttClient.Disconnect(1000)

	if err := device.Stop(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close device, err: %s\n", err)
		os.Exit(1)
	}
}
