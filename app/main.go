package main

import (
	"encoding/json"
	slog "log"
	"net/http"
	"sync"

	"google.golang.org/appengine"

	"github.com/drewwells/chargerstore/api"
	"github.com/drewwells/chargerstore/math"
	"github.com/drewwells/chargerstore/store"
)

func main() {
	appengine.Main()
}

func init() {
	http.HandleFunc("/api/v1/id/chargerate", rateHandler)
	http.HandleFunc("/api/v1/id/battery", batteryStatusHandler)

	router, err := api.New()
	if err != nil {
		slog.Fatal(err)
	}
	http.Handle("/", router)
	slog.Println("Web Server started")
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

const devID = "520041000351353337353037"

func batteryStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	stat, err := store.GetCarStatus(ctx, devID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	marshal(w, math.BatteryCharging(
		stat.LastSOC,
		stat.LastPower,
	))
}

func rateHandler(w http.ResponseWriter, r *http.Request) {
	marshal(w, struct {
		Amps  float64 `json:"amps"`
		Volts float64 `json:"volts"`
		Power float64 `json:"power"`
	}{
		Amps:  store.LastAmps.Data,
		Volts: store.LastVolts.Data,
		Power: store.LastPower.Data,
	})
}

const maxMessages = 10

var (
	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string
)
