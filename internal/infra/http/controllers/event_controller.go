package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
)

type EventController struct {
	eventService app.EventService
	deviceRepo   database.DeviceRepository
}

func NewEventController(es app.EventService, dr database.DeviceRepository) *EventController {
	return &EventController{
		eventService: es,
		deviceRepo:   dr,
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

		deviceDomain, err := c.deviceRepo.Find(event.DeviceId)
		if err != nil {
			log.Printf("EventController: Error fetching device: %s", err)
			http.Error(w, "Failed to fetch device", http.StatusInternalServerError)
			return
		}

		event.RoomId = deviceDomain.RoomId

		createdEvent, err := c.eventService.Save(event)
		if err != nil {
			log.Printf("EventController: %s", err)
			http.Error(w, "Failed to save event", http.StatusInternalServerError)
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

func (c *EventController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := c.eventService.FindAll()
		if err != nil {
			log.Printf("EventController: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		eventDtos := make([]resources.EventDto, len(events))
		for i, event := range events {
			eventDtos[i] = resources.EventDto{}.DomainToDto(event)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(eventDtos)
	}
}

func (c *EventController) GetPowerConsumptionByRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomIDParam := chi.URLParam(r, "roomId")
		if roomIDParam == "" {
			log.Printf("EventController: Room ID is empty")
			http.Error(w, "Room ID is empty", http.StatusBadRequest)
			return
		}

		roomID, err := strconv.ParseUint(roomIDParam, 10, 64)
		if err != nil {
			log.Printf("EventController: Error parsing room ID: %s", err)
			http.Error(w, "Invalid room ID", http.StatusBadRequest)
			return
		}

		startDateParam := r.URL.Query().Get("startDate")
		endDateParam := r.URL.Query().Get("endDate")

		startDate, err := time.Parse("2006-01-02", startDateParam)
		if err != nil {
			log.Printf("EventController: Error parsing start date: %s", err)
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateParam)
		if err != nil {
			log.Printf("EventController: Error parsing end date: %s", err)
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}

		totalPowerConsumption, err := c.eventService.GetPowerConsumptionByRoom(roomID, startDate, endDate)
		if err != nil {
			log.Printf("EventController: Error calculating power consumption by room: %s", err)
			http.Error(w, "Failed to calculate power consumption by room", http.StatusInternalServerError)
			return
		}

		response := map[string]float64{"total_power_consumption": totalPowerConsumption}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
