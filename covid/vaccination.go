package covid

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type VaccinationData struct {
	StateVaccinationData
	LastUpdate string                          `json:"lastUpdate"`
	States     map[string]StateVaccinationData `json:"states"`
}

type StateVaccinationData struct {
	Total      uint    `json:"total"`
	Vaccinated uint    `json:"vaccinated"`
	Quote      float64 `json:"quote"`
}

func GetVaccinationData() (*VaccinationData, error) {
	response, err := resty.New().
		SetHeader("accept", "application/json").
		SetHostURL("https://rki-vaccination-data.vercel.app").
		R().
		Get("/api")
	if err != nil {
		return nil, err
	}

	result := new(VaccinationData)
	if err := json.Unmarshal(response.Body(), result); err != nil {
		return nil, err
	}

	return result, err
}
