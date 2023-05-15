package app

import (
	"log"
	"net/http"

	"github.com/MakaroffAV/thesis-blockchain-node-root/internal/rts"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

func RunRootNode() {

	for _, route := range rts.GetRoutes() {
		http.HandleFunc(route.Path, route.Handler)
	}

	log.Fatal(http.ListenAndServe("0.0.0.0:2605", nil))

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
