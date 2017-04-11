package main

import (
	"encoding/json"
	slog "log"
	"net/http"
	"sync"

	"google.golang.org/appengine"

	"github.com/drewwells/chargerstore"
	"github.com/drewwells/chargerstore/api"
	"github.com/drewwells/chargerstore/math"
	"github.com/drewwells/chargerstore/store"
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
	http.HandleFunc("/api/v1/car/id/laststatus", lastStatusHandler)
	http.HandleFunc("/api/v1/car/id/chargerate", rateHandler)
	http.HandleFunc("/api/v1/car/id/battery", batteryStatusHandler)

	// http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
	// 	ctx := appengine.NewContext(r)
	// 	w.Write([]byte("to test\n"))
	// 	test(ctx)
	// 	w.Write([]byte("wrote it"))
	// })

	router, err := api.New()
	if err != nil {
		slog.Fatal(err)
	}
	http.Handle("/", router)
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

func lastStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	// TODO: read deviceid from url or account

	stat, err := store.GetCarStatus(ctx, devID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make(map[string]types.LastMsg)
	resp["staticamps"] = chargerstore.LastAmps
	resp["staticvolts"] = chargerstore.LastVolts
	resp["staticsoc"] = chargerstore.LastSOC
	resp["amps"] = stat.LastAmps
	resp["volts"] = stat.LastVolts
	resp["soc"] = stat.LastSOC
	marshal(w, resp)
}

func batteryStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	stat, err := store.GetCarStatus(ctx, devID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	marshal(w, math.BatteryCharging(
		stat.LastSOC,
		stat.LastAmps,
		stat.LastVolts,
	))
}

func rateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	// TODO: read deviceid from url or account
	devID := "520041000351353337353037"
	stat, err := store.GetCarStatus(ctx, devID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p := math.Power(stat.LastVolts.Data, stat.LastAmps.Data)
	marshal(w, struct {
		Amps  float64 `json:"amps"`
		Volts float64 `json:"volts"`
		Power float64 `json:"power"`
	}{
		Power: p,
		Amps:  chargerstore.LastAmps.Data,
		Volts: chargerstore.LastVolts.Data,
	})
}

const maxMessages = 10

var (
	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string
)
