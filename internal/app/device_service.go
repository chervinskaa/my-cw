package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type DeviceService interface {
	Save(o domain.Device) (domain.Device, error)
	FindAll() ([]domain.Device, error)
	FindById(id uint64) (domain.Device, error)
	Update(o domain.Device) (domain.Device, error)
	InstallDevice(deviceId uint64, roomId uint64) error
	UninstallDevice(deviceId uint64) error
	Delete(id uint64) error
}

type deviceService struct {
	deviceRepo database.DeviceRepository
}

func NewDeviceService(dr database.DeviceRepository) DeviceService {
	return &deviceService{
		deviceRepo: dr,
	}
}

func (s *deviceService) Save(r domain.Device) (domain.Device, error) {
	device, err := s.deviceRepo.Save(r)
	if err != nil {
		log.Printf("DeviceService: %s", err)
		return domain.Device{}, err
	}

	return device, nil
}

func (s *deviceService) FindAll() ([]domain.Device, error) {
	devices, err := s.deviceRepo.FindAll()
	if err != nil {
		log.Printf("DeviceService: %s", err)
		return nil, err
	}

	return devices, nil
}

func (s *deviceService) FindById(id uint64) (domain.Device, error) {
	device, err := s.deviceRepo.FindById(id)
	if err != nil {
		log.Printf("DeviceService: %s", err)
		return domain.Device{}, err
	}

	return device, nil
}

func (s *deviceService) Update(r domain.Device) (domain.Device, error) {
	device, err := s.deviceRepo.Update(r)
	if err != nil {
		log.Printf("DeviceService: %s", err)
		return domain.Device{}, err
	}

	return device, nil
}

func (s *deviceService) InstallDevice(deviceId uint64, roomId uint64) error {
	err := s.deviceRepo.InstallDevice(deviceId, roomId)
	if err != nil {
		log.Printf("DeviceService: %s", err)
		return err
	}

	return nil
}

func (s *deviceService) UninstallDevice(deviceId uint64) error {
	err := s.deviceRepo.Delete(deviceId)
	if err != nil {
		log.Printf("DeviceService: %s", err)
		return err
	}

	return nil
}

func (s *deviceService) Delete(id uint64) error {
	err := s.deviceRepo.Delete(id)
	if err != nil {
		log.Printf("DeviceService: %s", err)
		return err
	}

	return nil
}
