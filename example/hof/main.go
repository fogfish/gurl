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

type IO struct {
	*gurl.IOCat
}

type headers struct {
	Host string `json:"Host,omitempty"`
}

type httpbin struct {
	URL     string  `json:"url,omitempty"`
	Headers headers `json:"headers,omitempty"`
}

func main() {
	io := &IO{gurl.IO()}
	input := io.readHttpBin()
	val := io.writeHttpBin(input)

	if io.Fail != nil {
		fmt.Printf("fail %v\n", io.Fail)
	}
	fmt.Printf("==> %v\n", val)
}

func (io *IO) readHttpBin() (val httpbin) {
	io.GET("https://httpbin.org/get").
		With("Accept", "application/json").
		Code(200).
		Head("Content-Type", "application/json").
		Recv(&val)

	return
}

func (io *IO) writeHttpBin(json httpbin) (val httpbin) {
	io.POST("https://httpbin.org/post").
		With("Content-Type", "application/json").
		With("Accept", "application/json").
		Send(json).
		Code(200).
		Recv(&val)

	return
}
