package model

type EmialContent struct {
	Name        string `json:"name"`
	AwayMatch   bool   `json:"awaymatch"`
	Location    string `json:"location"`
	Date        string `json:"date"`
	Emailadress string `json:"emailadress"`
}
