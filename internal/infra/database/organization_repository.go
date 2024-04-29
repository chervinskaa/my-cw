package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const OrganizationsTableName = "organizations"

type organization struct {
	Id          uint64     `db:"id"`
	UserId      uint64     `db:"user_id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	City        string     `db:"city"`
	Address     string     `db:"address"`
	Lat         float64    `db:"lat"`
	Lon         float64    `db:"lon"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date"`
}

type OrganizationRepository interface {
	Save(o domain.Organization) (domain.Organization, error)
}

type organizationRepository struct {
	coll db.Collection
	sess db.Session
}

func NewOrganizationRepository(dbSession db.Session) OrganizationRepository {
	return organizationRepository{
		coll: dbSession.Collection(OrganizationsTableName),
		sess: dbSession,
	}
}

func (r organizationRepository) Save(o domain.Organization) (domain.Organization, error) {
	org := r.mapDomainToModel(o)
	err := r.coll.InsertReturning(&org)
	if err != nil {
		return domain.Organization{}, err
	}
	o = r.mapModelToDomain(org)
	return o, nil
}

func (r organizationRepository) mapDomainToModel(d domain.Organization) organization {
	return organization{
		Id:          d.Id,
		UserId:      d.UserId,
		Name:        d.Name,
		Description: d.Description,
		City:        d.City,
		Address:     d.Address,
		Lat:         d.Lat,
		Lon:         d.Lon,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r organizationRepository) mapModelToDomain(d organization) domain.Organization {
	return domain.Organization{
		Id:          d.Id,
		UserId:      d.UserId,
		Name:        d.Name,
		Description: d.Description,
		City:        d.City,
		Address:     d.Address,
		Lat:         d.Lat,
		Lon:         d.Lon,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}