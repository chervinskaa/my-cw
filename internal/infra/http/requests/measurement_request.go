package requests

import (
	"errors"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type MeasurementRequest struct {
	DeviceId uint64  `json:"device_id" validate:"required"`
	Value    float64 `json:"value" validate:"required"`
}

func (r MeasurementRequest) ToDomainModel(roomId uint64) (domain.Measurement, error) {

	if r.DeviceId == 0 {
		return domain.Measurement{}, errors.New("device_id is required")
	}
	if roomId == 0 {
		return domain.Measurement{}, errors.New("room_id is required")
	}
	if r.Value == 0 {
		return domain.Measurement{}, errors.New("value is required")
	}

	return domain.Measurement{
		DeviceId:    r.DeviceId,
		RoomId:      roomId,
		Value:       r.Value,
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}, nil
}
