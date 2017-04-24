package store

import (
	"fmt"
	"time"

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
	q := qry.Order("-CreatedAt")

	// Battery lookup is special, just keep looking backwards in time for it
	if field == "Battery" {
		q = q.Limit(10000) // 2 days
	} else {
		q = q.Limit(20) // 4x records per minute
	}

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
			// 0 is not a valid battery value, can bus loves to send it anyways
			if cm.LastSOC.Data > 0 {
				break endloop
			}
		case "ChargerVolts":
			if cm.LastVolts.Data > -1 {
				break endloop
			}
		case "ChargerAmps":
			if cm.LastAmps.Data > -1 {
				break endloop
			}
		case "ChargerPower":
			if cm.LastPower.Data > 0 {
				break endloop
			}
		}
	}
	return &cm, nil
}

// GetCarStatus fetches the current deviceID's vehicle status. It
// will return the cache unless it is found to be missing
var GetCarStatus = func(ctx context.Context, deviceID string) (*types.CarStatus, error) {

	qry := aedatastore.NewQuery(statusbucket).Filter("DeviceID =", deviceID)

	if LastSOC.PublishTime.IsZero() {
		bat, err := getLastField(ctx, qry, "Battery")
		if err != nil {
			log.Errorf(ctx, "%s", err)
		} else {
			LastSOC = bat.LastSOC
		}
	}

	if LastVolts.PublishTime.IsZero() {
		volts, err := getLastField(ctx, qry, "ChargerVolts")
		if err != nil {
			log.Errorf(ctx, "%s", err)
		} else {
			LastVolts = volts.LastVolts
		}
	}

	if LastAmps.PublishTime.IsZero() {
		amps, err := getLastField(ctx, qry, "ChargerAmps")
		if err != nil {
			log.Errorf(ctx, "%s", err)
		} else {
			LastAmps = amps.LastAmps
		}
	}

	if LastPower.PublishTime.IsZero() {
		power, err := getLastField(ctx, qry, "ChargerPower")
		if err != nil {
			log.Errorf(ctx, "%s", err)
		} else {
			LastPower = power.LastPower
		}
	}

	return &types.CarStatus{
		DeviceID:  deviceID,
		LastSOC:   LastSOC,
		LastAmps:  LastAmps,
		LastVolts: LastVolts,
		LastPower: LastPower,
		CreatedAt: time.Now(),
	}, nil
}
