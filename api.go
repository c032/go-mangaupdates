package mangaupdates

import (
	"bytes"
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

type SeriesID int64

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

type SType string

const (
	STypeTitle       SType = "title"
	STypeDescription SType = "description"
)

type Series struct {
	SeriesID int    `json:"series_id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Type     string `json:"type"`
}

type SeriesSearchRequest struct {
	Search string `json:"search"`
	SType  SType  `json:"stype"`
}

type SeriesSearchResponseResult struct {
	Record *Series `json:"record"`
}

type SeriesSearchResponse struct {
	TotalHits int                          `json:"total_hits"`
	Page      int                          `json:"page"`
	PerPage   int                          `json:"per_page"`
	Results   []SeriesSearchResponseResult `json:"results"`
}

func (muc *Client) SeriesSearch(request SeriesSearchRequest) (*SeriesSearchResponse, error) {
	var (
		err             error
		seriesSearchURL *url.URL
	)

	seriesSearchURL, err = baseURL.Parse("/v1/series/search")
	if err != nil {
		return nil, fmt.Errorf("could not parse url: %w", err)
	}

	var rawRequestBody []byte

	rawRequestBody, err = json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("could not encode request body: %w", err)
	}

	reqBody := bytes.NewBuffer(rawRequestBody)

	var req *nethttp.Request
	req, err = nethttp.NewRequest(nethttp.MethodPost, seriesSearchURL.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	hc := muc.httpClient()

	var resp *nethttp.Response

	resp, err = hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get HTTP response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// TODO: Return information about error.
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	d := json.NewDecoder(resp.Body)

	ssr := &SeriesSearchResponse{}

	err = d.Decode(&ssr)
	if err != nil {
		return nil, fmt.Errorf("could not decode JSON: %w", err)
	}

	return ssr, nil
}

func (muc *Client) SeriesByID(seriesID SeriesID) (*Series, error) {
	var (
		err       error
		seriesURL *url.URL
	)

	seriesURL, err = baseURL.Parse(fmt.Sprintf("/v1/series/%d", seriesID))
	if err != nil {
		return nil, fmt.Errorf("could not parse url: %w", err)
	}

	var req *nethttp.Request
	req, err = nethttp.NewRequest(nethttp.MethodGet, seriesURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	hc := muc.httpClient()

	var resp *nethttp.Response

	resp, err = hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get HTTP response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// TODO: Return information about error.
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	d := json.NewDecoder(resp.Body)

	s := &Series{}

	err = d.Decode(&s)
	if err != nil {
		return nil, fmt.Errorf("could not decode JSON: %w", err)
	}

	return s, nil
}
