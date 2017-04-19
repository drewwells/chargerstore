// +build !appengine

package api

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
)

func init() {
	configureContext = func(ctx context.Context) context.Context {
		aeCtx, _, _ := aetest.NewContext()
		return aeCtx
	}

	newContext = func(r *http.Request) context.Context {
		aeCtx, _, _ := aetest.NewContext()
		return aeCtx
	}

}
