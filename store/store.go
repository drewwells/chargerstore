package store

import (
	"context"
	"fmt"

	"github.com/drewwells/chargerstore/types"
	"google.golang.org/appengine/datastore"
)

const (
	carbucket    = "car"
	statusbucket = "status"
)

func PutCarMsg(ctx context.Context, cm *types.CarMsg) error {
	k := datastore.NewKey(ctx, carbucket, cm.ID, 0, nil)
	if _, err := datastore.Put(ctx, k, cm); err != nil {
		return fmt.Errorf("failed to save %s: %s", cm.ID, err)
	}
	return nil
}

// Walk through db looking for last good status
func PutCarStatus(ctx context.Context, status *types.CarStatus) error {
	k := datastore.NewIncompleteKey(ctx, statusbucket, nil)
	_, err := datastore.Put(ctx, k, status)
	return err
}

func GetCarStatus(ctx context.Context, deviceID string) (*types.CarMsg, error) {
	q := datastore.NewQuery(statusbucket).Filter("DeviceID =", deviceID).Order("-CreatedAt").Limit(1)

	it := q.Run(ctx)
	var cm types.CarMsg
	_, err := it.Next(&cm)
	return &cm, err
}
