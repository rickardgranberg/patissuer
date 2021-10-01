package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

type AuthClient struct {
	tenantId string
	clientId string
	client   public.Client
}

const (
	LoginMethodInteractive = "interactive"
	LoginMethodDeviceCode  = "devicecode"
	LoginMethodBearerToken = "bearertoken"
)

const loginHost = "login.microsoftonline.com"
const aadInstance = "https://" + loginHost + "/%s/v2.0"

var scopes = []string{"499b84ac-1321-427f-aa17-267ca6975798/user_impersonation"} //Constant value to target Azure DevOps. Do not change

func NewAuthClient(tenantId, clientId string) (*AuthClient, error) {
	cl := &AuthClient{
		tenantId: tenantId,
		clientId: clientId,
	}

	http := &http.Client{}
	client, err := public.New(clientId,
		public.WithHTTPClient(http),
		public.WithAuthority(fmt.Sprintf(aadInstance, tenantId)))

	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	cl.client = client

	return cl, nil
}

func (a *AuthClient) Login(ctx context.Context, method, token string) (string, error) {

	switch method {
	case LoginMethodBearerToken:
		return a.loginBearerToken(ctx, token)
	case LoginMethodInteractive:
		return a.loginInteractive(ctx)
	case LoginMethodDeviceCode:
		return a.loginDeviceCode(ctx)
	default:
		return "", fmt.Errorf("unsupported login method provided: %s", method)
	}
}

func (a *AuthClient) loginBearerToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("use of the '%s' login method requires a token to be provided", LoginMethodBearerToken)
	}
	return token, nil
}

func (a *AuthClient) loginInteractive(ctx context.Context) (string, error) {
	// This is done as an attempt to avoid DNS resolution errors during login:
	_, err := net.LookupHost(loginHost)

	if err != nil {
		return "", fmt.Errorf("name lookup error: %w", err)
	}

	accounts := a.client.Accounts()
	if len(accounts) > 0 {
		// Assuming the user wanted the first account
		userAccount := accounts[0]
		// found a cached account, now see if an applicable token has been cached
		result, err := a.client.AcquireTokenSilent(ctx, scopes, public.WithSilentAccount(userAccount))
		if err != nil {
			return "", fmt.Errorf("aquire token silent failed: %w", err)
		}
		return result.AccessToken, nil
	}

	result, err := a.client.AcquireTokenInteractive(ctx, scopes, public.WithRedirectURI("http://localhost"))
	if err != nil {
		return "", fmt.Errorf("aquire token interactive failed: %w", err)
	}
	return result.AccessToken, nil
}

func (a *AuthClient) loginDeviceCode(ctx context.Context) (string, error) {
	// This is done as an attempt to avoid DNS resolution errors during login:
	_, err := net.LookupHost(loginHost)

	if err != nil {
		return "", fmt.Errorf("name lookup error: %w", err)
	}

	accounts := a.client.Accounts()
	if len(accounts) > 0 {
		// Assuming the user wanted the first account
		userAccount := accounts[0]
		// found a cached account, now see if an applicable token has been cached
		result, err := a.client.AcquireTokenSilent(ctx, scopes, public.WithSilentAccount(userAccount))
		if err != nil {
			return "", err
		}
		return result.AccessToken, fmt.Errorf("aquire token silent failed: %w", err)
	}

	code, err := a.client.AcquireTokenByDeviceCode(ctx, scopes)
	if err != nil {
		return "", fmt.Errorf("aquire token device code failed: %w", err)
	}

	fmt.Println(code.Result.Message)
	result, err := code.AuthenticationResult(ctx)
	if err != nil {
		return "", fmt.Errorf("auth result failed: %w", err)
	}
	return result.AccessToken, nil
}
