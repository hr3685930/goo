package errors

import (
	"goo/pkg/grpcgw"
	"google.golang.org/grpc/codes"
	"runtime/debug"
)

func GWRequestTimeout(msg string) error {
	return grpcgw.Err(codes.Canceled, msg, string(debug.Stack()))
}

func GWInternalServer(msg string) error {
	return grpcgw.Err(codes.Unknown, msg, string(debug.Stack()))
}

func GWValidationFailed(msg string) error {
	return grpcgw.Err(codes.InvalidArgument, msg, string(debug.Stack()))
}

func GWBadRequest(msg string) error {
	return grpcgw.Err(codes.OutOfRange, msg, string(debug.Stack()))
}

func GWResourceNotFound(msg string) error {
	return grpcgw.Err(codes.NotFound, msg, string(debug.Stack()))
}

func GWConflict(msg string) error {
	return grpcgw.Err(codes.AlreadyExists, msg, string(debug.Stack()))
}

func GWAuthorizationFailed(msg string) error {
	return grpcgw.Err(codes.PermissionDenied, msg, string(debug.Stack()))
}

func GWAuthenticationFailed(msg string) error {
	return grpcgw.Err(codes.Unauthenticated, msg, string(debug.Stack()))
}

func GWTooManyRequests(msg string) error {
	return grpcgw.Err(codes.ResourceExhausted, msg, string(debug.Stack()))
}

func GWServiceUnavailable(msg string) error {
	return grpcgw.Err(codes.Unavailable, msg, string(debug.Stack()))
}