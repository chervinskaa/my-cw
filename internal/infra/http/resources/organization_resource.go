package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type OrgsDto struct {
	Organizations []OrgDto `json:"organizations"`
}

type OrgDto struct {
	Id          uint64    `json:"id"`
	UserId      uint64    `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	City        string    `json:"city"`
	Address     string    `json:"address"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	Rooms       []RoomDto `json:"rooms"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}

func (d OrgDto) DomainToDto(o domain.Organization) OrgDto {
	var rooms []RoomDto
	for _, r := range o.Rooms {
		rDto := RoomDto{}.DomainToDto(r)
		rooms = append(rooms, rDto)

	}
	return OrgDto{
		Id:          o.Id,
		UserId:      o.UserId,
		Name:        o.Name,
		Description: o.Description,
		City:        o.City,
		Address:     o.Address,
		Lat:         o.Lat,
		Lon:         o.Lon,
		Rooms:       rooms,
		CreatedDate: o.CreatedDate,
		UpdatedDate: o.UpdatedDate,
	}
}

func (d OrgsDto) DomainToDto(orgs []domain.Organization) OrgsDto {
	var organizations []OrgDto
	for _, o := range orgs {
		var oDto OrgDto
		org := oDto.DomainToDto(o)
		organizations = append(organizations, org)
	}
	response := OrgsDto{
		Organizations: organizations,
	}
	return response
}
