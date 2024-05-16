package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const RoomTableName = "rooms"

type room struct {
	Id             uint64     `db:"id,omitempty"`
	OrganizationId uint64     `db:"organization_id"`
	Name           string     `db:"name"`
	Description    string     `db:"description"`
	CreatedDate    time.Time  `db:"created_date"`
	UpdatedDate    time.Time  `db:"updated_date"`
	DeletedDate    *time.Time `db:"deleted_date"`
}

type RoomRepository interface {
	Save(o domain.Room) (domain.Room, error)
	FindForOrganization(oId uint64) ([]domain.Room, error)
	FindById(id uint64) (domain.Room, error)
	Update(o domain.Room) (domain.Room, error)
	Delete(id uint64) error
}

type roomRepository struct {
	coll db.Collection
	sess db.Session
}

func NewRoomRepository(dbSession db.Session) RoomRepository {
	return &roomRepository{
		coll: dbSession.Collection(RoomTableName),
		sess: dbSession,
	}
}

func (r *roomRepository) Save(o domain.Room) (domain.Room, error) {
	room := r.mapDomainToModel(o)
	room.CreatedDate, room.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&room)
	if err != nil {
		return domain.Room{}, err
	}
	o = r.mapModelToDomain(room)
	return o, nil
}

func (r *roomRepository) FindForOrganization(oId uint64) ([]domain.Room, error) {
	var rooms []room
	err := r.coll.Find(db.Cond{"organization_id": oId, "deleted_date": nil}).All(&rooms)
	if err != nil {
		return nil, err
	}
	res := r.mapModelToDomainCollection(rooms)
	return res, nil
}

func (r *roomRepository) FindById(id uint64) (domain.Room, error) {
	var room room
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&room)
	if err != nil {
		return domain.Room{}, err
	}
	o := r.mapModelToDomain(room)
	return o, nil
}

func (r *roomRepository) Update(o domain.Room) (domain.Room, error) {
	room := r.mapDomainToModel(o)
	room.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": room.Id, "deleted_date": nil}).Update(&room)
	if err != nil {
		return domain.Room{}, err
	}
	o = r.mapModelToDomain(room)
	return o, nil
}

func (r *roomRepository) Delete(id uint64) error {
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
	return err
}

func (r *roomRepository) mapDomainToModel(d domain.Room) room {
	return room{
		Id:             d.Id,
		OrganizationId: d.OrganizationId,
		Name:           d.Name,
		Description:    d.Description,
		CreatedDate:    d.CreatedDate,
		UpdatedDate:    d.UpdatedDate,
		DeletedDate:    d.DeletedDate,
	}
}

func (r *roomRepository) mapModelToDomain(d room) domain.Room {
	return domain.Room{
		Id:             d.Id,
		OrganizationId: d.OrganizationId,
		Name:           d.Name,
		Description:    d.Description,
		CreatedDate:    d.CreatedDate,
		UpdatedDate:    d.UpdatedDate,
		DeletedDate:    d.DeletedDate,
	}
}

func (r *roomRepository) mapModelToDomainCollection(roomModels []room) []domain.Room {
	var rooms []domain.Room
	for _, roomModel := range roomModels {
		room := r.mapModelToDomain(roomModel)
		rooms = append(rooms, room)
	}
	return rooms
}
