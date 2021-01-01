package covid

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type GovernmentDiseaseData struct {
	Updated                         uint64 `json:"updated"`
	Province                        string `json:"province"`
	Cases                           uint   `json:"cases"`
	CasePreviousDayChange           uint   `json:"casePreviousDayChange"`
	CasesPerHundredThousand         uint   `json:"casesPerHundredThousand"`
	SevenDayCasesPerHundredThousand uint   `json:"sevenDayCasesPerHundredThousand"`
	Deaths                          uint   `json:"deaths"`
}

func GetGovernmentData(country string, allowNull bool) ([]GovernmentDiseaseData, error) {
	response, err := resty.New().
		SetHeader("accept", "application/json").
		SetQueryParam("allowNull", strconv.FormatBool(allowNull)).
		SetHostURL("https://disease.sh/v3/covid-19").
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
