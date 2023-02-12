package go_hidrive

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	StratoHiDriveAPIV21   = "https://api.hidrive.strato.com/2.1"      // Default HiDrive API endpoint
	StratoHiDriveAuthURL  = "https://my.hidrive.com/client/authorize" // Default HiDrive authentication URL
	StratoHiDriveTokenURL = "https://my.hidrive.com/oauth2/token"     // Default HiDrive token operations URL
)

// Api - TODO
type Api struct {
	APIEndpoint string
	HTTPClient  *http.Client
}

func NewApi(client *http.Client, endpoint string) Api {
	if endpoint == "" {
		endpoint = StratoHiDriveAPIV21
	}
	return Api{
		APIEndpoint: endpoint,
		HTTPClient:  client,
	}
}

func (a Api) NewHTTPRequest(ctx context.Context, method, uri string, r io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, strings.Join([]string{a.APIEndpoint, uri}, "/"), r)
}

func (a Api) DoGET(ctx context.Context, uri string, params url.Values) (*http.Response, error) {
	return a.DoHTTPRequest(ctx, "GET", uri, params, nil)
}

func (a Api) DoDELETE(ctx context.Context, uri string, params url.Values) (*http.Response, error) {
	return a.DoHTTPRequest(ctx, "DELETE", uri, params, nil)
}

func (a Api) DoPOST(ctx context.Context, uri string, params url.Values, body io.ReadCloser) (*http.Response, error) {
	return a.DoHTTPRequest(ctx, "POST", uri, params, body)
}

func (a Api) DoPUT(ctx context.Context, uri string, params url.Values, body io.ReadCloser) (*http.Response, error) {
	return a.DoHTTPRequest(ctx, "PUT", uri, params, body)
}

func (a Api) DoHTTPRequest(ctx context.Context, method, uri string, params url.Values, body io.ReadCloser) (*http.Response, error) {
	var (
		req *http.Request
		res *http.Response
	)

	{
		var err error
		if req, err = a.NewHTTPRequest(ctx, method, uri, body); err != nil {
			return nil, err
		}
	}

	req.URL.RawQuery = params.Encode()

	{
		var err error
		if res, err = a.HTTPClient.Do(req); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (a Api) checkHTTPStatus(desiredCode int, res *http.Response) error {
	var err error
	var body []byte

	if res.StatusCode != desiredCode {
		hdErr := &HiDriveError{}
		if body, err = io.ReadAll(res.Body); err != nil {
			return err
		}
		if err := json.Unmarshal(body, hdErr); err != nil {
			return err
		}
		return hdErr
	}

	return nil
}
