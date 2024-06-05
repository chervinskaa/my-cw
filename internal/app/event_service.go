package app

import (
	"log"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type EventService interface {
	Save(de domain.Event) (domain.Event, error)
	Find(id uint64) (interface{}, error)
}

type eventService struct {
	eventRepo database.EventRepository
}

func NewEventService(er database.EventRepository) EventService {
	return &eventService{
		eventRepo: er,
	}
}

func (s *eventService) Save(de domain.Event) (domain.Event, error) {
	de.CreatedDate = time.Now()
	createdEvent, err := s.eventRepo.Save(de)
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
		log.Printf("EventService: %s", err)
		return domain.Event{}, err
	}

	return event, nil
}
