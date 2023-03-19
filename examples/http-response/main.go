//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package main

import (
	"context"
	"fmt"

	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

// data types used by HTTP payload(s)
type Headers struct {
	UserAgent string `json:"User-Agent,omitempty"`
}

type HTTPBin struct {
	URL     string  `json:"url,omitempty"`
	Origin  string  `json:"origin,omitempty"`
	Headers Headers `json:"headers,omitempty"`
}

// combinator validates HTTP response
func (bin *HTTPBin) validate(*http.Context) error {
	if bin.Headers.UserAgent == "" {
		return fmt.Errorf("User-Agent is not defined")
	}

	if bin.Headers.UserAgent != "gurl" {
		return fmt.Errorf("User-Agent is not valid")
	}

	return nil
}

func request() (*HTTPBin, http.Arrow) {
	var data HTTPBin

	return &data, http.GET(
		// HTTP Request
		ø.URI("https://httpbin.org/get"),
		ø.Accept.JSON,
		ø.UserAgent.Set("gurl"),

		// HTTP Response
		ƒ.Status.OK,
		ƒ.ContentType.JSON,
		ƒ.Body(&data),

		// asserts
		data.validate,
	)
}

func main() {
	// instance of http stack
	stack := http.New(http.WithDebugPayload())

	data, lazy := request()

	// executes http I/O
	err := stack.IO(context.Background(), lazy)
	if err != nil {
		panic(err)
	}

	fmt.Printf("==> %+v\n", data)
}
