//Package cache provides caching utilities.
package cache

import (
	"sync"
	"time"

	"github.com/kolbis/corego/errors"
)

const (
	// DefaultExpiration : Time by default it will remain an item in the cache
	DefaultExpiration = 350 * time.Second

	// DefaultTimeout : Time maxium allow to take a refreshing function to take
	DefaultTimeout = 30 * time.Second

	// DefaultRegion : Region used when no region is specified
	DefaultRegion = "default"
)

// RefreshFunction function type that will renew the data
type RefreshFunction func() (interface{}, error)

// Cache used to store information in memory divided by regions
type Cache struct {
	regions      map[string]*region
	mu           sync.RWMutex
	Expiration   time.Duration
	Timeout      time.Duration
	JitterFactor float64
}

// Config used to config the cache
type Config struct {
	Expiration   time.Duration
	Timeout      time.Duration
	JitterFactor float64
}

// DefaultConfig returns a default configuration of the cache
func DefaultConfig() Config {
	return Config{
		Expiration:   DefaultExpiration,
		Timeout:      DefaultTimeout,
		JitterFactor: DefaultJitterFactor,
	}
}

// NewCache creates a new cache instance
func NewCache(config Config) *Cache {
	return &Cache{
		regions:      make(map[string]*region),
		Expiration:   Jitter(config.Expiration, config.JitterFactor),
		Timeout:      config.Timeout,
		JitterFactor: config.JitterFactor,
	}
}

// Get an item
func (c *Cache) Get(region, key string) (interface{}, error) {
	reg, ok := c.regions[region]
	if !ok {
		err := errors.New("Get return error")
		return nil, errors.NewNotFoundError(err, "Cannot find region "+region)
	}
	v, ok := reg.items.Get(key)
	if !ok {
		err := errors.New("Get return error")
		return nil, errors.NewNotFoundError(err, "Key "+key+" not found")
	}
	return v, nil
}

// GetDefault return an item from the default region
func (c *Cache) GetDefault(key string) (interface{}, error) {
	return c.Get(DefaultRegion, key)
}

// GetOrCreate looks for a value by region and key
// if do not find, call the function to generate the entry and assign it
func (c *Cache) GetOrCreate(
	region, key string,
	expiration time.Duration,
	refreshFunction RefreshFunction,
) (interface{}, error) {
	reg := c.getOrCreateRegion(region)
	value, ok := reg.items.Get(key)
	if ok {
		return value, nil
	}

	valueChan := make(chan interface{})
	errChan := make(chan error)
	go func() {
		value, err := refreshFunction()
		if err != nil {
			errChan <- err
		}
		valueChan <- value
	}()

	select {
	case value = <-valueChan:
		expiration = Jitter(expiration, c.JitterFactor)
		reg.items.Set(key, value, expiration)
		return value, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(c.Timeout):
		err := errors.New("GetOrCreate error")
		return nil, errors.NewTimeoutError(err, "Timedout request")
	}
}

// GetOrCreateDefault looks for a value by key
// If do not find, call the function to generate the entry and assign it
func (c *Cache) GetOrCreateDefault(key string, refreshFunction RefreshFunction) (interface{}, error) {
	return c.GetOrCreate(key, DefaultRegion, DefaultExpiration, refreshFunction)
}

// Set assign a value to a region and key for a period of time
func (c *Cache) Set(region, key string, value interface{}, expiration time.Duration) {
	reg := c.getOrCreateRegion(region)
	expiration = Jitter(expiration, c.JitterFactor)
	reg.items.Set(key, value, expiration)
}

// SetDefault assign a value to a key on the default region for a period of time
func (c *Cache) SetDefault(key string, value interface{}, expiration time.Duration) {
	c.Set(DefaultRegion, key, value, expiration)
}

// InvalidateRegion flush the entries of that region
func (c *Cache) InvalidateRegion(region string) error {
	reg, err := c.getRegion(region)
	if err != nil {
		return err
	}
	reg.items.Flush()
	return nil
}

// Invalidate flush all entries from all regions
func (c *Cache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.regions = make(map[string]*region)
}

func (c *Cache) getOrCreateRegion(region string) *region {
	c.mu.RLock()
	reg, found := c.regions[region]
	c.mu.RUnlock()
	if !found {
		reg = newRegion(c.Expiration)
		c.mu.Lock()
		c.regions[region] = reg
		c.mu.Unlock()
	}
	return reg
}

func (c *Cache) getRegion(region string) (*region, error) {
	c.mu.RLock()
	r, ok := c.regions[region]
	c.mu.RUnlock()
	if !ok {
		err := errors.New("Get return error")
		return nil, errors.NewNotFoundError(err, "Cannot find region "+region)
	}
	return r, nil
}
