package go_hidrive

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

/*
File - structure represents a set of methods for interacting with HiDrive `/file` API endpoint.
*/
type File struct {
	Api
}

/*
NewFile - create new instance of File.

Accepts http.Client and API endpoint as input parameters.
If `endpoint` is empty string, then default `StratoHiDriveAPIV21` value is used.
*/
func NewFile(client *http.Client, endpoint string) File {
	api := NewApi(client, endpoint)
	return File{api}
}

/*
Get - This method retrieves a given file from the HiDrive.

Usage details:
Both, the `pid` and `path` parameters identify a filesystem object, at least one of them is always mandatory.
It is allowed to use both together, in which case pid addresses a parent directory and the value of path is then considered relative to that directory. (<pid>/<path>)

Status codes:
  - 200 - OK
  - 304 - Not Modified
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (password required)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (ID does not exist or given path is not shared).
  - 416 - Requested Range Not Satisfiable
  - 500 - Internal Error

Supported parameters:
  - path ([Parameters.SetPath])
  - pid ([Parameters.SetPid])

Returns an io.ReadCloser object to read file contents using standard Go mechanisms.
*/
func (f File) Get(ctx context.Context, params url.Values) (io.ReadCloser, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = f.doGET(ctx, "file", params, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	return res.Body, nil
}

/*
Upload -This method can be used to create a new file and store uploaded content.

Using POST guarantees that existing files will not be overwritten.

Usage Details:

The dir_id and dir parameters must be specified as query parameters of the URI.
They refer to the file's target directory on the HiDrive storage, at least one of the parameters is mandatory.
If both are given, dir_id addresses a directory and the value of dir is taken as a relative path to that directory.

The size of the request body to upload is limited to 2147483648 bytes (2G). The size of the complete request, including
header, and after possible decoding of chunked encoding and decompression is limited to 3206545408 bytes (3058MB).
Larger requests are rejected with 413 Request Entity Too Large.

As existence of the target file will be checked only after the upload is complete, the target file may have sprung into
existence during the upload. To avoid losing the uploaded content in this case, the optional on_exist parameter can be
set to "autoname". When set, the returned data will contain a name that differs from the provided name parameter.

Status codes:
  - 201 - Created
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (password required)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (ID does not exist or given path is not shared).
  - 409 - Conflict
  - 413 - Request Entity Too Large
  - 415 - Unsupported Media Type
  - 422 - Unprocessable Entity (e.g. name too long)
  - 500 - Internal Error
  - 507 - Insufficient Storage

Supported parameters:
  - dir ([Parameters.SetDir])
  - dir_id ([Parameters.SetDirId])
  - name ([Parameters.SetName])
  - on_exist ([Parameters.SetOnExist])
  - mtime ([Parameters.SetMTime])
  - parent_mtime ([Parameters.SetParentMTime])

Returns [Object] with information about uploaded file.
*/
func (f File) Upload(ctx context.Context, params url.Values, fileBody io.ReadCloser) (*Object, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = f.doPOST(ctx, "file", params, []int{http.StatusCreated}, fileBody); err != nil {
		return nil, err
	}

	obj := &Object{}
	if err := f.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
Delete - Delete a given file.
Both, the pid and path parameters identify a filesystem object, at least one of them is always mandatory.
It is allowed to use both together, in which case pid addresses a parent directory and the value of path is then considered relative to that directory. (<pid>/<path>)

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
  - parent_mtime ([Parameters.SetParentMTime])
*/
func (f File) Delete(ctx context.Context, params url.Values) error {
	if _, err := f.doDELETE(ctx, "file", params, []int{http.StatusNoContent}); err != nil {
		return err
	}
	return nil
}

/*
Update - Update a file by overwriting the target file with uploaded content.
If the target file does not exist it will be created.

If you wish to create a file without overwriting data, use [File.Upload].

Usage Details:

The dir_id and dir parameters must be specified as query parameters of the URI.
They refer to the file's target directory on the HiDrive storage, at least one of the parameters is mandatory.
If both are given, dir_id addresses a directory and the value of dir is taken as a relative path to that directory.

The size of the request body to upload is limited to 2147483648 bytes (2G). The size of the complete request, including
header, and after possible decoding of chunked encoding and decompression is limited to 3206545408 bytes (3058MB).
Larger requests are rejected with 413 Request Entity Too Large.

As existence of the target file will be checked only after the upload is complete, the target file may have sprung into
existence during the upload. To avoid losing the uploaded content in this case, the optional on_exist parameter can be
set to "autoname". When set, the returned data will contain a name that differs from the provided name parameter.

Status codes:
  - 201 - Created
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (password required)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (ID does not exist or given path is not shared).
  - 409 - Conflict
  - 413 - Request Entity Too Large
  - 415 - Unsupported Media Type
  - 422 - Unprocessable Entity (e.g. name too long)
  - 500 - Internal Error
  - 507 - Insufficient Storage

Supported parameters:
  - dir ([Parameters.SetDir])
  - dir_id ([Parameters.SetDirId])
  - name ([Parameters.SetName])
  - mtime ([Parameters.SetMTime])
  - parent_mtime ([Parameters.SetParentMTime])

Returns [Object] with information about uploaded file.
*/
func (f File) Update(ctx context.Context, params url.Values, fileBody io.ReadCloser) (*Object, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = f.doPUT(ctx, "file", params, []int{http.StatusOK}, fileBody); err != nil {
		return nil, err
	}

	obj := &Object{}
	if err := f.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
Copy - copy a file.

The parameters `src` and `src_id` as well as `dst` and `dst_id` identify the source and destination for the operation.
At least one source identifier and `dst` are always mandatory.
It is allowed to use the related parameters together, in which case `src_id` and `dst_id` each address a parent
directory and the values of `src` and `dst` are considered relative to that directory (e.g.<src_id>/<src>).

Status codes:
  - 200 - OK
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (password required)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (ID does not exist or given path is not shared).
  - 409 - Conflict
  - 422 - Unprocessable Entity (e.g. name too long)
  - 500 - Internal Error

Supported parameters:
  - src
  - src_id
  - dst
  - dst_id
  - on_exist ([Parameters.SetOnExist]) (possible values: `autoname`, `overwrite`)
  - dst_parent_mtime
  - preserve_mtime
*/
func (f File) Copy(ctx context.Context, params url.Values) (*Object, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = f.doPOST(ctx, "file/copy", params, []int{http.StatusOK}, nil); err != nil {
		return nil, err
	}

	obj := &Object{}
	if err := f.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
Move - move a file.

The parameters `src` and `src_id` as well as `dst` and `dst_id` identify the source and destination for the operation.
At least one source identifier and `dst` are always mandatory.
It is allowed to use the related parameters together, in which case `src_id` and `dst_id` each address a parent
directory and the values of `src` and `dst` are considered relative to that directory (e.g.<src_id>/<src>).

Status codes:
  - 200 - OK
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (password required)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (ID does not exist or given path is not shared).
  - 409 - Conflict
  - 422 - Unprocessable Entity (e.g. name too long)
  - 500 - Internal Error

Supported parameters:
  - src
  - src_id
  - dst
  - dst_id
  - on_exist ([Parameters.SetOnExist]) (possible values: `autoname`, `overwrite`)
  - src_parent_mtime
  - dst_parent_mtime
*/
func (f File) Move(ctx context.Context, params url.Values) (*Object, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = f.doPOST(ctx, "file/move", params, []int{http.StatusOK}, nil); err != nil {
		return nil, err
	}

	obj := &Object{}
	if err := f.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

/*
Rename - rename a file.

Both, the `pid` and `path` parameters identify a filesystem object, at least one of them is always mandatory.
It is allowed to use both together, in which case `pid` addresses a parent directory and the value of `path` is then
considered relative to that directory (<pid>/<path>).

Status codes:
  - 201 - Created
  - 400 - Bad Request (e.g. invalid parameter)
  - 401 - Unauthorized (password required)
  - 403 - Forbidden (wrong password)
  - 404 - Not Found (ID does not exist or given path is not shared).
  - 409 - Conflict
  - 422 - Unprocessable Entity (e.g. name too long)
  - 500 - Internal Error

Supported parameters:
  - path ([Parameters.SetPath])
  - pid ([Parameters.SetPid])
  - name ([Parameters.SetName])
  - on_exist ([Parameters.SetOnExist]) (possible values: `autoname`, `overwrite`)
  - parent_mtime ([Parameters.SetParentMTime])
*/
func (f File) Rename(ctx context.Context, params url.Values) (*Object, error) {
	var (
		res *http.Response
		err error
	)

	if res, err = f.doPOST(ctx, "file/rename", params, []int{http.StatusCreated}, nil); err != nil {
		return nil, err
	}

	obj := &Object{}
	if err := f.unmarshalBody(res, obj); err != nil {
		return nil, err
	}

	return obj, nil
}
