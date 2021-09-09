package mock

import (
	"fmt"
	"net/http"
)

type GetResponse struct {
	RequestResponse *http.Response
	ErrResponse     bool
}

type HTTPClient struct {
	RequestResponse *http.Response
	ErrResponse     bool
}

func (h *HTTPClient) Get(url string) (*http.Response, error) {
	if h.ErrResponse {
		return nil, fmt.Errorf("Web request error")
	}

	return h.RequestResponse, nil
}
