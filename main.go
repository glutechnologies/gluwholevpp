package main

import (
	"flag"
	"gluwholevpp/pkg/app"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

func main() {
	// Load configuration
	configSrc := flag.String("config", "/etc/gluwholevpp.toml", "Config source path")
	flag.Parse()

	body, err := ioutil.ReadFile(*configSrc)

	if err != nil {
		log.Fatalf("Error loading configuration file")
	}

	var config app.Config
	_, err = toml.Decode(string(body), &config)

	if err != nil {
		log.Fatalf("Error decoding configuration file")
	}

	app.RunHttpServer(&config)
}
