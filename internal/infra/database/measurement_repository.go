package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const MeasurementsTableName = "measurements"

type measurement struct {
	Id          uint64     `db:"id,omitempty"`
	DeviceId    uint64     `db:"device_id"`
	RoomId      uint64     `db:"room_id"`
	Value       float64    `db:"value"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date"`
}

type MeasurementRepository interface {
	Save(m domain.Measurement) (domain.Measurement, error)
	FindByDeviceAndDate(deviceId uint64, startDate, endDate time.Time) ([]domain.Measurement, error)
	Find(id uint64) (domain.Measurement, error)
	Update(m domain.Measurement) (domain.Measurement, error)
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

func (r *measurementRepository) Save(m domain.Measurement) (domain.Measurement, error) {
	measurement := r.mapDomainToModel(m)
	measurement.CreatedDate, measurement.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&measurement)
	if err != nil {
		return domain.Measurement{}, err
	}
	m = r.mapModelToDomain(measurement)
	return m, nil
}

func (r *measurementRepository) FindByDeviceAndDate(deviceId uint64, startDate, endDate time.Time) ([]domain.Measurement, error) {
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
	m := r.mapModelToDomain(measurement)
	return m, nil
}

func (r *measurementRepository) Update(m domain.Measurement) (domain.Measurement, error) {
	measurement := r.mapDomainToModel(m)
	measurement.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": measurement.Id, "deleted_date": nil}).Update(&measurement)
	if err != nil {
		return domain.Measurement{}, err
	}
	m = r.mapModelToDomain(measurement)
	return m, nil
}

func (r *measurementRepository) Delete(id uint64) error {
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
	return err
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

func (r *measurementRepository) mapModelToDomain(d measurement) domain.Measurement {
	return domain.Measurement{
		Id:          d.Id,
		DeviceId:    d.DeviceId,
		RoomId:      d.RoomId,
		Value:       d.Value,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r *measurementRepository) mapModelToDomainCollection(measurementModels []measurement) []domain.Measurement {
	var measurements []domain.Measurement
	for _, measurementModel := range measurementModels {
		measurement := r.mapModelToDomain(measurementModel)
		measurements = append(measurements, measurement)
	}
	return measurements
}
