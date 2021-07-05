package server

import (
	"goo/internal/errors"
	"goo/internal/repo/user"
	"goo/internal/svc"
	pb "goo/proto"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
)

type User struct {
	pb.UnimplementedUserServer

	svc *svc.ServiceContext
}

func NewUser(svc *svc.ServiceContext) *User {
	return &User{svc: svc}
}

func (u *User) GetCurrentUserInfo(ctx context.Context, in *pb.Empty) (*pb.CurrentUserResponse, error) {
	result := &pb.CurrentUserResponse{}
	md, _ := metadata.FromIncomingContext(ctx)
	userInfo,err := u.svc.Helper.ValidaUser(md.Get("x-authenticated-userid"))
	if err != nil {
		return result, errors.GWAuthorizationFailed("用户未找到")
	}

	adminMember := &user.AdminMember{}
	adminMember.ID = userInfo.ID
	err = u.svc.UserRepo.GetCurrentUserInfo(ctx, adminMember)
	if err != nil {
		return result,errors.GWResourceNotFound("用户未找到")
	}
	result.Email = adminMember.Email
	result.Username = adminMember.Username
	result.Password = adminMember.Password
	return result, nil
}

func (u *User) UpdateCurrentUserInfo(ctx context.Context, in *pb.CurrentUserRequest) (*pb.CurrentUserResponse, error) {
	fmt.Print(in.String())

	result := &pb.CurrentUserResponse{}

	return result, nil
}


