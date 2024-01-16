package models

type RequestEnrich struct {
	Name string `json:"name"`
}

type ResponseEnrich struct {
	RequestEnrich
	Age     int    `json:"age"`
	Gender  string `json:"sex"`
	Country string `json:"country"`
}
