package main

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/hiromaily/hatena-fake-detector/pkg/args"
	"github.com/hiromaily/hatena-fake-detector/pkg/envs"
	"github.com/hiromaily/hatena-fake-detector/pkg/registry"
)

// value is passed when building application
var CommitID string

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Parse arguments
	args, _, appCode := args.Parse()
	if args.Version {
		fmt.Println(CommitID)
		return
	}

	// Parse Environment Variables
	var cfg envs.Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	// Register for initialization of dependencies
	reg := registry.NewRegistry(&cfg, appCode, CommitID, args.URLs)
	app, err := reg.InitializeApp()
	if err != nil {
		panic(err)
	}

	// Execute application
	err = app.Run()
	if err != nil {
		reg.Logger().Error("failed to run application", "error", err)
	}
}
