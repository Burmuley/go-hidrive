package go_hidrive

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

/*
Dir - structure represents a set of methods for interacting with HiDrive `/dir` API endpoint.
*/
type Dir struct {
	Api
}

/*
NewDir - create new instance of [Dir].

Accepts http.Client and API endpoint as input parameters.
If `endpoint` is empty string, then default [StratoHiDriveAPIV21] value is used.
*/
func NewDir(client *http.Client, endpoint string) Dir {
	api := NewApi(client, endpoint)
	return Dir{api}
}

/*
Get - this method allows to query information about a given directory and all its contents.

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

Returns [Object] with information about given directory.
*/
func (d Dir) Get(ctx context.Context, params url.Values) (*Object, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = d.doGET(ctx, "dir", params, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	obj := &Object{}
	if err := d.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
Create - this method creates a new directory.

Provide a `path` and optional superior `pid` to create a new directory.
Both, the `pid` and `path` parameters identify a filesystem object, at least one of them is always mandatory.
It is allowed to use both together, in which case pid addresses a parent directory and the value of path is then
considered relative to that directory.

Note: This will not create all missing parent directories - the parent directory within path has to exist!
To create all missing parent directories in one shot use [Dir.CreatePath] method.

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

Returns [Object] with information about the directory created.
*/
func (d Dir) Create(ctx context.Context, params url.Values) (*Object, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = d.doPOST(ctx, "dir", params, []int{http.StatusCreated}, nil); err != nil {
		return nil, err
	}

	obj := &Object{}
	if err := d.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
CreatePath - this method performs the same action as [Dir.Create] and also creates all parent directories
if they are missing.

Note: method does not support `pid` parameter, only `path` can be used

Returns [Object] with information about the directory created.
*/
func (d Dir) CreatePath(ctx context.Context, params url.Values) (*Object, error) {
	path := params.Get("path")
	if len(path) < 1 {
		return nil, fmt.Errorf("path: %w", ErrShouldNotBeEmpty)
	}

	dirs := strings.Split(path, "/")
	for k := range dirs[:len(dirs)-1] {
		dir := fmt.Sprintf("/%s", strings.Join(dirs[1:k+1], "/"))
		tmpp := NewParameters().SetMembers([]string{"none"}).SetFields([]string{"path"}).SetPath(dir)
		if _, err := d.Get(ctx, tmpp.Values); err == nil {
			continue
		}
		if _, err := d.Create(ctx, NewParameters().SetPath(dir).Values); err != nil {
			return nil, err
		}
	}
	return d.Create(ctx, NewParameters().SetPath(path).Values)
}

/*
Delete - this method deletes a given directory.

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
func (d Dir) Delete(ctx context.Context, params url.Values) error {
	if _, err := d.doDELETE(ctx, "dir", params, []int{http.StatusNoContent}); err != nil {
		return err
	}
	return nil
}
