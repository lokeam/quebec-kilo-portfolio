package interfaces

import "context"

type TokenRetriever interface {
	GetToken(
		ctx context.Context,
		clientID,
		clientSecret,
		authURL string,
		logger Logger,
	) (string, error)
}
