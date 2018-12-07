package main

import (
	"encoding/binary"
	"errors"

	"github.com/currantlabs/gatt"
)

type VersionBatteryResponse struct {
	firmware_version string // as "x.y.z"
	battery_level    uint8  // in percent 0-100
}

type SensorDataResponse struct {
	temperature  float64 // in degree C
	brightness   uint32  // in lux
	moisture     uint8   // in percent 0-100
	conductivity uint16  // in µS/cm
}

func mifloraGetModeChangeData() []byte {
	return []byte{0xa0, 0x1f}
}

var mifloraServiceUUID = gatt.MustParseUUID("00001204-0000-1000-8000-00805f9b34fb")
var mifloraCharModeChangeUUID = gatt.MustParseUUID("00001a00-0000-1000-8000-00805f9b34fb")
var mifloraCharReadSensorDataUUID = gatt.MustParseUUID("00001a01-0000-1000-8000-00805f9b34fb")
var mifloraCharVersionBatteryUUID = gatt.MustParseUUID("00001a02-0000-1000-8000-00805f9b34fb")

func findServiceByUUID(services []*gatt.Service, u gatt.UUID) *gatt.Service {
	for _, service := range services {
		if service.UUID().Equal(u) {
			return service
		}
	}
	return nil
}

func findCharacteristicByUUID(characteristics []*gatt.Characteristic, u gatt.UUID) *gatt.Characteristic {
	for _, characteristic := range characteristics {
		if characteristic.UUID().Equal(u) {
			return characteristic
		}
	}
	return nil
}

func mifloraRequestVersionBattery(p gatt.Peripheral) (VersionBatteryResponse, error) {
	mifloraService := findServiceByUUID(p.Services(), mifloraServiceUUID)
	if mifloraService == nil {
		return VersionBatteryResponse{}, errors.New("Failed to get the miflora service")
	}

	mifloraVersionBatteryChar := findCharacteristicByUUID(mifloraService.Characteristics(), mifloraCharVersionBatteryUUID)
	if mifloraVersionBatteryChar == nil {
		return VersionBatteryResponse{}, errors.New("Failed to get the version battery characteristic")
	}

	bytes, err := p.ReadCharacteristic(mifloraVersionBatteryChar)
	if err != nil {
		return VersionBatteryResponse{}, err
	}

	return VersionBatteryResponse{string(bytes[2:]), uint8(bytes[0])}, nil
}

func mifloraRequestModeChange(p gatt.Peripheral) error {
	mifloraService := findServiceByUUID(p.Services(), mifloraServiceUUID)
	if mifloraService == nil {
		return errors.New("Failed to get the miflora service")
	}

	mifloraModeChangeChar := findCharacteristicByUUID(mifloraService.Characteristics(), mifloraCharModeChangeUUID)
	if mifloraModeChangeChar == nil {
		return errors.New("Failed to discover the mode change characteristic")
	}

	err := p.WriteCharacteristic(mifloraModeChangeChar, mifloraGetModeChangeData(), false)
	if err != nil {
		return err
	}

	return nil
}

func mifloraRequstSensorData(p gatt.Peripheral) (SensorDataResponse, error) {
	mifloraService := findServiceByUUID(p.Services(), mifloraServiceUUID)
	if mifloraService == nil {
		return SensorDataResponse{}, errors.New("Failed to get the miflora service")
	}

	mifloraSensorDataChar := findCharacteristicByUUID(mifloraService.Characteristics(), mifloraCharReadSensorDataUUID)
	if mifloraSensorDataChar == nil {
		return SensorDataResponse{}, errors.New("Failed to discover the sensor data characteristic")
	}

	bytes, err := p.ReadCharacteristic(mifloraSensorDataChar)
	if err != nil {
		return SensorDataResponse{}, err
	}

	return SensorDataResponse{
		temperature:  float64(binary.LittleEndian.Uint16(bytes[0:2])) / 10.0,
		brightness:   binary.LittleEndian.Uint32(bytes[3:7]),
		moisture:     uint8(bytes[7]),
		conductivity: binary.LittleEndian.Uint16(bytes[8:10]),
	}, nil
}
