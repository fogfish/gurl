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
)

type headers struct {
	UserAgent string `json:"X-User-Agent,omitempty"`
}

type httpbin struct {
	URL     string  `json:"url,omitempty"`
	Origin  string  `json:"origin,omitempty"`
	Headers headers `json:"headers,omitempty"`
}

func request() (val httpbin, err error) {
	err = gurl.IO().
		GET("https://httpbin.org/get").
		With("Accept", "application/json").
		With("X-User-Agent", "gurl").
		Code(200).
		Head("Content-Type", "application/json").
		Recv(&val).
		Require(val.Headers.UserAgent, "gurl").
		Assert(validate(val)).
		Fail

	return
}

func validate(val httpbin) func() error {
	return func() error {
		return nil
	}
}

func main() {
	val, err := request()

	if err != nil {
		fmt.Printf("fail %v\n", err)
	}
	fmt.Printf("==> %v\n", val)
}
