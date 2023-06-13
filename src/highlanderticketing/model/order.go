package model

type Order struct {
	OrderType OrderType
	Amount    int32
	User      User
}
type OrderType string

const (
	MATCHTICKET OrderType = "MATCHTICKET"
	BUSTICKET   OrderType = "TRAVELTICKET"
)
