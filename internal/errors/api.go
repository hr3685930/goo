package errors

import (
    "goo/pkg/http/echo"
    "net/http"
    "runtime/debug"
)

func BadRequest(msg string) *echo.HttpError {
    err := &echo.HttpError{Msg: msg, Stack: debug.Stack()}
    err.HttpCode = http.StatusBadRequest
    err.Code = 4400
    return err
}

func ResourceNotFound(msg string) *echo.HttpError {
    err := &echo.HttpError{Msg: msg, Stack: debug.Stack()}
    err.HttpCode = http.StatusNotFound
    err.Code = 4404
    return err
}

func AuthenticationFailed(msg string) *echo.HttpError {
    err := &echo.HttpError{Msg: msg, Stack: debug.Stack()}
    err.HttpCode = http.StatusUnauthorized
    err.Code = 4401
    return err
}

func AuthorizationFailed(msg string) *echo.HttpError {
    err := &echo.HttpError{Msg: msg, Stack: debug.Stack()}
    err.HttpCode = http.StatusForbidden
    err.Code = 4403
    return err
}

func Conflict(msg string) *echo.HttpError {
    err := &echo.HttpError{Msg: msg, Stack: debug.Stack()}
    err.HttpCode = http.StatusMethodNotAllowed
    err.Code = 4405
    return err
}

func ValidationFailed(msg string) *echo.HttpError {
    err := &echo.HttpError{Msg: msg, Stack: debug.Stack()}
    err.HttpCode = http.StatusUnprocessableEntity
    err.Code = 4422
    return err
}

func InternalError(msg string) *echo.HttpError {
    err := &echo.HttpError{Msg: msg, Stack: debug.Stack()}
    err.HttpCode = http.StatusInternalServerError
    err.Code = 5500
    return err
}
