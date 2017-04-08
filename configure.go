package chargerstore

import (
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
)

type Options struct {
	ps                 *pubsub.Client
	subName, topicName string
	topic              *pubsub.Topic
}

func NewPS() (*Options, error) {
	ctx := context.Background()
	c, err := pubsub.NewClient(ctx, "particle-volt")
	if err != nil {
		return nil, err
	}
	return &Options{
		ps: c,
	}, nil
}

func (o *Options) Subscribe(subName string, topicName string) {
	ctx := context.Background()
	topic, err := o.getTopic(ctx, topicName)
	if err != nil {
		log.Fatal(err)
	}

	o.topic = topic
	sub, err := o.ps.CreateSubscription(ctx, subName, topic, 0, nil)
	if err != nil {
		log.Fatal(err)
	}

	go o.subscribe(sub)
}

func (o *Options) getTopic(ctx context.Context, topicName string) (*pubsub.Topic, error) {

	topic := o.ps.Topic(o.subName)
	if exists, err := o.ps.Topic(o.topicName).Exists(ctx); err != nil {
		return nil, err
	} else if !exists {
		if _, err := o.ps.CreateTopic(ctx, o.subName); err != nil {
			return nil, err
		}
	}
	return topic, nil
}

type CarMsgMeta struct {
	Data        CarMsg `json:"data"`
	TTL         int
	PublishedAt time.Time `json:"published_at"`
	Name        string    `json:"name"`
	CoreID      string    `json:"coreid"`
}

type CarMsg struct {
	VehicleSpeed float32 `json:"VEHICLE_SPEED"`
	AirTemp      float32 `json:"AMBIENT_AIR_TEMPERATURE"`
	CMV          float32 `json:"CONTROL_MODULE_VOLTAGE"`
	Fuel         float32 `json:"FUEL_TANK_LEVEL_INPUT"`
	ChargerAmps  float32 `json:"CHARGER_AMPS_IN"`
	ChargerVolts float32 `json:"CHARGER_VOLTS_IN"`
	Battery      float32 `json:"EXTENDED_HYBRID_BATTERY_REMAINING_LIFE"`
}

func (o *Options) subscribe(sub *pubsub.Subscription) {
	ctx := context.Background()
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {

		log.Printf("% #v\n", msg)
		var cm CarMsgMeta
		if err := json.Unmarshal(msg.Data, &cm); err != nil {
			log.Printf("could not decode message data: %#v", msg)
			msg.Ack()
			return
		}

		msg.Ack()
		log.Printf("received %#v\n", cm)
	})
	if err != nil {
		log.Fatal(err)
	}
}
