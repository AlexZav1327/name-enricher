package models

type RequestEnrich struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type ResponseEnrich struct {
	RequestEnrich
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Country string `json:"country"`
}

type ListingQueryParams struct {
	TextFilter   string
	ItemsPerPage int
	Offset       int
	Sorting      string
	Descending   bool
}
