package cache

import (
	"container/list"
	"context"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type Options struct {
	Capacity        int
	DefaultTTL      time.Duration
	CleanupInterval time.Duration
}

type Cache struct {
	mu         sync.Mutex
	capacity   int
	defaultTTL time.Duration
	items      map[string]*list.Element
	order      *list.List
	group      singleflight.Group

	cleanupInterval time.Duration
	stopCh          chan struct{}
	closeOnce       sync.Once
}

type entry struct {
	key       string
	value     any
	expiresAt time.Time
}

func New(options Options) *Cache {
	if options.Capacity < 0 {
		options.Capacity = 0
	}
	if options.CleanupInterval <= 0 {
		options.CleanupInterval = time.Minute
	}

	c := &Cache{
		capacity:        options.Capacity,
		defaultTTL:      options.DefaultTTL,
		items:           make(map[string]*list.Element),
		order:           list.New(),
		cleanupInterval: options.CleanupInterval,
		stopCh:          make(chan struct{}),
	}
	go c.cleanupLoop()
	return c
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	if key == "" {
		return
	}
	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		e := elem.Value.(*entry)
		e.value = value
		e.expiresAt = expiresAt
		c.order.MoveToFront(elem)
		return
	}

	e := &entry{key: key, value: value, expiresAt: expiresAt}
	elem := c.order.PushFront(e)
	c.items[key] = elem
	c.evictIfNeededLocked()
}

func (c *Cache) Get(key string) (any, bool) {
	if key == "" {
		return nil, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return nil, false
	}
	e := elem.Value.(*entry)
	if e.expiredAt(time.Now()) {
		c.removeElementLocked(elem)
		return nil, false
	}
	c.order.MoveToFront(elem)
	return e.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[key]; ok {
		c.removeElementLocked(elem)
	}
}

func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

func (c *Cache) GetOrLoad(ctx context.Context, key string, ttl time.Duration, loader func(context.Context) (any, error)) (any, error) {
	if value, ok := c.Get(key); ok {
		return value, nil
	}

	v, err, _ := c.group.Do(key, func() (any, error) {
		if value, ok := c.Get(key); ok {
			return value, nil
		}
		value, err := loader(ctx)
		if err != nil {
			return nil, err
		}
		c.Set(key, value, ttl)
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *Cache) Close() {
	c.closeOnce.Do(func() {
		close(c.stopCh)
	})
}

func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stopCh:
			return
		}
	}
}

func (c *Cache) removeExpired() {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, elem := range c.items {
		e := elem.Value.(*entry)
		if e.expiredAt(now) {
			c.removeElementLocked(elem)
		}
	}
}

func (c *Cache) evictIfNeededLocked() {
	if c.capacity <= 0 {
		return
	}
	for len(c.items) > c.capacity {
		back := c.order.Back()
		if back == nil {
			return
		}
		c.removeElementLocked(back)
	}
}

func (c *Cache) removeElementLocked(elem *list.Element) {
	e := elem.Value.(*entry)
	delete(c.items, e.key)
	c.order.Remove(elem)
}

func (e *entry) expiredAt(now time.Time) bool {
	if e.expiresAt.IsZero() {
		return false
	}
	return now.After(e.expiresAt)
}
