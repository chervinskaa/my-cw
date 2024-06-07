package app

import (
	"log"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type EventService interface {
	Save(event domain.Event) (domain.Event, error)
	Find(id uint64) (interface{}, error)
	FindAll() ([]domain.Event, error)
	GetPowerConsumptionByRoom(roomID uint64, startDate, endDate time.Time) (float64, error)
}

type eventService struct {
	eventRepo  database.EventRepository
	deviceRepo database.DeviceRepository
}

func NewEventService(eventRepo database.EventRepository, deviceRepo database.DeviceRepository) EventService {
	return &eventService{
		eventRepo:  eventRepo,
		deviceRepo: deviceRepo,
	}
}

func (s *eventService) Save(event domain.Event) (domain.Event, error) {
	event.CreatedDate = time.Now()
	createdEvent, err := s.eventRepo.Save(event)
	if err != nil {
		log.Printf("EventService: Error saving event: %s", err)
		return domain.Event{}, err
	}

	log.Printf("EventService: Event saved successfully: %+v", createdEvent)
	return createdEvent, nil
}

func (s *eventService) Find(id uint64) (interface{}, error) {
	event, err := s.eventRepo.Find(id)
	if err != nil {
		log.Printf("EventService: Error finding event: %s", err)
		return domain.Event{}, err
	}

	return event, nil
}

func (s *eventService) FindAll() ([]domain.Event, error) {
	events, err := s.eventRepo.FindAll()
	if err != nil {
		log.Printf("EventService: Error finding all events: %s", err)
		return nil, err
	}

	return events, nil
}

func (s *eventService) GetPowerConsumptionByRoom(roomID uint64, startDate, endDate time.Time) (float64, error) {
	events, err := s.eventRepo.FindByRoomAndDate(roomID, startDate, endDate)
	if err != nil {
		log.Printf("EventService: Error finding events by room: %s", err)
		return 0, err
	}

	totalPC, err := s.calculatePowerConsumption(events, startDate, endDate)
	if err != nil {
		log.Printf("EventService: Error calculating power consumption by room: %s", err)
		return 0, err
	}

	return totalPC, nil
}

func (s *eventService) calculatePowerConsumption(events []domain.Event, startDate, endDate time.Time) (float64, error) {
	totalPowerConsumption := 0.0
	currentTime := time.Now()
	consumptionByDevice := make(map[uint64]float64)
	deviceOnTimes := make(map[uint64]time.Time)

	for _, event := range events {
		if event.Action == domain.TurnOn {
			deviceOnTimes[event.DeviceId] = event.CreatedDate
		} else if event.Action == domain.TurnOff {
			onTime, exists := deviceOnTimes[event.DeviceId]
			if exists {
				device, err := s.deviceRepo.Find(event.DeviceId)
				if err != nil {
					log.Printf("EventService: Error fetching device: %s", err)
					return 0, err
				}

				onTimeVal := maxTime(onTime, startDate)
				offTimeVal := minTime(event.CreatedDate, endDate)

				if onTimeVal.Before(offTimeVal) {
					duration := offTimeVal.Sub(onTimeVal).Hours()
					consumptionByDevice[event.DeviceId] += duration * *device.PowerConsumption
					log.Printf("Device ID: %d, OnTime: %v, OffTime: %v, Duration: %f hours, PowerConsumption: %f",
						event.DeviceId, onTimeVal, offTimeVal, duration, *device.PowerConsumption)
				}
				delete(deviceOnTimes, event.DeviceId)
			}
		}
	}

	for deviceId, onTime := range deviceOnTimes {
		device, err := s.deviceRepo.Find(deviceId)
		if err != nil {
			log.Printf("EventService: Error fetching device: %s", err)
			return 0, err
		}

		onTimeVal := maxTime(onTime, startDate)
		offTimeVal := minTime(currentTime, endDate)

		if onTimeVal.Before(offTimeVal) {
			duration := offTimeVal.Sub(onTimeVal).Hours()
			consumptionByDevice[deviceId] += duration * *device.PowerConsumption

			log.Printf("Device ID: %d, OnTime: %v, OffTime: %v, Duration: %f hours, PowerConsumption: %f",
				deviceId, onTimeVal, offTimeVal, duration, *device.PowerConsumption)
		}
	}

	for _, consumption := range consumptionByDevice {
		totalPowerConsumption += consumption
	}

	return totalPowerConsumption, nil
}

// Допоміжні функції для визначення мінімального та максимального часу
func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
