package services

import (
	"context"
	authv1 "github.com/KBcHMFollower/blog_user_service/api/protos/gen/auth"
	"log/slog"
	dependencies "test-plate/internal/dependencies/userservice"
	"test-plate/internal/domain/models"
)

type AuthService struct {
	userGrpc *dependencies.UsersGrpcClient
	log      *slog.Logger
}

func NewAuthService(userGrpc *dependencies.UsersGrpcClient, log *slog.Logger) *AuthService {
	return &AuthService{
		userGrpc: userGrpc,
		log:      log,
	}
}

func (as *AuthService) Register(ctx context.Context, registerData *models.RegisterData) (*models.TokenData, error) {
	op := "AuthService.Register"
	log := as.log.With(
		slog.String("op", op))

	res, err := as.userGrpc.AuthApi.Register(ctx, &authv1.RegisterDTO{
		Email:    registerData.Email,
		Password: registerData.Password,
	})
	if err != nil {
		log.Error("Failed to register user", err)
		return nil, err
	}

	return &models.TokenData{
		Token: res.GetToken(),
	}, nil
}

func (as *AuthService) Login(ctx context.Context, loginData *models.LoginData) (*models.TokenData, error) {
	op := "AuthService.Login"
	log := as.log.With(
		slog.String("op", op))

	res, err := as.userGrpc.AuthApi.Login(ctx, &authv1.LoginDTO{
		Email:    loginData.Email,
		Password: loginData.Password,
	})
	if err != nil {
		log.Error("Failed to login user", err)
		return nil, err
	}

	return &models.TokenData{
		Token: res.GetToken(),
	}, nil
}

func (as *AuthService) CheckAuth(ctx context.Context, token *models.TokenData) (*models.TokenData, error) {
	op := "AuthService.CheckAuth"
	log := as.log.With(
		slog.String("op", op))

	res, err := as.userGrpc.AuthApi.CheckAuth(ctx, &authv1.CheckAuthDTO{
		Token: token.Token,
	})
	if err != nil {
		log.Error("Failed to checkAuth user", err)
		return nil, err
	}

	return &models.TokenData{
		Token: res.GetToken(),
	}, nil
}
