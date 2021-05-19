package auth

import (
	"context"
	"fmt"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

type AuthClient struct {
	tenantId  string
	clientId  string
	clientApp public.Client
}

const aadInstance = "https://login.microsoftonline.com/%s/v2.0"

var scopes = []string{"499b84ac-1321-427f-aa17-267ca6975798/user_impersonation"} //Constant value to target Azure DevOps. Do not change

func NewAuthClient(tenantId, clientId string) (*AuthClient, error) {
	cl := &AuthClient{
		tenantId: tenantId,
		clientId: clientId,
	}

	app, err := public.New(clientId,
		public.WithAuthority(fmt.Sprintf(aadInstance, tenantId)))
	if err != nil {
		return nil, err
	}

	cl.clientApp = app
	return cl, nil
}

func (a *AuthClient) Login(ctx context.Context) (string, error) {
	accounts := a.clientApp.Accounts()
	if len(accounts) > 0 {
		// Assuming the user wanted the first account
		userAccount := accounts[0]
		// found a cached account, now see if an applicable token has been cached
		result, err := a.clientApp.AcquireTokenSilent(ctx, scopes, public.WithSilentAccount(userAccount))
		if err != nil {
			return "", err
		}
		return result.AccessToken, nil
	}
	result, err := a.clientApp.AcquireTokenInteractive(ctx, scopes, public.WithRedirectURI("http://localhost"))
	if err != nil {
		return "", err
	}
	return result.AccessToken, nil
}
