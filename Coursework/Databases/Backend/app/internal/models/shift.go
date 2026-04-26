package models

import "time"

type ShiftEmployee struct {
	ID       string
	FullName string
	Login    string
}

type Shift struct {
	ID       string
	OpenedAt time.Time
	ClosedAt *time.Time
	IsActive bool
	Duty     *ShiftEmployee
}
