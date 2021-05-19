package devops

import (
	"fmt"

	"github.com/rickardgranberg/patissuer/pkg/devops/auth"
)

type Client struct {
	authClient *auth.AuthClient
}

func NewClient(tenantId, clientId string) (*Client, error) {
	ac, err := auth.NewAuthClient(tenantId, clientId)

	if err != nil {
		return nil, err
	}
	return &Client{authClient: ac}, nil
}

func (c *Client) IssuePat(scopes []string) (string, error) {
	_, err := c.authClient.Login()

	if err != nil {
		return "", fmt.Errorf("failed to login: %w", err)
	}

	return "", nil
}
