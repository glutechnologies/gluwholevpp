package api

import (
	"gluwholevpp/pkg/repository"
	"gluwholevpp/pkg/vpp"
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	router  http.Handler
	storage repository.Storage
	vpp     vpp.Client
}

type Server interface {
	Router() http.Handler
	Close() error
}

func New() Server {
	a := &Api{}
	r := mux.NewRouter()

	// Add endpoints
	r.HandleFunc("/bitstreams", a.GetBitstreamsHandler).Methods(http.MethodGet)
	r.HandleFunc("/bitstream", a.CreateBitstreamHandler).Methods(http.MethodPost)

	// Init DB
	a.storage.Init("data.db")
	a.storage.OpenDB()

	// Init VPP
	a.vpp.Init("/var/run/vpp/api.sock")

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
