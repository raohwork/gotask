// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import (
	"fmt"
	"testing"
	"time"
)

func (c *integralTest) BenchmarkRun(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		b.Run("3", c.bRunSync(3))
		b.Run("30", c.bRunSync(30))
		b.Run("100", c.bRunSync(100))
	})
	b.Run("async", func(b *testing.B) {
		b.Run("3", c.bRunAsync(3))
		b.Run("30", c.bRunAsync(30))
		b.Run("100", c.bRunAsync(100))
	})
	b.Run("syncLong", func(b *testing.B) {
		b.Run("3", c.bRunSyncLong(3))
		b.Run("30", c.bRunSyncLong(30))
		b.Run("100", c.bRunSyncLong(100))
	})
	b.Run("asyncLong", func(b *testing.B) {
		b.Run("3", c.bRunAsyncLong(3))
		b.Run("30", c.bRunAsyncLong(30))
		b.Run("100", c.bRunAsyncLong(100))
	})
}

func (c *integralTest) bRunSync(n int) func(*testing.B) {
	return func(b *testing.B) {
		r := c.New()
		for i := 0; i < n; i++ {
			r.QuickAdd(fmt.Sprintf("%d", i), nop)
		}
		r.Parse()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r.RunSync()
		}
	}
}

func (c *integralTest) bRunAsyncLong(n int) func(*testing.B) {
	return func(b *testing.B) {
		r := c.New()
		for i := 0; i < n; i++ {
			r.QuickAdd(fmt.Sprintf("%d", i), func() error {
				time.Sleep(5 * time.Millisecond)
				return nil
			})
		}
		r.Parse()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r.Run(6)
		}
	}
}

func (c *integralTest) bRunSyncLong(n int) func(*testing.B) {
	return func(b *testing.B) {
		r := c.New()
		for i := 0; i < n; i++ {
			r.QuickAdd(fmt.Sprintf("%d", i), func() error {
				time.Sleep(5 * time.Millisecond)
				return nil
			})
		}
		r.Parse()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r.RunSync()
		}
	}
}

func (c *integralTest) bRunAsync(n int) func(*testing.B) {
	return func(b *testing.B) {
		r := c.New()
		for i := 0; i < n; i++ {
			r.QuickAdd(fmt.Sprintf("%d", i), nop)
		}
		r.Parse()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r.Run(6)
		}
	}
}
