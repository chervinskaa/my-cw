package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type MeasurementDto struct {
	Id          uint64    `json:"id"`
	DeviceId    uint64    `json:"device_id"`
	RoomId      *uint64   `json:"room_id"`
	Value       float64   `json:"value"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

func (d MeasurementDto) DomainToDto(m domain.Measurement) MeasurementDto {
	return MeasurementDto{
		Id:          m.Id,
		DeviceId:    m.DeviceId,
		RoomId:      m.RoomId,
		Value:       m.Value,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
	}
}

func (d MeasurementDto) DomainToDtoCollection(measurements []domain.Measurement) []MeasurementDto {
	var measurementDtos []MeasurementDto
	for _, m := range measurements {
		measurementDto := MeasurementDto{}.DomainToDto(m)
		measurementDtos = append(measurementDtos, measurementDto)
	}
	return measurementDtos
}
