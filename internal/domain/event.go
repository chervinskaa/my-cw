package domain

import "time"

type EventAction string

const (
	TurnOn  EventAction = "ON"
	TurnOff EventAction = "OFF"
)

type Event struct {
	Id          uint64      `db:"id,omitempty"`
	DeviceId    uint64      `db:"device_id"`
	RoomId      *uint64     `db:"room_id"`
	Action      EventAction `db:"action"`
	CreatedDate time.Time   `db:"created_date"`
	UpdatedDate time.Time   `db:"updated_date"`
	DeletedDate *time.Time  `db:"deleted_date"`
}
