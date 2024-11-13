package redislock

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redislock.Client
	locker *redislock.Lock
}

type IClient interface {
	Lock(
		ctx context.Context, key string, ttl time.Duration, opt *redislock.Options,
	) error
	Unlock(ctx context.Context) error
	Refresh(
		ctx context.Context, ttl time.Duration, opt *redislock.Options,
	) error
	GetTTL(ctx context.Context) (time.Duration, error)
	GetToken(ctx context.Context) string
	GetMetadata(ctx context.Context) string
}

func NewClient(redisClient *redis.Client) *Client {
	return &Client{
		client: redislock.New(redisClient),
	}
}

// Lock tries to obtain a new lock using a key with the given TTL.
// May return ErrNotObtained if not successful.
func (c *Client) Lock(
	ctx context.Context,
	key string,
	ttl time.Duration,
	opt *redislock.Options,
) error {
	var err error
	c.locker, err = c.client.Obtain(ctx, key, ttl, opt)
	return err
}

// Unlock manually releases the lock. May return ErrLockNotHeld.
func (c *Client) Unlock(ctx context.Context) (err error) {
	err = c.locker.Release(ctx)
	return err
}

// Refresh extends the lock with a new TTL. May return ErrNotObtained if refresh is unsuccessful.
func (c *Client) Refresh(
	ctx context.Context,
	ttl time.Duration,
	opt *redislock.Options,
) error {
	err := c.locker.Refresh(ctx, ttl, opt)
	return err
}

// GetTTL returns the remaining time-to-live. Returns 0 if the lock has expired.
func (c *Client) GetTTL(ctx context.Context) (time.Duration, error) {
	duration, err := c.locker.TTL(ctx)
	return duration, err
}

// GetToken returns the token value set by the lock.
func (c *Client) GetToken(ctx context.Context) string {
	return c.locker.Token()
}

// GetMetadata returns the metadata of the lock.
func (c *Client) GetMetadata(ctx context.Context) string {
	return c.locker.Metadata()
}
