package database

import (
	"log"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const MeasurementsTableName = "measurements"

type measurement struct {
	Id          uint64     `db:"id,omitempty"`
	DeviceId    uint64     `db:"device_id"`
	RoomId      *uint64    `db:"room_id"`
	Value       float64    `db:"value"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date"`
}

type MeasurementRepository interface {
	Save(dm domain.Measurement) (domain.Measurement, error)
	FindByDeviceAndDate(deviceId uint64, startDate, endDate time.Time) (interface{}, error)
	Find(id uint64) (domain.Measurement, error)
	Update(dm domain.Measurement) (domain.Measurement, error)
	Delete(id uint64) error
}

type measurementRepository struct {
	coll db.Collection
	sess db.Session
}

func NewMeasurementRepository(dbSession db.Session) MeasurementRepository {
	return &measurementRepository{
		coll: dbSession.Collection(MeasurementsTableName),
		sess: dbSession,
	}
}

func (r *measurementRepository) Save(dm domain.Measurement) (domain.Measurement, error) {
	measurement := r.mapDomainToModel(dm)
	now := time.Now()
	measurement.CreatedDate, measurement.UpdatedDate = now, now
	log.Printf("MeasurementRepository: Saving measurement %+v", measurement)
	err := r.coll.InsertReturning(&measurement)
	if err != nil {
		log.Printf("MeasurementRepository: Error saving measurement: %s", err)
		return domain.Measurement{}, err
	}
	dm = r.mapModelToDomain(measurement)
	log.Printf("MeasurementRepository: Saved measurement %+v", dm)
	return dm, nil
}

func (r *measurementRepository) FindByDeviceAndDate(deviceId uint64, startDate, endDate time.Time) (interface{}, error) {
	var measurements []measurement
	err := r.coll.Find(db.Cond{"device_id": deviceId, "created_date >=": startDate, "created_date <=": endDate, "deleted_date": nil}).All(&measurements)
	if err != nil {
		return nil, err
	}
	res := r.mapModelToDomainCollection(measurements)
	return res, nil
}

func (r *measurementRepository) Find(id uint64) (domain.Measurement, error) {
	var measurement measurement
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&measurement)
	if err != nil {
		return domain.Measurement{}, err
	}
	return r.mapModelToDomain(measurement), nil
}

func (r *measurementRepository) Update(dm domain.Measurement) (domain.Measurement, error) {
	measurement := r.mapDomainToModel(dm)
	measurement.UpdatedDate = time.Now()
	log.Printf("MeasurementRepository: Updating measurement %+v", measurement)
	err := r.coll.Find(db.Cond{"id": measurement.Id, "deleted_date": nil}).Update(&measurement)
	if err != nil {
		log.Printf("MeasurementRepository: Error updating measurement: %s", err)
		return domain.Measurement{}, err
	}
	dm = r.mapModelToDomain(measurement)
	log.Printf("MeasurementRepository: Updated measurement %+v", dm)
	return dm, nil
}

func (r *measurementRepository) Delete(id uint64) error {
	log.Printf("MeasurementRepository: Deleting measurement with id %d", id)
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
	if err != nil {
		log.Printf("MeasurementRepository: Error deleting measurement with id %d: %s", id, err)
		return err
	}
	log.Printf("MeasurementRepository: Deleted measurement with id %d", id)
	return nil
}

func (m measurementRepository) mapDomainToModel(d domain.Measurement) measurement {
	return measurement{
		Id:          d.Id,
		DeviceId:    d.DeviceId,
		RoomId:      d.RoomId,
		Value:       d.Value,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (mr measurementRepository) mapModelToDomain(m measurement) domain.Measurement {
	return domain.Measurement{
		Id:          m.Id,
		DeviceId:    m.DeviceId,
		RoomId:      m.RoomId,
		Value:       m.Value,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (mr measurementRepository) mapModelToDomainCollection(measurs []measurement) []domain.Measurement {
	var measurements []domain.Measurement
	for _, m := range measurs {
		measurs := mr.mapModelToDomain(m)
		measurements = append(measurements, measurs)
	}
	return measurements
}
