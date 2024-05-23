package domain

import "time"

type Measurement struct {
	Id          uint64
	DeviceId    uint64
	RoomId      uint64
	Value       float64
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}
