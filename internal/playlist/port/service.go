package port

import (
	"context"

	"github.com/melody-mood/internal/playlist/payload"
)

type IPlaylistService interface {
	GeneratePlaylists(ctx context.Context, req payload.GeneratePlaylistReq) (res payload.GeneratePlaylistResp, err error)
	CreateUserSpotifyPlaylist(ctx context.Context, req payload.CreateUserSpotifyPlaylistReq) error
}
