package cache

import (
	"time"

	goCache "github.com/patrickmn/go-cache"
)

const (
	defaultCleanupInterval = DefaultExpiration * 2
)

type region struct {
	items *goCache.Cache
}

func newRegion(expiration time.Duration) *region {
	items := goCache.New(expiration, defaultCleanupInterval)
	return &region{items: items}
}
