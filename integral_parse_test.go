// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import (
	"testing"
)

func (c *integralTest) TestParse(t *testing.T) {
	t.Run("ok", c.parseOK)
	t.Run("missing", c.parseMissing)
	t.Run("cyclic", func(t *testing.T) {
		t.Run("2", c.parseCyclic2)
		t.Run("3", c.parseCyclic3)
	})
}

func (c *integralTest) parseMissing(t *testing.T) {
	r := c.New()
	t.Run("addTaskA", ensureAdd(r, New("a", nop, []string{"b"})))

	err := r.Parse()
	if err == nil {
		t.Fatal("expected ErrMissing, got nothing")
	}
	if err != ErrMissing("b") {
		t.Fatalf("expected missing b, got %s", err)
	}
}

func (c *integralTest) parseCyclic2(t *testing.T) {
	r := c.New()
	t.Run("addTaskA", ensureAdd(r, New("a", nop, []string{"b"})))
	t.Run("addTaskB", ensureAdd(r, New("b", nop, []string{"a"})))

	err := r.Parse()
	if err == nil {
		t.Fatal("expected ErrCyclic, got nothing")
	}
	if err != ErrCyclic {
		t.Fatalf("expected ErrCyclic, got %s", err)
	}
}

func (c *integralTest) parseCyclic3(t *testing.T) {
	r := c.New()
	t.Run("addTaskA", ensureAdd(r, New("a", nop, []string{"b"})))
	t.Run("addTaskB", ensureAdd(r, New("b", nop, []string{"c"})))
	t.Run("addTaskC", ensureAdd(r, New("c", nop, []string{"a"})))

	err := r.Parse()
	if err == nil {
		t.Fatal("expected ErrCyclic, got nothing")
	}
	if err != ErrCyclic {
		t.Fatalf("expected ErrCyclic, got %s", err)
	}
}

func (c *integralTest) parseOK(t *testing.T) {
	r := c.New()
	t.Run("addTaskA", ensureAdd(r, New("a", nop, []string{"b"})))
	t.Run("addTaskB", ensureAdd(r, New("b", nop, []string{"c"})))
	t.Run("addTaskC", ensureAdd(r, New("c", nop, []string{})))

	err := r.Parse()
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}
}
