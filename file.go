package go_hidrive

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

/*
FileApi - structure represents a set of methods for interacting with HiDrive `/file` API endpoint.
*/
type FileApi struct {
	Api
}

/*
NewFileApi - create new instance of FileApi.

Accepts http.Client and API endpoint as input parameters.
If `endpoint` is empty string, then default `StratoHiDriveAPIV21` value is used.
*/
func NewFileApi(client *http.Client, endpoint string) FileApi {
	api := NewApi(client, endpoint)
	return FileApi{api}
}

/*
GetFile - This method retrieves a given file from the HiDrive.

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
func (f FileApi) GetFile(ctx context.Context, params url.Values) (io.ReadCloser, error) {
	var res *http.Response

	{
		var err error
		if res, err = f.doGET(ctx, "file", params); err != nil {
			return nil, err
		}
	}

	if err := f.checkHTTPStatus(http.StatusOK, res); err != nil {
		return nil, err
	}

	return res.Body, nil
}

/*
UploadFile -This method can be used to create a new file and store uploaded content.

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

Returns [HiDriveObject] with information about uploaded file.
*/
func (f FileApi) UploadFile(ctx context.Context, params url.Values, fileBody io.ReadCloser) (*HiDriveObject, error) {
	var (
		res  *http.Response
		body []byte
	)

	{
		var err error
		if res, err = f.doPOST(ctx, "file", params, fileBody); err != nil {
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

/*
DeleteFile - Delete a given file.
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
func (f FileApi) DeleteFile(ctx context.Context, params url.Values) error {
	var res *http.Response

	{
		var err error
		if res, err = f.doDELETE(ctx, "file", params); err != nil {
			return err
		}
	}

	if err := f.checkHTTPStatus(http.StatusNoContent, res); err != nil {
		return err
	}

	return nil
}
