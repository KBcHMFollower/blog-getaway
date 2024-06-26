package models

import (
	usersv1 "github.com/KBcHMFollower/blog_user_service/api/protos/gen/users"
)

type User struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	FName       string `json:"fname"`
	LName       string `json:"lname"`
	Avatar      string `json:"avatar"`
	AvatarMin   string `json:"avatar_min"`
	IsDeleted   bool   `json:"is_deleted"`
	PassHash    []byte `json:"pass_hash"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}

func ConvertUserFromProto(u *usersv1.User) *User {
	return &User{
		Id:          u.GetId(),
		Email:       u.GetEmail(),
		FName:       u.GetFname(),
		LName:       u.GetLname(),
		Avatar:      u.GetAvatar(),
		AvatarMin:   u.GetAvatarMin(),
		IsDeleted:   u.GetIsDeleted(),
		CreatedDate: u.GetCreatedDate(),
		UpdatedDate: u.GetUpdatedDate(),
	}
}

func UsersArrayFromProto(usersProto []*usersv1.User) []*User {
	users := make([]*User, 0)

	for _, userProto := range usersProto {
		users = append(users, ConvertUserFromProto(userProto))
	}

	return users
}
