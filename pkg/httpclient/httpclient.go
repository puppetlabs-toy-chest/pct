package httpclient

import "net/http"

type HTTPClientI interface {
	Get(url string) (resp *http.Response, err error)
}

type HTTPClient struct {
	Client *http.Client
}
