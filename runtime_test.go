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

func TestingBusyFn() {
	j := 1
	for i := 0; i < 5000000000; i++ {
		j = j - i + 42
	}
}

func TestingLazyFn() {
	time.Sleep(time.Millisecond * 100)
}
