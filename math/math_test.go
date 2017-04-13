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

func TestCharge_cardisplay(t *testing.T) {
	t.Skip()
	now := time.Now()
	var testMap = []struct {
		Deficit        float64
		battery, power types.LastMsg

		eCurrent      int
		eV120Standard int
		eV120Max      int
		eV240         int
		eEstimate     bool
	}{
		/*{
			battery: types.LastMsg{
				Data:        0.4945,
				PublishTime: now,
			},
			amps: types.LastMsg{
				Data: 8,
			},
			volts: types.LastMsg{
				Data: 120,
			},
			eEstimate:     true,
			eCurrent:      3000,
			eV120Standard: 276, //22:30 - 17:54,
			//eV120Max:      2000,
			eV240: 66, //1900 - 1754,

			// system guess
			// v240 36.5mins
			// v120max 84mins
			// v120std 125mins
		},*/
		{
			battery: types.LastMsg{
				Data:        0.65882355,
				PublishTime: now,
			},
			power: types.LastMsg{
				Data: 0.700,
			},
			eEstimate:     true,
			eCurrent:      225, // 7:00 - 10:45
			eV120Standard: 225, // 7:00 - 10:45
			//eV120Max:      2000,
			eV240: 60, // 7:00 - 8:00,

			// system guess
			// current 81mins
			// v120std 76mins
			// v120max 51mins
			// v240 22mins
		},
	}

	for i, tm := range testMap {
		bc := BatteryCharging(tm.battery, tm.power)
		if e := tm.eEstimate; e != bc.Estimate {
			t.Errorf("%d got: %t wanted: %t", i, bc.Estimate, e)
		}

		sec := int(bc.Current.Duration.Seconds())
		if e := int(tm.eCurrent); e != sec {
			t.Errorf("%d got: %d wanted: %d", i, sec, e)
		}

		sec = int(bc.V120Standard.Duration.Seconds())
		if e := int(tm.eV120Standard); e != sec {
			t.Errorf("%d got: %d wanted: %d", i, sec, e)
		}

		sec = int(bc.V120Max.Duration.Seconds())
		if e := int(tm.eV120Max); tm.eV120Max > 0 && e != sec {
			t.Errorf("%d got: %d wanted: %d", i, sec, e)
		}

		sec = int(bc.V240.Duration.Seconds())
		if e := int(tm.eV240); e != sec {
			t.Errorf("%d got: %d wanted: %d", i, sec, e)
		}

	}

}

func TestCharge_BatteryCharging(t *testing.T) {
	now := time.Now()
	var testMap = []struct {
		battery, power types.LastMsg

		eCurrent      int
		eV120Standard int
		eV120Max      int
		eV240         int
		eEstimate     bool
	}{
		{
			battery: types.LastMsg{
				Data:        MIN_PCT,
				PublishTime: now,
			},
			power: types.LastMsg{
				Data: 0.711,
			},
			eEstimate:     false,
			eCurrent:      56757,
			eV120Standard: 56749,
			eV120Max:      37832,
			eV240:         15133,
		},
	}

	for i, tm := range testMap {
		bc := BatteryCharging(tm.battery, tm.power)
		if e := tm.eEstimate; e != bc.Estimate {
			t.Errorf("%d got: %t wanted: %t", i, bc.Estimate, e)
		}

		sec := int(bc.Current.Duration.Seconds())
		if e := int(tm.eCurrent); e != sec {
			t.Errorf("%d got: %d wanted: %d", i, sec, e)
		}

		sec = int(bc.V120Standard.Duration.Seconds())
		if e := int(tm.eV120Standard); e != sec {
			t.Errorf("%d got: %d wanted: %d", i, sec, e)
		}

		sec = int(bc.V120Max.Duration.Seconds())
		if e := int(tm.eV120Max); e != sec {
			t.Errorf("%d got: %d wanted: %d", i, sec, e)
		}

		sec = int(bc.V240.Duration.Seconds())
		if e := int(tm.eV240); e != sec {
			t.Errorf("%d got: %d wanted: %d", i, sec, e)
		}

	}
}

func TestCharge_regain(t *testing.T) {
	t.Skip()
	now := time.Now()
	var testMap = []struct {
		battery, power types.LastMsg
		reg            float64
	}{
		{
			// missing 1.5532605
			// guessometer probably reads 30miles
			battery: types.LastMsg{
				Data:        0.7490196,
				PublishTime: now.Add(-2 * time.Hour),
			},
			power: types.LastMsg{
				Data: 0.700,
			},
			reg: 1.92,
		},
		{
			battery: types.LastMsg{
				Data:        MIN_PCT,
				PublishTime: now.Add(-10 * time.Minute),
			},
			power: types.LastMsg{
				Data: 0.700,
			},
			reg: 0.16,
		},
		{
			battery: types.LastMsg{
				Data:        MIN_PCT,
				PublishTime: now.Add(-10 * time.Hour),
			},
			power: types.LastMsg{
				Data: 0.700,
			},
			reg: 9.6,
		},
		{
			battery: types.LastMsg{
				Data:        MIN_PCT,
				PublishTime: now.Add(-12 * time.Hour),
			},
			power: types.LastMsg{
				Data: 0.700,
			},
			reg: float64(MAX_ENERGY),
		},
	}

	for i, tm := range testMap {
		reg := guessRecharged(tm.battery, tm.power)

		if e := tm.reg; !float64eq(e, reg) {
			t.Errorf("%d got: %f wanted: %f", i, reg, e)
		}
	}
}
