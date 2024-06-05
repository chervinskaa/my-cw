package requests

import (
	"errors"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type EventRequest struct {
	DeviceId uint64  `json:"device_id"`
	RoomId   *uint64 `json:"room_id"`
	Action   string  `json:"action" validate:"required,oneof='ON' 'OFF'"`
}

func (r EventRequest) ToDomainModel() (domain.Event, error) {
	var action domain.EventAction
	switch r.Action {
	case "ON":
		action = domain.TurnOn
	case "OFF":
		action = domain.TurnOff
	default:
		return domain.Event{}, errors.New("invalid action")
	}

	return domain.Event{
		DeviceId:    r.DeviceId,
		RoomId:      r.RoomId,
		Action:      action,
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}, nil
}
