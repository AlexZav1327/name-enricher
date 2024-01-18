package models

type AgeEnriched struct {
	Age int `json:"age"`
}

type GenderEnriched struct {
	Gender string `json:"gender"`
}

type CountryEnriched struct {
	CountryID   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type CountryEnrichedList struct {
	Country []CountryEnriched `json:"country"`
}
