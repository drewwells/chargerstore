package math

import (
	"math"
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

func float64eq(a, b float64) bool {
	return math.Abs(a-b) < 0.00001
}

func TestCharge_BatteryCharging(t *testing.T) {
	now := time.Now()
	var testMap = []struct {
		Deficit              float64
		battery, amps, volts types.LastMsg

		eCurrent      int
		eV120Standard int
		eV120Max      int
		eV240         int
		eEstimate     bool
	}{
		{
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
			eEstimate: true,
		},
		{
			battery: types.LastMsg{
				Data:        0.7490196,
				PublishTime: now.Add(-10 * time.Minute),
			},
			amps: types.LastMsg{
				Data: 8,
			},
			volts: types.LastMsg{
				Data: 120,
			},
			eEstimate:     true,
			eCurrent:      3000,
			eV120Standard: 3000,
			eV120Max:      2000,
			eV240:         869,
		},
	}

	for i, tm := range testMap {
		bc := BatteryCharging(tm.battery, tm.amps, tm.volts)
		if e := tm.eEstimate; e != bc.Estimate {
			t.Errorf("%d got: %t wanted: %t", i, bc.Estimate, e)
		}

		if e := int(tm.eCurrent); e != int(bc.Current.Duration.Seconds()) {
			t.Errorf("%d got: %d wanted: %d", i, bc.Current.Duration, e)
		}

		if e := int(tm.eV120Standard); e != int(bc.V120Standard.Duration.Seconds()) {
			t.Errorf("%d got: %d wanted: %d", i, int(bc.V120Standard.Duration.Seconds()), e)
		}

		if e := int(tm.eV120Max); e != int(bc.V120Max.Duration.Seconds()) {
			t.Errorf("%d got: %d wanted: %d", i, bc.V120Max.Duration, e)
		}

		if e := int(tm.eV240); e != int(bc.V240.Duration.Seconds()) {
			t.Errorf("%d got: %d wanted: %d", i, bc.V240.Duration, e)
		}

	}
}

func TestCharge_regain(t *testing.T) {

	now := time.Now()
	var testMap = []struct {
		battery types.LastMsg
		amps    types.LastMsg
		volts   types.LastMsg
		reg     float64
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
			reg: 1.92,
		},
		{
			battery: types.LastMsg{
				Data:        MIN_PCT,
				PublishTime: now.Add(-10 * time.Minute),
			},
			amps: types.LastMsg{
				Data: 8,
			},
			volts: types.LastMsg{
				Data: 120,
			},
			reg: 0.16,
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
			reg: 9.6,
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
			reg: float64(MAX_ENERGY),
		},
	}

	for i, tm := range testMap {
		reg := guessRecharged(tm.battery, tm.amps, tm.volts)

		if e := tm.reg; !float64eq(e, reg) {
			t.Errorf("%d got: %f wanted: %f", i, reg, e)
		}
	}
}
