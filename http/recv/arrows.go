//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

// Package recv defines a pure computations to compose HTTP response receivers
package recv

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/fogfish/gurl"
)

/*

Code is a mandatory statement to match expected HTTP Status Code against
received one. The execution fails with BadMatchCode if service responds
with other value then specified one.
*/
func Code(code ...int) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.Unsafe()
		if io.Fail != nil {
			return io
		}

		status := io.HTTP.Ingress.StatusCode
		if !hasCode(code, status) {
			io.Fail = &gurl.BadMatchCode{Expect: code, Actual: status}
		}
		return io
	}
}

func hasCode(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/*

Header matches presence of header in the response or match its entire content.
The execution fails with BadMatchHead if the matched value do not meet expectations.
Use wildcard string ("*") to match any header value
*/
func Header(header, value string) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		h := io.HTTP.Ingress.Header.Get(header)
		if h == "" {
			io.Fail = &gurl.BadMatchHead{Header: header, Expect: value}
		} else if value != "*" && !strings.HasPrefix(h, value) {
			io.Fail = &gurl.BadMatchHead{Header: header, Expect: value, Actual: h}
		}

		return io
	}
}

/*

HeaderString matches a header value to closed variable of string type.
*/
func HeaderString(header string, value *string) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		val := io.HTTP.Ingress.Header.Get(header)
		if val == "" {
			io.Fail = &gurl.BadMatchHead{Header: header}
		} else {
			*value = val
		}

		return io
	}
}

// Served is a syntax sugar of Header("Content-Type", ...)
func Served(mime string) gurl.Arrow {
	return Header("Content-Type", mime)
}

// ServedJSON is a syntac sugar of Head("Content-Type", "application/json")
func ServedJSON() gurl.Arrow {
	return Header("Content-Type", "application/json")
}

/*

Recv applies auto decoders for response and returns either binary or
native Go data structure. The Content-Type header give a hint to decoder.
Supply the pointer to data target data structure.
*/
func Recv(out interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.Fail = decode(
			io.HTTP.Ingress.Header.Get("Content-Type"),
			io.HTTP.Ingress.Body,
			&out,
		)
		io.HTTP.Ingress.Body.Close()
		io.Body = out
		io.HTTP.Ingress = nil
		return io
	}
}

func decode(content string, stream io.ReadCloser, data interface{}) error {
	switch {
	case strings.HasPrefix(content, "application/json"):
		return json.NewDecoder(stream).Decode(&data)
	default:
		return &gurl.BadMatchHead{
			Header: "Content-Type",
			Actual: content,
		}
	}
}

/*

Bytes receive raw binary from HTTP response
*/
func Bytes(val *[]byte) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		*val, io.Fail = ioutil.ReadAll(io.HTTP.Ingress.Body)
		io.Fail = io.HTTP.Ingress.Body.Close()
		io.HTTP.Ingress = nil
		return io
	}
}

/*

FMap applies clojure to matched HTTP request.
*/
func FMap(f func() error) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.Fail = f()
		return io
	}
}

/*

Defined checks if the value is defined, use a pointer to the value.
*/
func Defined(value interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		va := reflect.ValueOf(value)
		if va.Kind() == reflect.Ptr {
			va = va.Elem()
		}

		if !va.IsValid() {
			io.Fail = &gurl.Undefined{Type: va.Type().Name()}
		}

		if va.IsValid() && va.IsZero() {
			io.Fail = &gurl.Undefined{Type: va.Type().Name()}
		}
		return io
	}
}

/*

Require checks if the value equals to defined one.
Supply the pointer to actual value
*/
func Require(actual interface{}, expect interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		if !isEqual(actual, expect) {
			io.Fail = &gurl.BadMatch{Expect: expect, Actual: actual}
		}
		return io
	}
}

func isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Kind() == reflect.Ptr {
		va = va.Elem()
	}
	if vb.Kind() == reflect.Ptr {
		vb = vb.Elem()
	}

	if !va.Type().Comparable() {
		return false
	}
	return va.Interface() == vb.Interface()
}
