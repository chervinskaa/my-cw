package requests

import "github.com/BohdanBoriak/boilerplate-go-back/internal/domain"

type DeviceRequest struct {
	RoomId           *uint64  `json:"room_id"`
	Characteristics  string   `json:"characteristics" validate:"required"`
	PowerConsumption *float64 `json:"power_consumption" validate:"omitempty"`
	Units            *string  `json:"units" validate:"omitempty"`
}

func (r DeviceRequest) ToDomainModel() (domain.Device, error) {
	return domain.Device{
		RoomId:           r.RoomId,
		Characteristics:  r.Characteristics,
		PowerConsumption: r.PowerConsumption,
		Units:            r.Units,
	}, nil
}
