package math

import "time"

// 234.0000 volts
// 13.20000 amps
// 84.313728
// 16.5 kWh (10.9 kWh usable)

const (
	MAX_POWER = 10.9 // only seen 10.2 on my vehicle

	// As reported from car
	MAX_PCT = 0.84313728
	// Derived from usable numbers from Wikipedia, needs verification
	MIN_PCT = MAX_PCT - 0.660606061
	MAX_KWH = 16.5 * MAX_PCT
	MIN_KWH = 16.5 * MIN_PCT
)

// Power provided in kwh
func Power(amp float32, volt float32) float32 {
	return amp * volt / float32(1000)
}

// Remaining determines the amount of time remaining to charge to full
func Remaining(currentPct float32, power float32) time.Duration {
	if power == 0 {
		return time.Duration(-1 * time.Minute)
	}
	deficitPwr := (MAX_PCT - currentPct) * MAX_POWER
	hours := (deficitPwr * 1000) / power
	return time.Duration(hours) * 60 * time.Minute
}
