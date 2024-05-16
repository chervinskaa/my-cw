package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type RoomService interface {
	Save(r domain.Room) (domain.Room, error)
	FindForOrganization(oId uint64) ([]domain.Room, error)
	Find(id uint64) (domain.Room, error)
	Update(r domain.Room) (domain.Room, error)
	Delete(id uint64) error
}

type roomService struct {
	roomRepo database.RoomRepository
}

func NewRoomService(rr database.RoomRepository) RoomService {
	return &roomService{
		roomRepo: rr,
	}
}

func (s *roomService) Save(r domain.Room) (domain.Room, error) {
	room, err := s.roomRepo.Save(r)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	return room, nil
}

func (s *roomService) FindForOrganization(oId uint64) ([]domain.Room, error) {
	rooms, err := s.roomRepo.FindForOrganization(oId)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return nil, err
	}

	return rooms, nil
}

func (s *roomService) Find(id uint64) (domain.Room, error) {
	room, err := s.roomRepo.FindById(id)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	return room, nil
}

func (s *roomService) Update(r domain.Room) (domain.Room, error) {
	room, err := s.roomRepo.Update(r)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	return room, nil
}

func (s *roomService) Delete(id uint64) error {
	err := s.roomRepo.Delete(id)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return err
	}

	return nil
}
