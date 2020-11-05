package cache_test

import (
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	cache "github.com/kolbis/corego/cache"
)

func runDoubleCacheBenchmark(b *testing.B, keys, regions []string, refreshFunction cache.RefreshFunction, percentageWrites uint64, initializeCache bool) {
	b.ReportAllocs()

	keySize := len(keys)
	regionSize := len(regions)
	rc := uint64(0)
	keyMask := keySize - 1
	regionMask := regionSize - 1

	config := cache.DefaultDoubleCacheConfig()
	cache := cache.NewDoubleCache(config)

	if initializeCache {
		for i := 0; i < keySize; i++ {
			for j := 0; j < regionSize; j++ {
				cache.Set(regions[j], keys[i], "data", randomExpiration())
			}
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		keyIndex := rand.Int() & keyMask
		regionIndex := rand.Int() & regionMask

		mc := atomic.AddUint64(&rc, 1)

		if percentageWrites*mc/100 != percentageWrites*(mc-1)/100 {
			for pb.Next() {
				cache.Set(regions[regionIndex&regionMask], keys[keyIndex&keyMask], "data", randomExpiration())
				regionIndex++
				keyIndex++
			}
		} else {
			for pb.Next() {
				_, _ = cache.GetOrCreate(
					regions[regionIndex&regionMask],
					keys[keyIndex&keyMask],
					randomExpiration(),
					refreshFunction,
				)
				keyIndex++
				regionIndex++
			}
		}
	})
}

func BenchmarkDoubleCache(b *testing.B) {
	zipfKeys := zipfKeyList(workloadSize)
	zipfRegions := zipfKeyList(regionSize)

	oneEntryKeys := oneEntryKeyList(workloadSize)
	oneEntryRegions := oneEntryKeyList(regionSize)

	minSecondTime := 1
	maxSecondTime := 30

	resfreshFunction := func() (interface{}, error) {
		// We simulate a request time range [1...30] seconds
		randTime := rand.Intn(maxSecondTime-minSecondTime+1) + minSecondTime
		t := time.NewTicker(time.Duration(randTime) * time.Second)

		return t.C, nil
	}

	benchmarks := []struct {
		name            string
		keys            []string
		regions         []string
		refreshFunction cache.RefreshFunction
		percentageWrite uint64
		initializeCache bool
	}{
		{"Zipf keys", zipfKeys, zipfRegions, resfreshFunction, 25, true},
		{"One entry Key", oneEntryKeys, oneEntryRegions, resfreshFunction, 25, true},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			runDoubleCacheBenchmark(b, bm.keys, bm.regions, bm.refreshFunction, bm.percentageWrite, bm.initializeCache)
		})
	}
}
