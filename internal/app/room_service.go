package app

import (
	"errors"
	"fmt"
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type RoomService interface {
	Save(r domain.Room) (domain.Room, error)
	FindByOrgId(orgId uint64) ([]domain.Room, error)
	Find(id uint64) (interface{}, error)
	Update(r domain.Room) (domain.Room, error)
	Delete(id uint64) error
}

type roomService struct {
	roomRepo   database.RoomRepository
	orgRepo    database.OrganizationRepository
	deviceRepo database.DeviceRepository
}

func NewRoomService(rr database.RoomRepository, or database.OrganizationRepository, dr database.DeviceRepository) (RoomService, error) {
	if rr == nil {
		return nil, errors.New("room repository is nil")
	}
	if or == nil {
		return nil, errors.New("organization repository is nil")
	}

	log.Printf("NewRoomService: room repository and organization repository are initialized")

	return &roomService{
		roomRepo:   rr,
		orgRepo:    or,
		deviceRepo: dr,
	}, nil
}

func (s *roomService) Save(r domain.Room) (domain.Room, error) {
	createdRoom, err := s.roomRepo.Save(r)
	if err != nil {
		log.Printf("RoomService: Error saving room: %s", err)
		return domain.Room{}, err
	}

	log.Printf("RoomService: Room saved successfully: %+v", createdRoom)
	return createdRoom, nil
}

func (s *roomService) FindByOrgId(orgId uint64) ([]domain.Room, error) {
	rooms, err := s.roomRepo.FindByOrgId(orgId)
	if err != nil {
		log.Printf("RoomService: Error finding rooms for organization ID %d: %s", orgId, err)
		return nil, err
	}

	if len(rooms) == 0 {
		log.Printf("RoomService: No rooms found for organization ID %d", orgId)
	}

	log.Printf("RoomService: Found %d rooms for organization ID %d", len(rooms), orgId)
	return rooms, nil
}

func (s *roomService) Find(id uint64) (interface{}, error) {
	log.Printf("RoomService: Finding room with ID %d", id)

	room, err := s.roomRepo.Find(id)
	if err != nil {
		log.Printf("RoomService: Error finding room: %s", err)
		return domain.Room{}, err
	}

	if room.Id == 0 {
		log.Printf("RoomService: No room found with ID %d", id)
		return domain.Room{}, fmt.Errorf("room not found")
	}

	log.Printf("RoomService: Found room successfully: %+v", room)
	return room, nil
}

func (s *roomService) Update(r domain.Room) (domain.Room, error) {
	log.Printf("RoomService: Updating room %+v", r)

	room, err := s.roomRepo.Update(r)
	if err != nil {
		log.Printf("RoomService: Error updating room: %s", err)
		return domain.Room{}, err
	}

	log.Printf("RoomService: Room updated successfully: %+v", room)
	return room, nil
}

func (s *roomService) Delete(id uint64) error {
	log.Printf("RoomService: Deleting room with ID %d", id)

	err := s.roomRepo.Delete(id)
	if err != nil {
		log.Printf("RoomService: Error deleting room: %s", err)
		return err
	}

	log.Printf("RoomService: Room deleted successfully")
	return nil
}
