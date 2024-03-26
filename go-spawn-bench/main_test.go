package main

import (
	"sync"
	"testing"
)

func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sum := 0

		for j := 0; j < 1000; j++ {
			sum += j
		}
	}

	b.ReportAllocs()
}

func BenchmarkSumPar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for j := 0; j < 50; j++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				sum := 0
				for j := 0; j < 20; j++ {
					sum += j
				}
			}()
		}

		wg.Wait()
	}

	b.ReportAllocs()
}

// ❯ go test -bench Sum
// goos: darwin
// goarch: arm64
// pkg: spawnbench
// BenchmarkSum-10       	3556826	      320.0 ns/op	      0 B/op	      0 allocs/op
// BenchmarkSumPar-10    	 126237	      9224 ns/op	    816 B/op	     51 allocs/op
// PASS
// ok  	spawnbench	3.024s
