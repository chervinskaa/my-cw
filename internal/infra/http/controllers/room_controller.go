package controllers

import (
	"encoding/json"
	"errors"
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
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if roomRequest.OrganizationId == 0 {
			err := errors.New("organizationId is required")
			log.Printf("RoomController: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = c.organizationService.Find(roomRequest.OrganizationId)
		if err != nil {
			log.Printf("RoomController: Error finding organization: %s", err)
			http.Error(w, "Organization not found", http.StatusBadRequest)
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
			http.Error(w, "Failed to save room", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdRoom)
	}
}

func (c RoomController) FindByOrgId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParam(r, "orgId")
		orgId, err := strconv.ParseUint(orgIdParam, 10, 64)
		if err != nil {
			log.Printf("RoomController: Invalid organization ID: %s", err)
			http.Error(w, "Invalid organization ID", http.StatusBadRequest)
			return
		}

		rooms, err := c.roomService.FindByOrgId(orgId)
		if err != nil {
			log.Printf("RoomController: Error finding rooms: %s", err)
			http.Error(w, "Failed to retrieve rooms", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rooms)
	}
}

func (c RoomController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := r.Context().Value(OrgKey).(domain.Organization)
		ro := r.Context().Value(RoKey).(domain.Room)

		if ro.OrganizationId != organization.Id {
			err := errors.New("access denied")
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		var roomDto resources.RoomDto
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(roomDto.DomainToDto(ro))
	}
}

func (c RoomController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := r.Context().Value(OrgKey).(domain.Organization)
		ro, err := requests.Bind(r, requests.RoomRequest{}, domain.Room{})
		if err != nil {
			log.Printf("RoomController: %s", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		room := r.Context().Value(RoKey).(domain.Room)
		if room.OrganizationId != organization.Id {
			err := errors.New("access denied")
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		room.Name = ro.Name
		room.Description = ro.Description
		updatedRoom, err := c.roomService.Update(room)
		if err != nil {
			log.Printf("RoomController: %s", err)
			http.Error(w, "Failed to update room", http.StatusInternalServerError)
			return
		}

		var roomDto resources.RoomDto
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(roomDto.DomainToDto(updatedRoom))
	}
}

func (c RoomController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := r.Context().Value(OrgKey).(domain.Organization)
		ro := r.Context().Value(RoKey).(domain.Room)

		if ro.OrganizationId != organization.Id {
			err := errors.New("access denied")
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		err := c.roomService.Delete(ro.Id)
		if err != nil {
			log.Printf("RoomController: %s", err)
			http.Error(w, "Failed to delete room", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
