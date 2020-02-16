//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

// Package http defines a pure computations to compose HTTP networking
package http

import (
	"bytes"
	"encoding/json"
	"net/url"
	"reflect"
	"strings"

	"github.com/fogfish/gurl"
)

//-----------------------------------------------------------------------------
//
// HTTP request
//
//-----------------------------------------------------------------------------

// URL defines a mandatory parameters to the request such as
// HTTP method and destination URL
func URL(method, uri string) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.HTTP = nil
		io.URL, io.Fail = url.Parse(uri)

		switch io.URL.Scheme {
		case "http", "https":
			io.HTTP = &gurl.IOSpec{
				Method:  method,
				Header:  make(map[string]string),
				Payload: bytes.NewBuffer(nil),
				Ingress: nil,
			}
		default:
			io.Fail = &gurl.BadSchema{io.URL.Scheme}
		}
		return io
	}
}

// GET is syntax sugar of URL("GET", ...)
func GET(uri string) gurl.Arrow {
	return URL("GET", uri)
}

// POST is syntax sugar of URL("POST", ...)
func POST(uri string) gurl.Arrow {
	return URL("POST", uri)
}

// PUT is syntax sugar of URL("PUT", ...)
func PUT(uri string) gurl.Arrow {
	return URL("PUT", uri)
}

// DELETE is syntax sugar of URL("DELETE", ...)
func DELETE(uri string) gurl.Arrow {
	return URL("DELETE", uri)
}

// With defines output HTTP headers, you can add as many headers as needed using With syntax.
func With(header, value string) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.HTTP.Header[header] = value
		return io
	}
}

// Accept is syntax sugar of With("Accept", ...)
func Accept(mime string) gurl.Arrow {
	return With("Accept", mime)
}

// AcceptJSON is syntax sugar of With("Accept", "application/json")
func AcceptJSON() gurl.Arrow {
	return With("Accept", "application/json")
}

// Content is syntax sugar of With("Content-Type", ...)
func Content(mime string) gurl.Arrow {
	return With("Content-Type", mime)
}

// ContentJSON is syntax sugar of With("Content-Type", "application/json")
func ContentJSON() gurl.Arrow {
	return With("Content-Type", "application/json")
}

// KeepAlive is a syntax sugar of With("Connection", "keep-alive")
func KeepAlive() gurl.Arrow {
	return With("Connection", "keep-alive")
}

// Authorization is syntax sugar of With("Authorization", ...)
func Authorization(token string) gurl.Arrow {
	return With("Authorization", token)
}

// Send payload to destination URL. You can also use native Go data types
// (e.g. maps, struct, etc) as egress payload. The library implicitly encodes
// input structures to binary using Content-Type as a hint
func Send(data interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.HTTP.Payload, io.Fail = encode(io.HTTP.Header["Content-Type"], data)
		return io
	}
}

//-----------------------------------------------------------------------------
//
// HTTP response
//
//-----------------------------------------------------------------------------

// Code is a mandatory statement to match expected HTTP Status Code against
// received one. The execution fails with BadMatchCode if service responds
// with other value then specified one.
func Code(code ...int) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.Unsafe()
		if io.Fail != nil {
			return io
		}

		status := io.HTTP.Ingress.StatusCode
		if !hasCode(code, status) {
			io.Fail = &gurl.BadMatchCode{code, status}
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

// Head matches presence of header in the response or match its entire content.
// The execution fails with BadMatchHead if the matched value do not meet expectations.
// Use wildcard string ("*") to match any header value
func Head(header, value string) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		h := io.HTTP.Ingress.Header.Get(header)
		if h == "" {
			io.Fail = &gurl.BadMatchHead{header, value, h}
		} else if value != "*" && !strings.HasPrefix(h, value) {
			io.Fail = &gurl.BadMatchHead{header, value, h}
		}

		return io
	}
}

// Served is a syntax sugar of Head("Content-Type", ...)
func Served(mime string) gurl.Arrow {
	return Head("Content-Type", mime)
}

// ServedJSON is a syntac sugar of Head("Content-Type", "application/json")
func ServedJSON() gurl.Arrow {
	return Head("Content-Type", "application/json")
}

// Recv applies auto decoders for response and returns either binary or
// native Go data structure. The Content-Type header give a hint to decoder.
// Supply the pointer to data target data structure.
func Recv(out interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.Fail = json.NewDecoder(io.HTTP.Ingress.Body).Decode(&out)
		io.Fail = io.HTTP.Ingress.Body.Close()
		io.Body = out
		io.HTTP.Ingress = nil
		return io
	}
}

// Defined checks if the value is defined.
// Supply the pointer to the value
func Defined(value interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		va := reflect.ValueOf(value)
		if va.Kind() == reflect.Ptr {
			va = va.Elem()
		}

		if !va.IsValid() {
			io.Fail = &gurl.Undefined{va.Type().Name()}
		}

		if va.IsValid() && va.IsZero() {
			io.Fail = &gurl.Undefined{va.Type().Name()}
		}
		return io
	}
}

// Require checks if the value equals to defined one.
// Supply the pointer to actual value
func Require(actual interface{}, expect interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		if !isEqual(actual, expect) {
			io.Fail = &gurl.BadMatch{expect, actual}
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

// Test evaluates the assert function in the context of IO category
func Test(check func() error) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.Fail = check()
		return io
	}
}
