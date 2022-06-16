package httpClient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type HTTP struct {
	client  *http.Client // HTTP client that is reused across requests.
	headers map[string]string
	params  map[string]string
	name    string
	uri     string
	method  string
	body    []byte
}

func NewHTTP() *HTTP {

	return &HTTP{
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
			},
		},
	}

}

func (h *HTTP) FetchResponse() (*http.Response, error) {

	var reader io.Reader
	if h.body != nil {
		reader = bytes.NewReader(h.body)
	}
	req, err := http.NewRequest(h.method, h.uri, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	res, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}

	return res, nil
}

func (h *HTTP) FetchContent() ([]byte, error) {

	res, err := h.FetchResponse()
	if err != nil {
		return nil, err
	}
	if res != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP error %d in %s: %s", res.StatusCode, h.name, res.Status)
	}

	return ioutil.ReadAll(res.Body)
}

func (h *HTTP) SetMethod(method string) {
	h.method = method
}

func (h *HTTP) SetBody(body []byte) {
	h.body = body
}

func (h *HTTP) GetURI() string {
	return h.uri
}

func (h *HTTP) SetURI(uri string) {
	h.uri = uri
}

func (h *HTTP) SetTransport(tran http.RoundTripper) {
	h.client.Transport = tran
}
