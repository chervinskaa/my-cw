package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type RoomsDto struct {
	Rooms []RoomDto `json:"rooms"`
}

type RoomDto struct {
	Id             uint64    `json:"id"`
	OrganizationId uint64    `json:"organizationId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CreatedDate    time.Time `json:"createdDate"`
	UpdatedDate    time.Time `json:"updatedDate"`
}

func (d RoomDto) DomainToDto(o domain.Room) RoomDto {
	return RoomDto{
		Id:             o.Id,
		OrganizationId: o.OrganizationId,
		Name:           o.Name,
		Description:    o.Description,
		CreatedDate:    o.CreatedDate,
		UpdatedDate:    o.UpdatedDate,
	}
}

func (d RoomsDto) DomainToDto(rooms []domain.Room) RoomsDto {
	var roomDtos []RoomDto
	for _, o := range rooms {
		roomDto := RoomDto{}.DomainToDto(o)
		roomDtos = append(roomDtos, roomDto)
	}
	return RoomsDto{Rooms: roomDtos}
}
