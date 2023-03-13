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

// NotSupported is returned if communication schema is not supported.
type NotSupported struct{ URL string }

func (e *NotSupported) Error() string {
	return fmt.Sprintf("Not supported: %s", e.URL)
}

// Mismatch is returned by api if expectation at body value is failed
type NoMatch struct {
	ID       string // unique ID of failed combinator
	Protocol any    // protocol primitive caused failure
	Diff     string // human readable difference between expected & actual values
	Expect   any    // expected value
	Actual   any    // actual value
}

func (e *NoMatch) Error() string { return e.Diff }
