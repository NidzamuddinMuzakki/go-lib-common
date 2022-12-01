package cache

import (
	"context"
	"sync"
	"time"

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
