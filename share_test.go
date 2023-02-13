//go:build integration
// +build integration

package go_hidrive

import (
	"context"
	"net/url"
	"testing"
)

func TestShareApi_CreateShare(t *testing.T) {
	type fields struct {
		Api Api
	}
	type args struct {
		params url.Values
	}

	client, err := createTestHTTPClient()
	if err != nil {
		t.Errorf("error setting up HTTP client: %s", err.Error())
		return
	}
	shareApi := NewApi(client, StratoHiDriveAPIV21)
	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *HiDriveShareObject
		wantErr bool
	}{
		{
			name:    "Create share for non-existing directory",
			fields:  fields{Api: shareApi},
			wantErr: true,
			args: args{
				params: NewParameters().SetPath("/public/i_do_not_exist_skdjcb3974hfc.wvjh4o397fbn").Values,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ShareApi{
				Api: tt.fields.Api,
			}
			_, err := s.CreateShare(ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateShare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
