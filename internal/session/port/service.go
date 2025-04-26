package port

import (
	"context"

	"github.com/melody-mood/internal/session/payload"
)

type ISessionService interface {
	GenerateSessionID(ctx context.Context) (res payload.SessionResponse, err error)
	GenerateAuthSpotifyURL(ctx context.Context, sessionID string) (string, error)
}
