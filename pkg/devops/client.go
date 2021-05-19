package devops

import (
	"context"
	"fmt"
	"time"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/rickardgranberg/patissuer/pkg/devops/pat"
)

type Client struct {
	orgUrl string
	token  string
}

func NewClient(orgUrl, token string) (*Client, error) {
	return &Client{orgUrl: orgUrl, token: token}, nil
}

func (c *Client) IssuePat(ctx context.Context, name string, scopes []string, validTo time.Time) (pat.PatToken, error) {

	conn := createConnection(c.orgUrl, c.token)
	client, err := pat.NewPatClient(ctx, conn, c.orgUrl)
	if err != nil {
		return pat.PatToken{}, err
	}

	t, err := client.CreateToken(ctx, name, scopes, validTo)

	if err != nil {
		return pat.PatToken{}, err
	}

	return t, nil
}

func (c *Client) ListPats(ctx context.Context) ([]pat.PatToken, error) {

	conn := createConnection(c.orgUrl, c.token)
	client, err := pat.NewPatClient(ctx, conn, c.orgUrl)
	if err != nil {
		return []pat.PatToken{}, err
	}

	t, err := client.ListTokens(ctx)

	if err != nil {
		return []pat.PatToken{}, err
	}

	return t, nil
}

func createConnection(orgUrl, token string) *azuredevops.Connection {
	conn := azuredevops.NewAnonymousConnection(orgUrl)
	conn.AuthorizationString = fmt.Sprintf("Bearer %s", token)
	return conn
}
