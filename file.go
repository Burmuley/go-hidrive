package go_hidrive

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HDFileApi struct {
	*HDApi
}

func (f *HDFileApi) GetFile(path string) (io.ReadCloser, error) {
	var (
		query url.Values
		req   *http.Request
		cli   *http.Client
		res   *http.Response
	)

	{
		var err error
		if req, err = f.NewGETRequest("file"); err != nil {
			return nil, err
		}
	}

	query = req.URL.Query()
	query.Add("path", path)
	req.URL.RawQuery = query.Encode()

	{
		var err error
		if cli, err = f.Authenticator.Client(); err != nil {
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
		var body []byte

		if res.StatusCode != http.StatusOK {
			hdErr := &HiDriveError{}
			if body, err = io.ReadAll(res.Body); err != nil {
				return nil, err
			}
			if err := json.Unmarshal(body, hdErr); err != nil {
				return nil, err
			}
			return nil, hdErr
		}
	}

	return res.Body, nil
}

func (f *HDFileApi) UploadFile(path string, fileBody io.ReadCloser) (*HiDriveObject, error) {
	var (
		query      url.Values
		req        *http.Request
		cli        *http.Client
		res        *http.Response
		body       []byte
		dir, fname string
	)

	{
		var err error
		if req, err = f.NewPOSTRequest("file"); err != nil {
			return nil, err
		}
	}

	elems := strings.Split(path, "/")
	fname = elems[len(elems)-1]
	dir = strings.Join(elems[:len(elems)-1], "/")
	query = req.URL.Query()
	query.Add("dir", dir)
	query.Add("name", fname)
	req.URL.RawQuery = query.Encode()
	req.Body = fileBody

	{
		var err error
		if cli, err = f.Authenticator.Client(); err != nil {
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

func (f *HDFileApi) DeleteFile(path string) error {
	var (
		query url.Values
		req   *http.Request
		cli   *http.Client
		res   *http.Response
		body  []byte
	)

	{
		var err error
		if req, err = f.NewDELETERequest("file"); err != nil {
			return err
		}
	}

	query = req.URL.Query()
	query.Add("path", path)
	req.URL.RawQuery = query.Encode()

	{
		var err error
		if cli, err = f.Authenticator.Client(); err != nil {
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
