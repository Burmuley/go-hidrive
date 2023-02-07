package go_hidrive

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (c *ClosingBuffer) Close() (err error) {
	return
}

func TestHDFileApi_UploadFile(t *testing.T) {
	type fields struct {
		HDApi *HDApi
	}
	type args struct {
		path     string
		fileBody io.ReadCloser
	}

	a, err := CreateTestAuthenticator()
	if err != nil {
		t.Errorf("error setting up authenticator: %s", err.Error())
		return
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *HiDriveObject
		wantErr bool
	}{
		{
			name:    "upload random file to /public",
			want:    nil,
			wantErr: false,
			args: args{
				path: fmt.Sprintf("/public/%s", uuid.New().String()),
				fileBody: func() io.ReadCloser {
					f, err := os.Open("go.mod")
					if err != nil {
						t.Errorf("error opening file go.mod for reading")
						return nil
					}

					return f
				}(),
			},
			fields: fields{HDApi: &HDApi{
				Authenticator: a,
				Endpoint:      DefaultEndpointPrefix,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &HDFileApi{
				HDApi: tt.fields.HDApi,
			}
			_, err := f.UploadFile(tt.args.path, tt.args.fileBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("UploadFile() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestHDFileApi_DeleteFile(t *testing.T) {
	type fields struct {
		HDApi *HDApi
	}
	type args struct {
		path string
	}

	a, err := CreateTestAuthenticator()
	if err != nil {
		t.Errorf("error setting up authenticator: %s", err.Error())
		return
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "delete non-existent file",
			wantErr: true,
			args: args{
				path: fmt.Sprintf("/public/%s", uuid.New().String()),
			},
			fields: fields{HDApi: &HDApi{
				Authenticator: a,
				Endpoint:      DefaultEndpointPrefix,
			}},
		},
		{
			name:    "delete existing file",
			wantErr: false,
			args: args{
				path: func() string {
					buf := &ClosingBuffer{
						&bytes.Buffer{},
					}
					for i := 0; i < 500; i++ {
						rand.Seed(time.Now().UnixNano())
						rndByte := rand.Intn(127)
						_ = buf.WriteByte(byte(rndByte))
					}
					f := &HDFileApi{
						HDApi: &HDApi{
							Authenticator: a,
							Endpoint:      DefaultEndpointPrefix,
						},
					}
					fname := fmt.Sprintf("/public/%s", uuid.New().String())
					if _, err := f.UploadFile(fname, buf); err != nil {
						return ""
					}

					return fname
				}(),
			},
			fields: fields{HDApi: &HDApi{
				Authenticator: a,
				Endpoint:      DefaultEndpointPrefix,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &HDFileApi{
				HDApi: tt.fields.HDApi,
			}
			if err := f.DeleteFile(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
