//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package gurl

import "fmt"

// BadSchema is returned by api if protocol is not supported.
type BadSchema struct {
	Schema string
}

func (e *BadSchema) Error() string {
	return fmt.Sprintf("Bad protocol argument %v", e.Schema)
}

// BadMatchCode is returned by api if HTTP status code expectation is failed
type BadMatchCode struct {
	Expect []int
	Actual int
}

func (e *BadMatchCode) Error() string {
	return fmt.Sprintf("Mismatch of http status code %v, required one of %v.", e.Actual, e.Expect)
}

// BadMatchHead is returned by api if HTTP header expectation is failed
type BadMatchHead struct {
	Header string
	Expect string
	Actual string
}

func (e *BadMatchHead) Error() string {
	return fmt.Sprintf("Mismatch of http header %v value %v, required %v.", e.Header, e.Actual, e.Expect)
}

// Undefined is returned by api if expectation at body value is failed
type Undefined struct {
	Type string
}

func (e *Undefined) Error() string {
	return fmt.Sprintf("Value of type %v is not defined.", e.Type)
}

// BadMatch is returned by api if expectation at body value is failed
type BadMatch struct {
	Expect interface{}
	Actual interface{}
}

func (e *BadMatch) Error() string {
	return fmt.Sprintf("Mismatch of value %v, required %v.", e.Actual, e.Expect)
}
