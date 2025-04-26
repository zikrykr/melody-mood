package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/melody-mood/constants"
	"github.com/melody-mood/internal/playlist/payload"
	"github.com/melody-mood/internal/playlist/port"
	"github.com/melody-mood/pkg"
	"github.com/redis/go-redis/v9"
	openai "github.com/sashabaranov/go-openai"
)

type PlaylistService struct {
	rds          *redis.Client
	openAIClient *openai.Client
}

func NewPlaylistService(openAIClient *openai.Client, rds *redis.Client) port.IPlaylistService {
	return PlaylistService{
		rds:          rds,
		openAIClient: openAIClient,
	}
}

func (r PlaylistService) GeneratePlaylists(ctx context.Context, req payload.GeneratePlaylistReq) (res payload.GeneratePlaylistResp, err error) {
	if len(req.PickedSongs) == 0 {
		return res, fmt.Errorf("no picked songs provided")
	}

	var (
		playlistTracks []payload.PlaylistTrack
		cacheKey       = fmt.Sprintf(constants.PLAYLIST_CACHE_KEY, req.SessionID)
	)
	// Try read from cache
	// Per user only can have one playlist
	cached, err := r.rds.Get(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		return res, err
	}

	if !req.IsRegenerate && cached != "" {
		playlistTracks, err = pkg.ParseToStruct[[]payload.PlaylistTrack](cached)
		if err != nil {
			return res, err
		}

		res.Tracks = playlistTracks

		return res, nil
	}

	// invalidate cache
	if err := r.rds.Del(ctx, cacheKey).Err(); err != nil {
		return res, err
	}

	playlistResp, err := pkg.CreateOpenAIMessage(ctx, r.openAIClient, fmt.Sprintf(constants.GeneratePlaylistPrompt, composePickedSongsPrompt(req.PickedSongs), req.Genre))
	if err != nil {
		return res, err
	}

	playlistTracks, err = pkg.ParseToStruct[[]payload.PlaylistTrack](playlistResp)
	if err != nil {
		return res, err
	}

	// search to spotify to retrieve the song data
	var composedPlaylistTracks []payload.PlaylistTrack
	for _, track := range playlistTracks {
		songData, err := pkg.SpotifySearch(ctx, r.rds, track.SongName)
		if err != nil {
			return res, err
		}

		if len(songData.Tracks.Items) > 0 {
			song := songData.Tracks.Items[0]
			composedPlaylistTracks = append(composedPlaylistTracks, payload.PlaylistTrack{
				SpotifyTrackID:  song.ID,
				SongName:        song.Name,
				SongArtist:      song.Artists[0].Name,
				SongAlbum:       song.Album.Name,
				ReleaseDate:     song.Album.ReleaseDate,
				SpotifyCoverArt: song.Album.Images[0].URL,
			})
		}
	}

	jsonData, err := json.Marshal(composedPlaylistTracks)
	if err != nil {
		return res, err
	}

	// Save to Redis
	err = r.rds.Set(ctx, cacheKey, jsonData, 2*time.Hour).Err()
	if err != nil {
		return res, err
	}

	res.Tracks = composedPlaylistTracks

	return res, nil
}

func composePickedSongsPrompt(pickedSongs []payload.PickedSongReq) (res string) {
	for _, song := range pickedSongs {
		res += fmt.Sprintf("%s by %s\n", song.SongName, song.SongArtist)
	}
	return res
}
