package api

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	configs "orderly/internal/config"
	api "orderly/internal/services/api/handlers"
	"orderly/pkg/db"
)

type App struct {
	Config *configs.Config
	DB     *db.DB
	Logger *slog.Logger
}

type AppDeps struct {
	Config *configs.Config
	DB     *db.DB
	Logger *slog.Logger
}

func NewApp(deps *AppDeps) *App {
	return &App{
		Config: deps.Config,
		DB:     deps.DB,
		Logger: deps.Logger,
	}
}

func (app *App) Run() {
	router := chi.NewMux()
	apiService := NewApiService(app.Config.JWTSecret)
	api.NewHandler(router, &api.HandlerDeps{
		ApiService: apiService,
		Config:     app.Config,
		Logger:     app.Logger,
	})
	server := http.Server{
		Handler: router,
		Addr:    app.Config.ApiAddress,
	}
	app.Logger.Info("Api run!", slog.String("Address", app.Config.ApiAddress))
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
