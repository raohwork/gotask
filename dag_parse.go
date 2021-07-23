// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

func (g *dag) Reset() {
	if !g.parsed {
		return
	}

	for _, v := range g.vertexes {
		v.depth = 0
	}

	for idx, arr := range g.depthMap {
		g.depthMap[idx] = make([]*vertex, 0, len(arr))
	}

	g.errs = nil
	g.parsed = false
}

func (g *dag) Parse() (err error) {
	if g.parsed {
		return
	}

	for _, v := range g.vertexes {
		for _, dep := range v.t.Depends() {
			err = g.connect(dep, v)
		}
	}

	g.depthMap = [][]*vertex{}
	for _, v := range g.vertexes {
		g.ensureMapDepth(v.depth + 1)
		g.depthMap[v.depth] = append(g.depthMap[v.depth], v)
	}

	g.parsed = true

	// release memory
	for _, v := range g.vertexes {
		v.up = nil
		v.down = nil
		v.parents = map[string]bool{}
	}
	return
}

func (g *dag) connect(up string, down *vertex) (err error) {
	upv := g.vertexes[up]
	if upv == nil {
		return ErrMissing(up)
	}

	if upv.parents[down.name] {
		return ErrCyclic
	}

	down.parents[up] = true
	down.up = append(down.up, up)
	upv.down = append(upv.down, down.name)

	for u := range upv.parents {
		down.parents[u] = true
	}

	depth := upv.depth + 1
	if down.depth < depth {
		down.depth = depth
	}

	g.updateChildren(down)

	return
}
