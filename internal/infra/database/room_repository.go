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
	FindByOrgId(oId uint64) ([]domain.Room, error)
	Find(id uint64) (domain.Room, error)
	Update(r domain.Room) (domain.Room, error)
	Delete(id uint64) error
}

type roomRepository struct {
	coll db.Collection
}

func NewRoomRepository(sess db.Session) RoomRepository {
	return &roomRepository{
		coll: sess.Collection(RoomsTableName),
	}
}

func (r *roomRepository) Save(dr domain.Room) (domain.Room, error) {
	room := r.mapDomainToModel(dr)
	now := time.Now()
	room.CreatedDate, room.UpdatedDate = now, now
	log.Printf("RoomRepository: Saving room %+v", room)
	err := r.coll.InsertReturning(&room)
	if err != nil {
		log.Printf("RoomRepository: Error saving room: %s", err)
		return domain.Room{}, err
	}
	dr = r.mapModelToDomain(room)
	log.Printf("RoomRepository: Saved room %+v", dr)
	return dr, nil
}

func (r *roomRepository) FindByOrgId(orgId uint64) ([]domain.Room, error) {
	var rooms []room
	err := r.coll.Find(db.Cond{"organization_id": orgId}).All(&rooms)
	if err != nil {
		if err == db.ErrNoMoreRows {
			log.Printf("RoomRepository: No rooms found for organization ID %d", orgId)
			return []domain.Room{}, nil
		}
		log.Printf("RoomRepository: Error finding rooms for organization ID %d: %s", orgId, err)
		return nil, err
	}

	log.Printf("RoomRepository: Found %d rooms for organization ID %d", len(rooms), orgId)

	return r.mapModelToDomainCollection(rooms), nil
}

func (r *roomRepository) Find(id uint64) (domain.Room, error) {
	var roomModel room
	err := r.coll.Find(db.Cond{"id": id}).One(&roomModel)
	if err != nil {
		if err == db.ErrNoMoreRows {
			log.Printf("RoomRepository: No room found with ID %d", id)
			return domain.Room{}, nil
		}
		log.Printf("RoomRepository: Error finding room with ID %d: %s", id, err)
		return domain.Room{}, err
	}
	domainRoom := r.mapModelToDomain(roomModel)
	log.Printf("RoomRepository: Found room with ID %d: %+v", id, domainRoom)
	return domainRoom, nil
}

func (r *roomRepository) Update(dr domain.Room) (domain.Room, error) {
	room := r.mapDomainToModel(dr)
	room.UpdatedDate = time.Now()
	log.Printf("RoomRepository: Updating room %+v", room)
	err := r.coll.Find(db.Cond{"id": room.Id, "deleted_date": nil}).Update(&room)
	if err != nil {
		log.Printf("RoomRepository: Error updating room: %s", err)
		return domain.Room{}, err
	}
	dr = r.mapModelToDomain(room)
	log.Printf("RoomRepository: Updated room %+v", dr)
	return dr, nil
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

func (r *roomRepository) mapModelToDomain(m room) domain.Room {
	return domain.Room{
		Id:             m.Id,
		OrganizationId: m.OrganizationId,
		Name:           m.Name,
		Description:    m.Description,
		CreatedDate:    m.CreatedDate,
		UpdatedDate:    m.UpdatedDate,
		DeletedDate:    m.DeletedDate,
	}
}

func (r *roomRepository) mapModelToDomainCollection(roomModels []room) []domain.Room {
	rooms := make([]domain.Room, len(roomModels))
	for i, roomModel := range roomModels {
		rooms[i] = r.mapModelToDomain(roomModel)
	}
	return rooms
}
