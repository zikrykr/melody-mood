package pkg

import (
	"bytes"
	"context"
	"encoding/base64"
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
			URI   string `json:"uri"`
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

type GenerateSpotifyAccessTokenReq struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
}

type CreateSpotifyPlaylistReq struct {
	Name string `json:"name"`
}

type CreateSpotifyPlaylistResp struct {
	ID string `json:"id"`
}

type GetUserProfileResp struct {
	ID string `json:"id"`
}

const (
	GRANT_TYPE_CLIENT_CREDENTIALS = "client_credentials"
	GRANT_TYPE_AUTHORIZATION_CODE = "authorization_code"

	SPOTIFY_ACCESS_TOKEN_CACHE_KEY      = "spotify:access_token"
	SPOTIFY_ACCESS_TOKEN_USER_CACHE_KEY = "spotify:%s:access_token"        // by session ID
	SPOTIFY_PLAYLIST_USER_CACHE_KEY     = "spotify:playlist:%s:session:%s" // by spotify playlist ID, session ID

	SCOPE_CREATE_PLAYLIST = "playlist-modify-public playlist-modify-private playlist-read-private playlist-read-collaborative user-read-private user-read-email"

	SPOTIFY_AUTH_URL = "https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s"
)

func GenerateSpotifyAccessToken(ctx context.Context, payload GenerateSpotifyAccessTokenReq) (SpotifyTokenResponse, error) {
	conf := config.GetConfig()
	form := url.Values{}
	form.Add("grant_type", payload.GrantType)

	if payload.GrantType == GRANT_TYPE_CLIENT_CREDENTIALS {
		form.Add("client_id", conf.Spotify.ClientID)
		form.Add("client_secret", conf.Spotify.ClientSecret)
	}

	if payload.Code != "" {
		form.Add("code", payload.Code)
	}

	if payload.RedirectURI != "" {
		form.Add("redirect_uri", payload.RedirectURI)
	}

	var tokenResp SpotifyTokenResponse

	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return tokenResp, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if payload.GrantType == GRANT_TYPE_AUTHORIZATION_CODE {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(conf.Spotify.ClientID+":"+conf.Spotify.ClientSecret)))
	}

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
			tokenResp, errToken := GenerateSpotifyAccessToken(ctx, GenerateSpotifyAccessTokenReq{
				GrantType: GRANT_TYPE_CLIENT_CREDENTIALS,
			})
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

func GetUserSpotifyAccessToken(ctx context.Context, rds *redis.Client, sessionID string) (string, error) {
	token, err := rds.Get(ctx, fmt.Sprintf(SPOTIFY_ACCESS_TOKEN_USER_CACHE_KEY, sessionID)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("user has not authenticate their spotify account: %w", err)
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

func GetUserProfile(ctx context.Context, accessToken string) (GetUserProfileResp, error) {
	var res GetUserProfileResp

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.spotify.com/v1/me", nil)
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
		return res, fmt.Errorf("spotify get user profile error: %s", string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return res, fmt.Errorf("failed to decode response: %w", err)
	}

	return res, nil
}

func CreateSpotifyPlaylist(ctx context.Context, accessToken, userID, name string) (CreateSpotifyPlaylistResp, error) {
	var res CreateSpotifyPlaylistResp

	payload := map[string]string{
		"name":        name,
		"description": "🎶 MelodyMood - Music that feels you.",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return res, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.spotify.com/v1/users/"+userID+"/playlists", bytes.NewBuffer(jsonData))
	if err != nil {
		return res, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return res, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return res, fmt.Errorf("spotify create playlist error: %s", string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return res, fmt.Errorf("failed to decode response: %w", err)
	}

	return res, nil
}

func AddTracksToSpotifyPlaylist(ctx context.Context, accessToken, playlistID string, trackURIs []string) error {
	payload := map[string][]string{
		"uris": trackURIs,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.spotify.com/v1/playlists/"+playlistID+"/tracks", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("spotify add tracks to playlist error: %s", string(body))
	}

	return nil
}
