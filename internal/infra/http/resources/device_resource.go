package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type DevicesDto struct {
	Devices []DeviceDto `json:"devices"`
}

type DeviceDto struct {
	Id              uint64    `json:"id"`
	OrganizationId  uint64    `json:"organizationId"`
	GUID            string    `json:"guid"`
	InventoryNumber string    `json:"inventory_number"`
	SerialNumber    string    `json:"serial_number"`
	Category        string    `json:"category"`
	CreatedDate     time.Time `json:"createdDate"`
	UpdatedDate     time.Time `json:"updatedDate"`
}

func (d DeviceDto) DomainToDto(o domain.Device) DeviceDto {
	return DeviceDto{
		Id:              o.Id,
		OrganizationId:  o.OrganizationId,
		GUID:            o.GUID,
		InventoryNumber: o.InventoryNumber,
		SerialNumber:    o.SerialNumber,
		Category:        string(o.Category),
		CreatedDate:     o.CreatedDate,
		UpdatedDate:     o.UpdatedDate,
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
