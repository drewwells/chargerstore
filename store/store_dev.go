// +build !appengine

package store

import (
	"encoding/json"

	"golang.org/x/net/context"

	"github.com/drewwells/chargerstore/types"
)

func init() {
	GetCarStatus = func(ctx context.Context, deviceID string) (*types.CarStatus, error) {

		bs := []byte(`{
  "amps": {
    "Data": 8,
    "PublishTime": "2017-04-11T03:16:36.421Z"
  },
  "soc": {
    "Data": 0.63,
    "PublishTime": "2017-04-11T04:02:14.364Z"
  },
  "volts": {
    "Data": 122,
    "PublishTime": "2017-04-11T04:01:43.03Z"
  }
}`)
		var cs types.CarStatus
		err := json.Unmarshal(bs, &cs)
		return &cs, err
	}

}
