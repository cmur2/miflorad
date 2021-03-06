package ble

import (
	"miflorad/common"

	"github.com/go-ble/ble"
	"github.com/pkg/errors"
)

func FindServiceByUUID(services []*ble.Service, uuid string) *ble.Service {
	u := ble.MustParse(uuid)
	for _, service := range services {
		if service.UUID.Equal(u) {
			return service
		}
	}
	return nil
}

func FindCharacteristicByUUID(characteristics []*ble.Characteristic, uuid string) *ble.Characteristic {
	u := ble.MustParse(uuid)
	for _, characteristic := range characteristics {
		if characteristic.UUID.Equal(u) {
			return characteristic
		}
	}
	return nil
}

func RequestVersionBattery(client ble.Client) (common.VersionBatteryResponse, error) {
	mifloraService := FindServiceByUUID(client.Profile().Services, common.MifloraServiceUUID)
	if mifloraService == nil {
		return common.VersionBatteryResponse{}, errors.New("Failed to get the miflora service")
	}

	mifloraVersionBatteryChar := FindCharacteristicByUUID(mifloraService.Characteristics, common.MifloraCharVersionBatteryUUID)
	if mifloraVersionBatteryChar == nil {
		return common.VersionBatteryResponse{}, errors.New("Failed to get the version battery characteristic")
	}

	bytes, err := client.ReadCharacteristic(mifloraVersionBatteryChar)
	if err != nil {
		return common.VersionBatteryResponse{}, errors.Wrap(err, "can't read version battery")
	}

	return common.ParseVersionBattery(bytes), nil
}

func RequestModeChange(client ble.Client) error {
	mifloraService := FindServiceByUUID(client.Profile().Services, common.MifloraServiceUUID)
	if mifloraService == nil {
		return errors.New("Failed to get the miflora service")
	}

	mifloraModeChangeChar := FindCharacteristicByUUID(mifloraService.Characteristics, common.MifloraCharModeChangeUUID)
	if mifloraModeChangeChar == nil {
		return errors.New("Failed to discover the mode change characteristic")
	}

	err := client.WriteCharacteristic(mifloraModeChangeChar, common.MifloraGetModeChangeData(), false)
	if err != nil {
		return errors.Wrap(err, "can't change mode")
	}

	return nil
}

func RequestSensorData(client ble.Client) (common.SensorDataResponse, error) {
	mifloraService := FindServiceByUUID(client.Profile().Services, common.MifloraServiceUUID)
	if mifloraService == nil {
		return common.SensorDataResponse{}, errors.New("Failed to get the miflora service")
	}

	mifloraSensorDataChar := FindCharacteristicByUUID(mifloraService.Characteristics, common.MifloraCharReadSensorDataUUID)
	if mifloraSensorDataChar == nil {
		return common.SensorDataResponse{}, errors.New("Failed to discover the sensor data characteristic")
	}

	bytes, err := client.ReadCharacteristic(mifloraSensorDataChar)
	if err != nil {
		return common.SensorDataResponse{}, errors.Wrap(err, "can't read sensor data")
	}

	return common.ParseSensorData(bytes), nil
}
