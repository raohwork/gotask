// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import "testing"

func (c *integralTest) TestAdd(t *testing.T) {
	t.Run("ok", c.addOK)
	t.Run("afterParse", c.addAfterParse)
	t.Run("afterReset", c.addAfterReset)
	t.Run("cyclic", c.addCyclic)
	t.Run("duplicated", c.addDuplicated)
}

func (c *integralTest) addOK(t *testing.T) {
	r := c.New()
	if err := r.Add(New("a", nop, []string{})); err != nil {
		t.Fatal("unexpected error: ", err)
	}
}

func (c *integralTest) addAfterParse(t *testing.T) {
	r := c.New()
	t.Run("addTaskA", ensureAdd(r, New("a", nop, []string{})))
	if err := r.Parse(); err != nil {
		t.Fatal("unexpected parse error: ", err)
	}

	err := r.QuickAdd("b", nop)
	if err == nil {
		t.Fatal("expected to have ErrParsed, got nothing")
	}
	if err != ErrParsed {
		t.Fatal("expected to have ErrParsed, got ", err)
	}
}

func (c *integralTest) addAfterReset(t *testing.T) {
	r := c.New()
	t.Run("addTaskA", ensureAdd(r, New("a", nop, []string{})))
	if err := r.Parse(); err != nil {
		t.Fatal("unexpected parse error: ", err)
	}
	r.Reset()

	err := r.QuickAdd("b", nop)
	if err != nil {
		t.Fatal("unexpected error ", err)
	}
}

func (c *integralTest) addCyclic(t *testing.T) {
	r := c.New()
	err := r.Add(New("a", nop, []string{"a"}))

	if err == nil {
		t.Fatal("expected ErrCyclic, got nothing")
	}
	if err != ErrCyclic {
		t.Fatalf("expected ErrCyclic, got %s", err)
	}

}

func (c *integralTest) addDuplicated(t *testing.T) {
	r := c.New()
	t.Run("AddTaskA", ensureAdd(r, New("a", nop, []string{})))

	err := r.Add(New("a", nop, []string{"b"}))
	if err == nil {
		t.Fatal("expected ErrDuplicated, got nothing")
	}
	if err != ErrDuplicated("a") {
		t.Fatalf("expected ErrDuplicated(a), got %s", err)
	}
}
