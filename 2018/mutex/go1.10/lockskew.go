package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"whitehouse.id.au/talks/2018/mutex/go1.10/sync"
)

var n = flag.Int("n", 10, "number of locks to acquire")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		runtime.SetMutexProfileFraction(100)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	done := make(chan bool, 1)

	// START OMIT
	var mu sync.Mutex

	// Goroutine A
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				mu.Lock()                          // HL
				time.Sleep(100 * time.Microsecond) // HL
				mu.Unlock()                        // HL
			}
		}
	}()

	// Goroutine B
	for i := 0; i < *n; i++ {
		time.Sleep(100 * time.Microsecond) // HL
		start := time.Now()
		mu.Lock() // HL
		fmt.Println("done in", time.Since(start))
		mu.Unlock() // HL
	}
	// END OMIT
}
