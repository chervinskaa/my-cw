package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type EventsDto struct {
	Events []EventDto `json:"events"`
}

type EventDto struct {
	Id          uint64    `json:"id"`
	DeviceId    uint64    `json:"device_id"`
	RoomId      *uint64   `json:"room_id"`
	Action      string    `json:"action"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}

func (d EventDto) DomainToDto(o domain.Event) EventDto {
	return EventDto{
		Id:          o.Id,
		DeviceId:    o.DeviceId,
		RoomId:      o.RoomId,
		Action:      string(o.Action),
		CreatedDate: o.CreatedDate,
		UpdatedDate: o.UpdatedDate,
	}
}

func (d EventsDto) DomainToDto(events []domain.Event) EventsDto {
	var eventDtos []EventDto
	for _, o := range events {
		eventDto := EventDto{}.DomainToDto(o)
		eventDtos = append(eventDtos, eventDto)
	}
	return EventsDto{Events: eventDtos}
}
