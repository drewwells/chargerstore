// +build !appengine

package api

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/aetest"
)

func init() {
	configureContext = func(ctx context.Context) context.Context {
		aeCtx, _, _ := aetest.NewContext()
		return aeCtx
	}
}
