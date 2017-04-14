package chargerstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/appengine/log"

	"github.com/drewwells/chargerstore/store"
	"github.com/drewwells/chargerstore/types"

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
		// slog.Fatal(err)
	}

	o.topic = topic
	sub := o.ps.Subscription(subName)

	go o.subscribe(sub)
}

func (o *Options) getTopic(ctx context.Context, subName, topicName string) (*pubsub.Topic, error) {

	return o.ps.Topic(subName), nil
}

var (
	muCount sync.RWMutex
	count   int64
)

func Count() int64 {
	muCount.RLock()
	defer muCount.RUnlock()
	return count
}

func Process(ctx context.Context, msg *pubsub.Message) (types.CarMsg, error) {
	var cm types.CarMsg
	if msg == nil {
		return cm, errors.New("empty message")
	}
	log.Infof(ctx, "data %s", string(msg.Data))
	if err := json.Unmarshal(msg.Data, &cm); err != nil {
		err := fmt.Errorf("could not decode message data: %#v", msg)
		log.Errorf(ctx, err.Error())
		return cm, err
	}
	cm.ID = msg.ID
	cm.PublishTime = msg.PublishTime
	cm.Event = msg.Attributes["event"]
	cm.DeviceID = msg.Attributes["device_id"]

	muCount.Lock()
	count++
	muCount.Unlock()

	log.Infof(ctx, "received: %#v\n", cm)
	store.PutCarMsg(ctx, &cm)

	processLastMsg(cm)

	err := store.PutCarStatus(ctx, &types.CarStatus{
		DeviceID:  msg.Attributes["device_id"],
		LastSOC:   LastSOC,
		LastAmps:  LastAmps,
		LastVolts: LastVolts,
		LastPower: LastPower,
		CreatedAt: time.Now(),
	})

	return cm, err
}

func processLastMsg(cm types.CarMsg) {
	// battery tends to report 0, probably an error on C side
	if cm.Battery > 0 {
		LastSOC = types.LastMsg{
			Data:        cm.Battery / 100, // battery is sent as pct * 100
			PublishTime: cm.PublishTime,
		}
	}

	// Power constantly reports 0, always ignore these unless it's been
	// 5 minutes since a non-zero update has been published
	if cm.ChargerPower > -1 || time.Since(LastPower.PublishTime) > 5*time.Minute {
		LastPower = types.LastMsg{
			Data:        cm.ChargerPower / 1000,
			PublishTime: cm.PublishTime,
		}
	}

	if cm.ChargerAmps > -1 {
		LastAmps = types.LastMsg{
			Data:        cm.ChargerAmps,
			PublishTime: cm.PublishTime,
		}
	}

	if cm.ChargerVolts > -1 {
		LastVolts = types.LastMsg{
			Data:        cm.ChargerVolts,
			PublishTime: cm.PublishTime,
		}
	}
}

func (o *Options) subscribe(sub *pubsub.Subscription) {
	ctx := context.Background()
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {

		_, err := Process(ctx, msg)
		if err != nil {
			log.Errorf(ctx, "failed to process msg: %s", err)
		}

	})
	if err != nil {
		log.Errorf(ctx, "error receiving event: %s", err)
	}
}
