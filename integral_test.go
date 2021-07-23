// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import "testing"

type integralTest struct {
	New func() Runner
}

func nop() error { return nil }

func ensureAdd(r Runner, task Task) func(*testing.T) {
	return func(t *testing.T) {
		if err := r.Add(task); err != nil {
			t.Fatalf("failed to add task %s: %v", task.Name(), err)
		}
	}
}

func (c *integralTest) Test(t *testing.T) {
	t.Run("add", c.TestAdd)
	t.Run("parse", c.TestParse)
	t.Run("run", c.TestRun)

}

func (c *integralTest) Benchmark(b *testing.B) {
	b.Run("parse", c.BenchmarkParse)
	b.Run("run", c.BenchmarkRun)
}

func TestRunner(t *testing.T) {
	c := &integralTest{New: NewRunner}
	t.Run("integral", c.Test)
}

func BenchmarkRunner(b *testing.B) {
	c := &integralTest{New: NewRunner}
	b.Run("integral", c.Benchmark)
}
