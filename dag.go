// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

type vertex struct {
	name    string
	t       Task
	depth   uint
	up      []string
	down    []string
	parents map[string]bool
	skip    bool
}

type dag struct {
	vertexes map[string]*vertex
	parsed   bool
	depthMap [][]*vertex
	errs     []error
}

// NewRunner creates a Runner instance
func NewRunner() (ret Runner) {
	return &dag{
		vertexes: map[string]*vertex{},
	}
}

func (g *dag) Tasks() (ret []string) {
	ret = make([]string, 0, len(g.vertexes))
	for n := range g.vertexes {
		ret = append(ret, n)
	}

	return
}

func (g *dag) Errors() []error {
	return g.errs
}

func (g *dag) updateChildren(v *vertex) {
	depth := v.depth + 1

	for _, name := range v.down {
		g.updateChild(v, name, depth)
	}
}

func (g *dag) updateChild(up *vertex, name string, depth uint) {
	me := g.vertexes[name]
	me.parents[up.name] = true
	for dep := range up.parents {
		me.parents[dep] = true
	}
	if me.depth < depth {
		me.depth = depth
	}
	g.updateChildren(me)
}

func (g *dag) ensureMapDepth(d uint) {
	l := len(g.depthMap)
	if uint(l) >= d {
		return
	}
	arr := make([][]*vertex, d-uint(l))
	g.depthMap = append(g.depthMap, arr...)
}
