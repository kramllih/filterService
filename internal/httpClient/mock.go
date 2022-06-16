package httpClient

import (
	"net/http"
)

func MockHTTP() *HTTP {
	return &HTTP{
		client: &http.Client{},
	}
}
