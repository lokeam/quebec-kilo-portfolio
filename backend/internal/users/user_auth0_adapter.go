package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/types"
)

// Auth0Adapter handles Auth0 Management API calls for user metadata
// Now uses AppContext for config and logging
//``
type Auth0Adapter struct {
	appCtx     *appcontext.AppContext
	HTTPClient *http.Client
}

func NewAuth0Adapter(appCtx *appcontext.AppContext) (*Auth0Adapter, error) {
	// Validate Auth0 config
	cfg := appCtx.Config.Auth0
	if cfg.Domain == "" {
		return nil, fmt.Errorf("auth0 domain is not configured")
	}
	if cfg.ClientID == "" {
		return nil, fmt.Errorf("auth0 client ID is not configured")
	}
	if cfg.ClientSecret == "" {
		return nil, fmt.Errorf("auth0 client secret is not configured")
	}
	if cfg.Audience == "" {
		return nil, fmt.Errorf("auth0 audience is not configured")
	}

	appCtx.Logger.Debug("Auth0Adapter: config validation passed", map[string]any{
		"domain": cfg.Domain,
		"clientID": cfg.ClientID,
		"audience": cfg.Audience,
	})

	return &Auth0Adapter{
		appCtx:     appCtx,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// getManagementToken fetches a Management API token using client credentials
func (a *Auth0Adapter) getManagementToken(ctx context.Context) (string, error) {
	cfg := a.appCtx.Config.Auth0
	url := fmt.Sprintf("https://%s/oauth/token", cfg.Domain)
	payload := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     cfg.ClientID,
		"client_secret": cfg.ClientSecret,
		"audience":      cfg.ManagementAudience,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		a.appCtx.Logger.Error("Auth0Adapter: failed to marshal token request", map[string]any{"error": err})
		return "", fmt.Errorf("failed to marshal token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		a.appCtx.Logger.Error("Auth0Adapter: failed to create token request", map[string]any{"error": err})
		return "", fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		a.appCtx.Logger.Error("Auth0Adapter: token request failed", map[string]any{"error": err})
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error response body for better debugging
		var errorBody []byte
		if resp.Body != nil {
			errorBody, _ = json.Marshal(resp.Body)
		}
		a.appCtx.Logger.Error("Auth0Adapter: token request failed", map[string]any{
			"status": resp.Status,
			"body":   string(errorBody),
		})
		return "", fmt.Errorf("auth0 token request failed with status %s: %s", resp.Status, string(errorBody))
	}

	var tokenResp types.Auth0TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		a.appCtx.Logger.Error("Auth0Adapter: failed to decode token response", map[string]any{"error": err})
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("auth0 returned empty access token")
	}

	return tokenResp.AccessToken, nil
}


// PatchUserMetadata updates user_metadata for a user in Auth0
func (a *Auth0Adapter) PatchUserMetadata(
	ctx context.Context,
	userID string,
	metadata map[string]any,
) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	if metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}

	token, err := a.getManagementToken(ctx)
	if err != nil {
		a.appCtx.Logger.Error("Auth0Adapter: failed to get management token", map[string]any{"error": err})
		return fmt.Errorf("failed to get management token: %w", err)
	}

	cfg := a.appCtx.Config.Auth0
	url := fmt.Sprintf("https://%s/api/v2/users/%s", cfg.Domain, userID)
	reqBody := types.Auth0UserMetadataPatchRequest{
		UserMetadata: metadata,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		a.appCtx.Logger.Error("Auth0Adapter: failed to marshal PATCH request", map[string]any{"error": err})
		return fmt.Errorf("failed to marshal PATCH request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"PATCH",
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		a.appCtx.Logger.Error("Auth0Adapter: failed to create PATCH request", map[string]any{"error": err})
		return fmt.Errorf("failed to create PATCH request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		a.appCtx.Logger.Error("Auth0Adapter: PATCH request failed", map[string]any{"error": err})
		return fmt.Errorf("PATCH request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error response body for better debugging
		var errorBody []byte
		if resp.Body != nil {
			errorBody, _ = json.Marshal(resp.Body)
		}
		a.appCtx.Logger.Error("Auth0Adapter: PATCH request failed", map[string]any{
			"status": resp.Status,
			"userID": userID,
			"body":   string(errorBody),
		})
		return fmt.Errorf("auth0 PATCH failed with status %s for user %s: %s", resp.Status, userID, string(errorBody))
	}

	// Optionally parse response if you want to use it
	// var patchResp types.Auth0UserMetadataPatchResponse
	// if err := json.NewDecoder(resp.Body).Decode(&patchResp); err != nil {
	//     return err
	// }

	a.appCtx.Logger.Info("Auth0Adapter: successfully updated user metadata", map[string]any{
		"userID": userID,
		"status": resp.Status,
	})
	return nil
}