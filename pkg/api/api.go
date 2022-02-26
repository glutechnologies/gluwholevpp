package api

import (
	"gluwholevpp/pkg/repository"
	"gluwholevpp/pkg/vpp"
	"net/http"

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

func New(vppEnabled bool, srcInterface int, srcDatabase string) Server {
	a := &Api{srcInterface: srcInterface}
	r := mux.NewRouter()

	// Add endpoints
	r.HandleFunc("/bitstreams", a.GetBitstreamsHandler).Methods(http.MethodGet)
	r.HandleFunc("/bitstream", a.CreateBitstreamHandler).Methods(http.MethodPost)

	// Init DB
	a.storage.Init(srcDatabase)
	a.storage.OpenDB()

	// Init VPP
	a.vpp.Init("/var/run/vpp/api.sock", vppEnabled)

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
