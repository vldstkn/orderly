package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	configs "orderly/internal/config"
	"orderly/internal/di"
	api "orderly/internal/services/api/dto"
	"orderly/internal/services/api/middleware"
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
		r.Group(func(r chi.Router) {
			r.Use(middleware.IsAuthed(handler.Config.JWTSecret))
			r.Get("/login/access-token", handler.GetNewTokens())
		})
	})
	router.Route("/account", func(r chi.Router) {
		r.Use(middleware.IsAuthed(handler.Config.JWTSecret))
		r.Get("/profile", handler.GetProfile())
		r.Put("/role", handler.ChangeRole())
	})
}

func (handler *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[api.RegisterRequest](&w, r)
		if err != nil {
			handler.Logger.Error("req.HandleBody", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
		body, err := req.HandleBody[api.LoginRequest](&w, r)
		if err != nil {
			handler.Logger.Error("req.HandleBody", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		response, err := handler.AccountClient.Login(context.Background(), &pb.LoginRequest{
			Email:    body.Email,
			Password: body.Password,
		})
		if err != nil {
			http_error.BadRequest(w, handler.Logger, "AccountClient.Login", err)
			return
		}
		handler.ApiService.AddCookie(&w, "refresh_token", response.RefreshToken, 3600)
		res.Json(w, api.LoginResponse{
			Id:          int(response.Id),
			AccessToken: response.AccessToken,
		}, http.StatusCreated)
	}
}

func (handler *Handler) GetNewTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value("authData").(middleware.AuthData)

		if data.Id <= 0 {
			handler.Logger.Error("r.Context()", slog.String("err", "id is missing"))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		accRes, err := handler.AccountClient.GetNewTokens(context.Background(), &pb.GetNewTokensRequest{
			Id:   int64(data.Id),
			Role: data.Role,
		})
		if err != nil {
			handler.Logger.Error("AccountClient.GetNewTokens", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		handler.ApiService.AddCookie(&w, "refresh_token", accRes.RefreshToken, 3600)
		res.Json(w, api.GetNewTokensResponse{
			AccessToken: accRes.AccessToken,
		}, http.StatusOK)
	}
}

func (handler *Handler) GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value("authData").(middleware.AuthData)

		if data.Id <= 0 {
			handler.Logger.Error("r.Context()", slog.String("err", "id is missing"))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		accRes, err := handler.AccountClient.GetProfileById(context.Background(), &pb.GetProfileRequest{
			Id: int64(data.Id),
		})
		if err != nil {
			handler.Logger.Error("AccountClient.GetProfile", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		res.Json(w, api.GetProfileResponse{
			Id:        int(accRes.Profile.Id),
			Name:      accRes.Profile.Name,
			Role:      accRes.Profile.Role,
			Email:     accRes.Profile.Email,
			CreatedAt: accRes.Profile.CreatedAt,
			UpdatedAt: accRes.Profile.UpdatedAt,
		}, http.StatusOK)
	}
}

func (handler *Handler) ChangeRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[api.ChangeRoleRequest](&w, r)
		data := r.Context().Value("authData").(middleware.AuthData)
		if err != nil || data.Id <= 0 {
			handler.Logger.Error("req.HandleBody", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		accRes, err := handler.AccountClient.ChangeRole(context.Background(), &pb.ChangeRoleRequest{
			Role: body.Role,
			Id:   int64(data.Id),
		})
		// TODO:
		if err != nil {
			fmt.Println(err)
			handler.Logger.Error("AccountClient.ChangeRole", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		handler.ApiService.AddCookie(&w, "refresh_token", accRes.RefreshToken, 3600)

		res.Json(w, api.ChangeRoleResponse{
			Id:          int(accRes.Id),
			Role:        accRes.Role,
			AccessToken: accRes.AccessToken,
		}, http.StatusOK)
	}
}
