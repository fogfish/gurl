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
	c "github.com/fogfish/gurl/cats"
	"github.com/fogfish/gurl/http"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

// data types used by HTTP payload(s)
type headers struct {
	UserAgent string `json:"User-Agent,omitempty"`
}

type httpbin struct {
	URL     string  `json:"url,omitempty"`
	Origin  string  `json:"origin,omitempty"`
	Headers headers `json:"headers,omitempty"`
}

// basic declarative request
func request(val *httpbin) gurl.Arrow {
	return http.Join(
		// HTTP Request
		ø.GET.URL("https://httpbin.org/get"),
		ø.Accept.JSON,
		ø.UserAgent.Is("gurl"),
		// HTTP Response and its validation
		ƒ.Status.OK,
		ƒ.ContentType.JSON,
		ƒ.Recv(val),
	).Then(
		c.Defined(&val.Headers.UserAgent),
		c.Value(&val.Headers.UserAgent).String("gurl"),
		c.FMap(validate(val)),
	)
}

func validate(val *httpbin) func() error {
	return func() error {
		return nil
	}
}

func main() {
	var val httpbin
	req := request(&val)
	cat := http.DefaultIO(gurl.Logging(3))

	if err := req(cat).Fail; err != nil {
		fmt.Printf("fail %v\n", err)
	}
	fmt.Printf("==> %v\n", val)
}
