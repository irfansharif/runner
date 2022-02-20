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

package runner_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/irfansharif/runner"
	"github.com/stretchr/testify/assert"
)

// TestRunDuration tests that run durations behaves as expected.
func TestRunDuration(t *testing.T) {
	for _, work := range []string{"busy", "lazy"} {
		work := work
		t.Run(fmt.Sprintf("loop=%s", work), func(t *testing.T) {
			var wg sync.WaitGroup
			r := runner.New()

			start := time.Now()
			for i := 0; i < 10; i++ {
				wg.Add(1)
				r.Run(func() {
					defer wg.Done()

					if work == "busy" {
						runner.TestingBusyFn()
					} else {
						runner.TestingLazyFn()
					}
				})
			}
			wg.Wait()

			walltime := time.Since(start)
			cputime := r.Duration()
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
