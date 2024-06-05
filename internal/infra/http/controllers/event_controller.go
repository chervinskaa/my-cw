package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type EventController struct {
	eventService app.EventService
}

func NewEventController(es app.EventService) *EventController {
	return &EventController{
		eventService: es,
	}
}

func (c *EventController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var eventRequest requests.EventRequest
		err := json.NewDecoder(r.Body).Decode(&eventRequest)
		if err != nil {
			log.Printf("EventController: Error decoding request body: %s", err)
			BadRequest(w, errors.New("invalid request payload"))
			return
		}

		event, err := eventRequest.ToDomainModel()
		if err != nil {
			log.Printf("EventController: Error converting to domain model: %s", err)
			BadRequest(w, err)
			return
		}

		createdEvent, err := c.eventService.Save(event)
		if err != nil {
			log.Printf("EventController: %s", err)
			InternalServerError(w, errors.New("failed to save event"))
			return
		}

		eventDto := resources.EventDto{}.DomainToDto(createdEvent)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(eventDto)
	}
}

func (c *EventController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event := r.Context().Value(EventKey).(domain.Event)

		eventDto := resources.EventDto{}.DomainToDto(event)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(eventDto)
	}
}
