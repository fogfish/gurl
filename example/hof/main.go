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
	ø "github.com/fogfish/gurl/http"
)

type headers struct {
	Host string `json:"Host,omitempty"`
}

type httpbin struct {
	URL     string  `json:"url,omitempty"`
	Headers headers `json:"headers,omitempty"`
	Data    string  `json:"data,omitempty"`
}

func main() {
	var val httpbin
	http := gurl.Join(
		read(&val),
		post(&val),
	)

	if err := http(gurl.IO()).Fail; err != nil {
		fmt.Printf("fail %v\n", err)
	}
	fmt.Printf("==> %v\n", val)
}

func read(val *httpbin) gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://httpbin.org/get"),
		ø.Accept("application/json"),
		ø.Code(200),
		ø.Served("application/json"),
		ø.Recv(val),
	)
}

func post(val *httpbin) gurl.Arrow {
	return gurl.HTTP(
		ø.POST("https://httpbin.org/post"),
		ø.Accept("application/json"),
		ø.Content("application/json"),
		ø.Send(val),
		ø.Code(200),
		ø.Recv(val),
	)
}
