package rts

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import (
	"net/http"

	"github.com/MakaroffAV/thesis-blockchain-node-root/internal/hds"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

var (

	// routes is array of web
	// server routes and handlers
	routes = []Route{
		{
			Path:    "/track",
			Handler: hds.BaseNodeTrack,
		},
		{
			Path:    "/sh-ip",
			Handler: hds.BaseNodeGetIp,
		},
	}
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

type (

	// Route is data structure for
	// describing the web server route and handler
	Route struct {
		Path    string
		Handler func(http.ResponseWriter, *http.Request)
	}
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// GetRoutes is global
// access point for fetching web server routes
func GetRoutes() []Route { return routes }

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
