package types

import (
	"encoding/json"
	"fmt"
	"testing"

	"cloud.google.com/go/pubsub"
)

func TestCarMsg(t *testing.T) {
	msg := &pubsub.Message{
		ID:   "",
		Data: []uint8{},
		Attributes: map[string]string{
			"device_id": "520041000351353337353037", "event": "CAR",
			"published_at": "2017-04-09T20:29:53.867Z",
		},
	}

	var cm CarMsg
	fmt.Println(string(msg.Data))
	err := json.Unmarshal(msg.Data, &cm)
	if err != nil {
		t.Fatal(err)
	}
}
