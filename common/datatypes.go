package common

import (
	"math"
	"strconv"
	"strings"
)

type VersionBatteryResponse struct {
	FirmwareVersion string // as "x.y.z"
	BatteryLevel    uint8  // in percent 0-100
}

type SensorDataResponse struct {
	Temperature  float64 // in degree C
	Brightness   uint32  // in lux
	Moisture     uint8   // in percent 0-100
	Conductivity uint16  // in µS/cm
}

func (res VersionBatteryResponse) NumericFirmwareVersion() int {
	// turns "2.3.4" into 20304
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
