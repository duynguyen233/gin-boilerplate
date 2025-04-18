package main

import (
	"log"
	"wine-be/config"
	"wine-be/internal/app"

	_ "ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
