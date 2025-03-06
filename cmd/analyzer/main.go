package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/hiromaily/hatena-analyzer/pkg/args"
	"github.com/hiromaily/hatena-analyzer/pkg/envs"
	"github.com/hiromaily/hatena-analyzer/pkg/registry"
)

// value is passed when building application
var CommitID string

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
		fmt.Println(err)
		os.Exit(1)
	}

	// Register for initialization of dependencies
	reg, err := registry.NewRegistry(&cfg, appCode, CommitID, args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app, err := reg.InitializeApp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Execute application
	err = app.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
