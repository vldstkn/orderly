package main

import (
	configs "orderly/internal/config"
	"orderly/internal/services/account"
	"orderly/pkg/db"
	"orderly/pkg/logger"
	"os"
)

func main() {
	conf := configs.LoadConfig()
	logger := logger.NewLogger(os.Stdout)
	database := db.NewDb(conf.Dsn)
	_ = database
	app := account.NewApp(&account.AppDeps{
		Config: conf,
		DB:     database,
		Logger: logger,
	})
	app.Run()
}
