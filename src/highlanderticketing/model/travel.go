package model

import (
	"google.golang.org/genproto/googleapis/type/date"
)

type Travel struct {
	ID                  uint
	TravelType          TravelType
	InitialSeatAmount   int32
	AvailableSeatAmount int32
	StartLocation       string
	EndLocation         string
	StartDate           date.Date
	Orders              []Order
}

type TravelType string

const (
	CAR   TravelType = "CAR"
	BUS   TravelType = "BUS"
	PLANE TravelType = "PLANE"
)
