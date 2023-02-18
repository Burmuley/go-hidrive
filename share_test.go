//go:build integration
// +build integration

package go_hidrive

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"os"
	"testing"
)

func TestShare_Create(t *testing.T) {
	type fields struct {
		Api Share
	}
	type args struct {
		params url.Values
	}

	client, err := createTestHTTPClient()
	if err != nil {
		t.Errorf("error setting up HTTP client: %s", err.Error())
		return
	}
	shareApi := NewShare(client, StratoHiDriveAPIV21)
	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ShareObject
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
			s := tt.fields.Api
			_, err := s.Create(ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestShare_Invite(t *testing.T) {
	type fields struct {
		Api Api
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
	shareApi := NewShare(client, StratoHiDriveAPIV21)
	dirAPi := NewDir(client, StratoHiDriveAPIV21)
	ctx := context.Background()

	testFilePath := fmt.Sprintf("/public/test_dir_for_sharing_%s", uuid.New().String())
	var dirObj *Object
	var shareObj *ShareObject

	{
		var err error
		if dirObj, err = dirAPi.Create(ctx, NewParameters().SetPath(testFilePath).Values); err != nil {
			t.Errorf("Create() error = %v", err)
			return
		}
	}

	{
		var err error
		if shareObj, err = shareApi.Create(ctx, NewParameters().SetPath(dirObj.Path).SetPassword("test@123!").Values); err != nil {
			t.Errorf("Create() error = %v", err)
			return
		}
	}

	{
		var err error
		recip, ok := os.LookupEnv("STRATO_INVITE_EMAIL")
		if !ok || recip == "" {
			t.Error("Environment variable STRATO_INVITE_EMAIL is not properly set", err)
		}
		if _, err = shareApi.Invite(ctx,
			NewParameters().SetId(shareObj.ID).SetRecipient(recip).SetMsg("Test invitation from good people").Values,
		); err != nil {
			t.Errorf("Invite() error = %v", err)
		}
	}

	{
		var err error
		if err = dirAPi.Delete(ctx, NewParameters().SetPath(testFilePath).Values); err != nil {
			t.Errorf("Delete() error = %v", err)
			return
		}
	}
}
