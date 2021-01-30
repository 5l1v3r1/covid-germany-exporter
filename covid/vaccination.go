package covid

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

const (
	rkiVaccinationDataAPIURL = "https://rki-vaccination-data.vercel.app"
)

// VaccinationData is a struct that is JSON-compatible with the RKI vaccination data API.
// It contains the current vaccination data grouped by province.
type VaccinationData struct {
	// StateVaccinationData is embedded to represent the vaccination data for Germany as a whole.
	StateVaccinationData
	// LastUpdate is the ISO-instant formatted time of the last data update.
	LastUpdate string `json:"lastUpdate"`
	// States contains a map with an entry for each province and the corresponding data for that province.
	States map[string]StateVaccinationData `json:"states"`
}

// StateVaccinationData is a struct that is JSON-compatible with the RKI vaccination data API.
// It contains the current vaccination data for a single province.
type StateVaccinationData struct {
	// VaccinationCountData is embedded to represent the vaccination counts for the province.
	VaccinationCountData
	// Total is the number of citizens that live in the province.
	Total uint `json:"total"`
	// Quote is the relative vaccination progression in the province.
	Quote float64 `json:"quote"`
	// VaccinationsPer1000Inhabitants is the number of vaccinated citizens per 1000 citizens.
	VaccinationsPer1000Inhabitants float64 `json:"vaccinations_per_1000_inhabitants"`
	// SecondVaccination is the vaccinations count for the second vaccination wave.
	SecondVaccination VaccinationCountData `json:"2nd_vaccination"`
}

// VaccinationCountData is a struct that is JSON-compatible with the RKI vaccination data API.
// It contains the actual vaccination counts and difference to the last day.
type VaccinationCountData struct {
	// Vaccinated is the number of citizens that already received a vaccination.
	Vaccinated uint `json:"vaccinated"`
	// DifferenceToThePreviousDay is the number of citizens that received a vaccination in the last 24h period.
	DifferenceToThePreviousDay uint `json:"difference_to_the_previous_day"`
}

// GetVaccinationData returns the vaccination stats for Germany from the RKI API endpoint.
// The data is usually updated once per week-day.
func GetVaccinationData() (*VaccinationData, error) {
	response, err := resty.New().
		SetHeader("accept", "application/json").
		SetHostURL(rkiVaccinationDataAPIURL).
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
