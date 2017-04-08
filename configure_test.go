package chargerstore

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
)

var testmap = [][]byte{
	[]byte(`{"data":"{\"VEHICLE_SPEED\":-1.000000,\"AMBIENT_AIR_TEMPERATURE\":20.000000,\"CONTROL_MODULE_VOLTAGE\":-1.000000,\"FUEL_TANK_LEVEL_INPUT\":-1.000000,\"CHARGE_AMPS_IN\":0.000000,\"CHARGER_VOLTS_IN\":0.000000,\"EXTENDED_HYBRID_BATTERY_PACK_REMAINING_LIFE\":-1.000000}","ttl":"60","published_at":"2017-04-08T04:29:37.004Z","coreid":"520041000351353337353037","name":"CAR"}`)}

func TestSubscribe(t *testing.T) {

	opts, err := NewPS()
	if err != nil {
		t.Fatal(err)
	}

	opts.Subscribe("carpull", "CAR")
	ctx := context.Background()

	for _, tm := range testmap {
		_, err := opts.topic.Publish(ctx, &pubsub.Message{Data: tm}).Get(ctx)
		if err != nil {
			t.Fatal(err)
		}
		//log.Printf("Published update to Pub/Sub for Book ID %d: %v", bookID, err)
	}

}
