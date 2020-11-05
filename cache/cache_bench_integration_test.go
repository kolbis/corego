package cache_test

import (
	"math/rand"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pingcap/go-ycsb/pkg/generator"
	tlecache "github.com/kolbis/corego/cache"
)

const (
	workloadSize     = 1 << 20
	regionSize       = 1 << 6
	miniumExpiration = 300
	maxExpiration    = 1000
)

// zipfKeyList will return an array of string following a Zipf distribution
// where some values are more offen than other like we would have in real life
func zipfKeyList(size int) []string {
	// To ensure repetition of keys in the array,
	// we are generating keys in the range from 0 to size/3.
	maxKey := int64(size) / 3

	// scrambled zipfian to ensure same keys are not together
	z := generator.NewScrambledZipfian(0, maxKey, generator.ZipfianConstant)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	keys := make([]string, size)
	for i := 0; i < size; i++ {
		keys[i] = strconv.Itoa(int(z.Next(r)))
	}

	return keys
}

// oneEntryKeyList will return a unique list of keys
func oneEntryKeyList(size int) []string {
	v := rand.Int() % (size / 3)
	s := strconv.Itoa(v)

	keys := make([]string, size)
	for i := 0; i < size; i++ {
		keys[i] = s
	}

	return keys
}

func randomExpiration() time.Duration {
	return time.Duration(rand.Intn(maxExpiration-miniumExpiration+1)+miniumExpiration) * time.Microsecond
}

func runCacheBenchmark(b *testing.B, keys, regions []string, refreshFunction tlecache.RefreshFunction, percentageWrites uint64, initializeCache bool) {
	b.ReportAllocs()

	keySize := len(keys)
	regionSize := len(regions)
	rc := uint64(0)
	keyMask := keySize - 1
	regionMask := regionSize - 1

	config := tlecache.DefaultConfig()
	cache := tlecache.NewCache(config)
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

func BenchmarkCache(b *testing.B) {
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
		refreshFunction tlecache.RefreshFunction
		percentageWrite uint64
		initializeCache bool
	}{
		{"Zipf keys", zipfKeys, zipfRegions, resfreshFunction, 25, true},
		{"One entry Key", oneEntryKeys, oneEntryRegions, resfreshFunction, 25, true},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			runCacheBenchmark(b, bm.keys, bm.regions, bm.refreshFunction, bm.percentageWrite, bm.initializeCache)
		})
	}
}
