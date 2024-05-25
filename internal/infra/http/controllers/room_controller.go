package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type RoomRequest struct {
}

type RoomController struct {
	roomService app.RoomService
}

func NewRoomController(rs app.RoomService) RoomController {
	if rs == nil {
		log.Fatal("NewRoomController: roomService is nil")
	}

	log.Printf("NewRoomController: roomService is initialized")

	return RoomController{
		roomService: rs,
	}
}

func (c RoomController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		room, err := requests.Bind(r, requests.RoomRequest{}, domain.Room{})
		if err != nil {
			log.Printf("OrganizationController: %s", err)
			BadRequest(w, err)
			return
		}

		room, err = c.roomService.Save(room, user.Id)
		if err != nil {
			log.Printf("OrganizationController: %s", err)
			if err.Error() == "access denied" {
				Forbidden(w, err)
			} else {
				InternalServerError(w, err)
			}
			return
		}

		var roomDto resources.RoomDto
		Created(w, roomDto.DomainToDto(room))
	}
}

func (c RoomController) FindForOrganization() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization, ok := r.Context().Value(OrgKey).(domain.Organization)
		if !ok || organization.Id == 0 {
			err := fmt.Errorf("organization not found in context")
			log.Printf("RoomController: %s", err)
			InternalServerError(w, err)
			return
		}

		rooms, err := c.roomService.FindForOrganization(organization.Id)
		if err != nil {
			log.Printf("RoomController: %s", err)
			InternalServerError(w, err)
			return
		}

		var roomDtos []resources.RoomDto
		for _, room := range rooms {
			roomDto := resources.RoomDto{}.DomainToDto(room)
			roomDtos = append(roomDtos, roomDto)
		}

		response := resources.RoomsDto{Rooms: roomDtos}
		Success(w, response)
	}
}

func (c RoomController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := r.Context().Value(OrgKey).(domain.Organization)
		ro := r.Context().Value(RoKey).(domain.Room)

		if ro.OrganizationId != organization.Id {
			err := fmt.Errorf("access denied")
			Forbidden(w, err)
			return
		}

		var roomDto resources.RoomDto
		Success(w, roomDto.DomainToDto(ro))
	}
}

func (c RoomController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := r.Context().Value(OrgKey).(domain.Organization)
		ro, err := requests.Bind(r, requests.RoomRequest{}, domain.Room{})
		if err != nil {
			log.Printf("RoomController: %s", err)
			BadRequest(w, err)
			return
		}

		room := r.Context().Value(RoKey).(domain.Room)
		if room.OrganizationId != organization.Id {
			err := fmt.Errorf("access denied")
			Forbidden(w, err)
			return
		}

		room.Name = ro.Name
		room.Description = ro.Description
		updatedRoom, err := c.roomService.Update(room)
		if err != nil {
			log.Printf("RoomController: %s", err)
			InternalServerError(w, err)
			return
		}

		var roomDto resources.RoomDto
		Success(w, roomDto.DomainToDto(updatedRoom))
	}
}

func (c RoomController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := r.Context().Value(OrgKey).(domain.Organization)
		ro := r.Context().Value(RoKey).(domain.Room)

		if ro.OrganizationId != organization.Id {
			err := fmt.Errorf("access denied")
			Forbidden(w, err)
			return
		}

		err := c.roomService.Delete(ro.Id)
		if err != nil {
			log.Printf("RoomController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
