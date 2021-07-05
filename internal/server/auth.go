package server

import (
	"goo/internal/errors"
	"goo/internal/repo/user"
	"goo/internal/types"
	pb "goo/proto"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	GrantTypePassword = "password"
	GrantTypeClients  = "client_credentials"
	GrantTypeRefresh  = "refresh_token"
)

func (u *User) AuthToken(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	result := &pb.AuthResponse{}

	options := url.Values{
		"client_id":     {in.ClientId},
		"client_secret": {in.ClientSecret},
	}

	if in.GetGrantType() == GrantTypePassword {

		am := &user.AdminMember{}
		am.Email = in.Username
		am.Password = in.Password
		userType, err := u.svc.UserRepo.VerifyAdminMember(ctx, am)
		if err != nil {
			return result, errors.GWValidationFailed("用户名密码错误")
		}

		authUserID := fmt.Sprintf("%s:%d", userType, am.ID)
		options.Add("grant_type", GrantTypePassword)
		options.Add("authenticated_userid", authUserID)
		options.Add("provision_key", viper.GetString("kong.provision_key"))
	} else if in.GetGrantType() == GrantTypeClients {
		options.Add("grant_type", GrantTypeClients)
	} else if in.GetGrantType() == GrantTypeRefresh {
		options.Add("grant_type", GrantTypeRefresh)
		options.Add("refresh_token", in.RefreshToken)
	} else {
		return result, errors.GWValidationFailed("没有此类型")
	}


	uri := viper.GetString("kong.server_url") + "/oauth2/token"
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}


	req, err := http.NewRequest("POST", uri, strings.NewReader(options.Encode()))
	if err != nil {
		return result, errors.GWInternalServer("请求失败")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Host = viper.GetString("kong.auth_host")

	resp, err := client.Do(req)
	if err != nil {
		return result, errors.GWInternalServer("请求失败")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, errors.GWInternalServer("请求返回失败")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, errors.GWInternalServer("解析失败")
	}

	if result.AccessToken == "" {
		kongErr := &types.KongError{}
		err = json.Unmarshal(body, &kongErr)
		if err != nil {
			return result, errors.GWInternalServer("解析失败")
		}

		return result, errors.GWInternalServer(kongErr.ErrorDescription)
	}

	return result, nil
}
