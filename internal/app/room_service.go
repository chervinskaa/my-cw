package app

import (
	"errors"
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type RoomService interface {
	Save(r domain.Room, uId uint64) (domain.Room, error)
	FindForOrganization(oId uint64) ([]domain.Room, error)
	Find(id uint64) (interface{}, error)
	Update(r domain.Room) (domain.Room, error)
	Delete(id uint64) error
}

type roomService struct {
	roomRepo database.RoomRepository
	orgRepo  database.OrganizationRepository
}

func NewRoomService(rr database.RoomRepository, or database.OrganizationRepository) (RoomService, error) {
	if rr == nil {
		return nil, errors.New("room repository is nil")
	}
	if or == nil {
		return nil, errors.New("organization repository is nil")
	}

	log.Printf("NewRoomService: room repository and organization repository are initialized")

	return &roomService{
		roomRepo: rr,
		orgRepo:  or,
	}, nil
}

func (s *roomService) Save(r domain.Room, uId uint64) (domain.Room, error) {
	log.Printf("RoomService: Saving room %+v for user with ID %d", r, uId)

	org, err := s.orgRepo.FindById(r.OrganizationId)
	if err != nil {
		log.Printf("RoomService: Error finding organization: %s", err)
		return domain.Room{}, err
	}

	if org.UserId != uId {
		err = errors.New("access denied")
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	room, err := s.roomRepo.Save(r)
	if err != nil {
		log.Printf("RoomService: Error saving room: %s", err)
		return domain.Room{}, err
	}

	log.Printf("RoomService: Room saved successfully: %+v", room)
	return room, nil
}

func (s *roomService) FindForOrganization(oId uint64) ([]domain.Room, error) {
	log.Printf("RoomService: Finding rooms for organization with ID %d", oId)

	rooms, err := s.roomRepo.FindForOrganization(oId)
	if err != nil {
		log.Printf("RoomService: Error finding rooms: %s", err)
		return nil, err
	}

	log.Printf("RoomService: Found rooms successfully: %+v", rooms)
	return rooms, nil
}

func (s *roomService) Find(id uint64) (interface{}, error) {
	log.Printf("RoomService: Finding room with ID %d", id)

	room, err := s.roomRepo.Find(id)
	if err != nil {
		log.Printf("RoomService: Error finding room: %s", err)
		return nil, err
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
