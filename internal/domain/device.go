package domain

import "time"

type DeviceCategory string

const (
	SENSOR   DeviceCategory = "SENSOR"
	ACTUATOR DeviceCategory = "ACTUATOR"
)

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
