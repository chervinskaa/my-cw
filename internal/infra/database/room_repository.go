package database

import (
	"log"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const RoomsTableName = "rooms"

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
	Save(r domain.Room) (domain.Room, error)
	FindForOrganization(oId uint64) ([]domain.Room, error)
	Find(id uint64) (domain.Room, error)
	Update(r domain.Room) (domain.Room, error)
	Delete(id uint64) error
}

type roomRepository struct {
	coll db.Collection
	sess db.Session
}

func NewRoomRepository(dbSession db.Session) RoomRepository {
	if dbSession == nil {
		log.Fatal("NewRoomRepository: dbSession is nil")
	}

	log.Printf("NewRoomRepository: dbSession is initialized")

	return &roomRepository{
		coll: dbSession.Collection(RoomsTableName),
		sess: dbSession,
	}
}

func (r *roomRepository) Save(o domain.Room) (domain.Room, error) {
	room := r.mapDomainToModel(o)
	room.CreatedDate, room.UpdatedDate = time.Now(), time.Now()
	log.Printf("RoomRepository: Saving room %+v", room)
	err := r.coll.InsertReturning(&room)
	if err != nil {
		log.Printf("RoomRepository: Error saving room: %s", err)
		return domain.Room{}, err
	}
	o = r.mapModelToDomain(room)
	log.Printf("RoomRepository: Saved room %+v", o)
	return o, nil
}

func (r *roomRepository) FindForOrganization(oId uint64) ([]domain.Room, error) {
	log.Printf("RoomRepository: Finding rooms for organization %d", oId)
	var rooms []room
	err := r.coll.Find(db.Cond{"organization_id": oId, "deleted_date": nil}).All(&rooms)
	if err != nil {
		log.Printf("RoomRepository: Error finding rooms for organization %d: %s", oId, err)
		return nil, err
	}
	res := r.mapModelToDomainCollection(rooms)
	log.Printf("RoomRepository: Found rooms for organization %d: %+v", oId, res)
	return res, nil
}

func (r *roomRepository) Find(id uint64) (domain.Room, error) {
	log.Printf("RoomRepository: Finding room with id %d", id)
	var room room
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&room)
	if err != nil {
		log.Printf("RoomRepository: Error finding room with id %d: %s", id, err)
		return domain.Room{}, err
	}
	o := r.mapModelToDomain(room)
	log.Printf("RoomRepository: Found room with id %d: %+v", id, o)
	return o, nil
}

func (r *roomRepository) Update(o domain.Room) (domain.Room, error) {
	room := r.mapDomainToModel(o)
	room.UpdatedDate = time.Now()
	log.Printf("RoomRepository: Updating room %+v", room)
	err := r.coll.Find(db.Cond{"id": room.Id, "deleted_date": nil}).Update(&room)
	if err != nil {
		log.Printf("RoomRepository: Error updating room: %s", err)
		return domain.Room{}, err
	}
	o = r.mapModelToDomain(room)
	log.Printf("RoomRepository: Updated room %+v", o)
	return o, nil
}

func (r *roomRepository) Delete(id uint64) error {
	log.Printf("RoomRepository: Deleting room with id %d", id)
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
	if err != nil {
		log.Printf("RoomRepository: Error deleting room with id %d: %s", id, err)
		return err
	}
	log.Printf("RoomRepository: Deleted room with id %d", id)
	return nil
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
