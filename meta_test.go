//go:build integration
// +build integration

package go_hidrive

import (
	"context"
	"net/url"
	"testing"
)

func TestMeta_Get(t *testing.T) {
	type fields struct {
		Api Meta
	}
	type args struct {
		ctx    context.Context
		params url.Values
	}

	client, err := createTestHTTPClient()
	if err != nil {
		t.Errorf("error setting up HTTP client: %s", err.Error())
		return
	}
	metaApi := NewMeta(client, StratoHiDriveAPIV21)
	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Object
		wantErr bool
	}{
		{
			name:    "Simple test for /public directory",
			wantErr: false,
			fields:  fields{Api: metaApi},
			args: args{
				ctx:    ctx,
				params: NewParameters().SetPath("/public").SetFields([]string{"rshare"}).Values,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.fields.Api
			_, err := m.Get(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
