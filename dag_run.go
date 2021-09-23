// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import (
	"math"
	"runtime"
	"sync"
)

func ensureN(n uint8) uint8 {
	if n != 0 {
		return n
	}

	x := runtime.NumCPU()
	if x > math.MaxUint8 {
		x = math.MaxUint8
	}
	return uint8(x)
}

func (g *dag) RunSync() (err error) {
	if err = g.Parse(); err != nil {
		return
	}
	g.errs = nil

	for _, tasks := range g.depthMap {
		for _, task := range tasks {
			if task.skip {
				continue
			}
			err = task.t.Run()
			if err != nil {
				g.errs = append(g.errs, err)
				return
			}
		}
	}

	return
}

func (g *dag) Run(n uint8) (err error) {
	if err = g.Parse(); err != nil {
		return
	}

	n = ensureN(n)
	wg := &sync.WaitGroup{}
	lock := &sync.Mutex{}

	for _, tasks := range g.depthMap {
		if err = g.runTasks(tasks, n, wg, lock); err != nil {
			return
		}
	}

	return
}

func (g *dag) runTasks(tasks []*vertex, n uint8, wg *sync.WaitGroup, lock *sync.Mutex) (err error) {
	wg.Add(int(n))

	tokens := make(chan *vertex)
	go func() {
		for _, t := range tasks {
			tokens <- t
		}
		close(tokens)
	}()

	safe := func() bool {
		lock.Lock()
		defer lock.Unlock()
		return err == nil
	}
	worker := func() {
		for t := range tokens {
			if !safe() {
				continue
			}
			if t.skip {
				continue
			}

			e := t.t.Run()

			if e != nil {
				lock.Lock()
				err = e
				g.errs = append(g.errs, err)
				lock.Unlock()
			}
		}
		wg.Done()
	}

	for i := uint8(0); i < n; i++ {
		go worker()
	}

	wg.Wait()

	return
}
