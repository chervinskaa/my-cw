package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type OrganizationService interface {
	Save(o domain.Organization) (domain.Organization, error)
	FindForUser(uId uint64) ([]domain.Organization, error)
	Find(id uint64) (interface{}, error)
	Update(o domain.Organization) (domain.Organization, error)
	Delete(id uint64) error
}

type organizationService struct {
	organizationRepo database.OrganizationRepository
	roomRepo         database.RoomRepository
}

func NewOrganizationService(or database.OrganizationRepository, rr database.RoomRepository) OrganizationService {
	return organizationService{
		organizationRepo: or,
		roomRepo:         rr,
	}
}

func (s organizationService) Save(o domain.Organization) (domain.Organization, error) {
	o, err := s.organizationRepo.Save(o)
	if err != nil {
		log.Printf("OrganizationService: %s", err)
		return domain.Organization{}, err
	}

	return o, nil
}

func (s organizationService) FindForUser(uId uint64) ([]domain.Organization, error) {
	orgs, err := s.organizationRepo.FindForUser(uId)
	if err != nil {
		log.Printf("OrganizationService: %s", err)
		return nil, err
	}

	return orgs, nil
}

func (s organizationService) Find(id uint64) (interface{}, error) {
	org, err := s.organizationRepo.FindById(id)
	if err != nil {
		log.Printf("OrganizationService: %s", err)
		return nil, err
	}

	org.Rooms, err = s.roomRepo.FindByOrgId(org.Id)
	if err != nil {
		log.Printf("OrganizationService: %s", err)
		return nil, err
	}

	return org, nil
}

func (s organizationService) Update(o domain.Organization) (domain.Organization, error) {
	org, err := s.organizationRepo.Update(o)
	if err != nil {
		log.Printf("OrganizationService: %s", err)
		return domain.Organization{}, err
	}

	return org, nil
}

func (s organizationService) Delete(id uint64) error {
	err := s.organizationRepo.Delete(id)
	if err != nil {
		log.Printf("OrganizationService: %s", err)
		return err
	}

	return nil
}
