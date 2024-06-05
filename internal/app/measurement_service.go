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
	Update(dm domain.Measurement) (domain.Measurement, error)
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

func (s *measurementService) Update(dm domain.Measurement) (domain.Measurement, error) {
	log.Printf("MeasurementService: Updating measurement %+v", dm)

	measurement, err := s.measurementRepo.Update(dm)
	if err != nil {
		log.Printf("MeasurementService: Error updating measurement: %s", err)
		return domain.Measurement{}, err
	}

	log.Printf("MeasurementService: Measurement updated successfully: %+v", measurement)
	return measurement, nil
}

func (s *measurementService) Delete(id uint64) error {
	log.Printf("MeasurementService: Deleting measurement with ID %d", id)

	err := s.measurementRepo.Delete(id)
	if err != nil {
		log.Printf("MeasurementService: Error deleting measurement: %s", err)
		return err
	}

	log.Printf("MeasurementService: Measurement deleted successfully")
	return nil
}
