package main

import (
	"log"

	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/http/router"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		log.Fatal("cannot initiate app: ", err)
	}

	routes, err := router.New(cfg.Router, app)
	if err != nil {
		log.Fatal("cannot initiate routes: ", err)
	}
	routes.HandleRequests(cfg.API)

}
