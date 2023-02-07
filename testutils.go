package go_hidrive

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

func CreateTestAuthenticator() (*Authenticator, error) {
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

	a := NewAuthenticator(envVars["STRATO_CLIENT_ID"], envVars["STRATO_CLIENT_SECRET"], "", "admin,rw")
	a.Token = &oauth2.Token{
		TokenType:    "Bearer",
		RefreshToken: envVars["STRATO_REFRESH_TOKEN"],
	}

	return a, nil
}
