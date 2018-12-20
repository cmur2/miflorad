package common

import (
	"math"
	"strconv"
	"strings"
)

// captures response when reading meta data of miflora device
type VersionBatteryResponse struct {
	BatteryLevel    uint8  // in percent 0-100
	FirmwareVersion string // as "x.y.z"
}

// captures response when reading sensor data of miflora device
type SensorDataResponse struct {
	Temperature  float64 // in degree C
	Brightness   uint32  // in lux
	Moisture     uint8   // in percent 0-100
	Conductivity uint16  // in µS/cm
}

// turns firmware version "2.3.4" into 20304
func (res VersionBatteryResponse) NumericFirmwareVersion() int {
	version := 0
	parts := strings.Split(res.FirmwareVersion, ".")
	for i, part := range parts {
		partNumber, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		version += int(math.Pow10((len(parts)-(i+1))*2)) * partNumber
	}
	return version
}

// for the newer models a magic number must be written before we can read the current data
func (res VersionBatteryResponse) RequiresModeChangeBeforeRead() bool {
	return res.FirmwareVersion >= "2.6.6"
}
