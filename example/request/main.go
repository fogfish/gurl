//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package main

/*

Example shows a basic usage of HTTP I/O.

*/

import (
	"fmt"

	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

// data types used by HTTP payload(s)
type headers struct {
	UserAgent string `json:"X-User-Agent,omitempty"`
}

type httpbin struct {
	URL     string  `json:"url,omitempty"`
	Origin  string  `json:"origin,omitempty"`
	Headers headers `json:"headers,omitempty"`
}

// basic declarative request
func request(val *httpbin) gurl.Arrow {
	return gurl.HTTP(
		// HTTP output
		ø.GET("https://httpbin.org/get"),
		ø.Header("Accept").Is("application/json"),
		ø.Header("X-User-Agent").Is("gurl"),
		// HTTP input and its validation
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Header("Content-Type").Is("application/json"),
		ƒ.Recv(val),
		ƒ.Defined(&val.Headers.UserAgent),
		ƒ.Value(&val.Headers.UserAgent).String("gurl"),
		ƒ.FMap(validate(val)),
	)
}

func validate(val *httpbin) func() error {
	return func() error {
		return nil
	}
}

func main() {
	var val httpbin
	http := request(&val)

	if err := http(gurl.IO(gurl.Verbose(3))).Fail; err != nil {
		fmt.Printf("fail %v\n", err)
	}
	fmt.Printf("==> %v\n", val)
}
