package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericFirmwareVersion(t *testing.T) {
	tables := []struct {
		metaData VersionBatteryResponse
		firmware int
	}{
		{VersionBatteryResponse{BatteryLevel: 99, FirmwareVersion: "1.0.0"}, 10000},
		{VersionBatteryResponse{BatteryLevel: 88, FirmwareVersion: "2.6.6"}, 20606},
		{VersionBatteryResponse{BatteryLevel: 77, FirmwareVersion: "0.1.0"}, 100},
		{VersionBatteryResponse{BatteryLevel: 66, FirmwareVersion: "1.x.5"}, 10005},
		{VersionBatteryResponse{BatteryLevel: 55, FirmwareVersion: "fubar"}, 0},
	}

	for _, table := range tables {
		assert.Equal(t, table.metaData.NumericFirmwareVersion(), table.firmware)
	}
}
