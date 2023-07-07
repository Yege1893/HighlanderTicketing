package model

type EmialContent struct {
	OrderID     string `json:"orderid"`
	Name        string `json:"name"`
	AwayMatch   bool   `json:"awaymatch"`
	Location    string `json:"location"`
	Date        string `json:"date"`
	Emailadress string `json:"emailadress"`
}
