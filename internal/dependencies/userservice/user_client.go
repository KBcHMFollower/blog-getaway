package userdepend

import (
	authv1 "github.com/KBcHMFollower/blog_user_service/api/protos/gen/auth"
	usersv1 "github.com/KBcHMFollower/blog_user_service/api/protos/gen/users"
	"google.golang.org/grpc"
)

type UsersGrpcClient struct {
	UsersApi usersv1.UsersServiceClient
	AuthApi  authv1.AuthClient
}

func NewUsersClient(addr string) (*UsersGrpcClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return &UsersGrpcClient{
		UsersApi: usersv1.NewUsersServiceClient(conn),
		AuthApi:  authv1.NewAuthClient(conn),
	}, nil
}
