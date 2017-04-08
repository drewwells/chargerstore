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

	opts, err := chargerstore.NewPS()
	if err != nil {
		log.Fatal(err)
	}

	opts.Subscribe("CAR", "carpull")
	serve(opts)
	//http.HandleFunc("/pubsub/push", pushHandler)
}

func serve(opts *chargerstore.Options) {
	// [START http]
	// Publish a count of processed requests to the server homepage.
	http.HandleFunc("/summary", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This worker has processed %d events.", opts.Count())
	})
}

type pushRequest struct {
	Message struct {
		Attributes map[string]string
		Data       []byte
		ID         string `json:"message_id"`
	}
	Subscription string
}

const maxMessages = 10

var (

	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string
)

func pushHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("path", r.URL)
	var msg pushRequest
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	messagesMu.Lock()
	defer messagesMu.Unlock()
	// Limit to ten.
	messages = append(messages, string(msg.Message.Data))
	if len(messages) > maxMessages {
		messages = messages[len(messages)-maxMessages:]
	}
}
