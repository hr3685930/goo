package grpcgw

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"net/http"
	"net/textproto"
)

type GW struct {
}

func (*GW) Register(rpcAddr, httpAddr string, rpcHandler func(serverOpt []grpc.UnaryServerInterceptor) *grpc.Server,
	httpHandler func(gwmux *runtime.ServeMux, endpoint string) (http.Handler, error)) error {
	lis, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		fmt.Println("failed to listen: %v", err)
	}
	serverInterceptor := []grpc.UnaryServerInterceptor{}
	serverInterceptor = append(serverInterceptor, ErrorHandlerUnaryServerInterceptor)
	s := rpcHandler(serverInterceptor)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	flag.Parse()
	echoEndpoint := flag.String("rpc_server", rpcAddr, "RPC Service")

	jsonPb := &runtime.JSONPb{}
	jsonPb.UseProtoNames = true
	jsonPb.EmitUnpopulated = true

	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonPb),
		runtime.WithIncomingHeaderMatcher(HeaderMatcher),
	)

	mux, err := httpHandler(gwmux, *echoEndpoint)
	if err != nil {
		return errors.New("gw服务开启失败")
	}

	return http.ListenAndServe(httpAddr, mux)

}

func ErrorHandlerUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		md, _ := metadata.FromIncomingContext(ctx)
		err = ErrorHandler(md, req, err)
	}
	return resp, err
}

func HeaderMatcher(key string) (string, bool) {
	k, ok := runtime.DefaultHeaderMatcher(key)
	if ok {
		return k, ok
	}
	key = textproto.CanonicalMIMEHeaderKey(key)
	return key, true
}
