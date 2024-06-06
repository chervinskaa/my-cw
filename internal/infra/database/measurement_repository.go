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
	FindAll() ([]domain.Measurement, error)
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

	deviceRepo := NewDeviceRepository(r.coll.Session())
	device, err := deviceRepo.Find(dm.DeviceId)
	if err != nil {
		log.Printf("MeasurementRepository: Error fetching device: %s", err)
		return domain.Measurement{}, err
	}

	if device.Category != domain.Sensor {
		err := errors.New("only sensors can have measurements")
		log.Printf("MeasurementRepository: Device ID %d is not a sensor", dm.DeviceId)
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
	log.Printf("MeasurementRepository: Measurement saved successfully: %+v", dm)
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
	var measurements []measurement
	err := r.coll.Find(db.Cond{"device_id": deviceId, "deleted_date": nil}).All(&measurements)
	if err != nil {
		if err == db.ErrNoMoreRows {
			log.Printf("MeasurementRepository: No measurements found for device ID %d", deviceId)
			return []domain.Measurement{}, nil
		}
		log.Printf("MeasurementRepository: Error finding measurements for device ID %d: %s", deviceId, err)
		return nil, err
	}

	log.Printf("MeasurementRepository: Found %d measurements for device ID %d", len(measurements), deviceId)

	return r.mapModelToDomainCollection(measurements), nil
}

func (r *measurementRepository) FindAll() ([]domain.Measurement, error) {
	var measurements []measurement
	err := r.coll.Find(db.Cond{"deleted_date": nil}).All(&measurements)
	if err != nil {
		return nil, err
	}
	res := r.mapModelToDomainCollection(measurements)
	return res, nil
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
