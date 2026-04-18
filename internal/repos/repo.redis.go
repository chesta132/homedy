package repos

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type redisRepo struct {
	rdb *redis.Client
}

func (r *redisRepo) hGetWithParse(ctx context.Context, key, field string, dest any) error {
	var dataStr string
	err := r.rdb.HGet(ctx, key, field).Scan(&dataStr)
	if err != nil {
		return err
	}
	if dataStr == "" {
		return redis.Nil
	}

	return json.Unmarshal([]byte(dataStr), &dest)
}

func (r *redisRepo) hSetWithParse(ctx context.Context, key string, sets map[string]any) error {
	setsStr := make(map[string]string, len(sets))
	for key, val := range sets {
		dataBytes, err := json.Marshal(val)
		if err != nil {
			return err
		}
		setsStr[key] = string(dataBytes)
	}
	return r.rdb.HSet(ctx, key).Err()
}
