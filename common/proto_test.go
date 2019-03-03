package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersionBattery(t *testing.T) {
	tables := []struct {
		bytes    []byte
		metaData VersionBatteryResponse
	}{
		{[]byte{0x64, 0x15, 0x32, 0x2e, 0x37, 0x2e, 0x30}, VersionBatteryResponse{BatteryLevel: 100, FirmwareVersion: "2.7.0"}},
		{[]byte{0x64, 0xee, 0x32, 0x2e, 0x37, 0x2e, 0x30}, VersionBatteryResponse{BatteryLevel: 100, FirmwareVersion: "2.7.0"}},
		{[]byte{0x64, 0x42, 0x32, 0x2e, 0x37, 0x2e, 0x30}, VersionBatteryResponse{BatteryLevel: 100, FirmwareVersion: "2.7.0"}},
		{[]byte{0x50, 0x15, 0x32, 0x2e, 0x37, 0x2e, 0x30}, VersionBatteryResponse{BatteryLevel: 80, FirmwareVersion: "2.7.0"}},
		{[]byte{0x64, 0x42, 0x32, 0x2e, 0x36, 0x2e, 0x36}, VersionBatteryResponse{BatteryLevel: 100, FirmwareVersion: "2.6.6"}},
	}

	for _, table := range tables {
		assert.Equal(t, table.metaData, ParseVersionBattery(table.bytes))
	}
}

func TestParseSensorData(t *testing.T) {
	tables := []struct {
		bytes      []byte
		sensorData SensorDataResponse
	}{
		// Source: https://www.open-homeautomation.com/de/2016/08/23/reverse-engineering-the-mi-plant-sensor/
		{
			[]byte{0xf2, 0x00, 0x00, 0x79, 0x00, 0x00, 0x00, 0x10, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			SensorDataResponse{Temperature: 24.2, Brightness: 121, Moisture: 16, Conductivity: 101},
		},
		{
			[]byte{0x25, 0x01, 0x00, 0xf7, 0x26, 0x00, 0x00, 0x28, 0x0e, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			SensorDataResponse{Temperature: 29.3, Brightness: 9975, Moisture: 40, Conductivity: 270},
		},
	}

	for _, table := range tables {
		assert.Equal(t, table.sensorData, ParseSensorData(table.bytes))
	}
}
