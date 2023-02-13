//go:build integration
// +build integration

package go_hidrive

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/url"
	"os"
	"testing"
	"time"
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (c *ClosingBuffer) Close() (err error) {
	return
}

func TestFileApi_UploadFile(t *testing.T) {
	type fields struct {
		Api FileApi
	}
	type args struct {
		params   url.Values
		fileBody io.ReadCloser
	}

	client, err := createTestHTTPClient()
	if err != nil {
		t.Errorf("error setting up HTTP client: %s", err.Error())
		return
	}
	fileApi := NewFileApi(client, StratoHiDriveAPIV21)
	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *HiDriveObject
		wantErr bool
	}{
		{
			name:    "upload this file to /public",
			want:    nil,
			wantErr: false,
			args: args{
				params: NewParameters().SetFilePath(fmt.Sprintf("/public/%s", uuid.New().String())).Values,
				fileBody: func() io.ReadCloser {
					f, err := os.Open("go.mod")
					if err != nil {
						t.Errorf("error opening file go.mod for reading")
						return nil
					}

					return f
				}(),
			},
			fields: fields{Api: fileApi},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fileApi.UploadFile(ctx, tt.args.params, tt.args.fileBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFileApi_DeleteFile(t *testing.T) {
	type fields struct {
		Api FileApi
	}
	type args struct {
		params url.Values
	}

	client, err := createTestHTTPClient()
	if err != nil {
		t.Errorf("error setting up HTTP client: %s", err.Error())
		return
	}
	fileApi := NewFileApi(client, StratoHiDriveAPIV21)
	ctx := context.Background()

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
				params: NewParameters().SetPath(fmt.Sprintf("/public/%s", uuid.New().String())).Values,
			},
			fields: fields{Api: fileApi},
		},
		{
			name:    "delete existing file",
			wantErr: false,
			args: args{
				params: NewParameters().SetPath(func() string {
					buf := &ClosingBuffer{
						&bytes.Buffer{},
					}
					for i := 0; i < 500; i++ {
						rand.Seed(time.Now().UnixNano())
						rndByte := rand.Intn(127)
						_ = buf.WriteByte(byte(rndByte))
					}
					path := fmt.Sprintf("/public/%s", uuid.New().String())
					prm := NewParameters().SetFilePath(path).Values
					if _, err := fileApi.UploadFile(ctx, prm, buf); err != nil {
						return ""
					}

					return path
				}()).Values,
			},
			fields: fields{Api: fileApi},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fileApi.DeleteFile(ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
