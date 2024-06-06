package database

import (
	"errors"
	"log"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const EventsTableName = "events"

type event struct {
	Id          uint64     `db:"id,omitempty"`
	DeviceId    uint64     `db:"device_id"`
	RoomId      *uint64    `db:"room_id"`
	Action      string     `db:"action"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date"`
}

type EventRepository interface {
	Save(de domain.Event) (domain.Event, error)
	Find(id uint64) (domain.Event, error)
	FindByDeviceId(deviceId uint64) ([]domain.Event, error)
	FindAll() ([]domain.Event, error)
}

type eventRepository struct {
	coll db.Collection
}

func NewEventRepository(sess db.Session) EventRepository {
	return &eventRepository{
		coll: sess.Collection(EventsTableName),
	}
}

func (r *eventRepository) Save(de domain.Event) (domain.Event, error) {
	if de.Action != domain.TurnOn && de.Action != domain.TurnOff {
		return domain.Event{}, errors.New("invalid action")
	}

	deviceRepo := NewDeviceRepository(r.coll.Session())
	device, err := deviceRepo.Find(de.DeviceId)
	if err != nil {
		log.Printf("EventRepository: Error fetching device: %s", err)
		return domain.Event{}, err
	}

	if device.Category != domain.Actuator {
		err := errors.New("only actuators can have events")
		log.Printf("EventRepository: %s", err)
		return domain.Event{}, err
	}

	event := r.mapDomainToModel(de)
	now := time.Now()
	event.CreatedDate, event.UpdatedDate = now, now
	log.Printf("EventRepository: Saving event %+v", event)
	err = r.coll.InsertReturning(&event)
	if err != nil {
		log.Printf("EventRepository: Error saving event: %s", err)
		return domain.Event{}, err
	}
	de = r.mapModelToDomain(event)
	log.Printf("EventRepository: Saved event %+v", de)
	return de, nil
}

func (r *eventRepository) Find(id uint64) (domain.Event, error) {
	var eventModel event
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&eventModel)
	if err != nil {
		return domain.Event{}, err
	}
	return r.mapModelToDomain(eventModel), nil
}

func (r *eventRepository) FindAll() ([]domain.Event, error) {
	var events []event
	err := r.coll.Find(db.Cond{"deleted_date": nil}).All(&events)
	if err != nil {
		return nil, err
	}
	res := r.mapModelToDomainCollection(events)
	return res, nil
}

func (r *eventRepository) FindByDeviceId(deviceId uint64) ([]domain.Event, error) {
	var events []event
	err := r.coll.Find(db.Cond{"device_id": deviceId, "deleted_date": nil}).All(&events)
	if err != nil {
		return nil, err
	}
	return r.mapModelToDomainCollection(events), nil
}

func (r *eventRepository) mapDomainToModel(e domain.Event) event {
	return event{
		Id:          e.Id,
		DeviceId:    e.DeviceId,
		RoomId:      e.RoomId,
		Action:      string(e.Action),
		CreatedDate: e.CreatedDate,
		UpdatedDate: e.UpdatedDate,
		DeletedDate: e.DeletedDate,
	}
}

func (r *eventRepository) mapModelToDomain(e event) domain.Event {
	return domain.Event{
		Id:          e.Id,
		DeviceId:    e.DeviceId,
		RoomId:      e.RoomId,
		Action:      domain.EventAction(e.Action),
		CreatedDate: e.CreatedDate,
		UpdatedDate: e.UpdatedDate,
		DeletedDate: e.DeletedDate,
	}
}

func (r *eventRepository) mapModelToDomainCollection(eves []event) []domain.Event {
	var events []domain.Event
	for _, e := range eves {
		domainEvent := r.mapModelToDomain(e)
		events = append(events, domainEvent)
	}
	return events
}
