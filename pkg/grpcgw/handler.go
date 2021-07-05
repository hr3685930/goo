package grpcgw

import (
	"github.com/ddliu/go-httpclient"
	"github.com/spf13/viper"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

func HTTPFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusUnprocessableEntity
	case codes.DeadlineExceeded:
		return http.StatusInternalServerError
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	}

	return http.StatusInternalServerError
}

var dontReport = []int{
	http.StatusUnauthorized,
	http.StatusForbidden,
	http.StatusMethodNotAllowed,
	http.StatusUnsupportedMediaType,
	http.StatusUnprocessableEntity,
}

func ErrorHandler(md metadata.MD, req interface{}, err error) error {
	s := status.Convert(err)
	st := HTTPFromCode(s.Code())

	var stack string
	for _, detail := range s.Details() {
		switch t := detail.(type) {
		case *errdetails.DebugInfo:
			for _, violation := range t.GetStackEntries() {
				stack = violation
			}
		}
	}

	ErrorReport(md, req, stack, st)
	return status.Error(s.Code(),s.Message())
}

func ErrorReport(md metadata.MD, req interface{}, stack string, code int) {
	isDontReport := false
	for _, value := range dontReport {
		if value == code {
			isDontReport = true
		}
	}
	errUrl := viper.GetString("error.report")
	if errUrl != "" && !isDontReport {
		request := map[string]interface{}{
			"header": md,
			"params": req,
		}

		app := map[string]string{
			"name":        viper.GetString("app.name"),
			"environment": viper.GetString("app.env"),
		}

		exception := map[string]interface{}{
			"code":  code,
			"trace": stack,
		}

		option := map[string]interface{}{
			"error_type": "error",
			"app":        app,
			"exception":  exception,
			"request":    request,
		}

		go func() {
			_, _ = httpclient.Begin().PostJson(errUrl, option)
		}()
	}
}

func Err(code codes.Code, msg string, stack string) error {
	s := status.New(code, msg)
	st, _ := s.WithDetails(&errdetails.DebugInfo{
		StackEntries: []string{stack},
	})
	return st.Err()
}
