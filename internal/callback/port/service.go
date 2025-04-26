package port

import (
	"context"
)

type ICallbackService interface {
	HandleSpotifyCallback(ctx context.Context, code, errMsg, state string) error
}
