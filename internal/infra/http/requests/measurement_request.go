package requests

import (
	"errors"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type MeasurementRequest struct {
	DeviceId uint64  `json:"device_id" validate:"required"`
	RoomId   *uint64 `json:"room_id"`
	Value    float64 `json:"value" validate:"required"`
}

func (r MeasurementRequest) ToDomainModel() (domain.Measurement, error) {
	if r.DeviceId == 0 {
		return domain.Measurement{}, errors.New("device_id is required")
	}
	if r.Value == 0 {
		return domain.Measurement{}, errors.New("value is required")
	}

	measurement := domain.Measurement{
		DeviceId:    r.DeviceId,
		Value:       r.Value,
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}

	if r.RoomId != nil {
		measurement.RoomId = r.RoomId
	}

	return measurement, nil
}
