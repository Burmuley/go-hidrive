package go_hidrive

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type FileApi struct {
	Api
}

func NewFileApi(client *http.Client, endpoint string) FileApi {
	api := NewApi(client, endpoint)
	return FileApi{api}
}

func (f FileApi) GetFile(ctx context.Context, params url.Values) (io.ReadCloser, error) {
	var res *http.Response

	{
		var err error
		if res, err = f.DoGET(ctx, "file", params); err != nil {
			return nil, err
		}
	}

	if err := f.checkHTTPStatus(http.StatusOK, res); err != nil {
		return nil, err
	}

	return res.Body, nil
}

func (f FileApi) UploadFile(ctx context.Context, params url.Values, fileBody io.ReadCloser) (*HiDriveObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = f.DoPOST(ctx, "file", params, fileBody); err != nil {
			return nil, err
		}
	}

	if err := f.checkHTTPStatus(http.StatusCreated, res); err != nil {
		return nil, err
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	hdObj := &HiDriveObject{}
	if err := hdObj.UnmarshalJSON(body); err != nil {
		return nil, err
	}

	return hdObj, nil
}

func (f FileApi) DeleteFile(ctx context.Context, params url.Values) error {
	var res *http.Response

	{
		var err error
		if res, err = f.DoDELETE(ctx, "file", params); err != nil {
			return err
		}
	}

	if err := f.checkHTTPStatus(http.StatusNoContent, res); err != nil {
		return err
	}

	return nil
}
