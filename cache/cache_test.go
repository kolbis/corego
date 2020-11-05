package cache_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/kolbis/corego/cache"
	tlecache "github.com/kolbis/corego/cache"
	"github.com/kolbis/corego/errors"
)

func Test_Cache(t *testing.T) {
	tests := map[string]struct {
		region     string
		key        string
		input      interface{}
		expiration time.Duration
		action     func(t *testing.T, cache *tlecache.Cache, regionName string, key string, input interface{}, expiration time.Duration)
		expected   interface{}
		err        error
	}{

		"Set": {
			region:     "region1",
			key:        "key1",
			input:      3,
			expiration: 1 * time.Second,
			action: func(t *testing.T, cache *tlecache.Cache, region, key string, input interface{}, expiration time.Duration) {
				cache.Set(region, key, input, expiration)
			},
			expected: 3,
			err:      nil,
		},
		"Set on default region": {
			region:     tlecache.DefaultRegion,
			key:        "key1",
			input:      3,
			expiration: 1 * time.Second,
			action: func(t *testing.T, cache *tlecache.Cache, region, key string, input interface{}, expiration time.Duration) {
				cache.SetDefault(key, input, expiration)
			},
			expected: 3,
			err:      nil,
		},
		"value not there, function called": {
			region:     "region1",
			key:        "key1",
			input:      3,
			expiration: 10 * time.Microsecond,
			action: func(t *testing.T, cache *tlecache.Cache, regionName string, key string, input interface{}, expiration time.Duration) {
				fn := func() (interface{}, error) {
					return 3, nil
				}
				_, _ = cache.GetOrCreate(regionName, key, expiration, fn)
			},
			expected: 3,
			err:      nil,
		},
		"value should be there, fn not called": {
			region:     "region1",
			key:        "key1",
			input:      3,
			expiration: 1 * time.Second,
			action: func(t *testing.T, cache *tlecache.Cache, regionName string, key string, input interface{}, expiration time.Duration) {
				cache.Set(regionName, key, input, expiration)
				fn := func() (interface{}, error) {
					return 4, nil
				}
				_, _ = cache.GetOrCreate(regionName, key, expiration, fn)
			},
			expected: 3,
			err:      nil,
		},
		"timeout": {
			region:     "region2",
			key:        "key2",
			input:      3,
			expiration: 1 * time.Second,
			action: func(t *testing.T, cache *tlecache.Cache, regionName string, key string, input interface{}, expiration time.Duration) {
				fn := func() (interface{}, error) {
					time.Sleep(cache.Timeout + time.Second)
					return 4, nil
				}
				v, err := cache.GetOrCreate(regionName, key, expiration, fn)
				assert.Equal(t, nil, v)
				assert.True(t, strings.Contains(err.Error(), "Timedout request"))
			},
			expected: nil,
			err:      errors.NewNotFoundError(errors.New("Get return error"), "Key key2 not found"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			config := tlecache.Config{
				Expiration:   tlecache.DefaultExpiration,
				Timeout:      1 * time.Second,
				JitterFactor: 1.0,
			}
			cache := tlecache.NewCache(config)

			tc.action(t, cache, tc.region, tc.key, tc.input, tc.expiration)

			v, err := cache.Get(tc.region, tc.key)
			assert.Equal(t, tc.expected, v)
			if tc.err != nil && err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			}
		})
	}
}

func TestJitter(t *testing.T) {
	dur := time.Duration(5)
	expectedDuration := time.Duration(4)
	duration := cache.Jitter(dur, 1)
	if duration != expectedDuration {
		t.Errorf("Jitter return %d; expect %d ", duration, expectedDuration)
	}
}
