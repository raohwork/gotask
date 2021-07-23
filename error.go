// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import "errors"

// ErrMissing indicates a dependency is missing
type ErrMissing string

func (e ErrMissing) Error() string {
	return "missing dependency: " + string(e)
}

// ErrDuplicated indicates the task has been registered
type ErrDuplicated string

func (e ErrDuplicated) Error() string {
	return "task has been registered: " + string(e)
}

var (
	// ErrCyclic indicates cyclic dependencies are detected
	ErrCyclic = errors.New("cyclic dependencies detected")
	// ErrParsed indicates Runner.Add() is called after Runner.Parse()
	ErrParsed = errors.New("cannot add task after parsed")
)
