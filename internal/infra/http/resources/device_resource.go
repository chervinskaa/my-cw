package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type DevicesDto struct {
	Devices []DeviceDto `json:"devices"`
}

type DeviceDto struct {
	Id               uint64           `json:"id"`
	OrganizationId   uint64           `json:"organizationId"`
	RoomId           *uint64          `json:"room_id"`
	GUID             string           `json:"guid"`
	InventoryNumber  string           `json:"inventoryNumber"`
	SerialNumber     string           `json:"serialNumber"`
	Characteristics  string           `json:"characteristics"`
	Category         string           `json:"category"`
	Units            *string          `json:"units"`
	Measurements     []MeasurementDto `json:"measurements"`
	PowerConsumption *float64         `json:"power_consumption"`
	Events           []EventDto       `json:"events"`
	CreatedDate      time.Time        `json:"createdDate"`
	UpdatedDate      time.Time        `json:"updatedDate"`
}

func (d DeviceDto) DomainToDto(o domain.Device) DeviceDto {
	var measurements []MeasurementDto
	for _, dm := range o.Measurements {
		mDto := MeasurementDto{}.DomainToDto(dm)
		measurements = append(measurements, mDto)
	}
	var events []EventDto
	for _, de := range o.Events {
		eDto := EventDto{}.DomainToDto(de)
		events = append(events, eDto)
	}
	return DeviceDto{
		Id:               o.Id,
		OrganizationId:   o.OrganizationId,
		RoomId:           o.RoomId,
		GUID:             o.GUID,
		InventoryNumber:  o.InventoryNumber,
		SerialNumber:     o.SerialNumber,
		Characteristics:  o.Characteristics,
		Category:         string(o.Category),
		Units:            o.Units,
		PowerConsumption: o.PowerConsumption,
		CreatedDate:      o.CreatedDate,
		UpdatedDate:      o.UpdatedDate,
	}
}

func (d DevicesDto) DomainToDto(devices []domain.Device) DevicesDto {
	var deviceDtos []DeviceDto
	for _, o := range devices {
		deviceDto := DeviceDto{}.DomainToDto(o)
		deviceDtos = append(deviceDtos, deviceDto)
	}
	return DevicesDto{Devices: deviceDtos}
}
