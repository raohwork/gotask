// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gotask

import "fmt"

func ExampleRunner() {
	r := NewRunner()
	r.Add(New("a", func() error {
		fmt.Println("a")
		return nil
	}, []string{"b"}))
	r.Add(New("b", func() error {
		fmt.Println("b")
		return nil
	}, []string{"c"}))
	r.Add(New("c", func() error {
		fmt.Println("c")
		return nil
	}, []string{}))

	r.Run(3)

	// output: c
	// b
	// a
}
