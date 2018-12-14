package common

import (
	"encoding/binary"
)

const (
	MifloraServiceUUID            = "00001204-0000-1000-8000-00805f9b34fb"
	MifloraCharModeChangeUUID     = "00001a00-0000-1000-8000-00805f9b34fb"
	MifloraCharReadSensorDataUUID = "00001a01-0000-1000-8000-00805f9b34fb"
	MifloraCharVersionBatteryUUID = "00001a02-0000-1000-8000-00805f9b34fb"
)

func MifloraGetModeChangeData() []byte {
	return []byte{0xa0, 0x1f}
}

func ParseVersionBattery(bytes []byte) VersionBatteryResponse {
	return VersionBatteryResponse{
		BatteryLevel:    uint8(bytes[0]),
		FirmwareVersion: string(bytes[2:]),
	}
}

func ParseSensorData(bytes []byte) SensorDataResponse {
	return SensorDataResponse{
		Temperature:  float64(binary.LittleEndian.Uint16(bytes[0:2])) / 10.0,
		Brightness:   binary.LittleEndian.Uint32(bytes[3:7]),
		Moisture:     uint8(bytes[7]),
		Conductivity: binary.LittleEndian.Uint16(bytes[8:10]),
	}
}
