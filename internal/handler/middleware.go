package handler

import (
	"goo/internal/errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"google.golang.org/grpc"
)

func Middleware() []grpc.UnaryServerInterceptor {
	var grpcOpts []grpc.UnaryServerInterceptor
	RecoveryInterceptor := func() grpc_recovery.Option {
		return grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			return errors.GWInternalServer(fmt.Sprintf("panic triggered: %v", p))
		})
	}

	grpcOpts = append(grpcOpts,
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_recovery.UnaryServerInterceptor(RecoveryInterceptor()),
		grpc_validator.UnaryServerInterceptor(),
	)

	return grpcOpts
}
