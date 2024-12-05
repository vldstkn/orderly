package account

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	configs "orderly/internal/config"
	"orderly/internal/di"
	"orderly/internal/domain"
	pb "orderly/pkg/api/account"
	"orderly/pkg/jwt"
)

type Handler struct {
	pb.UnsafeAccountServer
	Service di.IAccountService
	Config  *configs.Config
	Logger  *slog.Logger
}

type HandlerDeps struct {
	Service di.IAccountService
	Config  *configs.Config
	Logger  *slog.Logger
}

func NewHandler(deps *HandlerDeps) *Handler {
	return &Handler{
		Service: deps.Service,
		Config:  deps.Config,
		Logger:  deps.Logger,
	}
}

func (handler *Handler) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	id, err := handler.Service.Register(r.Email, r.Password, r.Name)
	if err != nil {
		handler.Logger.Error(err.Error(), slog.String("email", r.Email), slog.String("name", r.Name))
		return nil, err
	}
	accessT, refreshT, err := handler.Service.IssueTokens(handler.Config.JWTSecret, jwt.Data{
		Id:   id,
		Role: string(domain.Customer),
	})
	if err != nil {
		handler.Logger.Error(err.Error())
		return nil, err
	}
	return &pb.RegisterResponse{
		RefreshToken: refreshT,
		Id:           int64(id),
		AccessToken:  accessT,
	}, nil
}

func (handler *Handler) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	id, role, err := handler.Service.Login(r.Email, r.Password)
	if err != nil {
		handler.Logger.Error(err.Error(), slog.String("email", r.Email))
		return nil, err
	}
	accessT, refreshT, err := handler.Service.IssueTokens(handler.Config.JWTSecret, jwt.Data{
		Id:   id,
		Role: role,
	})
	if err != nil {
		handler.Logger.Error(err.Error())
		return nil, err
	}
	return &pb.LoginResponse{
		RefreshToken: refreshT,
		Id:           int64(id),
		AccessToken:  accessT,
	}, nil
}

func (handler *Handler) GetProfileById(ctx context.Context, r *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	profile, err := handler.Service.GetProfileById(int(r.Id))
	if err != nil {
		return nil, err
	}
	return &pb.GetProfileResponse{
		Profile: &pb.PublicUserProfile{
			Id:        int64(profile.Id),
			Email:     profile.Email,
			Name:      profile.Name,
			Role:      profile.Role,
			CreatedAt: profile.CreatedAt,
			UpdatedAt: profile.UpdatedAt,
		},
	}, nil
}
func (handler *Handler) GetNewTokens(ctx context.Context, r *pb.GetNewTokensRequest) (*pb.GetNewTokensResponse, error) {
	accessT, refreshT, err := handler.Service.IssueTokens(handler.Config.JWTSecret, jwt.Data{
		Id:   int(r.Id),
		Role: r.Role,
	})
	if err != nil {
		return nil, err
	}
	return &pb.GetNewTokensResponse{
		AccessToken:  accessT,
		RefreshToken: refreshT,
	}, nil
}
func (handler *Handler) UpdateById(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	return &pb.UpdateUserResponse{
		IsSuccess: false,
	}, nil
}

func (handler *Handler) ChangeRole(ctx context.Context, r *pb.ChangeRoleRequest) (*pb.ChangeRoleResponse, error) {
	err := handler.Service.ChangeRole(int(r.Id), r.Role)
	if err != nil {
		handler.Logger.Error("Service.ChangeRole", slog.String("err", err.Error()))
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	accessT, refreshT, err := handler.Service.IssueTokens(handler.Config.JWTSecret, jwt.Data{
		Id:   int(r.Id),
		Role: r.Role,
	})

	if err != nil {
		handler.Logger.Error("Service.IssueTokens", slog.String("err", err.Error()))
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return &pb.ChangeRoleResponse{
		Role:         r.Role,
		Id:           r.Id,
		AccessToken:  accessT,
		RefreshToken: refreshT,
	}, nil
}
