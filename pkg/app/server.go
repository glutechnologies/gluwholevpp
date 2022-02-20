package app

import (
	"fmt"
	"gluwholevpp/pkg/api"
	"net/http"
)

func RunHttpServer(addr string) {
	// New gluwholevpp API
	s := api.New()

	fmt.Println("Listening on", addr)
	err := http.ListenAndServe(addr, s.Router())
	defer s.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
}
