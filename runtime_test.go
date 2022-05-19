// Copyright 2022 Irfan Sharif.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package runner

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestInternalDuration tests internal APIs measuring goroutine running time.
func TestInternalDuration(t *testing.T) {
	for _, work := range []string{"busy", "lazy"} {
		t.Run(fmt.Sprintf("loop=%s", work), func(t *testing.T) {
			tstart := time.Now()
			var totalnanos int64

			var wg sync.WaitGroup
			for g := 0; g < 10; g++ {
				wg.Add(1)
				go func(wg *sync.WaitGroup, total *int64) {
					start := grunningnanos()
					if work == "busy" {
						TestingBusyFn()
					} else {
						TestingLazyFn()
					}
					end := grunningnanos()

					atomic.AddInt64(total, end-start)
					wg.Done()
				}(&wg, &totalnanos)
			}
			wg.Wait()

			walltime := time.Since(tstart)
			cputime := time.Duration(totalnanos)
			mult := float64(cputime.Nanoseconds()) / float64(walltime.Nanoseconds())

			if work == "busy" {
				minexp := float64(runtime.GOMAXPROCS(-1)) - 1
				assert.Greaterf(t, mult, minexp,
					"expected multiplier > %f, got %f", minexp, mult)
			} else {
				maxexp := float64(0.1)
				assert.Lessf(t, mult, maxexp,
					"expected approximately zero multiplier, got %f", mult)
			}
		})
	}
}

// TestEquivalentGoroutines is a variant of the "parallel test" in
// https://github.com/golang/go/issues/36821. It tests whether goroutines that
// (should) spend the same amount of time on-CPU have similar measured on-CPU
// time.
func TestEquivalentGoroutines(t *testing.T) {
	mu := struct {
		sync.Mutex
		nanos map[int]int64
	}{}
	mu.nanos = make(map[int]int64)

	f := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()

		var sum int
		for i := 0; i < 500000000; i++ {
			sum -= i / 2
			sum *= i
			sum /= i/3 + 1
			sum -= i / 4
		}

		nanos := grunningnanos()
		mu.Lock()
		mu.nanos[id] = nanos
		mu.Unlock()
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		i := i // copy loop variable
		wg.Add(1)
		go f(&wg, i)
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	total := int64(0)
	for _, nanos := range mu.nanos {
		total += nanos
	}

	minexp, maxexp := float64(0.085), float64(0.115)
	for i, nanos := range mu.nanos {
		got := float64(nanos) / float64(total)

		assert.Greaterf(t, got, minexp,
			"expected proportion > %f, got %f", minexp, got)
		assert.Lessf(t, got, maxexp,
			"expected proportion < %f, got %f", maxexp, got)

		t.Logf("%d's got %0.2f%% of total time", i, got*100)
	}
}

// TestProportionalGoroutines is a variant of the "serial test" in
// https://github.com/golang/go/issues/36821. It tests whether goroutines that
// (should) spend a proportion of time on-CPU have proportionate measured on-CPU
// time.
func TestProportionalGoroutines(t *testing.T) {
	f := func(wg *sync.WaitGroup, v uint64, trip uint64, result *int64) {
		defer wg.Done()

		ret := v
		for i := trip; i > 0; i-- {
			ret += i
			ret = ret ^ (i + 0xcafebabe)
		}

		nanos := grunningnanos()
		atomic.AddInt64(result, nanos)
	}

	results := make([]int64, 10, 10)
	var wg sync.WaitGroup

	for iters := 0; iters < 10000; iters++ {
		for i := uint64(0); i < 10; i++ {
			i := i // copy loop variable
			wg.Add(1)
			go f(&wg, i+1, (i+1)*100000, &results[i])
		}
	}

	wg.Wait()

	total := int64(0)
	for _, result := range results {
		total += result
	}

	initial := float64(results[0]) / float64(total)
	maxdelta := float64(0.5)
	for i, result := range results {
		got := float64(result) / float64(total)
		mult := got / initial
		assert.InDelta(t, float64(i+1), mult, maxdelta)

		t.Logf("%d's got %0.2f%% of total time (%fx)", i, got*100, mult)
	}
}

// BenchmarkMetricNanos measures how costly it is to read the current
// goroutine's running time when going through the exported runtime metric.
func BenchmarkMetricNanos(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = metricnanos()
	}
}

// BenchmarkGRunningNanos measures how costly it is to read the current
// goroutine's running time when going through an internal runtime API.
func BenchmarkGRunningNanos(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = grunningnanos()
	}
}

// BenchmarkNanotime measures how costly it is to read the current nanotime.
func BenchmarkNanotime(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = nanotime()
	}
}

func TestingBusyFn() {
	j := 1
	for i := 0; i < 5000000000; i++ {
		j = j - i + 42
	}
}

func TestingLazyFn() {
	time.Sleep(time.Millisecond * 100)
}
