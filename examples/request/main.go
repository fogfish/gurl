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
	"context"
	"fmt"

	"github.com/fogfish/gurl/http"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

// data types used by HTTP payload(s)
type tHeaders struct {
	UserAgent string `json:"User-Agent,omitempty"`
}

type tHTTPBin struct {
	URL     string   `json:"url,omitempty"`
	Origin  string   `json:"origin,omitempty"`
	Headers tHeaders `json:"headers,omitempty"`
}

func (bin *tHTTPBin) validate(*http.Context) error {
	if bin.Headers.UserAgent == "" {
		return fmt.Errorf("User-Agent is not defined")
	}

	if bin.Headers.UserAgent != "gurl" {
		return fmt.Errorf("User-Agent is not valid")
	}

	return nil
}

// basic declarative request
func request(cat http.Stack) (*tHTTPBin, error) {
	var data tHTTPBin

	err := cat.IO(context.TODO(),
		// HTTP Request
		ø.GET.URL("https://httpbin.org/get"),
		ø.Accept.JSON,
		ø.UserAgent.Is("gurl"),

		// HTTP Response
		ƒ.Status.OK,
		ƒ.ContentType.JSON,
		ƒ.Recv(&data),

		// asserts
		data.validate,
	)

	return &data, err
}

func main() {
	cat := http.New(http.LogPayload())

	val, err := request(cat)
	if err != nil {
		fmt.Printf("fail %v\n", err)
	}

	fmt.Printf("==> %v\n", val)
}
