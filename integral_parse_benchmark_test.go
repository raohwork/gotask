// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import (
	"fmt"
	"testing"
)

func (c *integralTest) BenchmarkParse(b *testing.B) {
	b.Run("3", c.bParse(3))
	b.Run("30", c.bParse(30))
	b.Run("100", c.bParse(100))
}

func (c *integralTest) bParse(n int) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			r := c.New()
			for i := 0; i < n; i++ {
				r.QuickAdd(fmt.Sprintf("%d", i), nop)
			}

			b.StartTimer()
			r.Parse()
		}
	}
}
