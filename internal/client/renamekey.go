package client

import (
    "goo/internal/svc"
    "github.com/urfave/cli"
)


type User struct {
    svc *svc.ServiceContext
}

func NewUser(svc *svc.ServiceContext) *User {
    return &User{svc: svc}
}

func(*User) RenameKey(clis *cli.Context) {

}
