package go_hidrive

import (
	"context"
	"net/http"
	"net/url"
)

/*
Share - structure represents a set of methods for interacting with HiDrive `/share` API endpoint.
*/
type Share struct {
	Api
}

/*
NewShare - create new instance of Share.

Accepts http.Client and API endpoint as input parameters.
If `endpoint` is empty string, then default `StratoHiDriveAPIV21` value is used.
*/
func NewShare(client *http.Client, endpoint string) Share {
	api := NewApi(client, endpoint)
	return Share{api}
}

/*
Get - Get information about either one (given "id", "path" or "pid" parameter) or every existing share created by the authenticated user.
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
func (s Share) Get(ctx context.Context, params url.Values) (*ShareObject, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = s.doGET(ctx, "share", params, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	obj := &ShareObject{}
	if err := s.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
Create - create a new share for a given directory anywhere inside your accessible HiDrive.
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
func (s Share) Create(ctx context.Context, params url.Values) (*ShareObject, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = s.doPOST(ctx, "share", params, []int{http.StatusCreated}, nil); err != nil {
		return nil, err
	}

	obj := &ShareObject{}
	if err := s.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
Delete - delete a given share, thus invalidating each existing share `access_token` immediately.

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
func (s Share) Delete(ctx context.Context, params url.Values) error {
	if _, err := s.doDELETE(ctx, "share", params, []int{http.StatusNoContent}); err != nil {
		return err
	}

	return nil
}

/*
Update - update a given share. Change `ttl“, `maxcount` and add or remove a share password.

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
func (s Share) Update(ctx context.Context, params url.Values) (*ShareObject, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = s.doPUT(ctx, "share", params, []int{http.StatusOK}, nil); err != nil {
		return nil, err
	}

	obj := &ShareObject{}
	if err := s.unmarshalBody(res, obj); err != nil {
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

Returns [ShareInviteResponse] object.

The returned object contains the keys `done` and `failed`. Each of these keys holds an array of objects describing
successfully and unsuccessfully processed recipients. Each object holds at least the key `to`, which stores the
recipient's e-mail address. Failure-objects contain an additional key `msg` which describes the encountered error.
If all processed recipients share the same status code, the code will be returned as HTTP status code.
Partial success or differing status codes are indicated by setting the HTTP status code to "207 Multi-Status".
Failure- and done-objects will then contain the individual status of each processed recipient.
*/
func (s Share) Invite(ctx context.Context, params url.Values) (*ShareInviteResponse, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = s.doPOST(ctx, "share/invite", params, []int{http.StatusOK, http.StatusMultiStatus}, nil); err != nil {
		return nil, err
	}

	obj := &ShareInviteResponse{}
	if err := s.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}
