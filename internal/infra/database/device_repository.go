package database

import (
	"errors"
	"log"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const DevicesTableName = "devices"

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
	Save(d domain.Device) (domain.Device, error)
	FindByRoomId(roomId uint64) ([]domain.Device, error)
	Find(id uint64) (domain.Device, error)
	Update(d domain.Device) (domain.Device, error)
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
		coll: dbSession.Collection(DevicesTableName),
		sess: dbSession,
	}
}

func (r *deviceRepository) Save(dd domain.Device) (domain.Device, error) {
	if dd.Category == domain.Actuator && dd.PowerConsumption == nil {
		return domain.Device{}, errors.New("power consumption is required for ACTUATOR")
	}
	if dd.Category == domain.Sensor && dd.Units == nil {
		return domain.Device{}, errors.New("units is required for SENSOR")
	}
	device := r.mapDomainToModel(dd)
	now := time.Now()
	device.CreatedDate, device.UpdatedDate = now, now
	log.Printf("DeviceRepository: Saving device %+v", device)
	err := r.coll.InsertReturning(&device)
	if err != nil {
		log.Printf("DeviceRepository: Error saving device: %s", err)
		return domain.Device{}, err
	}
	dd = r.mapModelToDomain(device)
	log.Printf("DeviceRepository: Saved device %+v", dd)
	return dd, nil
}

func (r *deviceRepository) FindByRoomId(roomId uint64) ([]domain.Device, error) {
	var devices []device
	err := r.coll.Find(db.Cond{"room_id": roomId}).All(&devices)
	if err != nil {
		if err == db.ErrNoMoreRows {
			log.Printf("DeviceRepository: No devices found for room ID %d", roomId)
			return []domain.Device{}, nil
		}
		log.Printf("DeviceRepository: Error finding devices for room ID %d: %s", roomId, err)
		return nil, err
	}

	log.Printf("DeviceRepository: Found %d devices for room ID %d", len(devices), roomId)

	return r.mapModelToDomainCollection(devices), nil
}

func (r *deviceRepository) Find(id uint64) (domain.Device, error) {
	var dev device
	err := r.coll.Find(db.Cond{"id": id}).One(&dev)
	if err != nil {
		if err == db.ErrNoMoreRows {
			log.Printf("DeviceRepository: No device found with ID %d", id)
			return domain.Device{}, nil
		}
		log.Printf("DeviceRepository: Error finding device with ID %d: %s", id, err)
		return domain.Device{}, err
	}
	dd := r.mapModelToDomain(dev)
	log.Printf("DeviceRepository: Found device with ID %d: %+v", id, dd)
	return dd, nil
}

func (r *deviceRepository) Update(dd domain.Device) (domain.Device, error) {
	device := r.mapDomainToModel(dd)
	device.UpdatedDate = time.Now()
	log.Printf("DeviceRepository: Updating device %+v", device)
	err := r.coll.Find(db.Cond{"id": device.Id, "deleted_date": nil}).Update(&device)
	if err != nil {
		log.Printf("DeviceRepository: Error updating device: %s", err)
		return domain.Device{}, err
	}
	dd = r.mapModelToDomain(device)
	log.Printf("DeviceRepository: Updated device %+v", dd)
	return dd, nil
}

func (r *deviceRepository) InstallDevice(deviceId uint64, roomId uint64) error {
	err := r.coll.Find(db.Cond{"id": deviceId, "deleted_date": nil}).Update(map[string]interface{}{
		"room_id": roomId,
	})
	return err
}

func (r *deviceRepository) UninstallDevice(deviceId uint64) error {
	query := "UPDATE devices SET room_id = NULL, updated_date = ? WHERE id = ?"
	_, err := r.sess.SQL().Exec(query, time.Now(), deviceId)
	return err
}

func (r *deviceRepository) Delete(id uint64) error {
	log.Printf("DeviceRepository: Deleting device with id %d", id)
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
	if err != nil {
		log.Printf("DeviceRepository: Error deleting device with id %d: %s", id, err)
		return err
	}
	log.Printf("DeviceRepository: Deleted device with id %d", id)
	return nil
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
	for _, d := range devs {
		dev := r.mapModelToDomain(d)
		devices = append(devices, dev)
	}
	return devices
}
