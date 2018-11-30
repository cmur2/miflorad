package common

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/currantlabs/gatt"
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

func MifloraGetModeChangeData() []byte {
	return []byte{0xa0, 0x1f}
}

var MifloraServiceUUID = gatt.MustParseUUID("00001204-0000-1000-8000-00805f9b34fb")
var MifloraCharModeChangeUUID = gatt.MustParseUUID("00001a00-0000-1000-8000-00805f9b34fb")
var MifloraCharReadSensorDataUUID = gatt.MustParseUUID("00001a01-0000-1000-8000-00805f9b34fb")
var MifloraCharVersionBatteryUUID = gatt.MustParseUUID("00001a02-0000-1000-8000-00805f9b34fb")

func FindServiceByUUID(services []*gatt.Service, u gatt.UUID) *gatt.Service {
	for _, service := range services {
		if service.UUID().Equal(u) {
			return service
		}
	}
	return nil
}

func FindCharacteristicByUUID(characteristics []*gatt.Characteristic, u gatt.UUID) *gatt.Characteristic {
	for _, characteristic := range characteristics {
		if characteristic.UUID().Equal(u) {
			return characteristic
		}
	}
	return nil
}

func MifloraRequestVersionBattery(p gatt.Peripheral) (VersionBatteryResponse, error) {
	mifloraService := FindServiceByUUID(p.Services(), MifloraServiceUUID)
	if mifloraService == nil {
		return VersionBatteryResponse{}, errors.New("Failed to get the miflora service")
	}

	mifloraVersionBatteryChar := FindCharacteristicByUUID(mifloraService.Characteristics(), MifloraCharVersionBatteryUUID)
	if mifloraVersionBatteryChar == nil {
		return VersionBatteryResponse{}, errors.New("Failed to get the version battery characteristic")
	}

	bytes, err := p.ReadCharacteristic(mifloraVersionBatteryChar)
	if err != nil {
		return VersionBatteryResponse{}, err
	}

	return VersionBatteryResponse{string(bytes[2:]), uint8(bytes[0])}, nil
}

func MifloraRequestModeChange(p gatt.Peripheral) error {
	mifloraService := FindServiceByUUID(p.Services(), MifloraServiceUUID)
	if mifloraService == nil {
		return errors.New("Failed to get the miflora service")
	}

	mifloraModeChangeChar := FindCharacteristicByUUID(mifloraService.Characteristics(), MifloraCharModeChangeUUID)
	if mifloraModeChangeChar == nil {
		return errors.New("Failed to discover the mode change characteristic")
	}

	err := p.WriteCharacteristic(mifloraModeChangeChar, MifloraGetModeChangeData(), false)
	if err != nil {
		return err
	}

	return nil
}

func MifloraRequstSensorData(p gatt.Peripheral) (SensorDataResponse, error) {
	mifloraService := FindServiceByUUID(p.Services(), MifloraServiceUUID)
	if mifloraService == nil {
		return SensorDataResponse{}, errors.New("Failed to get the miflora service")
	}

	mifloraSensorDataChar := FindCharacteristicByUUID(mifloraService.Characteristics(), MifloraCharReadSensorDataUUID)
	if mifloraSensorDataChar == nil {
		return SensorDataResponse{}, errors.New("Failed to discover the sensor data characteristic")
	}

	bytes, err := p.ReadCharacteristic(mifloraSensorDataChar)
	if err != nil {
		return SensorDataResponse{}, err
	}

	return SensorDataResponse{
		Temperature:  float64(binary.LittleEndian.Uint16(bytes[0:2])) / 10.0,
		Brightness:   binary.LittleEndian.Uint32(bytes[3:7]),
		Moisture:     uint8(bytes[7]),
		Conductivity: binary.LittleEndian.Uint16(bytes[8:10]),
	}, nil
}

func (res VersionBatteryResponse) NumericFirmwareVersion() int {
	// turns "2.3.4" into 20304
	version := 0
	parts := strings.Split(res.FirmwareVersion, ".")
	for i, part := range parts {
		partNumber, err := strconv.Atoi(part)
		if err != nil {
			version += int(math.Pow10((len(parts)-(i+1))*2)) * partNumber
		}
	}
	return version
}
