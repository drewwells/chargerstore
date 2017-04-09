package math

import (
	"math"
	"time"

	"github.com/drewwells/chargerstore/types"
)

// 234.0000 volts
// 13.20000 amps
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

	POWER_120V_8A  = 120 * 8 / 1000    // 0.96
	POWER_120V_12A = 120 * 12 / 1000   // 1.44
	POWER_240V     = 240 * 13.8 / 1000 // 3.312
)

// Power provided in kwh
func Power(amp float64, volt float64) float64 {
	return amp * volt / float64(1000)
}

// TimeToCharge determines the amount of time remaining to charge to full
func TimeToCharge(currentPct float64, power float64) time.Duration {
	if power == 0 {
		return time.Duration(-1 * time.Minute)
	}
	deficitPwr := (MAX_PCT - currentPct) * MAX_ENERGY
	hours := (deficitPwr * 1000) / power
	return time.Duration(hours) * 60 * time.Minute
}

// guessRecharged, volt is terrible at reporting SOC while off. So instead
// guess how much has been recharged
func guessRecharged(lastBattery types.LastMsg, lastAmps types.LastMsg, lastVolts types.LastMsg) float64 {
	since := time.Since(lastBattery.PublishTime)

	hrs := float64(since / time.Hour)
	if lastAmps.Data == 0 || lastVolts.Data == 0 {
		return 0
	}

	// energy (kwh) = power / 1000
	pwr := Power(lastAmps.Data, lastVolts.Data)
	return math.Min(pwr*hrs, MAX_ENERGY)
}

func BatteryCharging(lastBattery types.LastMsg, lastAmps types.LastMsg, lastVolts types.LastMsg) types.BatteryCharging {
	//func BatteryCharging(currentPct float32, currentPwr float32, lastPublished time.Time) types.BatteryCharging {

	// charge exceeds maximum, all 0s for charging
	if lastBattery.Data >= MAX_PCT {
		return types.BatteryCharging{Deficit: 0}
	}

	// SOC publishing usually stops shortly after car turns off
	// So up to date SOC is not reliable.
	//
	// Here we make assumptions.
	// 1. If car reported non zero volt & amps recently <5mins, then it has
	//    been charging since Car was turned off
	currentPct := lastBattery.Data
	deficit := (MAX_PCT - currentPct) * MAX_ENERGY
	regained := guessRecharged(lastBattery, lastAmps, lastVolts)
	if regained > deficit {
		return types.BatteryCharging{Estimate: true}
	}

	// no recent SOC data, indicate so
	estimate := time.Since(lastBattery.PublishTime) > 5*time.Minute

	bc := types.BatteryCharging{
		Deficit:  deficit,
		Estimate: estimate,
	}

	estMissing := deficit - regained
	estPct := estMissing / MAX_ENERGY
	currentPwr := lastVolts.Data * lastAmps.Data / 1000

	current := TimeToCharge(estPct, currentPwr)
	bc.Current = types.Charger{
		Duration: current,
		Minutes:  current.Minutes(),
	}
	v120std := TimeToCharge(estPct, POWER_120V_8A)
	bc.V120Standard = types.Charger{
		Duration: v120std,
		Minutes:  v120std.Minutes(),
	}
	v120max := TimeToCharge(estPct, POWER_120V_12A)
	bc.V120Max = types.Charger{
		Duration: v120max,
		Minutes:  v120max.Minutes(),
	}
	v240 := TimeToCharge(estPct, POWER_240V)
	bc.V240 = types.Charger{
		Duration: v240,
		Minutes:  v240.Minutes(),
	}

	return bc
}
