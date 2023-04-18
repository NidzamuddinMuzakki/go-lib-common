package cache

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

type MemoryData struct {
	Value any
	TTL   time.Time
}

type InMemory struct {
	data map[Key]MemoryData
	mu   *sync.Mutex
}

var (
	ErrInMemNotFound = errors.New("inMemory: not found")
	ErrInMemExpired  = errors.New("inMemory: expired")
	ErrInMemCopy     = errors.New("inMemory: failed copying value to destination")
)

func NewInMemory() *InMemory {
	return &InMemory{
		data: make(map[Key]MemoryData),
		mu:   &sync.Mutex{},
	}
}

func (im *InMemory) Set(_ context.Context, data Data, duration time.Duration) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	im.data[data.Key] = MemoryData{
		Value: data.Value,
		TTL:   time.Now().Add(duration),
	}

	return nil
}

func (im *InMemory) Get(ctx context.Context, key Key, dest any) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	result, ok := im.data[key]
	if !ok {
		return ErrInMemNotFound
	}

	if result.TTL.Before(time.Now()) {
		err := im.Delete(ctx, key)
		if err != nil {
			return err
		}
		return ErrInMemExpired
	}

	err := copier.Copy(dest, result.Value)
	if err != nil {
		return ErrInMemCopy
	}

	return nil
}

func (im *InMemory) Delete(_ context.Context, key Key) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	delete(im.data, key)

	return nil
}

func (im *InMemory) BatchSet(_ context.Context, datas []Data, duration time.Duration) error {

	for _, data := range datas {
		im.mu.Lock()
		im.data[data.Key] = MemoryData{
			Value: data.Value,
			TTL:   time.Now().Add(duration),
		}
		im.mu.Unlock()
	}

	return nil
}

func (im *InMemory) BatchGet(ctx context.Context, keys []Key, dest any) error {
	switch v := dest.(type) {
	case map[string]struct{}:
		// only need its key is it available or not
		for _, key := range keys {
			var (
				err error
			)
			im.mu.Lock()
			if val, ok := im.data[key]; ok {

				if val.TTL.Before(time.Now()) {
					err = im.Delete(ctx, key)
					if err != nil {
						return err
					}
					continue
				}

				v[string(key)] = struct{}{}

			}
			im.mu.Unlock()

		}
	default:
		switch reflect.TypeOf(dest).Elem().Kind() {
		case reflect.Slice, reflect.Array:
			slicePtr := reflect.ValueOf(dest)
			sliceValuePtr := slicePtr.Elem()
			for _, key := range keys {
				var (
					err error
				)
				im.mu.Lock()
				if val, ok := im.data[key]; ok {

					if val.TTL.Before(time.Now()) {
						err = im.Delete(ctx, key)
						if err != nil {
							return err
						}
						continue
					}

					sliceValuePtr.Set(reflect.Append(sliceValuePtr, reflect.ValueOf(val.Value)))

				}
				im.mu.Unlock()

			}

		}

	}
	return nil
}

func (im *InMemory) Incr(ctx context.Context, key string) (*redis.IntCmd, error) {
	return nil, nil
}

func (im *InMemory) Expire(ctx context.Context, key string, ttl time.Duration) (*redis.BoolCmd, error) {
	return nil, nil
}
func (im *InMemory) Ttl(ctx context.Context, key string) (*redis.DurationCmd, error) {
	return nil, nil
}
func (im *InMemory) GetRedisInstance() *redis.Client {
	return nil
}

func (im *InMemory) SetNx(ctx context.Context, data Data, duration time.Duration) (isSuccessSet bool, err error) {
	if val, ok := im.data[data.Key]; ok {
		if val.TTL.Before(time.Now()) {
			err = im.Delete(ctx, data.Key)
			if err != nil {
				return false, err
			}
		} else {
			return false, nil
		}
	}

	err = im.Set(ctx, data, duration)
	if err != nil {
		return false, err
	}

	return true, nil
}
