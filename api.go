package go_hidrive

import (
	"io"
	"net/http"
	"strings"
)

var DefaultEndpointPrefix string = "https://api.hidrive.strato.com/2.1"

type Api interface {
	NewHTTPRequest(method, uri string, r io.Reader) (*http.Request, error)
	NewGETRequest(uri string) (*http.Request, error)
	NewPOSTRequest(uri string) (*http.Request, error)
	NewDELETERequest(uri string) (*http.Request, error)
	NewPUTRequest(uri string) (*http.Request, error)
}

type DirApi interface {
	Api
	GetDir(path string, params map[string]string) (*HiDriveObject, error)
	CreateDir(path string) (*HiDriveObject, error)
	CreatePath(path string) (*HiDriveObject, error)
	DeleteDir(path string, recursive bool) error
}

type FileApi interface {
	GetFile(path string) (io.ReadCloser, error)
	UploadFile(path string, fileBody io.ReadCloser) (*HiDriveObject, error)
	DeleteFile(path string) error
}

type HDApi struct {
	Authenticator *Authenticator
	Endpoint      string
}

func NewApi(authenticator *Authenticator, endpoint string) *HDApi {
	if endpoint == "" {
		endpoint = DefaultEndpointPrefix
	}
	return &HDApi{
		Authenticator: authenticator,
		Endpoint:      endpoint,
	}
}

func (a *HDApi) NewHTTPRequest(method, uri string, r io.Reader) (*http.Request, error) {
	return http.NewRequest(method, strings.Join([]string{a.Endpoint, uri}, "/"), r)
}

func (a *HDApi) NewGETRequest(uri string) (*http.Request, error) {
	return http.NewRequest("GET", strings.Join([]string{a.Endpoint, uri}, "/"), nil)
}

func (a *HDApi) NewPOSTRequest(uri string) (*http.Request, error) {
	return http.NewRequest("POST", strings.Join([]string{a.Endpoint, uri}, "/"), nil)
}

func (a *HDApi) NewDELETERequest(uri string) (*http.Request, error) {
	return http.NewRequest("DELETE", strings.Join([]string{a.Endpoint, uri}, "/"), nil)
}

func (a *HDApi) NewPUTRequest(uri string) (*http.Request, error) {
	return http.NewRequest("DELETE", strings.Join([]string{a.Endpoint, uri}, "/"), nil)
}
