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

type RoomController struct {
	roomService         app.RoomService
	organizationService app.OrganizationService
}

func NewRoomController(rs app.RoomService, os app.OrganizationService) *RoomController {
	return &RoomController{
		roomService:         rs,
		organizationService: os,
	}
}

func (c *RoomController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var roomRequest requests.RoomRequest
		err := json.NewDecoder(r.Body).Decode(&roomRequest)
		if err != nil {
			log.Printf("RoomController: Error decoding request body: %s", err)
			BadRequest(w, errors.New("invalid request payload"))
			return
		}

		if roomRequest.OrganizationId == 0 {
			err := errors.New("organizationId is required")
			log.Printf("RoomController: %s", err)
			BadRequest(w, err)
			return
		}

		_, err = c.organizationService.Find(roomRequest.OrganizationId)
		if err != nil {
			log.Printf("RoomController: Error finding organization: %s", err)
			BadRequest(w, errors.New("organization not found"))
			return
		}

		room := domain.Room{
			OrganizationId: roomRequest.OrganizationId,
			Name:           roomRequest.Name,
			Description:    roomRequest.Description,
		}

		createdRoom, err := c.roomService.Save(room)
		if err != nil {
			log.Printf("RoomController: %s", err)
			InternalServerError(w, errors.New("failed to save room"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdRoom)
	}
}

func (c *RoomController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		room := r.Context().Value(RoKey).(domain.Room)

		roomDto := resources.RoomDto{}.DomainToDto(room)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(roomDto)
	}
}

func (c *RoomController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := c.roomService.FindAll()
		if err != nil {
			log.Printf("RoomController: Error finding all rooms: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		roomDtos := make([]resources.RoomDto, len(rooms))
		for i, room := range rooms {
			roomDtos[i] = resources.RoomDto{}.DomainToDto(room)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(roomDtos); err != nil {
			log.Printf("RoomController: Error encoding response: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (c *RoomController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		room := r.Context().Value(RoKey).(domain.Room)

		var roomRequest requests.RoomRequest
		err := json.NewDecoder(r.Body).Decode(&roomRequest)
		if err != nil {
			log.Printf("RoomController: Error decoding request body: %s", err)
			BadRequest(w, errors.New("invalid request payload"))
			return
		}

		room.Name = roomRequest.Name
		room.Description = roomRequest.Description

		updatedRoom, err := c.roomService.Update(room)
		if err != nil {
			log.Printf("RoomController: Error updating room: %s", err)
			InternalServerError(w, errors.New("failed to update room"))
			return
		}

		roomDto := resources.RoomDto{}.DomainToDto(updatedRoom)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(roomDto)
	}
}

func (c *RoomController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ro := r.Context().Value(RoKey).(domain.Room)

		err := c.roomService.Delete(ro.Id)
		if err != nil {
			log.Printf("RoomController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
