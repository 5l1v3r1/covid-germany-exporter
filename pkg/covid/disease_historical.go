package covid

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type HistoricalData struct {
	Country  string   `json:"country"`
	Province []string `json:"province"`
	TimeLine Timeline `json:"timeline"`
}

type Timeline struct {
	Cases     map[string]uint `json:"cases"`
	Deaths    map[string]uint `json:"deaths"`
	Recovered map[string]uint `json:"recovered"`
}

func GetHistoricalData(country, lastdays string) (*HistoricalData, error) {
	response, err := resty.New().
		SetHeader("accept", "application/json").
		SetQueryParam("lastdays", lastdays).
		SetHostURL("https://disease.sh/v3/covid-19").
		R().
		Get(fmt.Sprintf("/historical/%s", country))
	if err != nil {
		return nil, err
	}

	result := new(HistoricalData)
	if err := json.Unmarshal(response.Body(), result); err != nil {
		return nil, err
	}

	return result, err
}
