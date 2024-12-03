package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	configs "orderly/internal/config"
	"orderly/internal/di"
	api "orderly/internal/services/api/dto"
	pb "orderly/pkg/api/account"
	grpc_conn "orderly/pkg/grpc-conn"
	http_error "orderly/pkg/http-error"
	"orderly/pkg/req"
	"orderly/pkg/res"
)

type HandlerDeps struct {
	ApiService di.IApiService
	Config     *configs.Config
	Logger     *slog.Logger
}

type Handler struct {
	ApiService    di.IApiService
	Config        *configs.Config
	AccountClient pb.AccountClient
	Logger        *slog.Logger
}

func NewHandler(router *chi.Mux, deps *HandlerDeps) {
	accountConn, err := grpc_conn.NewClientConn(deps.Config.AccountAddress)
	if err != nil {
		panic(err)
	}

	accountClient := pb.NewAccountClient(accountConn)

	handler := Handler{
		ApiService:    deps.ApiService,
		Config:        deps.Config,
		AccountClient: accountClient,
		Logger:        deps.Logger,
	}

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.Register())
		r.Post("/login", handler.Login())
		r.Get("/login/access-token", handler.GetNewTokens())
	})
	router.Route("/account", func(r chi.Router) {
		r.Get("/profile", handler.GetProfile())
	})
}

func (handler *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[api.RegisterRequest](&w, r)
		if err != nil {
			handler.Logger.Error("req.HandleBody", slog.String("err", err.Error()))
			return
		}
		response, err := handler.AccountClient.Register(context.Background(), &pb.RegisterRequest{
			Email:    body.Email,
			Password: body.Password,
			Name:     body.Name,
		})
		if err != nil {
			http_error.BadRequest(w, handler.Logger, "AccountClient.Register", err)
			return
		}
		handler.ApiService.AddCookie(&w, "refresh_token", response.RefreshToken, 3600)
		res.Json(w, api.RegisterResponse{
			Id:          int(response.Id),
			AccessToken: response.AccessToken,
		}, http.StatusCreated)
	}
}

func (handler *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handler *Handler) GetNewTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handler *Handler) GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
