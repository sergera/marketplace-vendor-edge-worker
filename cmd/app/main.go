package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sergera/marketplace-vendor-edge-worker/internal/api"
	"github.com/sergera/marketplace-vendor-edge-worker/internal/conf"
	"github.com/sergera/marketplace-vendor-edge-worker/internal/evt"
)

func main() {
	conf := conf.GetConf()

	mux := http.NewServeMux()

	orderAPI := api.NewOrderAPI()

	mux.HandleFunc("/update-order-status", orderAPI.UpdateOrderStatus)

	srv := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: mux,
	}

	listener := evt.NewOrderListener()
	go listener.Listen()

	fmt.Printf("starting application on port %s", conf.Port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
