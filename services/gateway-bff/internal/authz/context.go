package authz

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

var ErrUnauthenticated = errors.New("unauthenticated")

type ActorContext struct {
	ID          string
	DisplayName string
	Role        string
	Status      string
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (c Client) CurrentUser(ctx context.Context, source *http.Request, correlationID string) (ActorContext, error) {
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/auth/current", nil)
	if err != nil {
		return ActorContext{}, err
	}
	req.Header.Set("X-Correlation-Id", correlationID)
	for _, cookie := range source.Cookies() {
		req.AddCookie(cookie)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return ActorContext{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ActorContext{}, ErrUnauthenticated
	}
	var body struct {
		User ActorContext `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return ActorContext{}, err
	}
	if body.User.ID == "" || body.User.Status != "Active" {
		return ActorContext{}, ErrUnauthenticated
	}
	return body.User, nil
}
