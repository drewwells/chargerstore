package store

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/drewwells/chargerstore/types"
	aedatastore "google.golang.org/appengine/datastore"
)

const (
	projectid    = "particle-volt"
	carbucket    = "car"
	statusbucket = "status"
)

func PutCarMsg(ctx context.Context, cm *types.CarMsg) error {
	k := aedatastore.NewKey(ctx, carbucket, cm.ID, 0, nil)
	if _, err := aedatastore.Put(ctx, k, cm); err != nil {
		return fmt.Errorf("failed to save %s: %s", cm.ID, err)
	}
	return nil
}

// Walk through db looking for last good status
func PutCarStatus(ctx context.Context, status *types.CarStatus) error {
	k := aedatastore.NewIncompleteKey(ctx, statusbucket, nil)
	_, err := aedatastore.Put(ctx, k, status)
	return err
}

func getLastField(ctx context.Context, qry *aedatastore.Query, field string) (*types.CarStatus, error) {
	q := qry.Order("-CreatedAt").
		Filter(field+" >", 0).
		Order("-CreatedAt").
		Limit(1)

	it := q.Run(ctx)
	var cm types.CarStatus
	_, err := it.Next(&cm)
	return &cm, err
}

func GetCarStatus(ctx context.Context, deviceID string) (*types.CarStatus, error) {

	qry := aedatastore.NewQuery(statusbucket).Filter("DeviceID =", deviceID)
	bat, err := getLastField(ctx, qry, "Battery")
	if err != nil {
		return nil, err
	}
	volts, err := getLastField(ctx, qry, "ChargerVolts")
	if err != nil {
		return nil, err
	}
	amps, err := getLastField(ctx, qry, "ChargerAmps")
	if err != nil {
		return nil, err
	}

	return &types.CarStatus{
		DeviceID:  deviceID,
		LastSOC:   bat.LastSOC,
		LastAmps:  amps.LastAmps,
		LastVolts: volts.LastVolts,
		CreatedAt: bat.CreatedAt,
	}, nil
}
