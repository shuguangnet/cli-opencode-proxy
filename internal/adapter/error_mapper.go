package adapter

import (
	"fmt"
	"net/http"

	"opencode-cli-proxy/internal/domain"
)

func MapError(status int, err error) (int, domain.OpenAIErrorResponse) {
	if err == nil {
		err = fmt.Errorf("unknown error")
	}

	typeName := "server_error"
	code := fmt.Sprintf("%d", status)

	switch status {
	case http.StatusBadRequest:
		typeName = "invalid_request_error"
	case http.StatusUnauthorized:
		typeName = "invalid_api_key"
	case http.StatusForbidden:
		typeName = "authentication_error"
	case http.StatusTooManyRequests:
		typeName = "rate_limit_exceeded"
	case http.StatusGatewayTimeout:
		typeName = "request_timeout"
	case http.StatusBadGateway, http.StatusInternalServerError, http.StatusServiceUnavailable:
		typeName = "server_error"
	}

	return status, domain.OpenAIErrorResponse{Error: domain.OpenAIErrorBody{
		Message: err.Error(),
		Type:    typeName,
		Param:   nil,
		Code:    code,
	}}
}
