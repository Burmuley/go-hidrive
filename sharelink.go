package go_hidrive

import (
	"context"
	"net/http"
	"net/url"
)

/*
Sharelink - structure represents a set of methods for interacting with HiDrive `/sharelink` API endpoint.
*/
type Sharelink struct {
	Api
}

/*
NewSharelink - create new instance of Sharelink.

Accepts http.Client and API endpoint as input parameters.
If `endpoint` is empty string, then default `StratoHiDriveAPIV21` value is used.
*/
func NewSharelink(client *http.Client, endpoint string) Sharelink {
	api := NewApi(client, endpoint)
	return Sharelink{api}
}

/*
Create - create a new sharelink for a given file.

Both, the "pid" and "path" parameters refer to the file, at least one of them is mandatory.
If both are given, `pid` addresses a parent directory, and the value of `path` is a relative path to the actual file.

Specific values might be limited by package-feature settings:
Parameters `ttl` and `maxcount` are required, if the tariff defines a maximum limit for these values.
The password protection feature is not available in all tariffs.

Status codes:
  - 201 - Created
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (no authentication)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (e.g. ID not existing)
  - 500 - Internal Error

Supported parameters:
  - maxcount ([Parameters.SetMaxCount])
  - path ([Parameters.SetPath])
  - pid ([Parameters.SetPid])
  - password ([Parameters.SetPassword])
  - ttl ([Parameters.SetTTL])
  - type - always set by the method to value `file`
*/
func (sl Sharelink) Create(ctx context.Context, params url.Values) (*ShareObject, error) {
	var (
		res *http.Response
		err error
	)

	params.Set("type", "file")
	if res, err = sl.doPOST(ctx, "sharelink", params, []int{http.StatusCreated}, nil); err != nil {
		return nil, err
	}

	obj := &ShareObject{}
	if err := sl.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
Get - if no "id" parameter is given, a list of all sharelink objects of the user is returned.
With a given "id" only the corresponding `sharelink` object is returned, if that exists.

Status codes:
  - 200 - OK
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (password required)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (ID does not exist or given path is not shared).
  - 500 - Internal Error

Supported parameters:
  - id ([Parameters.SetId])
  - fields ([Parameters.SetFields])
*/
func (sl Sharelink) Get(ctx context.Context, params url.Values) (*ShareObject, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = sl.doGET(ctx, "sharelink", params, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	obj := &ShareObject{}
	if err := sl.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
Update - update values for a given `sharelink` (not available for all tariffs).

Specific values might be limited due to package-feature settings:
  - The password protection feature is not available in all tariffs.
  - Parameters `ttl` and `maxcount` are required, if the tariff defines a maximum limit for these values.
  - The new value for `maxcount` must be equal to greater than the current count and the difference
    between the value of the current `maxcount` and the new `maxcount` may be limited, depending on the tariff.

Status codes:
  - 200 - OK
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (no authentication)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (e.g. ID not existing)
  - 500 - Internal Error

Supported parameters:
  - maxcount ([Parameters.SetMaxCount])
  - id ([Parameters.SetId])
  - password ([Parameters.SetPassword])
  - ttl ([Parameters.SetTTL])
*/
func (sl Sharelink) Update(ctx context.Context, params url.Values) (*ShareObject, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = sl.doPUT(ctx, "sharelink", params, []int{http.StatusOK}, nil); err != nil {
		return nil, err
	}

	obj := &ShareObject{}
	if err := sl.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
Delete - remove `sharelink`.

Status codes:
  - 204 - No Content
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized
  - 403 - Forbidden
  - 404 - Not Found (e.g. ID not existing)
  - 500 - Internal Error

Supported parameters:
  - id ([Parameters.SetId])
*/
func (sl Sharelink) Delete(ctx context.Context, params url.Values) error {
	if _, err := sl.doDELETE(ctx, "sharelink", params, []int{http.StatusNoContent}); err != nil {
		return err
	}

	return nil
}
