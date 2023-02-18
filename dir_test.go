//go:build integration
// +build integration

package go_hidrive

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"testing"
)

func TestDir_Get(t *testing.T) {
	type fields struct {
		Api Dir
	}
	type args struct {
		params url.Values
	}
	client, err := createTestHTTPClient()
	if err != nil {
		t.Errorf("error setting up HTTP client: %s", err.Error())
		return
	}
	dirApi := NewDir(client, StratoHiDriveAPIV21)
	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Object
		wantErr bool
	}{
		{
			name:    "test /public directory",
			wantErr: false,
			want:    nil,
			args: args{
				params: NewParameters().SetPath("/public").SetMembers([]string{"none"}).SetFields([]string{"path"}).Values,
			},
			fields: fields{
				Api: dirApi,
			},
		},
		{
			name:    "test root directory",
			wantErr: false,
			want:    nil,
			args: args{
				params: NewParameters().SetPath("/").SetMembers([]string{"none"}).SetFields([]string{"path"}).Values,
			},
			fields: fields{Api: dirApi},
		},
		{
			name:    "test non-existent directory",
			wantErr: true,
			want:    nil,
			args: args{
				params: NewParameters().SetPath("/some_dir_that_does_not_exist_jhsbcv8374r").SetMembers([]string{"none"}).SetFields([]string{"path"}).Values,
			},
			fields: fields{Api: dirApi},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dirApi.Get(ctx, tt.args.params)
			ddParams := NewParameters().SetPath(tt.args.params.Get("path")).SetRecursive(false).Values
			_ = dirApi.Delete(ctx, ddParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDir_Delete(t *testing.T) {
	type fields struct {
		Api Dir
	}
	type args struct {
		params url.Values
	}

	client, err := createTestHTTPClient()
	if err != nil {
		t.Errorf("error setting up HTTP client: %s", err.Error())
		return
	}
	dirApi := NewDir(client, StratoHiDriveAPIV21)
	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test non-existent directory",
			wantErr: true,
			args: args{
				params: NewParameters().SetPath("/some_dir_that_does_not_exist_jhsbcv8374r").SetRecursive(false).Values,
			},
			fields: fields{Api: dirApi},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := dirApi.Delete(ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDir_Create(t *testing.T) {
	type fields struct {
		Api Dir
	}
	type args struct {
		params url.Values
	}

	client, err := createTestHTTPClient()
	if err != nil {
		t.Errorf("error setting up HTTP client: %s", err.Error())
		return
	}
	dirApi := NewDir(client, StratoHiDriveAPIV21)
	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Object
		wantErr bool
	}{
		{
			name:    "test creating directory",
			wantErr: false,
			args: args{
				params: NewParameters().SetPath(fmt.Sprintf("/public/%s", uuid.New().String())).Values,
			},
			fields: fields{Api: dirApi},
		},
		{
			name:    "test creating existing directory",
			wantErr: true,
			args: args{
				params: NewParameters().SetPath("/public").Values,
			},
			fields: fields{Api: dirApi},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dirApi.Create(ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDir_CreatePath(t *testing.T) {
	type fields struct {
		Api Dir
	}
	type args struct {
		params url.Values
	}

	client, err := createTestHTTPClient()
	if err != nil {
		t.Errorf("error setting up HTTP client: %s", err.Error())
		return
	}
	dirApi := NewDir(client, StratoHiDriveAPIV21)
	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Object
		wantErr bool
	}{
		{
			name:    "test creating path",
			wantErr: false,
			args: args{
				params: NewParameters().SetPath(fmt.Sprintf("/public/%s/%s/%s", uuid.New().String(), uuid.New().String(), uuid.New().String())).Values,
			},
			fields: fields{Api: dirApi},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dirApi.CreatePath(ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
