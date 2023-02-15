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

func TestShareApi_Invite(t *testing.T) {
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
	shareApi := NewShareApi(client, StratoHiDriveAPIV21)
	dirAPi := NewDirApi(client, StratoHiDriveAPIV21)
	ctx := context.Background()

	testFilePath := fmt.Sprintf("/public/test_dir_for_sharing_%s", uuid.New().String())
	var dirObj *HiDriveObject
	var shareObj *HiDriveShareObject

	{
		var err error
		if dirObj, err = dirAPi.CreateDir(ctx, NewParameters().SetPath(testFilePath).Values); err != nil {
			t.Errorf("CreateDir() error = %v", err)
			return
		}
	}

	{
		var err error
		if shareObj, err = shareApi.CreateShare(ctx, NewParameters().SetPath(dirObj.Path).SetPassword("test@123!").Values); err != nil {
			t.Errorf("CreateShare() error = %v", err)
			return
		}
	}
	fmt.Println(shareObj.ID)

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
		if err = dirAPi.DeleteDir(ctx, NewParameters().SetPath(testFilePath).Values); err != nil {
			t.Errorf("DeleteDir() error = %v", err)
			return
		}
	}

	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		s := ShareApi{
	//			Api: tt.fields.Api,
	//		}
	//		_, err := s.Invite(tt.args.ctx, tt.args.params)
	//		if (err != nil) != tt.wantErr {
	//			t.Errorf("Invite() error = %v, wantErr %v", err, tt.wantErr)
	//			return
	//		}
	//		//if !reflect.DeepEqual(got, tt.want) {
	//		//	t.Errorf("Invite() got = %v, want %v", got, tt.want)
	//		//}
	//	})
	//}
}
