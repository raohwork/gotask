// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

func (g *dag) Only(name ...string) Runner {
	ret := NewRunner()
	m := map[string]bool{}
	for _, n := range name {
		m[n] = true
	}

	for _, v := range g.vertexes {
		if !m[v.t.Name()] {
			continue
		}

		ret.Add(v.t)
	}

	return ret
}

func (g *dag) Without(name ...string) Runner {
	ret := NewRunner()
	m := map[string]bool{}
	for _, n := range name {
		m[n] = true
	}

	for _, v := range g.vertexes {
		if m[v.t.Name()] {
			continue
		}

		ret.Add(v.t)
	}

	return ret
}

func (g *dag) With(name ...string) Runner {
	ret := NewRunner()
	cur := map[string]bool{}
	next := map[string]bool{}
	for _, n := range name {
		cur[n] = true
	}

	for len(cur) > 0 {
		for _, v := range g.vertexes {
			if !cur[v.t.Name()] {
				continue
			}

			ret.Add(v.t)
			for _, d := range v.t.Depends() {
				next[d] = true
			}
		}
		cur = next
		next = map[string]bool{}
	}

	return ret
}
