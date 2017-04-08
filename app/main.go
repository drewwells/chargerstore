package main

import (
	"log"

	"github.com/drewwells/chargerstore"
	"google.golang.org/appengine"
)

func main() {
	appengine.Main()

	cli, err := chargerstore.NewPS()
	if err != nil {
		log.Fatal(err)
	}

	cli.Subscribe("CAR", "carpull")
}
