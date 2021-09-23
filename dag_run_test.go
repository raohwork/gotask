// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import "testing"

func TestRunSyncSkip(t *testing.T) {
	msg := ""
	f := func(n string) func() error {
		return func() error {
			msg += n
			return nil
		}
	}
	r := NewRunner()

	r.QuickAdd("a", f("a"))
	r.Add(New("b", f("b"), []string{"a"}))
	r.Add(New("c", f("c"), []string{"b"}))
	r.Skip("b")
	if err := r.RunSync(); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if msg != "ac" {
		t.Fatalf("expected ac, got %s", msg)
	}
}

func TestRunSkip(t *testing.T) {
	msg := make(chan string, 1)
	msg <- ""

	f := func(n string) func() error {
		return func() error {
			msg <- (<-msg) + n
			return nil
		}
	}
	r := NewRunner()

	r.QuickAdd("a", f("a"))
	r.Add(New("b", f("b"), []string{"a"}))
	r.Add(New("c", f("c"), []string{"b"}))
	r.Skip("b")
	if err := r.RunSync(); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if x := <-msg; x != "ac" {
		t.Fatalf("expected ac, got %s", x)
	}
}

func TestSubtasks(t *testing.T) {
	msg := ""
	f := func(s string) func() error {
		return func() error {
			msg += s
			return nil
		}
	}
	r := NewRunner()
	r.QuickAdd("a", f("a"))
	r.QuickAdd("b", f("b"))
	r.Add(New("c", f("c"), []string{"b"}))
	r.QuickAdd("d", f("d"))

	part := r.With("c")
	r.Skip(part.Tasks()...)

	part.RunSync()
	if msg != "bc" {
		t.Fatal("unexpected part result: ", msg)
	}
	if err := r.RunSync(); err != nil {
		t.Fatal("unexpected main error:", err)
	}
	if msg != "bcad" && msg != "bcda" {
		t.Fatal("unexpected main result: ", msg)
	}
}
