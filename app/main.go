package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"google.golang.org/appengine"

	"github.com/drewwells/chargerstore"
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
	serve()
	http.HandleFunc("/pubsub/push", pushHandler)
}

func serve() {
	// [START http]
	// Publish a count of processed requests to the server homepage.
	http.HandleFunc("/summary", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This worker has processed %d events.", chargerstore.Count())
	})
}

const maxMessages = 10

var (
	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string
)

func pushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var req chargerstore.PushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	msg, err := chargerstore.Process(ctx, req.Message)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to process msg: %s", err), http.StatusInternalServerError)
		return
	}
	log.Println("processed message from", msg.DeviceID)
}
