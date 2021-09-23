// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

// Runner manages task dependencies and runs the tasks
//
// Currently it caches all dependencies for each task in map while parsing, which
// cost you some extra memory space (should be few KBs in real world).
//
// Runner is not thread-safe, you MUST NOT share same instance amoung multiple
// goroutines.
type Runner interface {
	// Add a task
	//
	// This might returns following errors:
	//
	//   - ErrDuplicated: a task with same name has been registered
	//   - ErrParsed: calling Add() after Parse()
	//   - ErrCyclic: the task depends itself
	Add(Task) error
	// Helper to add a task without dependency
	QuickAdd(name string, f func() error) error
	// Reset internal state, make it possible to Add() after Parse()
	Reset()
	// Check dependencies and sort tasks. It might return following errors:
	//
	//   - ErrMissing: the depdency is missing
	//   - ErrCyclic: there's a loop in dependency graph
	Parse() error
	// Run tasks one-by-one (synchronously). It almost identical to Run(1), but
	// frees you from race condition.
	//
	// As scheduler is implemented in  map, the order of task is unspecified,
	// only dependencies are ensured.
	//
	// It calls Parse() if needed, so same error applies in addition to errors
	// returned from task.
	//
	// While RunSync() has less overhead compares to Run(). It is suggested only
	// if your tasks are not time-consuming.
	RunSync() error
	// Run the tasks in n goroutines. Set n to 1 will use only 1 groutine, which
	// is roughly same as synchronized executing. Set n to 0 = runtime.NumCPU()
	//
	// As scheduler is implemented in  map, the order of task is unspecified,
	// only dependencies are ensured.
	//
	// It calls Parse() if needed, so same error applies in addition to errors
	// returned from task.
	//
	// Since tasks are run in separated goroutines, you have to take care race
	// conditions yourself.
	//
	// Tasks are grouped by their "depth", tasks in next group will not run
	// before all tasks in current group done. For example
	//
	//    task  | depends | depth
	//   -------+---------+-------------------
	//    a     |         | 0
	//    b     |         | 0
	//    c     | a       | 1 (a+1)
	//    d     | b, c    | 2 (max(b, c) + 1)
	//
	// In this case, c will not run until both a and b are done, even if you set
	// n to 2.
	//
	// To prevent race-conditions, Run() has much more overhead compares to
	// RunSync(). However, RunSync() barely beats Run() in real world scenario
	// as most tasks are time-consuming.
	Run(n uint8) error
	// Gets known errors. Multiple errors might occur concurrently in Run(), but
	// only one is returned. If you need to know all errors, just use this.
	Errors() []error
	// Only creates a new Runner which contains only specified tasks.
	//
	// Non-exist tasks are ignored silently. Say you have a Runner contains four
	// tasks: a, b, c and d. Calling Only("a", "b", "f") returns a new Runner
	// contains a and b.
	Only(name ...string) Runner
	// Without creates a new Runner without specified tasks.
	//
	// Non-exist tasks are ignored silently. Say you have a Runner contains four
	// tasks: a, b, c and d. Calling Without("a", "b", "f") returns a new Runner
	// contains c and d.
	Without(name ...string) Runner
	// With creates a new Runner contains specified tasks and their deps.
	//
	// Non-exist tasks are ignored silently. Say you have a Runner contains four
	// tasks: a, b, c (depends b) and d. Calling With("a", "c", "f") returns a
	// new Runner contains a, b and c.
	With(name ...string) Runner
}

// taskState indicates state of a task
type taskState int

const (
	statePending  taskState = iota // task is waiting for it's deps
	stateQueueing                  // task is scheduled to run
	stateRunning                   // task is running
)
