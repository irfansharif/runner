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
	"runtime/metrics"
	_ "unsafe" // for go:linkname
)

// grunningnanos returns the running time observed by the current goroutine by
// linking to a private symbol in the runtime package.
//
//go:linkname grunningnanos runtime.grunningnanos
func grunningnanos() int64

// metricnanos returns the running time observed by the current goroutine, but
// doing it through an exported metric. It's an alternative to grunningnanos
// that doesn't require the go:linkname directive, though ~5x slower.
func metricnanos() int64 {
	const metric = "/goroutine/running:nanoseconds" // from the modified go runtime

	sample := []metrics.Sample{
		{Name: metric},
	}
	metrics.Read(sample)
	return int64(sample[0].Value.Uint64())
}
