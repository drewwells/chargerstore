package math

import (
	"math"
	"time"

	"github.com/drewwells/chargerstore/types"
)

// 234.0000 volts
// 13.40000 amps
// 2.915 power
// 0.69
// 70mins
//
// 0.84313728
// 16.5 kWh (10.9 kWh usable)
// MAX 13.91176512kwh
// MIN 3.7117

const (
	TOTAL_ENERGY = 16.5
	MAX_ENERGY   = 10.2 // only seen 10.2 on my vehicle

	// As reported from car
	MAX_PCT = 0.84313728
	// Derived from usable numbers from Wikipedia, needs verification
	MIN_PCT = 0.224955462 // ((MAX_PCT * TOTAL_POWER)-MAX_POWER)*TOTAL_POWER
	MAX_KWH = 16.5 * MAX_PCT
	MIN_KWH = 16.5 * MIN_PCT

	// Power levels reduced due to output emitted by car
	POWER_120V_8A  = 0.40 // 0.96  // 120 * 8 / 1000
	POWER_120V_12A = 0.6  // 1.44  // 120 * 12 / 1000
	POWER_240V     = 1.5  // 3.312 // 240 * 13.8 / 1000

	POWER_FACTOR = 0.5625
)

// Power provided in kwh
func Power(amp float64, volt float64) float64 {
	return amp * volt / float64(1000)
}

func TimeToCharge(currentEnergy, power float64) time.Duration {
	if power == 0 {
		return time.Duration(-1 * time.Minute)
	}
	deficit := float64(MAX_ENERGY) - currentEnergy
	hours := deficit / power
	d := time.Duration(hours * float64(time.Hour))
	return d
}

// TimeToCharge determines the amount of time remaining to charge to full
func TimeToChargePCT(currentPct float64, power float64) time.Duration {
	return TimeToCharge(currentPct*TOTAL_ENERGY, power)
}

// guessRecharged, volt is terrible at reporting SOC while off. So instead
// guess how much has been recharged
func guessRecharged(lastBattery types.LastMsg, power float64) float64 {
	since := time.Since(lastBattery.PublishTime)

	hrs := float64(since) / float64(time.Hour)
	if power <= 0 {
		return 0
	}

	// energy (kwh) = power / 1000
	energy := power * hrs
	return math.Min(energy, MAX_ENERGY)
}

func BatteryCharging(lastBattery types.LastMsg, lastPower types.LastMsg) types.BatteryCharging {

	// TODO: report errors for the various errorneous conditions

	// charge exceeds maximum, all 0s for charging
	if lastBattery.Data >= MAX_PCT {
		return types.BatteryCharging{}
	}

	// SOC publishing usually stops shortly after car turns off
	// So up to date SOC is not reliable.
	//
	// Here we make assumptions.
	// 1. If car reported non zero volt & amps recently <5mins, then it has
	//    been charging since Car was turned off
	currentPct := lastBattery.Data
	deficit := float64((MAX_PCT - currentPct) * TOTAL_ENERGY)
	regained := guessRecharged(lastBattery, lastPower.Data /* like this makes any sense */)

	// no recent SOC data, indicate so
	estimate := time.Since(lastBattery.PublishTime) > 5*time.Minute

	bc := types.BatteryCharging{
		State: types.ChargeState{
			LastSOCTime: lastBattery.PublishTime,
			Percent:     lastBattery.Data,
			Deficit:     deficit,
			Regained:    regained,
		},
		Estimate: estimate,
	}

	if regained > deficit {
		// Short circuit if regained has exceeded deficit
		bc.Estimate = true
		return bc
	}

	estCurrent := MAX_ENERGY - deficit + regained
	// estimates are calculated very pessimistically
	reducedPower := lastPower.Data * 0.5625
	current := TimeToCharge(estCurrent, reducedPower)
	bc.Current = types.Charger{
		Duration: current,
		Minutes:  current.Minutes(),
	}
	v120std := TimeToCharge(estCurrent, POWER_120V_8A)
	bc.V120Standard = types.Charger{
		Duration: v120std,
		Minutes:  v120std.Minutes(),
	}
	v120max := TimeToCharge(estCurrent, POWER_120V_12A)
	bc.V120Max = types.Charger{
		Duration: v120max,
		Minutes:  v120max.Minutes(),
	}
	v240 := TimeToCharge(estCurrent, POWER_240V)
	bc.V240 = types.Charger{
		Duration: v240,
		Minutes:  v240.Minutes(),
	}

	return bc
}
