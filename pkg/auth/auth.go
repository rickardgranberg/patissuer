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
}

const (
	LoginMethodInteractive = "interactive"
	LoginMethodDeviceCode  = "devicecode"
	LoginMethodBearerToken = "bearertoken"
)

const loginHost = "login.microsoftonline.com"
const aadInstance = "https://" + loginHost + "/%s"

var scopes = []string{"499b84ac-1321-427f-aa17-267ca6975798/user_impersonation"} //Constant value to target Azure DevOps. Do not change

func NewAuthClient(tenantId, clientId string) (*AuthClient, error) {
	cl := &AuthClient{
		tenantId: tenantId,
		clientId: clientId,
	}

	return cl, nil
}

func (a *AuthClient) Login(ctx context.Context, method, token string) (string, error) {

	cl, err := a.createPublicClient()
	if err != nil {
		return "", fmt.Errorf("error creating client: %w", err)
	}

	switch method {
	case LoginMethodBearerToken:
		return a.loginBearerToken(ctx, token)
	case LoginMethodInteractive:
		return a.loginInteractive(ctx, cl)
	case LoginMethodDeviceCode:
		return a.loginDeviceCode(ctx, cl)
	default:
		return "", fmt.Errorf("unsupported login method provided: %s", method)
	}
}

func (a *AuthClient) createPublicClient() (public.Client, error) {
	http := &http.Client{}
	return public.New(a.clientId,
		public.WithHTTPClient(http),
		public.WithAuthority(fmt.Sprintf(aadInstance, a.tenantId)))
}

func (a *AuthClient) loginBearerToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("use of the '%s' login method requires a token to be provided", LoginMethodBearerToken)
	}
	return token, nil
}

func (a *AuthClient) loginInteractive(ctx context.Context, cl public.Client) (string, error) {
	accounts, err := cl.Accounts(ctx)
	if err != nil {
		return "", fmt.Errorf("account fetch failed: %w", err)
	}
	if len(accounts) > 0 {
		// Assuming the user wanted the first account
		userAccount := accounts[0]
		// found a cached account, now see if an applicable token has been cached
		result, err := cl.AcquireTokenSilent(ctx, scopes, public.WithSilentAccount(userAccount))
		if err != nil {
			return "", fmt.Errorf("aquire token silent failed: %w", err)
		}
		return result.AccessToken, nil
	}
	result, err := cl.AcquireTokenInteractive(ctx, scopes, public.WithRedirectURI("http://localhost"))
	if err != nil {
		return "", fmt.Errorf("aquire token interactive failed: %w", err)
	}
	return result.AccessToken, nil
}

func (a *AuthClient) loginDeviceCode(ctx context.Context, cl public.Client) (string, error) {
	// This is done as an attempt to avoid DNS resolution errors during login:
	_, err := net.LookupHost(loginHost)

	if err != nil {
		return "", fmt.Errorf("name lookup error: %w", err)
	}

	accounts, err := cl.Accounts(ctx)
	if err != nil {
		return "", fmt.Errorf("account fetch failed: %w", err)
	}
	if len(accounts) > 0 {
		// Assuming the user wanted the first account
		userAccount := accounts[0]
		// found a cached account, now see if an applicable token has been cached
		result, err := cl.AcquireTokenSilent(ctx, scopes, public.WithSilentAccount(userAccount))
		if err != nil {
			return "", err
		}
		return result.AccessToken, fmt.Errorf("aquire token silent failed: %w", err)
	}

	code, err := cl.AcquireTokenByDeviceCode(ctx, scopes)
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
