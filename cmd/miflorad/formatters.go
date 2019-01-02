package main

import (
	"fmt"
	"strings"
	"time"
)

func publishGraphite(metric mifloraMetric, publish chan string, metricsBase string) {
	timestamp := time.Now().Unix()
	prefix := fmt.Sprintf("%s.miflora.%s", metricsBase, metric.getPeripheralId())

	switch metric := metric.(type) {
	case mifloraDataMetric:
		publish <- fmt.Sprintf("%s.battery_level %d %d", prefix, metric.metaData.BatteryLevel, timestamp)
		publish <- fmt.Sprintf("%s.firmware_version %d %d", prefix, metric.metaData.NumericFirmwareVersion(), timestamp)
		publish <- fmt.Sprintf("%s.temperature %.1f %d", prefix, metric.sensorData.Temperature, timestamp)
		publish <- fmt.Sprintf("%s.brightness %d %d", prefix, metric.sensorData.Brightness, timestamp)
		publish <- fmt.Sprintf("%s.moisture %d %d", prefix, metric.sensorData.Moisture, timestamp)
		publish <- fmt.Sprintf("%s.conductivity %d %d", prefix, metric.sensorData.Conductivity, timestamp)
		publish <- fmt.Sprintf("%s.connect_time %.2f %d", prefix, metric.connectTime, timestamp)
		publish <- fmt.Sprintf("%s.readout_time %.2f %d", prefix, metric.readoutTime, timestamp)
		publish <- fmt.Sprintf("%s.rssi %d %d", prefix, metric.rssi, timestamp)
	case mifloraErrorMetric:
		publish <- fmt.Sprintf("%s.failed %d %d", prefix, metric.failed, timestamp)
	}
}

func publishInflux(metric mifloraMetric, publish chan string) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("miflora,id=%s ", metric.getPeripheralId()))
	switch metric := metric.(type) {
	case mifloraDataMetric:
		b.WriteString(fmt.Sprintf("battery_level=%d,", metric.metaData.BatteryLevel))
		b.WriteString(fmt.Sprintf("firmware_version=%d,", metric.metaData.NumericFirmwareVersion()))
		b.WriteString(fmt.Sprintf("temperature=%.1f,", metric.sensorData.Temperature))
		b.WriteString(fmt.Sprintf("brightness=%d,", metric.sensorData.Brightness))
		b.WriteString(fmt.Sprintf("moisture=%d,", metric.sensorData.Moisture))
		b.WriteString(fmt.Sprintf("conductivity=%d,", metric.sensorData.Conductivity))
		b.WriteString(fmt.Sprintf("connect_time=%.2f,", metric.connectTime))
		b.WriteString(fmt.Sprintf("readout_time=%.2f,", metric.readoutTime))
		b.WriteString(fmt.Sprintf("rssi=%d", metric.rssi))
	case mifloraErrorMetric:
		b.WriteString(fmt.Sprintf("failed=%d", metric.failed))
	}
	b.WriteString(fmt.Sprintf(" %d", time.Now().UnixNano()))
	publish <- b.String()
}
