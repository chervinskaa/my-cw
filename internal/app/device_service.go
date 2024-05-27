package app

import (
	"fmt"
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
	"github.com/google/uuid"
)

type DeviceService interface {
	Save(d domain.Device) (domain.Device, error)
	FindByRoomId(roomId uint64) ([]domain.Device, error)
	Find(id uint64) (interface{}, error)
	Update(d domain.Device) (domain.Device, error)
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
	r.GUID = uuid.New().String()
	createdDevice, err := s.deviceRepo.Save(r)
	if err != nil {
		log.Printf("DeviceService: Error saving device: %s", err)
		return domain.Device{}, err
	}

	log.Printf("DeviceService: Device saved successfully: %+v", createdDevice)
	return createdDevice, nil
}

func (s *deviceService) FindByRoomId(roomId uint64) ([]domain.Device, error) {
	devices, err := s.deviceRepo.FindByRoomId(roomId)
	if err != nil {
		log.Printf("DeviceService: Error finding devices for organization ID %d: %s", roomId, err)
		return nil, err
	}

	if len(devices) == 0 {
		log.Printf("DeviceService: No devices found for organization ID %d", roomId)
	}

	log.Printf("DeviceService: Found %d devices for organization ID %d", len(devices), roomId)
	return devices, nil
}

func (s *deviceService) Find(id uint64) (interface{}, error) {
	log.Printf("DeviceService: Finding device with ID %d", id)

	device, err := s.deviceRepo.Find(id)
	if err != nil {
		log.Printf("DeviceService: Error finding device: %s", err)
		return domain.Device{}, err
	}

	if device.Id == 0 {
		log.Printf("DeviceService: No device found with ID %d", id)
		return domain.Device{}, fmt.Errorf("device not found")
	}

	log.Printf("DeviceService: Found device successfully: %+v", device)
	return device, nil
}

func (s *deviceService) Update(r domain.Device) (domain.Device, error) {
	log.Printf("DeviceService: Updating device %+v", r)

	device, err := s.deviceRepo.Update(r)
	if err != nil {
		log.Printf("DeviceService: Error updating device: %s", err)
		return domain.Device{}, err
	}

	log.Printf("DeviceService: Device updated successfully: %+v", device)
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
	err := s.deviceRepo.UninstallDevice(deviceId)
	if err != nil {
		log.Printf("DeviceService: %s", err)
		return err
	}

	return nil
}

func (s *deviceService) Delete(id uint64) error {
	log.Printf("DeviceService: Deleting device with ID %d", id)

	err := s.deviceRepo.Delete(id)
	if err != nil {
		log.Printf("DeviceService: Error deleting device: %s", err)
		return err
	}

	log.Printf("DeviceService: Device deleted successfully")
	return nil
}
