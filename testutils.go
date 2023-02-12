package go_hidrive

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

func createTestHTTPClient() (*http.Client, error) {
	envVars := map[string]string{
		"STRATO_CLIENT_ID":     "",
		"STRATO_CLIENT_SECRET": "",
		"STRATO_REFRESH_TOKEN": "",
	}

	for k := range envVars {
		ok := false
		val := ""
		if val, ok = os.LookupEnv(k); !ok {
			return nil, fmt.Errorf("missing %s value", k)
		}
		envVars[k] = val
	}

	oa2config := oauth2.Config{
		ClientID:     envVars["STRATO_CLIENT_ID"],
		ClientSecret: envVars["STRATO_CLIENT_SECRET"],
		Endpoint: oauth2.Endpoint{
			AuthURL:   StratoHiDriveAuthURL,
			TokenURL:  StratoHiDriveTokenURL,
			AuthStyle: 0,
		},
		Scopes: []string{"admin", "rw"},
	}
	token := &oauth2.Token{
		RefreshToken: envVars["STRATO_REFRESH_TOKEN"],
	}
	client := oa2config.Client(context.Background(), token)
	return client, nil
}
