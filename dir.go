package go_hidrive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

/*
DirApi - structure represents a set of methods for interacting with HiDrive `/dir` API endpoint.
*/
type DirApi struct {
	Api
}

/*
NewDirApi - create new instance of [DirApi].

Accepts http.Client and API endpoint as input parameters.
If `endpoint` is empty string, then default [StratoHiDriveAPIV21] value is used.
*/
func NewDirApi(client *http.Client, endpoint string) DirApi {
	api := NewApi(client, endpoint)
	return DirApi{api}
}

/*
GetDir - this method allows to query information about a given directory and all its contents.

In short; A few things to be aware of:
- path and name values are returned as URL-encoded strings
- an implicit limit of 5000 is used by default
- this also works for snapshots

Usage details:
Both, the pid and path parameters identify a filesystem object, at least one of them is always mandatory.
It is allowed to use both together, in which case pid addresses a parent directory and the value of path is then
considered relative to that directory.

Status codes:
  - 200 - OK
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (password required)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (ID does not exist or given path is not shared).
  - 500 - Internal Error

Supported parameters:
  - path ([Parameters.SetPath])
  - pid ([Parameters.SetPid])
  - members ([Parameters.SetMembers])
  - limit ([Parameters.SetLimit])
  - fields ([Parameters.SetFields])
  - sort ([Parameters.SetSortBy])
  - sort_lang ([Parameters.SetSortLang])

Returns [HiDriveObject] with information about given directory.
*/
func (d DirApi) GetDir(ctx context.Context, params url.Values) (*HiDriveObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = d.doGET(ctx, "dir", params); err != nil {
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

/*
CreateDir - this method creates a new directory.

Provide a `path` and optional superior `pid` to create a new directory.
Both, the `pid` and `path` parameters identify a filesystem object, at least one of them is always mandatory.
It is allowed to use both together, in which case pid addresses a parent directory and the value of path is then
considered relative to that directory.

Note: This will not create all missing parent directories - the parent directory within path has to exist!
To create all missing parent directories in one shot use [DirApi.CreatePath] method.

Status codes:
  - 201 - Created
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (no authentication)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (e.g. ID not existing)
  - 500 - Internal Error

Supported parameters:
  - path ([Parameters.SetPath])
  - pid ([Parameters.SetPid])
  - on_exist ([Parameters.SetOnExist])
  - mtime ([Parameters.SetMTime])
  - parent_mtime ([Parameters.SetParentMTime])

Returns [HiDriveObject] with information about the directory created.
*/
func (d DirApi) CreateDir(ctx context.Context, params url.Values) (*HiDriveObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = d.doPOST(ctx, "dir", params, nil); err != nil {
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

/*
CreatePath - this method performs the same action as [DirApi.CreateDir] and also creates all parent directories
if they are missing.

Note: method does not support `pid` parameter, only `path` can be used

Returns [HiDriveObject] with information about the directory created.
*/
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

/*
DeleteDir - this method deletes a given directory.

The optional parameter `recursive` determines, whether the operation shall fail on non-empty directories
(which is also the default behavior), or continue deleting recursively all contents.

Both, the pid and path parameters identify a filesystem object, at least one of them is always mandatory.
It is allowed to use both together, in which case pid addresses a parent directory and the value of path is then
considered relative to that directory.

Enforce the HiDrive ACLs but adjusts deviant POSIX ACLs of the src and his parent when required.
This may be necessary if the Permissions have been changed using protocols like SMB or rsync.
The permissions will be restored if the operation completes successfully. Due to missing transactional semantics the
permissions may not be restored if the operation fails (e.g. with an exception)

Status codes:
  - 204 - No Content
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized
  - 403 - Forbidden
  - 404 - Not Found (e.g. ID not existing)
  - 500 - Internal Error

Supported parameters:
  - path ([Parameters.SetPath])
  - pid ([Parameters.SetPid])
  - recursive ([Parameters.SetRecursive])
  - parent_mtime ([Parameters.SetParentMTime])
*/
func (d DirApi) DeleteDir(ctx context.Context, params url.Values) error {
	var res *http.Response

	{
		var err error
		if res, err = d.doDELETE(ctx, "dir", params); err != nil {
			return err
		}
	}

	if err := d.checkHTTPStatus(http.StatusNoContent, res); err != nil {
		return err
	}

	return nil
}
