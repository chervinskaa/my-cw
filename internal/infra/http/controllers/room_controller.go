package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
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

func (c *RoomController) FindByOrgId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParam(r, "orgId")
		orgId, err := strconv.ParseUint(orgIdParam, 10, 64)
		if err != nil {
			log.Printf("RoomController: Invalid organization ID: %s", err)
			BadRequest(w, errors.New("invalid organization ID"))
			return
		}

		log.Printf("RoomController: Looking for rooms with organization ID: %d", orgId)

		rooms, err := c.roomService.FindByOrgId(orgId)
		if err != nil {
			log.Printf("RoomController: Error finding rooms: %s", err)
			InternalServerError(w, errors.New("failed to retrieve rooms"))
			return
		}

		if len(rooms) == 0 {
			log.Printf("RoomController: No rooms found for organization ID: %d", orgId)
		}

		var roomsDto []resources.RoomDto
		for _, room := range rooms {
			roomDto := resources.RoomDto{}.DomainToDto(room)
			roomsDto = append(roomsDto, roomDto)
		}

		log.Printf("RoomController: Found %d rooms for organization ID: %d", len(roomsDto), orgId)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(roomsDto)
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
		org, ok := r.Context().Value(OrgKey).(domain.Organization)
		if !ok {
			log.Printf("RoomController: failed to retrieve organization from context")
			InternalServerError(w, errors.New("failed to retrieve organization from context"))
			return
		}

		room, ok := r.Context().Value(RoKey).(domain.Room)
		if !ok {
			log.Printf("RoomController: failed to retrieve room from context")
			InternalServerError(w, errors.New("failed to retrieve room from context"))
			return
		}

		if room.OrganizationId != org.Id {
			err := fmt.Errorf("access denied")
			Forbidden(w, err)
			return
		}

		err := c.roomService.Delete(room.Id)
		if err != nil {
			log.Printf("RoomController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
