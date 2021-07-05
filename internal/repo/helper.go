package repo

import (
	"errors"
	"github.com/spf13/cast"
	"strings"
)

type Helper struct {
}

func NewHelper() *Helper {
	return &Helper{}
}

type AuthUserInfo struct {
	ID   uint
	Type string
}

func (*Helper) ValidaUser(headers []string) (info *AuthUserInfo, err error) {
	info = &AuthUserInfo{}
	sps := strings.Split(headers[0], ":")
	if len(sps) != 2 {
		return info, errors.New("error")
	}
	info.ID = cast.ToUint(sps[1])
	info.Type = sps[0]

	return info,nil
}
