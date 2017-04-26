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
    "Data": 0,
    "PublishTime": "2017-04-26T17:06:34.85Z"
  },
  "power": {
    "Data": 0,
    "PublishTime": "2017-04-26T17:06:34.85Z"
  },
  "soc": {
    "Data": 0.84313728,
    "PublishTime": "2017-04-26T16:57:40.95Z"
  },
  "volts": {
    "Data": 0,
    "PublishTime": "2017-04-26T17:06:34.85Z"
  }
}`)
		var cs types.CarStatus
		err := json.Unmarshal(bs, &cs)
		return &cs, err
	}

}
