package go_hidrive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HDDirApi struct {
	*HDApi
}

// GetDir retrieves directory information defined in the `path`.
func (a *HDDirApi) GetDir(path string, params map[string]string) (*HiDriveObject, error) {
	var (
		query url.Values
		req   *http.Request
		cli   *http.Client
		res   *http.Response
		body  []byte
	)

	{
		var err error
		if req, err = a.NewGETRequest("dir"); err != nil {
			return nil, err
		}
	}

	query = req.URL.Query()
	query.Add("path", path)
	for k, v := range params {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	{
		var err error
		if cli, err = a.Authenticator.Client(); err != nil {
			return nil, err
		}
	}

	{
		var err error
		if res, err = cli.Do(req); err != nil {
			return nil, err
		}
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	if res.StatusCode != http.StatusOK {
		hdErr := &HiDriveError{}
		if err := json.Unmarshal(body, hdErr); err != nil {
			return nil, err
		}
		return nil, hdErr
	}

	hdObj := &HiDriveObject{}
	if err := hdObj.UnmarshalJSON(body); err != nil {
		return nil, err
	}

	return hdObj, nil
}

// CreateDir creates directory defined in the `path`.
// This functions will not create all parent directories, they must exist.
func (a *HDDirApi) CreateDir(path string) (*HiDriveObject, error) {
	var (
		query url.Values
		req   *http.Request
		cli   *http.Client
		res   *http.Response
		body  []byte
	)

	{
		var err error
		if req, err = a.NewPOSTRequest("dir"); err != nil {
			return nil, err
		}
	}

	query = req.URL.Query()
	query.Add("path", path)
	req.URL.RawQuery = query.Encode()

	{
		var err error
		if cli, err = a.Authenticator.Client(); err != nil {
			return nil, err
		}
	}

	{
		var err error
		if res, err = cli.Do(req); err != nil {
			return nil, err
		}
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	if res.StatusCode != http.StatusCreated {
		hdErr := &HiDriveError{}
		if err := json.Unmarshal(body, hdErr); err != nil {
			return nil, err
		}
		return nil, hdErr
	}

	hdObj := &HiDriveObject{}
	if err := hdObj.UnmarshalJSON(body); err != nil {
		return nil, err
	}

	return hdObj, nil
}

// CreatePath creates directory defined in the `path` and all parent directories preceding it
func (a *HDDirApi) CreatePath(path string) (*HiDriveObject, error) {
	dirs := strings.Split(path, "/")
	for k := range dirs[:len(dirs)-1] {
		dir := fmt.Sprintf("/%s", strings.Join(dirs[1:k+1], "/"))
		if _, err := a.GetDir(dir, map[string]string{"members": "none", "fields": "path"}); err == nil {
			continue
		}
		if _, err := a.CreateDir(dir); err != nil {
			return nil, err
		}
	}

	return a.CreateDir(path)
}

// DeleteDir deletes the directory defined in the `path`
func (a *HDDirApi) DeleteDir(path string, recursive bool) error {
	var (
		query url.Values
		req   *http.Request
		cli   *http.Client
		res   *http.Response
		body  []byte
	)

	{
		var err error
		if req, err = a.NewDELETERequest("dir"); err != nil {
			return err
		}
	}

	query = req.URL.Query()
	query.Add("path", path)
	query.Add("recursive", fmt.Sprint(recursive))
	req.URL.RawQuery = query.Encode()

	{
		var err error
		if cli, err = a.Authenticator.Client(); err != nil {
			return err
		}
	}

	{
		var err error
		if res, err = cli.Do(req); err != nil {
			return err
		}
	}

	{
		var err error
		if body, err = io.ReadAll(res.Body); err != nil {
			return err
		}
	}

	if res.StatusCode != http.StatusNoContent {
		hdErr := &HiDriveError{}
		if err := json.Unmarshal(body, hdErr); err != nil {
			return err
		}
		return hdErr
	}

	return nil
}
