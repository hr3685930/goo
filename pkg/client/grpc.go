package client

import (
	"google.golang.org/grpc"
)

func NewCli(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}