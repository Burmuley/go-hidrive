package go_hidrive

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

/*
ShareApi - structure represents a set of methods for interacting with HiDrive `/share` API endpoint.
*/
type ShareApi struct {
	Api
}

/*
NewShareApi - create new instance of ShareApi.

Accepts http.Client and API endpoint as input parameters.
If `endpoint` is empty string, then default `StratoHiDriveAPIV21` value is used.
*/
func NewShareApi(client *http.Client, endpoint string) ShareApi {
	api := NewApi(client, endpoint)
	return ShareApi{api}
}

/*
GetShare - Get information about either one (given "id", "path" or "pid" parameter) or every existing share created by the authenticated user.
You may customize the result set by adding optional "fields" values.

Status codes:
  - 200 - OK
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (password required)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (ID does not exist or given path is not shared).
  - 500 - Internal Error

Supported parameters:
  - id ([Parameters.SetId])
  - path ([Parameters.SetPath])
  - pid ([Parameters.SetPid])
  - fields ([Parameters.SetFields])
*/
func (s ShareApi) GetShare(ctx context.Context, params url.Values) (*HiDriveShareObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = s.doGET(ctx, "share", params); err != nil {
			return nil, err
		}
	}

	if err := s.checkHTTPStatusError([]int{http.StatusOK}, res); err != nil {
		return nil, err
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	obj := &HiDriveShareObject{}
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
CreateShare - create a new share for a given directory anywhere inside your accessible HiDrive.
You may limit the validity of a share to a given amount of time and protect it with a password.

Sharing a directory will allow anyone with knowledge of the specific (returned) share_id to access all data inside
that directory and all its children (read-only by default).

The path of the shared directory including 'root/' must not exceed 1000 bytes.

For ease of access, HiDrive also provides a share-gui to access the shared files comfortably.
The so-called "share_url" will be returned as well.

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
  - writable ([Parameters.SetWritable])
  - ttl ([Parameters.SetTTL])
  - salt ([Parameters.SetSalt])
  - share_access_key ([Parameters.SetShareAccessKey])
  - pw_sharekey ([Parameters.SetPwShareKey])
*/
func (s ShareApi) CreateShare(ctx context.Context, params url.Values) (*HiDriveShareObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = s.doPOST(ctx, "share", params, nil); err != nil {
			return nil, err
		}
	}

	if err := s.checkHTTPStatusError([]int{http.StatusCreated}, res); err != nil {
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

/*
DeleteShare - delete a given share, thus invalidating each existing share `access_token` immediately.

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
func (s ShareApi) DeleteShare(ctx context.Context, params url.Values) error {
	var res *http.Response

	{
		var err error
		if res, err = s.doDELETE(ctx, "share", params); err != nil {
			return err
		}
	}

	if err := s.checkHTTPStatusError([]int{http.StatusNoContent}, res); err != nil {
		return err
	}

	return nil
}

/*
UpdateShare - update a given share. Change `ttl“, `maxcount` and add or remove a share password.

Note:  It is not possible to change the target directory of an existing share!
Please create a new one, if you wish to share another directory.

Status codes:
  - 200 - OK
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (no authentication)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (e.g. ID not existing)
  - 500 - Internal Error

Supported parameters:
  - id ([Parameters.SetId])
  - maxcount ([Parameters.SetMaxCount])
  - password ([Parameters.SetPassword])
  - writable ([Parameters.SetWritable])
  - ttl ([Parameters.SetTTL])
  - salt ([Parameters.SetSalt])
  - share_access_key ([Parameters.SetShareAccessKey])
  - pw_sharekey ([Parameters.SetPwShareKey])
*/
func (s ShareApi) UpdateShare(ctx context.Context, params url.Values) (*HiDriveShareObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = s.doPUT(ctx, "share", params, nil); err != nil {
			return nil, err
		}
	}

	if err := s.checkHTTPStatusError([]int{http.StatusOK}, res); err != nil {
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

/*
Invite - Invite other people to a share via e-mail.

Status codes:
  - 200 - OK
  - 207 - Multi-Status (body contains multiple status messages)
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (no authentication)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (e.g. ID not existing)
  - 410 - Gone
  - 500 - Internal Error

Supported parameters:
  - id ([Parameters.SetId])
  - path ([Parameters.SetPath])
  - pid ([Parameters.SetPid])
  - fields ([Parameters.SetFields])

Returns [HiDriveShareInviteResponse] object.

The returned object contains the keys `done` and `failed`. Each of these keys holds an array of objects describing
successfully and unsuccessfully processed recipients. Each object holds at least the key `to`, which stores the
recipient's e-mail address. Failure-objects contain an additional key `msg` which describes the encountered error.
If all processed recipients share the same status code, the code will be returned as HTTP status code.
Partial success or differing status codes are indicated by setting the HTTP status code to "207 Multi-Status".
Failure- and done-objects will then contain the individual status of each processed recipient.
*/
func (s ShareApi) Invite(ctx context.Context, params url.Values) (*HiDriveShareInviteResponse, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = s.doPOST(ctx, "share/invite", params, nil); err != nil {
			return nil, err
		}
	}

	if err := s.checkHTTPStatusError([]int{http.StatusOK, http.StatusMultiStatus}, res); err != nil {
		return nil, err
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	obj := &HiDriveShareInviteResponse{}
	if err := json.Unmarshal(body, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
CreateShareLink - create a new sharelink for a given file.

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
func (s ShareApi) CreateShareLink(ctx context.Context, params url.Values) (*HiDriveShareObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	params.Set("type", "file")
	{
		var err error
		if res, err = s.doPOST(ctx, "sharelink", params, nil); err != nil {
			return nil, err
		}
	}

	if err := s.checkHTTPStatusError([]int{http.StatusCreated}, res); err != nil {
		return nil, err
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	obj := &HiDriveShareObject{}
	if err := json.Unmarshal(body, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
GetShareLink - if no "id" parameter is given, a list of all sharelink objects of the user is returned.
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
func (s ShareApi) GetShareLink(ctx context.Context, params url.Values) (*HiDriveShareObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = s.doGET(ctx, "sharelink", params); err != nil {
			return nil, err
		}
	}

	if err := s.checkHTTPStatusError([]int{http.StatusOK}, res); err != nil {
		return nil, err
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	obj := &HiDriveShareObject{}
	if err := json.Unmarshal(body, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
UpdateShareLink - update values for a given `sharelink` (not available for all tariffs).

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
func (s ShareApi) UpdateShareLink(ctx context.Context, params url.Values) (*HiDriveShareObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = s.doPUT(ctx, "sharelink", params, nil); err != nil {
			return nil, err
		}
	}

	if err := s.checkHTTPStatusError([]int{http.StatusOK}, res); err != nil {
		return nil, err
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	obj := &HiDriveShareObject{}
	if err := json.Unmarshal(body, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
DeleteShareLink - remove `sharelink`.

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
func (s ShareApi) DeleteShareLink(ctx context.Context, params url.Values) (*HiDriveShareObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = s.doDELETE(ctx, "sharelink", params); err != nil {
			return nil, err
		}
	}

	if err := s.checkHTTPStatusError([]int{http.StatusNoContent}, res); err != nil {
		return nil, err
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	obj := &HiDriveShareObject{}
	if err := json.Unmarshal(body, obj); err != nil {
		return nil, err
	}

	return obj, nil
}
