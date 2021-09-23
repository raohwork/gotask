// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import "testing"

func TestOnly(t *testing.T) {
	r := NewRunner()
	r.QuickAdd("a", nop)
	r.QuickAdd("b", nop)
	r.Add(New("c", nop, []string{"b"}))
	r.QuickAdd("d", nop)

	n := r.Only("a", "c").(*dag)
	if l := len(n.vertexes); l != 2 {
		t.Errorf("expected 2 tasks, got %d", l)
	}
	if _, ok := n.vertexes["a"]; !ok {
		t.Error("expected task a, got nothing")
	}
	if _, ok := n.vertexes["b"]; ok {
		t.Error("unexpected task b")
	}
	if _, ok := n.vertexes["c"]; !ok {
		t.Error("expected task c, got nothing")
	}
	if _, ok := n.vertexes["d"]; ok {
		t.Error("unexpected task d")
	}
}

func TestWithout(t *testing.T) {
	r := NewRunner()
	r.QuickAdd("a", nop)
	r.QuickAdd("b", nop)
	r.Add(New("c", nop, []string{"b"}))
	r.QuickAdd("d", nop)

	n := r.Without("a", "c").(*dag)
	if l := len(n.vertexes); l != 2 {
		t.Errorf("expected 2 tasks, got %d", l)
	}
	if _, ok := n.vertexes["a"]; ok {
		t.Error("unexpected task a")
	}
	if _, ok := n.vertexes["b"]; !ok {
		t.Error("expected task b, got nothing")
	}
	if _, ok := n.vertexes["c"]; ok {
		t.Error("unexpected task c")
	}
	if _, ok := n.vertexes["d"]; !ok {
		t.Error("expected task d, got nothing")
	}
}

func TestWith(t *testing.T) {
	r := NewRunner()
	r.QuickAdd("a", nop)
	r.Add(New("b", nop, []string{"a"}))
	r.Add(New("c", nop, []string{"b"}))
	r.QuickAdd("d", nop)

	n := r.With("c").(*dag)
	if l := len(n.vertexes); l != 3 {
		t.Errorf("expected 3 tasks, got %d", l)
	}
	if _, ok := n.vertexes["a"]; !ok {
		t.Error("expected task a, got nothing")
	}
	if _, ok := n.vertexes["b"]; !ok {
		t.Error("expected task b, got nothing")
	}
	if _, ok := n.vertexes["c"]; !ok {
		t.Error("expected task c, got nothing")
	}
	if _, ok := n.vertexes["d"]; ok {
		t.Error("unexpected task d")
	}
}
