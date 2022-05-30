//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package gurl

import (
	"fmt"
)

/*

Arrow is generic composable I/O
*/
type Arrow func() error

/*

Join composes I/O
*/
func Join(arrows ...Arrow) error {
	for _, f := range arrows {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

// NotSupported is returned if communication schema is not supported.
type NotSupported struct{ URL string }

func (e *NotSupported) Error() string {
	return fmt.Sprintf("Not supported: %s", e.URL)
}

// Mismatch is returned by api if expectation at body value is failed
type NoMatch struct {
	Diff    string
	Payload interface{}
}

func (e *NoMatch) Error() string {
	return e.Diff
}
