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

// Package runner is a library that's able to precisely measure on-CPU time for
// goroutines.
package runner

import (
	"sync/atomic"
	"time"
)

// TODO(irfansharif):
//
// - Make use of buildtags. Use smoketests to check if things are up to snuff?
//   Currently we fails with an opaque link error if not using the right go SDK.
// - Write out steps to download patched runtime.
// - Provide a default singleton runner for instances where plumbing in a runner
//   is inconvenient.
// - Maintain tags/task names and some bounded sketch of runner history.

// TODO(irfansharif): We only accumulate on-CPU time at the end of the invoked
// goroutine; there's no "inflight" view of ops. This post-hoc capture might be
// unsuitable for long-running goroutines unless we provide an API to
// periodically record on-CPU time, and/or perform taskgroup[1]-style tracking
// within the runtime.
//
// [1]: https://github.com/cockroachdb/cockroach/pull/60589

// Runner is able to measure precise on-CPU time for the set of all goroutines
// run through it.
type Runner struct {
	nanos int64 // accessed atomically
}

// New returns a new Runner.
func New() *Runner {
	return &Runner{}
}

// Run the given function in a new goroutine.
func (r *Runner) Run(f func()) {
	go func() {
		f()
		atomic.AddInt64(&r.nanos, grunningnanos()) // record the running nanoseconds for the goroutine
	}()
}

// Duration returns the total running (i.e. on CPU) duration observed by all
// goroutines run under the group.
func (r *Runner) Duration() time.Duration {
	return time.Duration(atomic.LoadInt64(&r.nanos))
}
