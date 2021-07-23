// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import (
	"sort"
	"testing"
)

func checkTasks(expect []string, actual []*vertex) (ok bool) {
	if len(expect) != len(actual) {
		return
	}

	m := map[string]bool{}
	for _, e := range expect {
		m[e] = false
	}

	for _, a := range actual {
		if _, x := m[a.name]; !x {
			return
		}
		m[a.name] = true
	}

	for _, v := range m {
		if !v {
			return
		}
	}

	return true
}

func checkDepthMap(expect [][]string, actual [][]*vertex) func(*testing.T) {
	return func(t *testing.T) {
		if x, y := len(expect), len(actual); x != y {
			t.Fatalf("expect %d groups, got %d", x, y)
		}

		for idx := range expect {
			if !checkTasks(expect[idx], actual[idx]) {
				t.Log("unexpected group ", idx)
				sort.Sort(sort.StringSlice(expect[idx]))
				arr := make([]string, len(actual[idx]))
				for i, v := range actual[idx] {
					arr[i] = v.name
				}
				sort.Sort(sort.StringSlice(arr))
				t.Logf("expect: %+v", expect[idx])
				t.Fatalf("actual: %+v", arr)
			}
		}
	}
}

func TestDAGParseEmpty(t *testing.T) {
	r := NewRunner()

	r.Parse()
	checkDepthMap([][]string{}, r.(*dag).depthMap)(t)
}

func TestDAGParseNoDep(t *testing.T) {
	r := NewRunner()
	r.QuickAdd("a", nop)
	r.QuickAdd("b", nop)
	r.QuickAdd("c", nop)

	r.Parse()
	checkDepthMap([][]string{
		{"a", "b", "c"},
	}, r.(*dag).depthMap)(t)
}

func TestDAGParseSimpleDep(t *testing.T) {
	r := NewRunner()
	r.Add(New("a", nop, []string{}))
	r.Add(New("b", nop, []string{"a"}))
	r.Add(New("c", nop, []string{"b"}))

	r.Parse()
	checkDepthMap([][]string{
		{"a"},
		{"b"},
		{"c"},
	}, r.(*dag).depthMap)(t)
}

func TestDAGParseComplexDep(t *testing.T) {
	r := NewRunner()
	r.Add(New("a1", nop, []string{}))
	r.Add(New("a2", nop, []string{}))
	r.Add(New("b", nop, []string{"a1"}))
	r.Add(New("c", nop, []string{"a2", "b"}))
	r.Add(New("d", nop, []string{"a2"}))

	r.Parse()
	checkDepthMap([][]string{
		{"a1", "a2"},
		{"b", "d"},
		{"c"},
	}, r.(*dag).depthMap)(t)
}
