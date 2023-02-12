package go_hidrive

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type ShareApi struct {
	Api
}

func NewShareApi(client *http.Client, endpoint string) ShareApi {
	api := NewApi(client, endpoint)
	return ShareApi{api}
}

func (s ShareApi) GetShare(ctx context.Context, params url.Values) ([]*HiDriveShareObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = s.DoGET(ctx, "share", params); err != nil {
			return nil, err
		}
	}

	if err := s.checkHTTPStatus(http.StatusOK, res); err != nil {
		return nil, err
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	obj := make([]*HiDriveShareObject, 0)
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s ShareApi) CreateShare(ctx context.Context, params url.Values) (*HiDriveShareObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = s.DoPOST(ctx, "share", params, nil); err != nil {
			return nil, err
		}
	}

	if err := s.checkHTTPStatus(http.StatusCreated, res); err != nil {
		return nil, err
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	obj := &HiDriveShareObject{}
	if err := obj.UnmarshalJSON(body); err != nil {
		return nil, err
	}

	return obj, nil

}

func (s ShareApi) DeleteShare(ctx context.Context, params url.Values) error {
	var res *http.Response

	{
		var err error
		if res, err = s.DoDELETE(ctx, "share", params); err != nil {
			return err
		}
	}

	if err := s.checkHTTPStatus(http.StatusNoContent, res); err != nil {
		return err
	}

	return nil
}

func (s ShareApi) UpdateShare(ctx context.Context, params url.Values) (*HiDriveShareObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = s.DoPUT(ctx, "share", params, nil); err != nil {
			return nil, err
		}
	}

	if err := s.checkHTTPStatus(http.StatusOK, res); err != nil {
		return nil, err
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	obj := &HiDriveShareObject{}
	if err := obj.UnmarshalJSON(body); err != nil {
		return nil, err
	}

	return obj, nil
}
