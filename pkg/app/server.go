package app

import (
	"fmt"
	"gluwholevpp/pkg/api"
	"net/http"
)

type Config struct {
	SrcDatabase  string
	Port         int
	Address      string
	SrcInterface int
	VppEnabled   bool
	SrcVppSocket string
}

func RunHttpServer(config *Config) {
	// New gluwholevpp API
	s := api.New(config.VppEnabled, config.SrcInterface, config.SrcDatabase, config.SrcVppSocket)
	addr := fmt.Sprintf("%v:%v", config.Address, config.Port)

	fmt.Println("Listening on", addr)
	err := http.ListenAndServe(addr, s.Router())
	defer s.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
}
