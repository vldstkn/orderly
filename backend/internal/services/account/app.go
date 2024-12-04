package account

import (
	"google.golang.org/grpc"
	"log/slog"
	"net"
	configs "orderly/internal/config"
	pb "orderly/pkg/api/account"
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
	var opts []grpc.ServerOption
	lis, err := net.Listen("tcp", app.Config.AccountAddress)
	if err != nil {
		panic(err)
	}
	handler := NewHandler(&HandlerDeps{
		Config: app.Config,
	})
	server := grpc.NewServer(opts...)
	pb.RegisterAccountServer(server, handler)

	app.Logger.Info("Account starts",
		slog.String("Address", app.Config.AccountAddress),
	)
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
