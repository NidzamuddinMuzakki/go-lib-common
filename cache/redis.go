package cache

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
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
	raw, err := json.Marshal(data.Value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, string(data.Key), raw, duration).Err()
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
		raw, err := json.Marshal(data.Value)
		if err != nil {
			return err
		}
		if err = pipe.Set(ctx, string(data.Key), raw, duration).Err(); err != nil {

			return err
		}
	}

	if _, err = pipe.Exec(ctx); err != nil {
		return err
	}

	return nil

}

func (r *Redis) BatchGet(ctx context.Context, keys []Key, dest any) error {
	pipeline := r.client.Pipeline()

	strCmds := make([]*redis.StringCmd, 0, len(keys))
	for _, key := range keys {
		strCmds = append(strCmds, pipeline.Get(ctx, string(key)))
	}

	if _, err := pipeline.Exec(ctx); err != nil {
		return err
	}

	switch v := dest.(type) {
	case map[string]struct{}:
		// only need its key is it available or not
		for idx, strCmd := range strCmds {
			if _, err := strCmd.Result(); err != nil && err != redis.Nil {
				return err
			} else {
				v[string(keys[idx])] = struct{}{}
			}
		}
	default:
		switch reflect.TypeOf(dest).Elem().Kind() {
		case reflect.Slice, reflect.Array:
			stringRes := make([]string, 0, len(strCmds))
			for _, strCmd := range strCmds {
				if res, err := strCmd.Result(); err != nil && err != redis.Nil {
					return err
				} else {
					stringRes = append(stringRes, res)
				}
			}

			stringJson := `[ ` + strings.Join(stringRes, ",") + ` ]`

			err := json.Unmarshal([]byte(stringJson), dest)
			if err != nil {
				return err
			}
		}

	}

	return nil

}
