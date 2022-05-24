//  GluWholeVPP API
//
//  Schemes: http
//  BasePath: /
//  Version: 1.0.0
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
//
// swagger:meta
package api

import (
	"gluwholevpp/pkg/repository"
	"gluwholevpp/pkg/vpp"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

type Api struct {
	router       http.Handler
	storage      repository.Storage
	vpp          vpp.Client
	srcInterface int
}

type Server interface {
	Router() http.Handler
	Close() error
}

func New(vppEnabled bool, srcInterface int, srcDatabase string, srcVppSocket string) Server {
	a := &Api{srcInterface: srcInterface}
	r := mux.NewRouter()

	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// Documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	// Add endpoints
	r.HandleFunc("/bitstreams", a.GetBitstreamsHandler).Methods(http.MethodGet)
	r.HandleFunc("/bitstream", a.CreateBitstreamHandler).Methods(http.MethodPost)
	r.HandleFunc("/bitstream/{id}", a.GetBitstreamHandler).Methods(http.MethodGet)
	r.HandleFunc("/bitstream/{id}", a.DeleteBitstreamHandler).Methods(http.MethodDelete)

	r.HandleFunc("/customers", a.GetCustomersHandler).Methods(http.MethodGet)
	r.HandleFunc("/customer", a.CreateCustomerHandler).Methods(http.MethodPost)
	r.HandleFunc("/customer/{id}", a.GetCustomerHandler).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id}", a.DeleteCustomerHandler).Methods(http.MethodDelete)
	r.HandleFunc("/customer/{id}", a.PatchCustomerHandler).Methods(http.MethodPatch)
	r.HandleFunc("/customer/{id}/bitstreams", a.GetBitstreamsFromCustomerHandler).Methods(http.MethodGet)

	// Init DB
	a.storage.Init(srcDatabase)
	a.storage.OpenDB()

	// Init VPP
	a.vpp.Init(srcVppSocket, vppEnabled)

	// Load Bitstreams from DB
	a.LoadBitstreamsStorage()

	a.router = r
	return a
}

func (a *Api) Close() error {
	a.vpp.Close()
	return a.storage.CloseDB()
}

func (a *Api) Router() http.Handler {
	return a.router
}
