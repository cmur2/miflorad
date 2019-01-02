package main

import (
	_ "fmt"
	"strconv"
	"strings"
	"testing"

	common "miflorad/common"

	"github.com/stretchr/testify/assert"
)

func TestPublishGraphite(t *testing.T) {
	tables := []struct {
		metric mifloraMetric
	}{
		{mifloraErrorMetric{peripheralId: "peri", failed: 1}},
		{mifloraDataMetric{
			peripheralId: "peri",
			metaData:     common.VersionBatteryResponse{BatteryLevel: 100, FirmwareVersion: "2.7.0"},
			sensorData:   common.SensorDataResponse{Temperature: 24.2, Brightness: 121, Moisture: 16, Conductivity: 101},
			connectTime:  3.42,
			readoutTime:  0.23,
			rssi:         -77,
		}},
	}

	for _, table := range tables {
		publish := make(chan string, 100)
		publishGraphite(table.metric, publish, "foo.base")
		close(publish)
		switch table.metric.(type) {
		case mifloraErrorMetric:
			for line := range publish {
				parts := strings.Split(line, " ")
				assert.Equal(t, len(parts), 3)
				assert.Equal(t, parts[0], "foo.base.miflora.peri.failed")
				assert.Equal(t, parts[1], "1")
				timestamp, err := strconv.ParseInt(parts[2], 10, 64)
				assert.NoError(t, err)
				assert.True(t, timestamp >= 0)
			}
		case mifloraDataMetric:
			for line := range publish {
				parts := strings.Split(line, " ")
				assert.Equal(t, len(parts), 3)
				assert.Equal(t, strings.Index(parts[0], "foo.base.miflora.peri"), 0)
				assert.True(t, len(parts[1]) > 0)
				timestamp, err := strconv.ParseInt(parts[2], 10, 64)
				assert.NoError(t, err)
				assert.True(t, timestamp >= 0)
			}
		}
	}
}

func TestPublishInflux(t *testing.T) {
	tables := []struct {
		metric mifloraMetric
	}{
		{mifloraErrorMetric{peripheralId: "peri", failed: 1}},
		{mifloraDataMetric{
			peripheralId: "peri",
			metaData:     common.VersionBatteryResponse{BatteryLevel: 100, FirmwareVersion: "2.7.0"},
			sensorData:   common.SensorDataResponse{Temperature: 24.2, Brightness: 121, Moisture: 16, Conductivity: 101},
			connectTime:  3.42,
			readoutTime:  0.23,
			rssi:         -77,
		}},
	}

	for _, table := range tables {
		publish := make(chan string, 100)
		publishInflux(table.metric, publish)
		close(publish)
		switch table.metric.(type) {
		case mifloraErrorMetric:
			line := <-publish
			parts := strings.Split(line, " ")
			assert.Equal(t, len(parts), 3)
			assert.Equal(t, parts[0], "miflora,id=peri")
			assert.Equal(t, parts[1], "failed=1")
			timestamp, err := strconv.ParseInt(parts[2], 10, 64)
			assert.NoError(t, err)
			assert.True(t, timestamp >= 0)
		case mifloraDataMetric:
			line := <-publish
			parts := strings.Split(line, " ")
			assert.Equal(t, len(parts), 3)
			assert.Equal(t, parts[0], "miflora,id=peri")
			assert.True(t, len(parts[1]) > 0)
			timestamp, err := strconv.ParseInt(parts[2], 10, 64)
			assert.NoError(t, err)
			assert.True(t, timestamp >= 0)
		}
	}
}
