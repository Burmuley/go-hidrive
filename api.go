/*
Package go_hidrive is a simple client SDK library for HiDrive cloud storage
(mainly provided by [Strato](https://www.strato.de/cloud-speicher/) provider)

Currently, the following implementation are available: [Dir], [File] and [Share].

All methods accept url.Values as a set of request parameters.
You can also use [Parameters] objects to simplify parameters gathering required for request.

Example reading file from HiDrive:

	package main

	import (
		"log"
		"context"
		"io"
		"fmt"

		"golang.org/x/oauth2"
		hidrive "github.com/Burmuley/go-hidrive"
	)

	func main() {
		oauth2config := oauth2.Config{
			ClientID:     "hi_drive_client_id",
			ClientSecret: "hi_drive_client_secret",
			Endpoint: oauth2.Endpoint{
				AuthURL:   hidrive.StratoHiDriveAuthURL,
				TokenURL:  hidrive.StratoHiDriveTokenURL,
				AuthStyle: 0,
			},
			Scopes: []string{"user", "rw"},
		}

		token := &oauth2.Token{
			RefreshToken: "hi_drive_oauth2_refresh_token",
		}

		client := oauth2config.Client(context.Background(), token)
		fileApi := hidrive.NewFile(client, hidrive.StratoHiDriveAPIV21)

		rdr, err := fileApi.Get(context.Background(), hidrive.NewParameters().SetPath("/public/test_file.txt").Values)

		if err != nil {
			log.Fatal(err)
		}

		contents, err := io.ReadAll(rdr)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(contents)
	}
*/
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

/*
Api - basic structure defining common logic for API interaction.

Property `HTTPClient` should be a [http.Client] type and retrieved from `oauth2` package,
i.e. it should be pre-configured to perform OAuth2 authentication against HiDrive API before
underlying method send any data.

Property `APIEndpoint` should be set to proper HiDrive API endpoint.
Use [NewApi] function to create new instances of this type, it supports empty `endpoint` and
injects default from [StratoHiDriveAPIV21] constant.
*/
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

func (a Api) newHTTPRequest(ctx context.Context, method, uri string, r io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, strings.Join([]string{a.APIEndpoint, uri}, "/"), r)
}

func (a Api) doGET(ctx context.Context, uri string, params url.Values, okCodes []int) (*http.Response, error) {
	return a.doHTTPRequest(ctx, "GET", uri, params, okCodes, nil)
}

func (a Api) doDELETE(ctx context.Context, uri string, params url.Values, okCodes []int) (*http.Response, error) {
	return a.doHTTPRequest(ctx, "DELETE", uri, params, okCodes, nil)
}

func (a Api) doPOST(ctx context.Context, uri string, params url.Values, okCodes []int, body io.ReadCloser) (*http.Response, error) {
	return a.doHTTPRequest(ctx, "POST", uri, params, okCodes, body)
}

func (a Api) doPUT(ctx context.Context, uri string, params url.Values, okCodes []int, body io.ReadCloser) (*http.Response, error) {
	return a.doHTTPRequest(ctx, "PUT", uri, params, okCodes, body)
}

func (a Api) doPATCH(ctx context.Context, uri string, params url.Values, okCodes []int, body io.ReadCloser) (*http.Response, error) {
	return a.doHTTPRequest(ctx, "PATCH", uri, params, okCodes, body)
}

func (a Api) doHTTPRequest(ctx context.Context, method, uri string, params url.Values, okCodes []int, body io.ReadCloser) (*http.Response, error) {
	var (
		req *http.Request
		res *http.Response
	)

	{
		var err error
		if req, err = a.newHTTPRequest(ctx, method, uri, body); err != nil {
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

	{
		var err error
		if err = a.checkHTTPStatusError(okCodes, res); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (a Api) checkHTTPStatusError(okCodes []int, res *http.Response) error {
	var err error
	var body []byte

	if !isItemInSlice(okCodes, res.StatusCode) {
		hdErr := &Error{}
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

func (a Api) unmarshalBody(res *http.Response, obj any) error {
	var body []byte
	var err error

	if body, err = io.ReadAll(res.Body); err != nil {
		return err
	}

	if err := json.Unmarshal(body, obj); err != nil {
		return err
	}

	return nil
}
