// +build !appengine

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func init() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("failed to listen", err)
	}
	fmt.Println("Serving", lis.Addr())
	http.Serve(lis, nil)
}
