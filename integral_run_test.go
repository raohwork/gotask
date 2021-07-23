// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func (c *integralTest) TestRun(t *testing.T) {
	t.Run("async", c.runAsync)
	t.Run("asyncEmpty", c.runAsyncEmpty)
	t.Run("asyncErrors", c.runAsyncErrors)
	t.Run("asyncAfterReset", c.runAsyncAfterReset)
	t.Run("async2Level", c.runAsync2level)
	t.Run("syncErr", c.runSyncErr)
	t.Run("sync", c.runSync)
	t.Run("syncEmpty", c.runSyncEmpty)
}

func (ct *integralTest) runAsyncEmpty(t *testing.T) {
	r := ct.New()
	if err := r.Run(0); err != nil {
		t.Fatal("unexpected error", err)
	}
}

func (ct *integralTest) runAsyncErrors(t *testing.T) {
	r := ct.New()
	errf := func(n string) func() error {
		return func() error {
			time.Sleep(10 * time.Millisecond)
			return errors.New(n)
		}
	}
	t.Run("addTaskA", ensureAdd(r, New("a", errf("a"), []string{})))
	t.Run("addTaskB", ensureAdd(r, New("b", errf("b"), []string{})))

	err := r.Run(2)
	if s := err.Error(); s != "a" && s != "b" {
		t.Fatal("unexpected error", err)
	}
	errs := r.Errors()
	if l := len(errs); l != 2 {
		t.Fatalf("expected 2 errors, get %+v", errs)
	}
	a, b := errs[0].Error(), errs[1].Error()
	if a == b {
		t.Fatal("expected 2 different errors, got two of", a)
	}
	if a != "a" && a != "b" {
		t.Fatal("unexpected error[0]", a)
	}
	if b != "a" && b != "b" {
		t.Fatal("unexpected error[1]", b)
	}
}

func (ct *integralTest) runAsync(t *testing.T) {
	r := ct.New()
	const (
		init int32 = iota
		started
		done
	)
	x, y, z := init, init, init
	a, b, c := &x, &y, &z
	t.Run("addTaskA", ensureAdd(r, New(
		"a",
		func() error {
			atomic.StoreInt32(a, started)
			time.Sleep(50 * time.Millisecond)
			atomic.StoreInt32(a, done)
			return nil
		},
		[]string{},
	)))
	t.Run("addTaskb", ensureAdd(r, New(
		"b",
		func() error {
			atomic.StoreInt32(b, started)
			time.Sleep(50 * time.Millisecond)
			atomic.StoreInt32(b, done)
			return nil
		},
		[]string{},
	)))
	t.Run("addTaskA", ensureAdd(r, New(
		"c",
		func() error {
			atomic.StoreInt32(c, started)
			time.Sleep(50 * time.Millisecond)
			atomic.StoreInt32(c, done)
			return nil
		},
		[]string{},
	)))

	go r.Run(3)
	time.Sleep(20 * time.Millisecond)
	if atomic.LoadInt32(a) != started {
		t.Fatal("a is not started at 20ms")
	}
	if atomic.LoadInt32(b) != started {
		t.Fatal("b is not started at 20ms")
	}
	if atomic.LoadInt32(c) != started {
		t.Fatal("c is not started at 20ms")
	}
	time.Sleep(35 * time.Millisecond)

	if atomic.LoadInt32(a) != done {
		t.Fatal("a is not done at 55ms")
	}
	if atomic.LoadInt32(b) != done {
		t.Fatal("b is not done at 55ms")
	}
	if atomic.LoadInt32(c) != done {
		t.Fatal("c is not done at 55ms")
	}
}

func (ct *integralTest) runAsyncAfterReset(t *testing.T) {
	r := ct.New()
	const (
		init int32 = iota
		started
		done
	)
	x, y, z := init, init, init
	a, b, c := &x, &y, &z
	t.Run("addTaskA", ensureAdd(r, New(
		"a",
		func() error {
			atomic.StoreInt32(a, started)
			time.Sleep(50 * time.Millisecond)
			atomic.StoreInt32(a, done)
			return nil
		},
		[]string{},
	)))
	t.Run("addTaskb", ensureAdd(r, New(
		"b",
		func() error {
			atomic.StoreInt32(b, started)
			time.Sleep(50 * time.Millisecond)
			atomic.StoreInt32(b, done)
			return nil
		},
		[]string{},
	)))
	t.Run("addTaskA", ensureAdd(r, New(
		"c",
		func() error {
			atomic.StoreInt32(c, started)
			time.Sleep(50 * time.Millisecond)
			atomic.StoreInt32(c, done)
			return nil
		},
		[]string{},
	)))

	if err := r.Run(3); err != nil {
		t.Fatal("unexpected error: ", err)
	}
	r.Reset()

	go r.Run(3)
	time.Sleep(20 * time.Millisecond)
	if atomic.LoadInt32(a) != started {
		t.Fatal("a is not started at 20ms")
	}
	if atomic.LoadInt32(b) != started {
		t.Fatal("b is not started at 20ms")
	}
	if atomic.LoadInt32(c) != started {
		t.Fatal("c is not started at 20ms")
	}
	time.Sleep(35 * time.Millisecond)

	if atomic.LoadInt32(a) != done {
		t.Fatal("a is not done at 55ms")
	}
	if atomic.LoadInt32(b) != done {
		t.Fatal("b is not done at 55ms")
	}
	if atomic.LoadInt32(c) != done {
		t.Fatal("c is not done at 55ms")
	}
}

func (ct *integralTest) runAsync2level(t *testing.T) {
	r := ct.New()
	const (
		init int32 = iota
		started
		done
	)
	x, y, z := init, init, init
	a, b, c := &x, &y, &z
	t.Run("addTaskA", ensureAdd(r, New(
		"a",
		func() error {
			atomic.StoreInt32(a, started)
			time.Sleep(50 * time.Millisecond)
			atomic.StoreInt32(a, done)
			return nil
		},
		[]string{},
	)))
	t.Run("addTaskA", ensureAdd(r, New(
		"b",
		func() error {
			atomic.StoreInt32(b, started)
			time.Sleep(50 * time.Millisecond)
			atomic.StoreInt32(b, done)
			return nil
		},
		[]string{"a"},
	)))
	t.Run("addTaskA", ensureAdd(r, New(
		"c",
		func() error {
			atomic.StoreInt32(c, started)
			time.Sleep(50 * time.Millisecond)
			atomic.StoreInt32(c, done)
			return nil
		},
		[]string{"a"},
	)))

	go r.Run(2)
	time.Sleep(20 * time.Millisecond)
	if atomic.LoadInt32(a) != started {
		t.Fatal("a is not started at 20ms")
	}
	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(a) != done {
		t.Fatal("a is not done at 70ms")
	}
	if atomic.LoadInt32(b) != started {
		t.Fatal("b is not started at 70ms")
	}
	if atomic.LoadInt32(c) != started {
		t.Fatal("c is not started at 70ms")
	}
	time.Sleep(35 * time.Millisecond)

	if atomic.LoadInt32(b) != done {
		t.Fatal("b is not done at 105ms")
	}
	if atomic.LoadInt32(c) != done {
		t.Fatal("c is not done at 105ms")
	}
}

func (ct *integralTest) runSyncErr(t *testing.T) {
	r := ct.New()
	e := errors.New("error")
	m := map[string]bool{}
	f := func(n string) func() error {
		return func() error {
			m[n] = true
			return nil
		}
	}
	errf := func(n string) func() error {
		return func() error {
			m[n] = true
			return e
		}
	}

	t.Run("addTaskA", ensureAdd(r, New("a", f("a"), []string{})))
	t.Run("addTaskB", ensureAdd(r, New("b", errf("b"), []string{"a"})))
	t.Run("addTaskC", ensureAdd(r, New("c", f("c"), []string{"b"})))

	err := r.RunSync()
	if err == nil {
		t.Error("expected error, got nothing")
	}
	if err != e {
		t.Error("unexpected error: ", err)
	}

	if !m["a"] {
		t.Error("a not ran")
	}
	if !m["b"] {
		t.Error("b not ran")
	}
	if m["c"] {
		t.Error("c has run")
	}
}

func (ct *integralTest) runSync(t *testing.T) {
	r := ct.New()

	// still need atomic since RunSync() is run in another goroutine
	a, b := int32(0), int32(0)
	running, done := &a, &b
	x := func(a, b int32) {
		if x := atomic.LoadInt32(running); x != a {
			t.Fatalf("expect %d running, got %d", a, x)
		}
		if x := atomic.LoadInt32(done); b != x {
			t.Fatalf("expect %d done, got %d", b, x)
		}
	}
	nop := func() error {
		atomic.AddInt32(running, 1)
		time.Sleep(50 * time.Millisecond)
		atomic.AddInt32(running, -1)
		atomic.AddInt32(done, 1)
		return nil
	}
	t.Run("addTaskA", ensureAdd(r, New("a", nop, []string{})))
	t.Run("addTaskB", ensureAdd(r, New("b", nop, []string{})))
	r.Parse()

	go r.RunSync()
	time.Sleep(20 * time.Millisecond)
	x(1, 0)
	time.Sleep(50 * time.Millisecond)
	x(1, 1)
	time.Sleep(35 * time.Millisecond)
	x(0, 2)
}

func (ct *integralTest) runSyncEmpty(t *testing.T) {
	r := ct.New()
	if err := r.RunSync(); err != nil {
		t.Fatal("unexpected error", err)
	}
}
