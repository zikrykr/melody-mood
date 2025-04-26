package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/melody-mood/config"
	redis "github.com/redis/go-redis/v9"
)

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type SpotifySearchResult struct {
	Tracks struct {
		Items []struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Album struct {
				Name   string `json:"name"`
				Images []struct {
					URL string `json:"url"`
				} `json:"images"`
				ReleaseDate string `json:"release_date"`
			} `json:"album"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"items"`
	} `json:"tracks"`
}

const (
	GRANT_TYPE_CLIENT_CREDENTIALS  = "client_credentials"
	SPOTIFY_ACCESS_TOKEN_CACHE_KEY = "spotify:access_token"
)

func GenerateSpotifyAccessToken(ctx context.Context) (SpotifyTokenResponse, error) {
	conf := config.GetConfig()

	form := url.Values{}
	form.Add("grant_type", GRANT_TYPE_CLIENT_CREDENTIALS)
	form.Add("client_id", conf.Spotify.ClientID)
	form.Add("client_secret", conf.Spotify.ClientSecret)

	var tokenResp SpotifyTokenResponse

	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return tokenResp, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return tokenResp, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return tokenResp, fmt.Errorf("spotify token error: %s", string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return tokenResp, fmt.Errorf("failed to decode response: %w", err)
	}

	return tokenResp, nil
}

func SaveSpotifyAccessToken(ctx context.Context, rds *redis.Client, token SpotifyTokenResponse) error {
	err := rds.Set(ctx, SPOTIFY_ACCESS_TOKEN_CACHE_KEY, token.AccessToken, time.Duration(token.ExpiresIn-600)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to save access token to Redis: %w", err)
	}
	return nil
}

func GetSpotifyAccessToken(ctx context.Context, rds *redis.Client) (string, error) {
	token, err := rds.Get(ctx, SPOTIFY_ACCESS_TOKEN_CACHE_KEY).Result()
	if err != nil {
		if err == redis.Nil {
			tokenResp, errToken := GenerateSpotifyAccessToken(ctx)
			if errToken != nil {
				return "", errToken
			}
			err = SaveSpotifyAccessToken(ctx, rds, tokenResp)
			if err != nil {
				return "", err
			}
			return tokenResp.AccessToken, nil
		}
		return "", fmt.Errorf("failed to get access token from Redis: %w", err)
	}
	return token, nil
}

func SpotifySearch(ctx context.Context, rds *redis.Client, query string) (SpotifySearchResult, error) {
	var res SpotifySearchResult
	accessToken, err := GetSpotifyAccessToken(ctx, rds)
	if err != nil {
		return res, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=1&offset=0&market=ID", url.QueryEscape(query)), nil)
	if err != nil {
		return res, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return res, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return res, fmt.Errorf("spotify search error: %s", string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return res, fmt.Errorf("failed to decode response: %w", err)
	}

	return res, nil
}
