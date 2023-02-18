package go_hidrive

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

/*
Meta - structure represents a set of methods for interacting with HiDrive `/meta` API endpoint.
*/
type Meta struct {
	Api
}

/*
NewMeta - create new instance of Meta.

Accepts http.Client and API endpoint as input parameters.
If `endpoint` is empty string, then default `StratoHiDriveAPIV21` value is used.
*/
func NewMeta(client *http.Client, endpoint string) Meta {
	api := NewApi(client, endpoint)
	return Meta{api}
}

/*
Get - get metadata of a storage object (dir, file or symlink).
At least one of the parameters path and pid is mandatory.
If both are given, pid must address a directory and the value of path is a relative path from that directory.

To get data from a snapshot, the name of the snapshot must be provided in the snapshot parameter. The name must be UTF-8 encoded.

Status codes:
  - 200 - OK
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (no authentication)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (e.g. ID not existing)
  - 500 - Internal Error

Supported parameters:
  - path ([Parameters.SetPath])
  - pid ([Parameters.SetPid])
  - fields ([Parameters.SetFields])
*/
func (m Meta) Get(ctx context.Context, params url.Values) (*Object, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = m.doGET(ctx, "meta", params, []int{http.StatusOK}); err != nil {
			return nil, err
		}
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	obj := &Object{}
	if err := json.Unmarshal(body, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
Update - modify metadata of a file or directory.

At the moment `mtime` is the only attribute that can be changed.
An attempt to modify metadata of a symlink causes 404 Not Found.

Status codes:
  - 200 - OK
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (no authentication)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (e.g. ID not existing)
  - 500 - Internal Error

Supported parameters:
  - path ([Parameters.SetPath])
  - pid ([Parameters.SetPid])
  - mtime ([Parameters.SetMTime])
*/
func (m Meta) Update(ctx context.Context, params url.Values) (*Object, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = m.doPATCH(ctx, "meta", params, []int{http.StatusOK, http.StatusNoContent}, nil); err != nil {
			return nil, err
		}
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	obj := &Object{}
	if err := json.Unmarshal(body, obj); err != nil {
		return nil, err
	}

	return obj, nil
}
