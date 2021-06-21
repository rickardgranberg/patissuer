package pat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
)

const apiVersion = "6.1-preview.1"

var resourceAreaId, _ = uuid.Parse("951917ac-a960-4999-8464-e3f0aa25b381")
var locationId, _ = uuid.Parse("2c1965e2-28da-4aa1-bc64-b8a58537ae9b")

type PatClient struct {
	apiClient *azuredevops.Client
}

func NewPatClient(ctx context.Context, conn *azuredevops.Connection, baseUrl string) (*PatClient, error) {
	c, err := conn.GetClientByResourceAreaId(ctx, resourceAreaId)

	if err != nil {
		return nil, err
	}

	return &PatClient{
		apiClient: c,
	}, nil
}

func (t *PatClient) ListTokens(ctx context.Context) ([]PatToken, error) {
	routeValues := createRouteValues()
	resp, err := t.apiClient.Send(
		ctx,
		http.MethodGet,
		locationId,
		apiVersion,
		routeValues,
		nil,
		nil,
		"",
		azuredevops.MediaTypeApplicationJson,
		nil)
	if err != nil {
		return nil, err
	}

	var tokens PagedPatTokens
	if err := t.apiClient.UnmarshalBody(resp, &tokens); err != nil {
		return nil, err
	}

	return tokens.PatTokens, nil
}

func (t *PatClient) CreateToken(ctx context.Context, displayName string, scopes []string, validTo time.Time) (PatToken, error) {
	request := &PatTokenCreateRequest{
		AllOrgs:     false,
		DisplayName: displayName,
		Scope:       strings.Join(scopes, " "),
		ValidTo:     validTo.Format(time.RFC3339),
	}
	body, marshalErr := json.Marshal(request)
	if marshalErr != nil {
		return PatToken{}, marshalErr
	}
	routeValues := createRouteValues()
	resp, err := t.apiClient.Send(
		ctx,
		http.MethodPost,
		locationId,
		apiVersion,
		routeValues,
		nil,
		bytes.NewReader(body),
		azuredevops.MediaTypeApplicationJson,
		azuredevops.MediaTypeApplicationJson,
		nil)
	if err != nil {
		return PatToken{}, err
	}

	var patResult PatTokenResult
	if err := t.apiClient.UnmarshalBody(resp, &patResult); err != nil {
		return PatToken{}, err
	}

	if patResult.PatTokenError != None {
		return PatToken{}, errors.New(patResult.PatTokenError)
	}

	return patResult.PatToken, nil
}

func createRouteValues() map[string]string {
	routeValues := make(map[string]string)
	routeValues["area"] = "tokens"
	routeValues["resource"] = "pats"
	return routeValues
}
