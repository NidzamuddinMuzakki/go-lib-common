package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(host, password string, db int) *Redis {
	r := &Redis{}
	r.client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	return r
}

func (r *Redis) Set(ctx context.Context, data Data, duration time.Duration) error {
	return r.client.Set(ctx, string(data.Key), data.Value, duration).Err()
}

func (r *Redis) Get(ctx context.Context, key Key, dest any) error {
	result, err := r.client.Get(ctx, string(key)).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(result), dest)
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Delete(ctx context.Context, key Key) error {
	return r.client.Del(ctx, string(key)).Err()
}

func (r *Redis) BatchSet(ctx context.Context, datas []Data, duration time.Duration) error {
	pipe := r.client.Pipeline()
	var err error
	defer func() {
		if err != nil {
			_ = pipe.Close()
		}
	}()
	for _, data := range datas {
		if err = pipe.Set(ctx, string(data.Key), data.Value, duration).Err(); err != nil {

			return err
		}
	}

	if _, err = pipe.Exec(ctx); err != nil {
		return err
	}

	return nil

}
