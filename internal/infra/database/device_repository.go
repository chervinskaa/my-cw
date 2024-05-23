package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const DeviceTableName = "devices"

type device struct {
	Id               uint64                `db:"id,omitempty"`
	OrganizationId   uint64                `db:"organization_id"`
	RoomId           *uint64               `db:"room_id"`
	GUID             string                `db:"guid_id"`
	InventoryNumber  string                `db:"inventory_number"`
	SerialNumber     string                `db:"serial_number"`
	Characteristics  string                `db:"characteristics"`
	Category         domain.DeviceCategory `db:"category"`
	Units            *string               `db:"units"`
	PowerConsumption *float64              `db:"power_consumption"`
	CreatedDate      time.Time             `db:"created_date"`
	UpdatedDate      time.Time             `db:"updated_date"`
	DeletedDate      *time.Time            `db:"deleted_date"`
}

type DeviceRepository interface {
	Save(o domain.Device) (domain.Device, error)
	FindAll() ([]domain.Device, error)
	Find(id uint64) (domain.Device, error)
	Update(o domain.Device) (domain.Device, error)
	InstallDevice(deviceId uint64, roomId uint64) error
	UninstallDevice(deviceId uint64) error
	Delete(id uint64) error
}

type deviceRepository struct {
	coll db.Collection
	sess db.Session
}

func NewDeviceRepository(dbSession db.Session) DeviceRepository {
	return &deviceRepository{
		coll: dbSession.Collection(DeviceTableName),
		sess: dbSession,
	}
}

func (r *deviceRepository) Save(o domain.Device) (domain.Device, error) {
	dev := r.mapDomainToModel(o)
	dev.CreatedDate, dev.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&dev)
	if err != nil {
		return domain.Device{}, err
	}
	o = r.mapModelToDomain(dev)
	return o, nil
}

func (r *deviceRepository) FindAll() ([]domain.Device, error) {
	var devs []device
	err := r.coll.Find(db.Cond{"deleted_date": nil}).All(&devs)
	if err != nil {
		return nil, err
	}
	res := r.mapModelToDomainCollection(devs)
	return res, nil
}

func (r *deviceRepository) Find(id uint64) (domain.Device, error) {
	var dev device
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&dev)
	if err != nil {
		return domain.Device{}, err
	}
	o := r.mapModelToDomain(dev)
	return o, nil
}

func (r *deviceRepository) Update(o domain.Device) (domain.Device, error) {
	dev := r.mapDomainToModel(o)
	dev.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": dev.Id, "deleted_date": nil}).Update(&dev)
	if err != nil {
		return domain.Device{}, err
	}
	o = r.mapModelToDomain(dev)
	return o, nil
}

func (r *deviceRepository) InstallDevice(deviceId uint64, roomId uint64) error {
	err := r.coll.Find(db.Cond{"id": deviceId, "deleted_date": nil}).Update(map[string]interface{}{
		"room_id": roomId,
	})
	return err
}

func (r *deviceRepository) UninstallDevice(deviceId uint64) error {
	err := r.coll.Find(db.Cond{"id": deviceId, "deleted_date": nil}).Update(map[string]interface{}{
		"room_id": nil,
	})
	return err
}

func (r *deviceRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r deviceRepository) mapDomainToModel(d domain.Device) device {
	return device{
		Id:               d.Id,
		OrganizationId:   d.OrganizationId,
		RoomId:           d.RoomId,
		GUID:             d.GUID,
		InventoryNumber:  d.InventoryNumber,
		SerialNumber:     d.SerialNumber,
		Characteristics:  d.Characteristics,
		Category:         d.Category,
		Units:            d.Units,
		PowerConsumption: d.PowerConsumption,
		CreatedDate:      d.CreatedDate,
		UpdatedDate:      d.UpdatedDate,
		DeletedDate:      d.DeletedDate,
	}
}

func (r deviceRepository) mapModelToDomain(d device) domain.Device {
	return domain.Device{
		Id:               d.Id,
		OrganizationId:   d.OrganizationId,
		RoomId:           d.RoomId,
		GUID:             d.GUID,
		InventoryNumber:  d.InventoryNumber,
		SerialNumber:     d.SerialNumber,
		Characteristics:  d.Characteristics,
		Category:         d.Category,
		Units:            d.Units,
		PowerConsumption: d.PowerConsumption,
		CreatedDate:      d.CreatedDate,
		UpdatedDate:      d.UpdatedDate,
		DeletedDate:      d.DeletedDate,
	}
}

func (r deviceRepository) mapModelToDomainCollection(devs []device) []domain.Device {
	var devices []domain.Device
	for _, o := range devs {
		dev := r.mapModelToDomain(o)
		devices = append(devices, dev)
	}
	return devices
}
