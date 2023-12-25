package main

import (
	"fmt"
	"log"
)

const gatewayPort = "80"

type application struct {
}

func main() {
	// setup application config
	app := application{}

	// getting router with gin engine
	router := app.routes()

	// Using gin to start api gateway server, exit status 1 if fail to start server
	if err := router.Run(fmt.Sprintf(":%s", gatewayPort)); err != nil {
		log.Fatal(err)
	}
}