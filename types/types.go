package types

import (
	"time"

	"cloud.google.com/go/pubsub"
)

// {
//   "message": {
//     "data": "eyJWRUhJQ0xFX1NQRUVEIjotMS4wMDAwMDAsIkFNQklFTlRfQUlSX1RFTVBFUkFUVVJFIjoyNy4wMDAwMDAsIkNPTlRST0xfTU9EVUxFX1ZPTFRBR0UiOi0xLjAwMDAwMCwiRlVFTF9UQU5LX0xFVkVMX0lOUFVUIjotMS4wMDAwMDAsIkNIQVJHRV9BTVBTX0lOIjo4LjAwMDAwMCwiQ0hBUkdFUl9WT0xUU19JTiI6MTE0LjAwMDAwMCwiRVhURU5ERURfSFlCUklEX0JBVFRFUllfUEFDS19SRU1BSU5JTkdfTElGRSI6LTEuMDAwMDAwfQ==",
//     "attributes": {
//       "device_id": "520041000351353337353037",
//       "event": "CAR",
//       "published_at": "2017-04-08T18:55:31.839Z"
//     },
//     "message_id": "68413786577982",
//     "messageId": "68413786577982",
//     "publish_time": "2017-04-08T18:55:32.029Z",
//     "publishTime": "2017-04-08T18:55:32.029Z"
//   },
//   "subscription": "projects/particle-volt/subscriptions/carpull"
// }
type PushRequest struct {
	Message      *pubsub.Message `json:"message"`
	Subscription string
}

// CarMsg is the format incoming from particle
// {
//   "VEHICLE_SPEED": -1,
//   "AMBIENT_AIR_TEMPERATURE": 26,
//   "CONTROL_MODULE_VOLTAGE": -1,
//   "FUEL_TANK_LEVEL_INPUT": -1,
//   "CHARGE_AMPS_IN": 0,
//   "CHARGER_VOLTS_IN": 0,
//   "EXTENDED_HYBRID_BATTERY_PACK_REMAINING_LIFE": 84.313728
// }
type CarMsg struct {
	ID                string    `json:"id"`
	VehicleSpeed      float64   `json:"vehicle_speed"`
	AirTemp           float64   `json:"temp"`
	CMV               float64   `json:"cmv"`
	Fuel              float64   `json:"fuel_tank"`
	ChargerAmps       float64   `json:"amps_in"`
	ChargerVolts      float64   `json:"volts_in"`
	Battery           float64   `json:"soc"`
	HVDischargeAmps   float64   `json:"hv_amps"`
	HVVolts           float64   `json:"hv_volts"`
	EVMilesThisCharge float64   `json:"ev_miles_cycle"`
	PublishTime       time.Time `json:"publish_time"`
	Event             string    `json:"event"`
	DeviceID          string    `json:"device_id"`
}

type CarStatus struct {
	DeviceID  string    `json:"device_id"`
	LastSOC   LastMsg   `json:"soc"`
	LastAmps  LastMsg   `json:"amps"`
	LastVolts LastMsg   `json:"volts"`
	CreatedAt time.Time `json:"created_at"`
}

type LastMsg struct {
	Data        float64
	PublishTime time.Time
}

type Charger struct {
	Duration time.Duration
	Minutes  float64
}

type ChargeState struct {
	LastSOCTime time.Time `json:"last_reported_soc"`
	Percent     float64   `json:"percent"`
	Deficit     float64   `json:"deficit_kwh"`
	Regained    float64   `json:"regained_kwh"`
}

type BatteryCharging struct {
	State ChargeState `json:"state"`

	Deficit       float64   `json:"deficit"`      // energy below maximum
	Estimate      bool      `json:"estimate"`     // indicate that we don't have SOC confirmation
	LastPublished time.Time `json:"published_at"` // last time SOC was reported
	Current       Charger   `json:"current"`      // current charging rate
	V120Standard  Charger   `json:"v120_standard"`
	V120Max       Charger   `json:"v120_max"`
	V240          Charger   `json:"v240"`
}
