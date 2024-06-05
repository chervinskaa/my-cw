package database

import (
	"errors"
	"log"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

type measurement struct {
	Id          uint64     `db:"id,omitempty"`
	DeviceId    uint64     `db:"device_id"`
	RoomId      *uint64    `db:"room_id"`
	Value       float64    `db:"value"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date"`
}

const MeasurementsTableName = "measurements"

type MeasurementRepository interface {
	Save(dm domain.Measurement) (domain.Measurement, error)
	FindByDeviceAndDate(deviceId uint64, startDate, endDate time.Time) ([]domain.Measurement, error)
	Find(id uint64) (domain.Measurement, error)
	FindByDeviceId(deviceId uint64) ([]domain.Measurement, error)
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
	var device domain.Device
	err := r.sess.Collection("devices").Find(db.Cond{"id": dm.DeviceId}).One(&device)
	if err != nil {
		log.Printf("MeasurementRepository: Error fetching device: %s", err)
		return domain.Measurement{}, err
	}

	if device.Category != domain.Sensor {
		err := errors.New("only sensors can have measurements")
		log.Printf("MeasurementRepository: %s", err)
		return domain.Measurement{}, err
	}

	measurement := r.mapDomainToModel(dm)
	now := time.Now()
	measurement.CreatedDate, measurement.UpdatedDate = now, now
	log.Printf("MeasurementRepository: Saving measurement %+v", measurement)
	err = r.coll.InsertReturning(&measurement)
	if err != nil {
		log.Printf("MeasurementRepository: Error saving measurement: %s", err)
		return domain.Measurement{}, err
	}
	dm = r.mapModelToDomain(measurement)
	log.Printf("MeasurementRepository: Saved measurement %+v", dm)
	return dm, nil
}

func (r *measurementRepository) FindByDeviceAndDate(deviceId uint64, startDate, endDate time.Time) ([]domain.Measurement, error) {
	var measurements []measurement
	err := r.coll.Find(db.Cond{"device_id": deviceId, "created_date >=": startDate, "created_date <=": endDate, "deleted_date": nil}).All(&measurements)
	if err != nil {
		return nil, err
	}
	return r.mapModelToDomainCollection(measurements), nil
}

func (r *measurementRepository) Find(id uint64) (domain.Measurement, error) {
	var measurement measurement
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&measurement)
	if err != nil {
		return domain.Measurement{}, err
	}
	return r.mapModelToDomain(measurement), nil
}

func (r *measurementRepository) FindByDeviceId(deviceId uint64) ([]domain.Measurement, error) {
	var measurements []domain.Measurement
	err := r.coll.Find(db.Cond{"device_id": deviceId}).All(&measurements)
	if err != nil {
		log.Printf("MeasurementRepository: Error finding measurements for device ID %d: %s", deviceId, err)
		return nil, err
	}
	return measurements, nil
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

func (r *measurementRepository) mapDomainToModel(d domain.Measurement) measurement {
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

func (r *measurementRepository) mapModelToDomain(m measurement) domain.Measurement {
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

func (r *measurementRepository) mapModelToDomainCollection(measurements []measurement) []domain.Measurement {
	var domainMeasurements []domain.Measurement
	for _, m := range measurements {
		domainMeasurement := r.mapModelToDomain(m)
		domainMeasurements = append(domainMeasurements, domainMeasurement)
	}
	return domainMeasurements
}
