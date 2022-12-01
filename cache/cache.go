package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Key string

type Data struct {
	Key   Key `json:"key"`
	Value any `json:"value"`
}

type Cacher interface {
	Set(ctx context.Context, data Data, duration time.Duration) error
	Get(ctx context.Context, key Key, dest any) error
	Delete(ctx context.Context, key Key) error
	BatchSet(ctx context.Context, datas []Data, duration time.Duration) error
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

func WitDriver(driver Driver) Option {
	return func(c *Cache) {
		c.driver = &driver
	}
}

func WitHost(host string) Option {
	return func(c *Cache) {
		c.host = host
	}
}

func WitPassword(password string) Option {
	return func(c *Cache) {
		c.password = password
	}
}

func WitDatabase(db string) Option {
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
