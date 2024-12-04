package main

import (
	configs "orderly/internal/config"
	"orderly/internal/services/api"
	"orderly/pkg/logger"
	"os"
)

func main() {
	conf := configs.LoadConfig()
	logger := logger.NewLogger(os.Stdout)
	app := api.NewApp(&api.AppDeps{
		Config: conf,
		Logger: logger,
	})

	app.Run()
}
