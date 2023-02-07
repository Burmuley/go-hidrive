package go_hidrive

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestHDDirApi_GetDir(t *testing.T) {
	type fields struct {
		Authenticator *Authenticator
		Endpoint      string
	}
	type args struct {
		path   string
		params map[string]string
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
			name:    "test /public directory",
			wantErr: false,
			want:    nil,
			args: args{
				path:   "/public",
				params: map[string]string{"members": "none", "fields": "path"},
			},
			fields: fields{
				Authenticator: a,
				Endpoint:      DefaultEndpointPrefix,
			},
		},
		{
			name:    "test root directory",
			wantErr: false,
			want:    nil,
			args: args{
				path:   "/",
				params: map[string]string{"members": "none", "fields": "path"},
			},
			fields: fields{
				Authenticator: a,
				Endpoint:      DefaultEndpointPrefix,
			},
		},
		{
			name:    "test non-existent directory",
			wantErr: true,
			want:    nil,
			args: args{
				path:   "/some_dir_that_does_not_exist_jhsbcv8374r",
				params: map[string]string{"members": "none", "fields": "path"},
			},
			fields: fields{
				Authenticator: a,
				Endpoint:      DefaultEndpointPrefix,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &HDDirApi{
				&HDApi{
					Authenticator: tt.fields.Authenticator,
					Endpoint:      tt.fields.Endpoint,
				},
			}
			_, err := a.GetDir(tt.args.path, tt.args.params)
			_ = a.DeleteDir(tt.args.path, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetDir() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestHDDirApi_DeleteDir(t *testing.T) {
	type fields struct {
		Authenticator *Authenticator
		Endpoint      string
	}
	type args struct {
		path      string
		recursive bool
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
			name:    "test non-existent directory",
			wantErr: true,
			args: args{
				path:      "/some_dir_that_does_not_exist_jhsbcv8374r",
				recursive: false,
			},
			fields: fields{
				Authenticator: a,
				Endpoint:      DefaultEndpointPrefix,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &HDDirApi{
				&HDApi{
					Authenticator: tt.fields.Authenticator,
					Endpoint:      tt.fields.Endpoint,
				},
			}
			if err := a.DeleteDir(tt.args.path, tt.args.recursive); (err != nil) != tt.wantErr {
				t.Errorf("DeleteDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHDDirApi_CreateDir(t *testing.T) {
	type fields struct {
		Authenticator *Authenticator
		Endpoint      string
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
		want    *HiDriveObject
		wantErr bool
	}{
		{
			name:    "test creating directory",
			wantErr: false,
			args: args{
				path: fmt.Sprintf("/public/%s", uuid.New().String()),
			},
			fields: fields{
				Authenticator: a,
				Endpoint:      DefaultEndpointPrefix,
			},
		},
		{
			name:    "test creating existing directory",
			wantErr: true,
			args: args{
				path: "/public",
			},
			fields: fields{
				Authenticator: a,
				Endpoint:      DefaultEndpointPrefix,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &HDDirApi{
				&HDApi{
					Authenticator: tt.fields.Authenticator,
					Endpoint:      tt.fields.Endpoint,
				},
			}
			_, err := a.CreateDir(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("CreateDir() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestHDDirApi_CreatePath(t *testing.T) {
	type fields struct {
		Authenticator *Authenticator
		Endpoint      string
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
		want    *HiDriveObject
		wantErr bool
	}{
		{
			name:    "test creating path",
			wantErr: false,
			args: args{
				path: fmt.Sprintf("/public/%s/%s/%s", uuid.New().String(), uuid.New().String(), uuid.New().String()),
			},
			fields: fields{
				Authenticator: a,
				Endpoint:      DefaultEndpointPrefix,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &HDDirApi{
				&HDApi{
					Authenticator: tt.fields.Authenticator,
					Endpoint:      tt.fields.Endpoint,
				},
			}
			_, err := a.CreatePath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("CreatePath() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
