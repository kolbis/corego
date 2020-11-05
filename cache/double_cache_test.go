package cache_test

import (
	"strings"
	"testing"
	"time"

	"github.com/kolbis/corego/errors"

	"github.com/stretchr/testify/assert"

	cache "github.com/kolbis/corego/cache"
)

func Test_DoubleCache(t *testing.T) {
	tests := map[string]struct {
		region     string
		key        string
		expiration time.Duration
		input      interface{}
		action     func(t *testing.T, doubleCache *cache.DoubleCache, regionName string, key string, input interface{}, expiration time.Duration)
		isExpired  bool
		pValue     interface{}
		isPOk      bool
		bValue     interface{}
		isBOk      bool
	}{
		"Set": {
			region:     "1",
			key:        "1",
			expiration: 1 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *cache.DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)
			},
			isExpired: false,
			pValue:    3,
			isPOk:     true,
			bValue:    3,
			isBOk:     true,
		},
		"Set default": {
			region:     cache.DefaultRegion,
			key:        "1",
			expiration: 1 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *cache.DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.SetDefault(key, input, expiration)
			},
			isExpired: false,
			pValue:    3,
			isPOk:     true,
			bValue:    3,
			isBOk:     true,
		},
		"Get from primary": {
			region:     "1",
			key:        "1",
			expiration: 2 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *cache.DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)

				fn := func() (interface{}, error) {
					return 4, nil
				}
				v, err := doubleCache.GetOrCreate(region, key, cache.DefaultExpiration, fn)

				assert.Equal(t, input, v)
				assert.Equal(t, nil, err)
			},
			isExpired: false,
			pValue:    3,
			isPOk:     true,
			bValue:    3,
			isBOk:     true,
		},
		"Value expire in primary, Get from backup": {
			region:     "2",
			key:        "2",
			expiration: 20 * time.Millisecond,
			input:      3,
			action: func(t *testing.T, doubleCache *cache.DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)

				// We wait for expiration * (1 + jitter)
				wait := time.Duration((1 + cache.DefaultJitterFactor) * float64(expiration))
				time.Sleep(wait + time.Millisecond)
			},
			isExpired: true,
			pValue:    nil,
			isPOk:     false,
			bValue:    3,
			isBOk:     true,
		},
		"Value expire in primary, Get from backup and call backend to renew data": {
			region:     "2",
			key:        "2",
			expiration: 1 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *cache.DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)

				// We wait for expiration * (1 + jitter)
				wait := time.Duration((1 + cache.DefaultJitterFactor) * float64(expiration))
				time.Sleep(wait + time.Millisecond)

				fn := func() (interface{}, error) {
					return 4, nil
				}
				v, err := doubleCache.GetOrCreate(region, key, expiration, fn)
				assert.Equal(t, 3, v)
				assert.Equal(t, nil, err)
				time.Sleep(500 * time.Millisecond)
			},
			isExpired: true,
			pValue:    4,
			isPOk:     false,
			bValue:    4,
			isBOk:     true,
		},
		"Value expire in primary, Get from backup and call backend to renew data with timeout": {
			region:     "2",
			key:        "2",
			expiration: 1 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *cache.DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)

				// We wait for expiration * (1 + jitter)
				wait := time.Duration((1 + cache.DefaultJitterFactor) * float64(expiration))
				time.Sleep(wait + time.Millisecond)

				fn := func() (interface{}, error) {
					time.Sleep(cache.DefaultTimeout + 1)
					return 4, nil
				}
				_, _ = doubleCache.GetOrCreate(region, key, expiration, fn)
			},
			isExpired: true,
			pValue:    nil,
			isPOk:     false,
			bValue:    3,
			isBOk:     true,
		},
		"Value missing in both": {
			region:     "2",
			key:        "2",
			expiration: 2 * time.Second,
			action: func(t *testing.T, doubleCache *cache.DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				fn := func() (interface{}, error) {
					return 4, errors.New("Cannot find")
				}
				v, err := doubleCache.GetOrCreate(region, key, expiration, fn)
				assert.Equal(t, nil, v)
				assert.True(t, strings.Contains(err.Error(), "Cannot find"))
			},
			isExpired: true,
			pValue:    nil,
			isPOk:     false,
			bValue:    nil,
			isBOk:     false,
		},
		"Value expire in both": {
			region:     "2",
			key:        "2",
			expiration: 20 * time.Millisecond,
			input:      3,
			action: func(t *testing.T, doubleCache *cache.DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)

				// We wait until the item is expire in both backups
				// Expiration time * (backup factor + jitter factor )
				wait := time.Duration(float64(expiration) * (doubleCache.BackupExpirationFactor + cache.DefaultJitterFactor))
				time.Sleep(wait + time.Second)
			},
			isExpired: true,
			pValue:    nil,
			isPOk:     false,
			bValue:    nil,
			isBOk:     false,
		},
		"Region not found": {
			region:     "2",
			key:        "2",
			expiration: 2 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *cache.DoubleCache, region, key string, input interface{}, expiration time.Duration) {
			},
			isExpired: false,
			pValue:    nil,
			isPOk:     false,
			bValue:    nil,
			isBOk:     false,
		},
		"Item not found": {
			region:     "2",
			key:        "2",
			expiration: 2 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *cache.DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, "fake", input, expiration)
			},
			isExpired: false,
			pValue:    nil,
			isPOk:     false,
			bValue:    nil,
			isBOk:     false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			config := cache.DefaultDoubleCacheConfig()
			config.Timeout = 500 * time.Millisecond

			c := cache.NewDoubleCache(config)
			tc.action(t, c, tc.region, tc.key, tc.input, tc.expiration)

			v, ok := c.Get(tc.region, tc.key)
			if !tc.isExpired {
				assert.Equal(t, tc.pValue, v)
				assert.Equal(t, tc.isPOk, ok)
			} else {
				assert.Equal(t, tc.bValue, v)
				assert.Equal(t, tc.isBOk, ok)
			}
		})
	}
}
