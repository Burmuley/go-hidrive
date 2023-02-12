package go_hidrive

//func TestShareApi_GetShare(t *testing.T) {
//	var (
//		shareApi *ShareApi
//		auth     *Authenticator
//	)
//	type fields struct {
//		Api *ShareApi
//	}
//	type args struct {
//		path   string
//		fields string
//	}
//
//	{
//		var err error
//		if auth, err = CreateTestAuthenticator(); err != nil {
//			t.Errorf("error setting up authenticator: %s", err.Error())
//			return
//		}
//		shareApi = NewShareApi(auth, DefaultEndpointPrefix)
//	}
//
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    []*HiDriveShareObject
//		wantErr bool
//	}{}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := shareApi.GetShare(tt.args.path, tt.args.fields)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetShare() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetShare() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
