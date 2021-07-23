// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

// Task represents a task to run
type Task interface {
	// task name, also used as dependency specification (CASE-SENSITIVE)
	Name() string
	// actually run the task
	Run() error
	// returns dependencies (CASE-SENSITIVE)
	Depends() []string
}

type task struct {
	name    string
	f       func() error
	depends []string
}

func (t *task) Name() string      { return t.name }
func (t *task) Run() error        { return t.f() }
func (t *task) Depends() []string { return t.depends }

// New is a helper to create a Task instance
func New(name string, f func() error, depends []string) (ret Task) {
	return &task{
		name:    name,
		f:       f,
		depends: depends,
	}
}
