package svc

import (
	"goo/internal/repo"
	"goo/internal/repo/user"
	"goo/pkg/db"
)

type ServiceContext struct {
	UserRepo user.UserRepo
	Helper *repo.Helper
}

func NewServiceContext() *ServiceContext {
	return &ServiceContext{
		UserRepo: user.NewUserDB(db.Orm),
		Helper: repo.NewHelper(),
	}
}
