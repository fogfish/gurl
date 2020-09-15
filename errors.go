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

ProtocolNotSupported is returned if handling of URL schema is not supported by the library
*/
type ProtocolNotSupported string

func (e ProtocolNotSupported) Error() string {
	return fmt.Sprintf("Not supported protocol: %s", string(e))
}

// Undefined is returned by api if expectation at body value is failed
type Undefined struct {
	Type string
}

func (e *Undefined) Error() string {
	return fmt.Sprintf("Value of type %v is not defined.", e.Type)
}

// Mismatch is returned by api if expectation at body value is failed
type Mismatch struct {
	Diff    string
	Payload interface{}
}

func (e *Mismatch) Error() string {
	return e.Diff
}
