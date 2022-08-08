package drone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Drone_Locate(t *testing.T) {
	tt := []struct {
		name     string
		drone    Drone
		sectorID float64
		expected float64
	}{
		{
			name: "HappyPath",
			drone: Drone{
				X:   123.12,
				Y:   456.56,
				Z:   789.89,
				Vel: 20.0,
			},
			sectorID: 1.0,
			expected: 1389.5700000000002,
		},
		{
			name: "HappyPathDifferentSector",
			drone: Drone{
				X:   123.12,
				Y:   456.56,
				Z:   789.89,
				Vel: 20.0,
			},
			sectorID: 2.0,
			expected: 2759.1400000000003,
		},
	}

	for _, c := range tt {
		t.Run(c.name, func(t *testing.T) {
			dr := c.drone
			loc := dr.Locate(c.sectorID)
			assert.Equal(t, c.expected, loc)
		})
	}
}
