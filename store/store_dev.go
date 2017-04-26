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
  "power": {
    "Data": 2.783,
    "PublishTime": "2017-04-26T16:23:55.598Z"
  },
  "soc": {
    "Data": 0.38039215,
    "PublishTime": "2017-04-26T13:59:48.227Z"
  },
  "volts": {
    "Data": 0,
    "PublishTime": "2017-04-26T16:23:55.746Z"
  }
}`)
		var cs types.CarStatus
		err := json.Unmarshal(bs, &cs)
		return &cs, err
	}

}
