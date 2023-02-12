package go_hidrive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type DirApi struct {
	Api
}

func NewDirApi(client *http.Client, endpoint string) DirApi {
	api := NewApi(client, endpoint)
	return DirApi{api}
}

// GetDir retrieves directory information defined in the `path`.
func (d DirApi) GetDir(ctx context.Context, params url.Values) (*HiDriveObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = d.DoGET(ctx, "dir", params); err != nil {
			return nil, err
		}
	}

	if err := d.checkHTTPStatus(http.StatusOK, res); err != nil {
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

// CreateDir creates directory defined in the `path`.
// This functions will not create all parent directories, they must exist.
func (d DirApi) CreateDir(ctx context.Context, params url.Values) (*HiDriveObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = d.DoPOST(ctx, "dir", params, nil); err != nil {
			return nil, err
		}
	}

	if err := d.checkHTTPStatus(http.StatusCreated, res); err != nil {
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

// CreatePath creates directory defined in the `path` and all parent directories preceding it
func (d DirApi) CreatePath(ctx context.Context, params url.Values) (*HiDriveObject, error) {
	path := params.Get("path")
	if len(path) < 1 {
		return nil, fmt.Errorf("path: %w", ErrShouldNotBeEmpty)
	}

	dirs := strings.Split(path, "/")
	for k := range dirs[:len(dirs)-1] {
		dir := fmt.Sprintf("/%s", strings.Join(dirs[1:k+1], "/"))
		tmpp := NewParameters().SetMembers([]string{"none"}).SetFields([]string{"path"}).SetPath(dir)
		if _, err := d.GetDir(ctx, tmpp.Values); err == nil {
			continue
		}
		if _, err := d.CreateDir(ctx, NewParameters().SetPath(dir).Values); err != nil {
			return nil, err
		}
	}
	return d.CreateDir(ctx, NewParameters().SetPath(path).Values)
}

// DeleteDir deletes the directory defined in the `path`
func (d DirApi) DeleteDir(ctx context.Context, params url.Values) error {
	var res *http.Response

	{
		var err error
		if res, err = d.DoDELETE(ctx, "dir", params); err != nil {
			return err
		}
	}

	if err := d.checkHTTPStatus(http.StatusNoContent, res); err != nil {
		return err
	}

	return nil
}
