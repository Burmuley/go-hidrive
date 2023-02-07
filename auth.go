package go_hidrive

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://my.hidrive.com/client/authorize",
	TokenURL:  "https://my.hidrive.com/oauth2/token",
	AuthStyle: 0,
}

type Authenticator struct {
	Token  *oauth2.Token
	Config oauth2.Config
}

func NewAuthenticator(clientId, clientSecret, redirectUrl, scopes string) *Authenticator {
	a := &Authenticator{
		Config: oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Endpoint:     Endpoint,
			RedirectURL:  redirectUrl,
			Scopes:       strings.Split(scopes, ","),
		},
	}

	return a
}

func (a *Authenticator) Exchange(code string) error {
	ctx := context.Background()
	token, err := a.Config.Exchange(ctx, code)

	if err != nil {
		return err
	}

	a.Token = token
	return nil
}

//func (a *Authenticator) RefreshToken() error {
//	if a.Token == nil {
//		return errors.Wrap(model.ErrorStorageAuth, "no initial token available")
//	}
//
//	ctx := context.Background()
//	cli := a.Config.Client(ctx, a.Token)
//}

func (a *Authenticator) Client() (*http.Client, error) {
	if a.Token == nil {
		return nil, fmt.Errorf("error creating authentication client: %w", ErrAuthNoToken)
	}

	return a.Config.Client(context.Background(), a.Token), nil
}

func (a *Authenticator) GetToken() *oauth2.Token {
	return a.Token
}

func (a *Authenticator) SetRefreshToken(token string) {
	a.Token.RefreshToken = token
}
