//
// Copyright (C) 2019 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package main

import (
	"context"

	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

// Declare the type, used for networking I/O.
type Payload struct {
	Origin string `json:"origin"`
	Url    string `json:"url"`
}

// declares http I/O
func request() http.Arrow {
	var data Payload

	return http.GET(
		// specify specify the request
		ø.URI("https://httpbin.org/get"),
		ø.Accept.ApplicationJSON,

		// specify requirements to the response
		ƒ.Status.OK,
		ƒ.ContentType.JSON,
		ƒ.Recv(&data),
	)
}

func main() {
	// instance of http stack
	stack := http.New(http.LogPayload())

	// declares http i/o
	lazy := request()

	// executes http I/O
	err := stack.IO(context.Background(), lazy)
	if err != nil {
		panic(err)
	}
}
