package main

import (
	"context"
	_ "embed"
	"log"

	handyconfig "github.com/FiyZou/handygo/config"
	"github.com/FiyZou/handygo/examples/internal/bootstrap"
	exampleconfig "github.com/FiyZou/handygo/examples/internal/config"
)

//go:embed manifest/config.yaml
var configYAML []byte

func main() {
	var cfg exampleconfig.AppConfig
	if err := handyconfig.LoadData("config", "yaml", configYAML, &cfg); err != nil {
		log.Fatalf("load embedded config: %v", err)
	}

	application, err := bootstrap.New(cfg)
	if err != nil {
		log.Fatalf("bootstrap application: %v", err)
	}
	if err := application.Run(context.Background()); err != nil {
		log.Fatalf("run application: %v", err)
	}
}
