package pkg

import (
	"context"
	"fmt"

	"github.com/melody-mood/constants"
	redis "github.com/redis/go-redis/v9"
)

func CheckValidSession(ctx context.Context, rds *redis.Client, sessionID string) (bool, error) {
	key := fmt.Sprintf(constants.SESSION_CACHE_KEY, sessionID)
	exists, err := rds.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
