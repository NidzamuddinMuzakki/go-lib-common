//go:generate mockery --name=Cacher
package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Key string

type Data struct {
	Key   Key `json:"key"`
	Value any `json:"value"`
}

type Cacher interface {
	GetRedisInstance() *redis.Client
	Set(ctx context.Context, data Data, duration time.Duration) error
	SetNx(ctx context.Context, data Data, duration time.Duration) (isSuccessSet bool, err error)
	Get(ctx context.Context, key Key, dest any) error
	Delete(ctx context.Context, key Key) error
	BatchSet(ctx context.Context, datas []Data, duration time.Duration) error
	BatchGet(ctx context.Context, keys []Key, dest any) error
	Incr(ctx context.Context, key string) (*redis.IntCmd, error)
	Expire(ctx context.Context, key string, ttl time.Duration) (*redis.BoolCmd, error)
}

type Driver string

// Drivers
const (
	InMemoryDriver = Driver("inMemory")
	RedisDriver    = Driver("redis")
)

type Cache struct {
	driver   *Driver
	host     string
	password string
	database string
}

type Option func(*Cache)

func WithDriver(driver Driver) Option {
	return func(c *Cache) {
		c.driver = &driver
	}
}

func WithHost(host string) Option {
	return func(c *Cache) {
		c.host = host
	}
}

func WithPassword(password string) Option {
	return func(c *Cache) {
		c.password = password
	}
}

func WithDatabase(db string) Option {
	return func(c *Cache) {
		c.database = db
	}
}

var (
	ErrDriverUnavailable = errors.New("cache: driver unavailable")
)

func NewCache(
	options ...Option,
) (Cacher, error) {
	c := Cache{}
	for _, option := range options {
		option(&c)
	}

	if c.driver == nil {
		return nil, ErrDriverUnavailable
	}

	switch *c.driver {
	case RedisDriver:
		db, err := strconv.Atoi(c.database)
		if err != nil {
			return nil, err
		}
		return NewRedis(c.host, c.password, db), nil
	case InMemoryDriver:
		return NewInMemory(), nil
	default:
		return nil, ErrDriverUnavailable
	}
}
