package covid

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetGovernmentData(t *testing.T) {
	data, err := GetGovernmentData("de", false)
	require.NoError(t, err)
	require.NotNil(t, data)

	if assert.Len(t, data, 17) {
		total := data[16]

		assert.Equal(t, "Total", total.Province)
		assert.True(t, total.Deaths > 0)
		assert.True(t, total.Cases > 0)
		assert.True(t, total.CasePreviousDayChange > 0)
		assert.True(t, total.CasesPerHundredThousand > 0)
		assert.True(t, total.SevenDayCasesPerHundredThousand > 0)
	}
}
