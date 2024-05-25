package resources

import (
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type RoomsDto struct {
	Rooms []RoomDto `json:"rooms"`
}

type RoomDto struct {
	Id             uint64 `json:"id"`
	OrganizationId uint64 `json:"organization_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
}

func (d RoomDto) DomainToDto(o domain.Room) RoomDto {
	return RoomDto{
		Id:             o.Id,
		OrganizationId: o.OrganizationId,
		Name:           o.Name,
		Description:    o.Description,
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
