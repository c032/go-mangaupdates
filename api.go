package mangaupdates

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"net/url"

	chttp "github.com/c032/go-http"
)

const baseURLStr = "https://api.mangaupdates.com/"

var (
	baseURL           *url.URL
	defaultHTTPClient chttp.Client
)

func init() {
	var err error

	baseURL, err = url.Parse(baseURLStr)
	if err != nil {
		panic(err)
	}

	defaultHTTPClient = nethttp.DefaultClient
}

type Client struct {
	BaseURL    *url.URL
	HTTPClient chttp.Client
}

func (muc *Client) baseURL() *url.URL {
	if muc.BaseURL != nil {
		return muc.BaseURL
	}

	return baseURL
}

func (muc *Client) httpClient() chttp.Client {
	c := muc.HTTPClient
	if c != nil {
		return c
	}

	return defaultHTTPClient
}

type TimeResponse struct {
	Timestamp int64  `json:"timestamp"`
	AsRFC3339 string `json:"as_rfc3339"`
	AsString  string `json:"as_string"`
}

func (muc *Client) Time() (*TimeResponse, error) {
	var (
		err     error
		timeURL *url.URL
	)

	timeURL, err = baseURL.Parse("/v1/misc/time")
	if err != nil {
		return nil, fmt.Errorf("could not parse url: %w", err)
	}

	var req *nethttp.Request

	req, err = nethttp.NewRequest(nethttp.MethodGet, timeURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %w", err)
	}

	hc := muc.httpClient()

	var resp *nethttp.Response

	resp, err = hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get HTTP response: %w", err)
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)

	tr := &TimeResponse{}
	err = d.Decode(&tr)
	if err != nil {
		return nil, fmt.Errorf("could not decode JSON: %w", err)
	}

	return tr, nil
}
