package account

import (
	"context"
	configs "orderly/internal/config"
	"orderly/internal/di"
	pb "orderly/pkg/api/account"
	"orderly/pkg/jwt"
)

type Handler struct {
	pb.UnsafeAccountServer
	Service di.IAccountService
	Config  *configs.Config
}

type HandlerDeps struct {
	Service di.IAccountService
	Config  *configs.Config
}

func NewHandler(deps *HandlerDeps) *Handler {
	return &Handler{
		Service: deps.Service,
		Config:  deps.Config,
	}
}

func (handler *Handler) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	id, err := handler.Service.Register(r.Email, r.Name, r.Password)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{
		RefreshToken: "refresh",
		Id:           int64(id),
		AccessToken:  "access",
	}, nil
}

func (handler *Handler) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {

	return &pb.LoginResponse{
		RefreshToken: "refresh",
		Id:           123,
		AccessToken:  "access",
	}, nil
}

func (handler *Handler) GetProfile(ctx context.Context, request *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {

	return &pb.GetProfileResponse{
		Profile: &pb.PublicUserProfile{
			Id:    int64(1),
			Email: int64(2),
			Name:  "name",
			Role:  "role",
		},
	}, nil
}
func (handler *Handler) GetNewTokens(ctx context.Context, r *pb.GetNewTokensRequest) (*pb.GetNewTokensResponse, error) {
	accessT, refreshT, err := handler.Service.IssueTokens(handler.Config.JWTSecret, jwt.Data{
		Id: int(r.Id),
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
