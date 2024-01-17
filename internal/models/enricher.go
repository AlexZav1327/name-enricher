package models

type AgeEnriched struct {
	Age int `json:"age"`
}

type GenderEnricher struct {
	Gender string `json:"gender"`
}

type CountryEnriched struct {
	CountryID   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type CountryEnrichedList struct {
	Country []CountryEnriched `json:"country"`
}

type ListingQueryParams struct {
	TextFilter   string
	ItemsPerPage int
	Offset       int
	Sorting      string
	Descending   bool
}
