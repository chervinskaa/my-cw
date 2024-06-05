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
	Find(id uint64) (interface{}, error)
	FindAll() ([]domain.Device, error)
	Update(d domain.Device) (domain.Device, error)
	InstallDevice(deviceId uint64, roomId uint64) error
	UninstallDevice(device domain.Device) (domain.Device, error)
	Delete(id uint64) error
}

type deviceService struct {
	deviceRepo      database.DeviceRepository
	measurementRepo database.MeasurementRepository
}

func NewDeviceService(dr database.DeviceRepository, mr database.MeasurementRepository) DeviceService {
	return &deviceService{
		deviceRepo:      dr,
		measurementRepo: mr,
	}
}

func (s *deviceService) Save(dd domain.Device) (domain.Device, error) {
	dd.GUID = uuid.New().String()
	createdDevice, err := s.deviceRepo.Save(dd)
	if err != nil {
		log.Printf("DeviceService: Error saving device: %s", err)
		return domain.Device{}, err
	}

	log.Printf("DeviceService: Device saved successfully: %+v", createdDevice)
	return createdDevice, nil
}

func (s *deviceService) FindAll() ([]domain.Device, error) {
	devices, err := s.deviceRepo.FindAll()
	if err != nil {
		log.Printf("DeviceService: %s", err)
		return nil, err
	}

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

func (s *deviceService) Update(dd domain.Device) (domain.Device, error) {
	log.Printf("DeviceService: Updating device %+v", dd)

	device, err := s.deviceRepo.Update(dd)
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

func (s *deviceService) UninstallDevice(device domain.Device) (domain.Device, error) {
	uninstalledDevice, err := s.deviceRepo.UninstallDevice(device)
	if err != nil {
		log.Printf("DeviceService: Error uninstalling device: %s", err)
		return domain.Device{}, err
	}
	log.Printf("DeviceService: Uninstalled device with ID %d", device.Id)
	return uninstalledDevice, nil
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
