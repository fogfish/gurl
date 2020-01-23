//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package gurl

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// Accept is a HTTP header literal "Accept"
const Accept = "Accept"

// ApplicationJson is Content-Type literal "application/json"
const ApplicationJson = "application/json"

// ApplicationForm is Content-Type literal "application/x-www-form-urlencoded"
const ApplicationForm = "application/x-www-form-urlencoded"

// ContentType is a HTTP header literal "Content-Type"
const ContentType = "Content-Type"

// IOCat defines the category or type for HTTP I/O. A composition of
// HTTP primitives within the category are written with the following syntax:
//
//   gurl.NewIO().Arrow1(). ... ArrowN()
//
// Here, each Arrow is a morphism applied to HTTP protocol, they composition
// is defined using dot (.). Effectively the implementation resembles the
// state monad. It defines an abstraction of environments and lenses to focus
// inside it. In other words, the category represents the environment as an
// "invisible" side-effect of the composition.
//
// The type IO implements http primitives ("arrow"). Fail is only exported
// value that provides the final status of the execution.
type IOCat struct {
	pool *http.Client
	uri  *url.URL
	http *httpio
	dur  time.Duration
	Body interface{}
	Fail error
}

// Http is used internally
type httpio struct {
	method  string
	head    map[string]string
	payload *bytes.Buffer
	ingress *http.Response
}

// Status contains summary about the IO
type Status struct {
	ID       string      `json:"id"`
	Status   string      `json:"status"`
	Duration int64       `json:"duration"`
	Expect   interface{} `json:"expect,omitempty"`
	Actual   interface{} `json:"actual,omitempty"`
}

// IO creates the instance of HTTP I/O category with default HTTP client.
// Please note that default client disables TLS verification.
// Use this only for testing.
func IO() *IOCat {
	return Use(
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	)
}

// Use creates the instance of HTTP I/O category with well-defined
// http client.
func Use(client *http.Client) *IOCat {
	return &IOCat{client, nil, nil, 0, nil, nil}
}

//-----------------------------------------------------------------------------
//
// up stream, output HTTP request
//
//-----------------------------------------------------------------------------

// URL defines a mandatory parameters to the request such as
// HTTP method and destination URL
func (io *IOCat) URL(method string, uri string) *IOCat {
	if io.Fail != nil {
		return io
	}
	io.http = nil
	io.uri, io.Fail = url.Parse(uri)

	switch io.uri.Scheme {
	case "http", "https":
		io.http = &httpio{method, make(map[string]string), bytes.NewBuffer(nil), nil}
	default:
		io.Fail = &BadSchema{io.uri.Scheme}
	}
	return io
}

// GET is warpper to URL("GET", ...)
func (io *IOCat) GET(uri string) *IOCat {
	return io.URL("GET", uri)
}

// POST is warpper to URL("POST", ...)
func (io *IOCat) POST(uri string) *IOCat {
	return io.URL("POST", uri)
}

// PUT is warpper to URL("PUT", ...)
func (io *IOCat) PUT(uri string) *IOCat {
	return io.URL("PUT", uri)
}

// DELETE is warpper to URL("DELETE", ...)
func (io *IOCat) DELETE(uri string) *IOCat {
	return io.URL("DELETE", uri)
}

// With defines output HTTP headers, you can add as many headers as needed using With syntax.
func (io *IOCat) With(head string, value string) *IOCat {
	if io.Fail != nil {
		return io
	}
	io.http.head[head] = value
	return io
}

// Send payload to destination URL. You can also use native Go data types
// (e.g. maps, struct, etc) as egress payload. The library implicitly encodes
// input structures to binary using Content-Type as a hint
func (io *IOCat) Send(data interface{}) *IOCat {
	if io.Fail != nil {
		return io
	}
	io.http.payload, io.Fail = encode(io.http.head[ContentType], data)
	return io
}

//-----------------------------------------------------------------------------
//
// down stream, input HTTP response
//
//-----------------------------------------------------------------------------

// Code is a mandatory statement to match expected HTTP Status Code against
// received one. The execution fails with BadMatchCode if service responds
// with other value then specified one.
func (io *IOCat) Code(code ...int) *IOCat {
	if io.Fail != nil {
		return io
	}
	io.unsafe()

	status := io.http.ingress.StatusCode
	if !hasInt(code, status) {
		io.Fail = &BadMatchCode{code, status}
	}
	return io
}

// Head matches presence of header in the response or match its entire content.
// The execution fails with BadMatchHead  if the matched value do not meet expectations.
func (io *IOCat) Head(head string, value string) *IOCat {
	if io.Fail != nil {
		return io
	}

	h := io.http.ingress.Header.Get(head)
	if h == "" {
		io.Fail = &BadMatchHead{head, value, h}
	} else if value != "*" && !strings.HasPrefix(h, value) {
		io.Fail = &BadMatchHead{head, value, h}
	}

	return io
}

// Recv applies auto decoders for response and returns either binary or
// native Go data structure. The Content-Type header give a hint to decoder.
func (io *IOCat) Recv(out interface{}) *IOCat {
	if io.Fail != nil {
		return io
	}
	defer io.http.ingress.Body.Close()
	io.Fail = json.NewDecoder(io.http.ingress.Body).Decode(&out)
	io.Body = out
	return io
}

// Defined checks if the value is defined
func (io *IOCat) Defined(x interface{}) *IOCat {
	if io.Fail != nil {
		return io
	}

	val := reflect.ValueOf(x)
	if !val.IsValid() {
		io.Fail = &Undefined{val.Type().Name()}
	}

	if val.IsValid() && val.IsZero() {
		io.Fail = &Undefined{val.Type().Name()}
	}
	return io

}

// Require checks if the value equals to defined one
func (io *IOCat) Require(actual interface{}, expect interface{}) *IOCat {
	if io.Fail != nil {
		return io
	}
	if !reflect.DeepEqual(actual, expect) {
		io.Fail = &BadMatch{expect, actual}
	}
	return io
}

// Assert implements custom assert logic on decoded value
func (io *IOCat) Assert(f func() error) *IOCat {
	if io.Fail != nil {
		return io
	}
	io.Fail = f()
	return io
}

//-----------------------------------------------------------------------------
//
// internal state
//
//-----------------------------------------------------------------------------

//
func (io *IOCat) unsafe() *IOCat {
	if io.Fail != nil {
		return io
	}

	var eg *http.Request
	eg, io.Fail = http.NewRequest(io.http.method, io.uri.String(), io.http.payload)
	if io.Fail != nil {
		return io
	}

	for head, value := range io.http.head {
		eg.Header.Set(head, value)
	}

	t := time.Now()
	io.http.ingress, io.Fail = io.pool.Do(eg)
	io.dur = time.Now().Sub(t)
	return io
}

//
func hasInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
