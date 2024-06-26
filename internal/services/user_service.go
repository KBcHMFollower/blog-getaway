package services

import (
	"context"
	usersv1 "github.com/KBcHMFollower/blog_user_service/api/protos/gen/users"
	"log/slog"
	dependencies "test-plate/internal/dependencies/userservice"
)

type UserService struct {
	userGrpc *dependencies.UsersGrpcClient
	log      *slog.Logger
}

func NewUserService(userGrpc *dependencies.UsersGrpcClient, log *slog.Logger) *UserService {
	return &UserService{
		userGrpc: userGrpc,
		log:      log,
	}
}

func (us *UserService) GetUser(ctx context.Context, userId string) (*usersv1.User, error) {
	op := "AuthService.Register"
	log := us.log.With(
		slog.String("op", op))

	res, err := us.userGrpc.UsersApi.GetUser(ctx, &usersv1.GetUserDTO{
		Id: userId,
	})
	if err != nil {
		log.Error("Failed to get user", "error", err)
		return nil, err
	}

	return res.User, nil
}

func (us *UserService) Subscribe(ctx context.Context, subscribeData *usersv1.SubscribeDTO) (bool, error) {
	op := "AuthService.Register"
	log := us.log.With(
		slog.String("op", op))

	_, err := us.userGrpc.UsersApi.Subscribe(ctx, subscribeData)
	if err != nil {
		log.Error("Failed to subscribe", "error", err)
		return false, err
	}

	return true, nil
}

func (us *UserService) Unsubscribe(ctx context.Context, unsubscribeData *usersv1.SubscribeDTO) (bool, error) {
	op := "AuthService.Register"
	log := us.log.With(
		slog.String("op", op))

	_, err := us.userGrpc.UsersApi.Unsubscribe(ctx, unsubscribeData)
	if err != nil {
		log.Error("Failed to unsubscribe", "error", err)
		return false, err
	}

	return true, nil
}

func (us *UserService) GetSubscribers(ctx context.Context, getData *usersv1.GetSubscribersDTO) (*usersv1.GetSubscribersRDO, error) {
	op := "AuthService.Register"
	log := us.log.With(
		slog.String("op", op))

	res, err := us.userGrpc.UsersApi.GetSubscribers(ctx, getData)
	if err != nil {
		log.Error("Failed to get subscribers", "error", err)
		return nil, err
	}

	return res, nil
}

func (us *UserService) GetSubscriptions(ctx context.Context, getData *usersv1.GetSubscriptionsDTO) (*usersv1.GetSubscriptionsRDO, error) {
	op := "AuthService.Register"
	log := us.log.With(
		slog.String("op", op))

	res, err := us.userGrpc.UsersApi.GetSubscriptions(ctx, getData)
	if err != nil {
		log.Error("Failed to get subscribers", "error", err)
		return nil, err
	}

	return res, nil

}

func (us *UserService) UpdateUser(ctx context.Context, updateData *usersv1.UpdateUserDTO) (*usersv1.UpdateUserRDO, error) {
	op := "AuthService.Register"
	log := us.log.With(
		slog.String("op", op))

	res, err := us.userGrpc.UsersApi.UpdateUser(ctx, updateData)
	if err != nil {
		log.Error("Failed to update user", "error", err)
		return nil, err
	}

	return res, nil
}

func (us *UserService) DeleteUser(ctx context.Context, userId string) (bool, error) {
	op := "AuthService.Register"
	log := us.log.With(
		slog.String("op", op))

	_, err := us.userGrpc.UsersApi.DeleteUser(ctx, &usersv1.DeleteUserDTO{
		Id: userId,
	})
	if err != nil {
		log.Error("Failed to delete user", "error", err)
		return false, err
	}

	return true, nil
}

func (us *UserService) UploadAvatar(ctx context.Context, uploadData *usersv1.UploadAvatarDTO) (*usersv1.UploadAvatarRDO, error) {
	op := "AuthService.Register"
	log := us.log.With(
		slog.String("op", op))

	res, err := us.userGrpc.UsersApi.UploadAvatar(ctx, uploadData)
	if err != nil {
		log.Error("Failed to upload avatar", "error", err)
		return nil, err
	}

	return res, nil
}

func (us *UserService) GetAvatar(ctx context.Context, imgName string) (*usersv1.GetAvatarRDO, error) {
	op := "AuthService.GetAvatar"
	log := us.log.With(
		slog.String("op", op))

	res, err := us.userGrpc.UsersApi.GetAvatar(ctx, &usersv1.GetAvatarDTO{
		FileName: imgName,
	})
	if err != nil {
		log.Error("Failed to load avatar", "error", err)
		return nil, err
	}

	return res, nil
}
