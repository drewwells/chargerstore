package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

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
	bs, _ := ioutil.ReadAll(r.Body)
	log.Infof(ctx, "incoming", string(bs))
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
