package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type RoomsDto struct {
	Rooms []RoomDto `json:"rooms"`
}

type RoomDto struct {
	Id             uint64      `json:"id"`
	OrganizationId uint64      `json:"organization_id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Devices        []DeviceDto `json:"devices"`
	CreatedDate    time.Time   `json:"createdDate"`
	UpdatedDate    time.Time   `json:"updatedDate"`
}

func (d RoomDto) DomainToDto(o domain.Room) RoomDto {
	var devices []DeviceDto
	for _, dd := range o.Devices {
		dDto := DeviceDto{}.DomainToDto(dd)
		devices = append(devices, dDto)
	}

	return RoomDto{
		Id:             o.Id,
		OrganizationId: o.OrganizationId,
		Name:           o.Name,
		Description:    o.Description,
		Devices:        devices,
		CreatedDate:    o.CreatedDate,
		UpdatedDate:    o.UpdatedDate,
	}
}

func (d RoomsDto) DomainToDto(rooms []domain.Room) RoomsDto {
	var roomDtos []RoomDto
	for _, room := range rooms {
		roomDto := RoomDto{}.DomainToDto(room)
		roomDtos = append(roomDtos, roomDto)
	}
	return RoomsDto{Rooms: roomDtos}
}
