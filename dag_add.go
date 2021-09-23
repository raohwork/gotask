// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

func (g *dag) Add(t Task) (err error) {
	if g.parsed {
		return ErrParsed
	}

	n := t.Name()
	if _, ok := g.vertexes[n]; ok {
		return ErrDuplicated(n)
	}
	for _, dep := range t.Depends() {
		if dep == n {
			return ErrCyclic
		}
	}

	v := &vertex{
		name:    n,
		t:       t,
		parents: map[string]bool{},
	}

	g.vertexes[v.name] = v

	return nil
}

func (g *dag) QuickAdd(name string, f func() error) error {
	return g.Add(&task{
		name: name,
		f:    f,
	})
}

func (g *dag) Skip(name ...string) {
	m := map[string]bool{}
	for _, n := range name {
		m[n] = true
	}
	for n, v := range g.vertexes {
		if !m[n] {
			continue
		}
		v.skip = true
	}
}
