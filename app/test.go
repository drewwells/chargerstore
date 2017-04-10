package main

import (
	"fmt"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/appengine/urlfetch"

	"golang.org/x/net/context"

	"github.com/drewwells/chargerstore/types"

	"cloud.google.com/go/datastore"
)

func test(ctx context.Context) {
	return // fuck appengine
	// Setx your Google Cloud Platform project ID.
	projectID := "particle-volt"
	log.Printf("%#v\n", ctx)
	httpCli := urlfetch.Client(ctx)
	log.Printf("% #v\n", httpCli)
	httpOpt := option.WithHTTPClient(httpCli)
	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID, httpOpt)
	if err != nil {
		log.Fatalf("datastore client creation failed: %v", err)
	}

	// Sets the kind for the new entity.
	kind := "test"
	// Sets the name/ID for the new entity.
	name := "sampletask1"
	// Creates a Key instance.
	taskKey := datastore.NameKey(kind, name, nil)

	// Creates a Task instance.
	task := types.CarStatus{
		DeviceID: "TEST",
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, taskKey, &task); err != nil {
		log.Fatalf("Failed to save task: %v", err)
	}

	fmt.Printf("Saved %v: %v\n", taskKey, task.DeviceID)
}
