package chargerstore

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"google.golang.org/appengine/datastore"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
)

const projectid = "particle-volt"

const (
	carbucket = "car"
)

type Options struct {
	ps      *pubsub.Client
	topic   *pubsub.Topic
	muCount sync.RWMutex
	count   int64
}

func (o *Options) Count() int64 {
	o.muCount.RLock()
	defer o.muCount.RUnlock()
	return o.count
}

func NewPS() (*Options, error) {
	ctx := context.Background()
	c, err := pubsub.NewClient(ctx, projectid)
	if err != nil {
		return nil, err
	}
	return &Options{
		ps: c,
	}, nil
}

// subscription created with
// -> % gcloud beta pubsub subscriptions create carpull --topic CAR --push-endpoint https://particle-volt.appspot.com/pubsub/push\?token\=
func (o *Options) Subscribe(subName string, topicName string) {
	ctx := context.Background()
	topic, err := o.getTopic(ctx, subName, topicName)
	if err != nil {
		log.Fatal(err)
	}

	o.topic = topic
	sub := o.ps.Subscription(subName)

	go o.subscribe(sub)
}

func (o *Options) getTopic(ctx context.Context, subName, topicName string) (*pubsub.Topic, error) {

	return o.ps.Topic(subName), nil
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

func (o *Options) subscribe(sub *pubsub.Subscription) {
	ctx := context.Background()
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {

		var cm CarMsg
		if err := json.Unmarshal(msg.Data, &cm); err != nil {
			log.Printf("could not decode message data: %#v", msg)
			msg.Ack()
			return
		}
		cm.PublishTime = msg.PublishTime
		cm.Event = msg.Attributes["event"]
		cm.DeviceID = msg.Attributes["device_id"]

		o.muCount.Lock()
		o.count++
		o.muCount.Unlock()

		msg.Ack()
		log.Printf("received %#v\n", cm)
		//k := datastore.NewKey(ctx, carbucket, msg.ID, 0, nil)
		k := datastore.NewIncompleteKey(ctx, carbucket, nil)
		if _, err := datastore.Put(ctx, k, &cm); err != nil {
			// Handle err
			log.Printf("failed to save %s: %s", msg.ID, err)
		}
	})
	if err != nil {
		log.Printf("error receiving event: %s", err)
	}
}
