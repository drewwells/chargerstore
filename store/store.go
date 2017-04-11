package store

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/drewwells/chargerstore/types"
	aedatastore "google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
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
		//Filter(field+" >", 0).
		Order("-CreatedAt") //.
		//Limit(1)

	var cm types.CarStatus
	it := q.Run(ctx)
endloop:
	for {
		_, err := it.Next(&cm)
		if err == aedatastore.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "fetching next Person: %v", err)
			break
		}

		switch field {
		case "Battery":
			if cm.LastSOC.Data > 0 {
				break endloop
			}
		case "ChargerVolts":
			if cm.LastVolts.Data > 0 {
				break endloop
			}
		case "ChargerAmps":
			if cm.LastAmps.Data > 0 {
				break endloop
			}
		}
	}
	return &cm, nil
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
