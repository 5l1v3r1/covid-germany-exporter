package covid

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetVaccinationData(t *testing.T) {
	data, err := GetVaccinationData()
	require.NoError(t, err)
	require.NotNil(t, data)

	if assert.Len(t, data.States, 16) {
		vaccinationData := data.States["Nordrhein-Westfalen"]
		require.NotNil(t, vaccinationData)

		assert.True(t, vaccinationData.Quote > 0)
		assert.True(t, vaccinationData.Vaccinated > 0)
		assert.True(t, vaccinationData.Total > 0)
	}
}
