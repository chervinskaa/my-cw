package app

import (
	"log"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type MeasurementService interface {
	Save(m domain.Measurement) (domain.Measurement, error)
	FindByDeviceAndDate(deviceId uint64, startDate, endDate time.Time) ([]domain.Measurement, error)
	Update(m domain.Measurement) (domain.Measurement, error)
	Delete(id uint64) error
}

type measurementService struct {
	measurementRepo database.MeasurementRepository
}

func NewMeasurementService(mr database.MeasurementRepository) MeasurementService {
	return &measurementService{
		measurementRepo: mr,
	}
}

func (s *measurementService) Save(m domain.Measurement) (domain.Measurement, error) {
	measurement, err := s.measurementRepo.Save(m)
	if err != nil {
		log.Printf("MeasurementService: %s", err)
		return domain.Measurement{}, err
	}

	return measurement, nil
}

func (s *measurementService) FindByDeviceAndDate(deviceId uint64, startDate, endDate time.Time) ([]domain.Measurement, error) {
	measurements, err := s.measurementRepo.FindByDeviceAndDate(deviceId, startDate, endDate)
	if err != nil {
		log.Printf("MeasurementService: %s", err)
		return nil, err
	}

	return measurements, nil
}

func (s *measurementService) Update(m domain.Measurement) (domain.Measurement, error) {
	measurement, err := s.measurementRepo.Update(m)
	if err != nil {
		log.Printf("MeasurementService: %s", err)
		return domain.Measurement{}, err
	}

	return measurement, nil
}

func (s *measurementService) Delete(id uint64) error {
	err := s.measurementRepo.Delete(id)
	if err != nil {
		log.Printf("MeasurementService: %s", err)
		return err
	}

	return nil
}
