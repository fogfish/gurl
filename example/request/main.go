//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package main

import (
	"fmt"

	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

type headers struct {
	UserAgent string `json:"X-User-Agent,omitempty"`
}

type httpbin struct {
	URL     string  `json:"url,omitempty"`
	Origin  string  `json:"origin,omitempty"`
	Headers headers `json:"headers,omitempty"`
}

func request(val *httpbin) gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://httpbin.org/get"),
		ø.Accept("application/json"),
		ø.Header("X-User-Agent", "gurl"),
		ƒ.Code(200),
		ƒ.Served("application/json"),
		ƒ.Recv(val),
		ƒ.Defined(&val.Headers.UserAgent),
		ƒ.Require(&val.Headers.UserAgent, "gurl"),
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

	if err := http(gurl.IO()).Fail; err != nil {
		fmt.Printf("fail %v\n", err)
	}
	fmt.Printf("==> %v\n", val)
}
