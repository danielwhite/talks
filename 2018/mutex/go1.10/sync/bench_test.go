package sync_test

import (
	"strconv"
	"testing"

	stdrand "math/rand"

	"whitehouse.id.au/talks/2018/mutex/go1.10/sync"
)

func BenchmarkRand(b *testing.B) {
	for i := 1; i < 256; i *= 2 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			benchmarkRand(b, i)
		})
	}
}

// START OMIT
func benchmarkRand(b *testing.B, n int) {
	b.SetParallelism(n)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Int()
		}
	})
}

// END OMIT

var rand = stdrand.New(&lockedSource{src: stdrand.NewSource(1)})

type lockedSource struct {
	lk  sync.Mutex
	src stdrand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}
