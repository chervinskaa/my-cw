package requests

import (
	"errors"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type DeviceRequest struct {
	OrganizationId   uint64   `json:"organizationId"`
	RoomId           *uint64  `json:"room_id"`
	InventoryNumber  string   `json:"inventoryNumber"`
	SerialNumber     string   `json:"serialNumber"`
	Characteristics  string   `json:"characteristics" validate:"required"`
	PowerConsumption *float64 `json:"power_consumption" validate:"omitempty"`
	Units            *string  `json:"units" validate:"omitempty"`
	Category         string   `json:"category" validate:"required"`
}

func (r DeviceRequest) ToDomainModel() (domain.Device, error) {
	if r.Category == "ACTUATOR" && r.PowerConsumption == nil {
		return domain.Device{}, errors.New("power consumption is required for ACTUATOR")
	}
	if r.Category == "SENSOR" && r.Units == nil {
		return domain.Device{}, errors.New("units is required for SENSOR")
	}

	category, err := domain.ParseDeviceCategory(r.Category)
	if err != nil {
		return domain.Device{}, err
	}

	return domain.Device{
		OrganizationId:   r.OrganizationId,
		RoomId:           r.RoomId,
		InventoryNumber:  r.InventoryNumber,
		SerialNumber:     r.SerialNumber,
		Characteristics:  r.Characteristics,
		PowerConsumption: r.PowerConsumption,
		Units:            r.Units,
		Category:         category,
	}, nil
}
