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
	"strings"
)

// Accept is a HTTP header literal "Accept"
const Accept = "Accept"

// ApplicationJson is Content-Type literal "application/json"
const ApplicationJson = "application/json"

// ApplicationForm is Content-Type literal "application/x-www-form-urlencoded"
const ApplicationForm = "application/x-www-form-urlencoded"

// ContentType is a HTTP header literal "Content-Type"
const ContentType = "Content-Type"

// IO defines the category or type for HTTP I/O. A composition of
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
type IO struct {
	pool *http.Client
	uri  *url.URL
	http *httpio
	Fail error
}

// Http is used internally
type httpio struct {
	method  string
	head    map[string]string
	payload *bytes.Buffer
	ingress *http.Response
}

// NewIO creates the instance of HTTP I/O category with default HTTP client.
// Please note that default client disables TLS verification.
// Use this only for testing.
func NewIO() *IO {
	return NewIOC(
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

// NewIOC creates the instance of HTTP I/O category with well-defined
// http client.
func NewIOC(client *http.Client) *IO {
	return &IO{client, nil, nil, nil}
}

// URL defines a mandatory parameters such as HTTP method and destination URL
func (io *IO) URL(method string, uri string) *IO {
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

// GET is warpper to IO("GET", ...)
func (io *IO) GET(uri string) *IO {
	return io.URL("GET", uri)
}

// POST is warpper to IO("POST", ...)
func (io *IO) POST(uri string) *IO {
	return io.URL("POST", uri)
}

// PUT is warpper to IO("PUT", ...)
func (io *IO) PUT(uri string) *IO {
	return io.URL("PUT", uri)
}

// DELETE is warpper to IO("DELETE", ...)
func (io *IO) DELETE(uri string) *IO {
	return io.URL("DELETE", uri)
}

// With defines HTTP headers, you can add as many headers as needed using With syntax.
func (io *IO) With(head string, value string) *IO {
	if io.Fail != nil {
		return io
	}
	io.http.head[head] = value
	return io
}

// Send payload to destination URL. You can also use native Go data types
// (e.g. maps, struct, etc) as egress payload. The library implicitly encodes
// input structures to binary using Content-Type as a hint
func (io *IO) Send(data interface{}) *IO {
	if io.Fail != nil {
		return io
	}
	io.http.payload, io.Fail = encode(io.http.head[ContentType], data)
	return io
}

// Code is a mandatory statement to match expected HTTP Status Code against
// received one. The execution fails with BadMatchCode if service responds
// with other value then specified one.
func (io *IO) Code(code ...int) *IO {
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
func (io *IO) Head(head string, value string) *IO {
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
func (io *IO) Recv(out interface{}) *IO {
	if io.Fail != nil {
		return io
	}
	defer io.http.ingress.Body.Close()
	json.NewDecoder(io.http.ingress.Body).Decode(&out)
	return io
}

//
func (io *IO) unsafe() *IO {
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

	io.http.ingress, io.Fail = io.pool.Do(eg)
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
