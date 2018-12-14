package gatt

import (
	"errors"

	"miflorad/common"

	"github.com/currantlabs/gatt"
)

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

func MifloraRequestVersionBattery(p gatt.Peripheral) (common.VersionBatteryResponse, error) {
	mifloraService := FindServiceByUUID(p.Services(), gatt.MustParseUUID(common.MifloraServiceUUID))
	if mifloraService == nil {
		return common.VersionBatteryResponse{}, errors.New("Failed to get the miflora service")
	}

	mifloraVersionBatteryChar := FindCharacteristicByUUID(mifloraService.Characteristics(), MifloraCharVersionBatteryUUID)
	if mifloraVersionBatteryChar == nil {
		return common.VersionBatteryResponse{}, errors.New("Failed to get the version battery characteristic")
	}

	bytes, err := p.ReadCharacteristic(mifloraVersionBatteryChar)
	if err != nil {
		return common.VersionBatteryResponse{}, err
	}

	return common.ParseVersionBattery(bytes), nil
}

func MifloraRequestModeChange(p gatt.Peripheral) error {
	mifloraService := FindServiceByUUID(p.Services(), gatt.MustParseUUID(common.MifloraServiceUUID))
	if mifloraService == nil {
		return errors.New("Failed to get the miflora service")
	}

	mifloraModeChangeChar := FindCharacteristicByUUID(mifloraService.Characteristics(), gatt.MustParseUUID(common.MifloraCharModeChangeUUID))
	if mifloraModeChangeChar == nil {
		return errors.New("Failed to discover the mode change characteristic")
	}

	err := p.WriteCharacteristic(mifloraModeChangeChar, common.MifloraGetModeChangeData(), false)
	if err != nil {
		return err
	}

	return nil
}

func MifloraRequstSensorData(p gatt.Peripheral) (common.SensorDataResponse, error) {
	mifloraService := FindServiceByUUID(p.Services(), gatt.MustParseUUID(common.MifloraServiceUUID))
	if mifloraService == nil {
		return common.SensorDataResponse{}, errors.New("Failed to get the miflora service")
	}

	mifloraSensorDataChar := FindCharacteristicByUUID(mifloraService.Characteristics(), gatt.MustParseUUID(common.MifloraCharReadSensorDataUUID))
	if mifloraSensorDataChar == nil {
		return common.SensorDataResponse{}, errors.New("Failed to discover the sensor data characteristic")
	}

	bytes, err := p.ReadCharacteristic(mifloraSensorDataChar)
	if err != nil {
		return common.SensorDataResponse{}, err
	}

	return common.ParseSensorData(bytes), nil
}
