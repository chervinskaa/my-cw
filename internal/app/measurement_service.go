package app

import (
	"log"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type MeasurementService interface {
	Save(m domain.Measurement) (domain.Measurement, error)
	FindByDeviceAndDate(deviceId uint64, startDate, endDate time.Time) (interface{}, error)
	Find(id uint64) (interface{}, error)
	FindAll() ([]domain.Measurement, error)
}

type measurementService struct {
	measurementRepo database.MeasurementRepository
}

func NewMeasurementService(mr database.MeasurementRepository) MeasurementService {
	return &measurementService{
		measurementRepo: mr,
	}
}

func (s *measurementService) Save(dm domain.Measurement) (domain.Measurement, error) {
	createdMeasurement, err := s.measurementRepo.Save(dm)
	if err != nil {
		log.Printf("MeasurementService: Error saving measurement: %s", err)
		return domain.Measurement{}, err
	}

	log.Printf("MeasurementService: Measurement saved successfully: %+v", createdMeasurement)
	return createdMeasurement, nil
}

func (s *measurementService) FindByDeviceAndDate(deviceId uint64, startDate, endDate time.Time) (interface{}, error) {
	measurements, err := s.measurementRepo.FindByDeviceAndDate(deviceId, startDate, endDate)
	if err != nil {
		log.Printf("MeasurementService: %s", err)
		return nil, err
	}

	return measurements, nil
}

func (s *measurementService) Find(id uint64) (interface{}, error) {
	measurement, err := s.measurementRepo.Find(id)
	if err != nil {
		log.Printf("MeasurementService: %s", err)
		return domain.Measurement{}, err
	}

	return measurement, nil
}

func (s *measurementService) FindAll() ([]domain.Measurement, error) {
	measurements, err := s.measurementRepo.FindAll()
	if err != nil {
		log.Printf("MeasurementService: %s", err)
		return nil, err
	}

	return measurements, nil
}
