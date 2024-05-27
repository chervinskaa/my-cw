package domain

import (
	"errors"
	"strings"
	"time"
)

type DeviceCategory string

const (
	Sensor   DeviceCategory = "SENSOR"
	Actuator DeviceCategory = "ACTUATOR"
)

func ParseDeviceCategory(category string) (DeviceCategory, error) {
	switch strings.ToUpper(category) {
	case "SENSOR":
		return Sensor, nil
	case "ACTUATOR":
		return Actuator, nil
	default:
		return "", errors.New("invalid device category")
	}
}

type Device struct {
	Id               uint64
	OrganizationId   uint64
	RoomId           *uint64
	GUID             string
	InventoryNumber  string
	SerialNumber     string
	Characteristics  string
	Category         DeviceCategory
	Units            *string
	PowerConsumption *float64
	CreatedDate      time.Time
	UpdatedDate      time.Time
	DeletedDate      *time.Time
}
