package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/drewwells/chargerstore"
	"github.com/drewwells/chargerstore/math"
	"github.com/drewwells/chargerstore/types"
)

func main() {
	appengine.Main()
}

func init() {

	// opts, err := chargerstore.NewPS()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//opts.Subscribe("CAR", "carpull")
	http.HandleFunc("/pubsub/push", pushHandler)
	http.HandleFunc("/summary", summaryHandler)
	http.HandleFunc("/car/id/laststatus", lastStatusHandler)
	http.HandleFunc("/car/id/chargerate", rateHandler)
	http.HandleFunc("/car/id/battery", batteryStatusHandler)
}

func marshal(w http.ResponseWriter, v interface{}) error {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
	return nil
}

func lastStatusHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]types.LastMsg)
	resp["amps"] = chargerstore.LastAmps
	resp["volts"] = chargerstore.LastVolts
	resp["battery"] = chargerstore.LastBattery
	marshal(w, resp)
}

func batteryStatusHandler(w http.ResponseWriter, r *http.Request) {
	p := math.Power(chargerstore.LastVolts.Data, chargerstore.LastAmps.Data)
	currentPct := chargerstore.LastBattery.Data / 100
	timeToCharge := math.TimeToChargePCT(
		currentPct,
		p,
	)
	marshal(w, math.BatteryCharging(currentPct, p, chargerstore.LastBattery.PublishTime))
}

func rateHandler(w http.ResponseWriter, r *http.Request) {
	p := math.Power(chargerstore.LastVolts.Data, chargerstore.LastAmps.Data)
	marshal(w, struct {
		Amps  float32 `json:"amps"`
		Volts float32 `json:"volts"`
		Power float32 `json:"power"`
	}{
		Power: p,
		Amps:  chargerstore.LastAmps.Data,
		Volts: chargerstore.LastVolts.Data,
	})
}

func summaryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This worker has processed %d events.", chargerstore.Count())
}

const maxMessages = 10

var (
	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string
)

func pushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var req types.PushRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err := fmt.Errorf("Could not decode body: %v", err)
		log.Errorf(ctx, "%s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Infof(ctx, "request % #v\n", req)
	msg, err := chargerstore.Process(ctx, req.Message)
	if err != nil {
		err := fmt.Errorf("failed to process msg: %s", err)
		log.Errorf(ctx, "%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, "processed message from: %s", msg.DeviceID)
}
