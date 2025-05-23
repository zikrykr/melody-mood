package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/melody-mood/config"
	"github.com/melody-mood/constants"
	"github.com/melody-mood/internal/session/payload"
	"github.com/melody-mood/internal/session/port"
	"github.com/melody-mood/pkg"
	"github.com/redis/go-redis/v9"
)

type SessionService struct {
	rds *redis.Client
}

func NewSessionService(rds *redis.Client) port.ISessionService {
	return SessionService{
		rds: rds,
	}
}

func (r SessionService) GenerateSessionID(ctx context.Context) (res payload.SessionResponse, err error) {
	sessionID := uuid.New().String()

	sessionCacheKey := fmt.Sprintf(constants.SESSION_CACHE_KEY, sessionID)
	sessionExpiresIn := time.Duration(constants.SESSION_EXPIRATION_TIME) * time.Second

	// Try read from cache
	err = r.rds.Set(ctx, sessionCacheKey, "authenticated", sessionExpiresIn).Err()
	if err != nil {
		return res, err
	}

	res.SessionID = sessionID
	res.ExpiresIn = constants.SESSION_EXPIRATION_TIME

	return res, nil
}

func (r SessionService) GenerateAuthSpotifyURL(ctx context.Context, sessionID string) (string, error) {
	conf := config.GetConfig()
	authURL := fmt.Sprintf(pkg.SPOTIFY_AUTH_URL, conf.Spotify.ClientID, conf.Spotify.RedirectURI, url.QueryEscape(pkg.SCOPE_CREATE_PLAYLIST), sessionID)

	return authURL, nil
}
