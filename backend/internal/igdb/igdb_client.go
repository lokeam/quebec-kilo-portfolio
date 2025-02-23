package igdb

import (
	"net/http"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
)

type IGDBClient struct {
	baseURL        string
	clientID       string
	token          string
	httpClient     *http.Client
	logger         interfaces.Logger
	appContext     *appcontext.AppContext
}

// NewIGDBClient creates a new IGDB client.
func NewIGDBClient(appContext *appcontext.AppContext, token string) *IGDBClient {
	if appContext == nil {
		panic("appContext is nil")
	}
	if token == "" {
			panic("token is empty")
	}
	if appContext.Config.IGDB.ClientID == "" {
			panic("ClientID is empty")
	}
	if appContext.Config.IGDB.BaseURL == "" {
			panic("BaseURL is empty")
	}

	return &IGDBClient{
		appContext:     appContext,
		clientID:       appContext.Config.IGDB.ClientID,
		token:          token,
		httpClient:     &http.Client{},
		baseURL:        appContext.Config.IGDB.BaseURL,
		logger:         appContext.Logger,
	}
}

// SetHTTPClient allows setting a custom HTTP client for testing purposes.
func (c *IGDBClient) SetHTTPClient(client *http.Client) {
	if client == nil {
			// Fallback to default client if nil is passed
			client = &http.Client{}
	}
	c.httpClient = client
}
