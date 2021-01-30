package covid

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

const (
	diseaseApiUrl      = "https://disease.sh/v3/covid-19"
	allowNullParameter = "allowNull"
)

// GovernmentDiseaseData is a struct that is JSON-compatible with the result of the disease.sh API.
type GovernmentDiseaseData struct {
	// Updated contains the UNIX timestamp of the last update to the data.
	Updated uint64 `json:"updated"`
	// Province contains the name of the province or state from which this data is originated.
	// May also be "Total" to represent the data from the whole country.
	Province string `json:"province"`
	// Cases contains the number of all confirmed COVID-19 cases.
	Cases uint `json:"cases"`
	// CasePreviousDayChange is the number of new cases within the previous day.
	CasePreviousDayChange uint `json:"casePreviousDayChange"`
	// CasesPerHundredThousand is the number of cases per 100.000 citizens.
	CasesPerHundredThousand uint `json:"casesPerHundredThousand"`
	// SevenDayCasesPerHundredThousand is the number of new cases per 100.000 citizens during the last week.
	SevenDayCasesPerHundredThousand uint `json:"sevenDayCasesPerHundredThousand"`
	// The absolute number of casualties related to COVID-19.
	Deaths uint `json:"deaths"`
}

// GetGovernmentData fetches the latest government-published data for a given country.
// The parameter country controls the geographic region for which data is fetched, while allowNull controls if missing
// data should be returned as '0' or as 'nil'.
func GetGovernmentData(country string, allowNull bool) ([]GovernmentDiseaseData, error) {
	response, err := resty.New().
		SetHeader("accept", "application/json").
		SetQueryParam(allowNullParameter, strconv.FormatBool(allowNull)).
		SetHostURL(diseaseApiUrl).
		R().
		Get(fmt.Sprintf("/gov/%s", country))
	if err != nil {
		return nil, err
	}

	var result []GovernmentDiseaseData
	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return nil, err
	}

	return result, err
}
