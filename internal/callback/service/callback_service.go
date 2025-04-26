package service

import (
	"context"
	"fmt"
	"time"

	"github.com/melody-mood/internal/callback/port"
	"github.com/melody-mood/pkg"
	"github.com/redis/go-redis/v9"
)

type CallbackService struct {
	rds *redis.Client
}

func NewCallbackService(rds *redis.Client) port.ICallbackService {
	return CallbackService{
		rds: rds,
	}
}

func (r CallbackService) HandleSpotifyCallback(ctx context.Context, code, errMsg, state string) error {
	if errMsg != "" {
		return fmt.Errorf("error in callback: %s", errMsg)
	}

	sessionID := state
	if sessionID == "" {
		return fmt.Errorf("session ID is empty")
	}

	// validate session ID
	valid, err := pkg.CheckValidSession(ctx, r.rds, sessionID)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("invalid session ID")
	}

	// generate user access token
	accessToken, err := pkg.GenerateSpotifyAccessToken(ctx, pkg.GenerateSpotifyAccessTokenReq{
		GrantType:   pkg.GRANT_TYPE_AUTHORIZATION_CODE,
		Code:        code,
		RedirectURI: "http://localhost:8080/v1/callback",
	})
	if err != nil {
		return err
	}

	// Save to Redis
	err = r.rds.Set(ctx, fmt.Sprintf(pkg.SPOTIFY_ACCESS_TOKEN_USER_CACHE_KEY, sessionID), accessToken.AccessToken, time.Duration(accessToken.ExpiresIn-600)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}
