package math

import (
	"testing"
	"time"

	"github.com/drewwells/chargerstore/types"
)

var testMap = []struct {
	pct float32
}{
	{
		pct: 84.313728,
	},
	{
		pct: 74.90196,
	},
}

func TestBatteryCharging(t *testing.T) {

}

func TestCharge_regain(t *testing.T) {

	now := time.Now()
	var testMap = []struct {
		battery types.LastMsg
		amps    types.LastMsg
		volts   types.LastMsg
		rem     float64
	}{
		{
			// missing 1.5532605
			// guessometer probably reads 30miles
			battery: types.LastMsg{
				Data:        0.7490196,
				PublishTime: now.Add(-2 * time.Hour),
			},
			amps: types.LastMsg{
				Data: 8,
			},
			volts: types.LastMsg{
				Data: 120,
			},
			rem: 1.92,
		},
		{
			battery: types.LastMsg{
				Data:        MIN_PCT,
				PublishTime: now.Add(-10 * time.Hour),
			},
			amps: types.LastMsg{
				Data: 8,
			},
			volts: types.LastMsg{
				Data: 120,
			},
			rem: 9.6,
		},
		{
			battery: types.LastMsg{
				Data:        MIN_PCT,
				PublishTime: now.Add(-12 * time.Hour),
			},
			amps: types.LastMsg{
				Data: 8,
			},
			volts: types.LastMsg{
				Data: 120,
			},
			rem: MAX_ENERGY,
		},
	}

	for i, tm := range testMap {
		rem := guessRecharged(tm.battery, tm.amps, tm.volts)

		if e := tm.rem; e != rem {
			t.Errorf("%d got: %f wanted: %f", i, rem, e)
		}
	}
}
