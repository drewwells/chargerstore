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
type CarMsg struct {
	VehicleSpeed float32   `json:"VEHICLE_SPEED"`
	AirTemp      float32   `json:"AMBIENT_AIR_TEMPERATURE"`
	CMV          float32   `json:"CONTROL_MODULE_VOLTAGE"`
	Fuel         float32   `json:"FUEL_TANK_LEVEL_INPUT"`
	ChargerAmps  float32   `json:"CHARGER_AMPS_IN"`
	ChargerVolts float32   `json:"CHARGER_VOLTS_IN"`
	Battery      float32   `json:"EXTENDED_HYBRID_BATTERY_REMAINING_LIFE"`
	PublishTime  time.Time `json:"publish_time"`
	Event        string    `json:"event"`
	DeviceID     string    `json:"device_id"`
}

type LastMsg struct {
	Data        float32
	PublishTime time.Time
}
