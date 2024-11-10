//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package main

/*

Example shows a composition of HTTP I/O.

*/

import (
	"context"
	"fmt"

	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

// id implements payload for https://httpbin.org/uuid
type ID struct {
	UUID string `json:"uuid,omitempty"`
}

// httpbin implements payload for https://httpbin.org/post
type HTTPBin struct {
	URL  string `json:"url,omitempty"`
	Data string `json:"data,omitempty"`
}

// context for HTTP I/O
type Heap struct {
	ID
	HTTPBin
}

// uuid declares HTTP I/O. Its result is returned via id variable.
func (hof *Heap) uuid() http.Arrow {
	return http.GET(
		ø.URI("https://httpbin.org/uuid"),
		ø.Accept.JSON,

		ƒ.Status.OK,
		ƒ.ContentType.JSON,
		ƒ.Body(&hof.ID),
	)
}

// post declares HTTP I/O. The HTTP request requires uuid.
// Its result is returned via doc variable.
func (hof *Heap) post() http.Arrow {
	return http.POST(
		ø.URI("https://httpbin.org/post"),
		ø.Accept.JSON,
		ø.ContentType.JSON,
		ø.Send(&hof.ID.UUID),

		ƒ.Status.OK,
		ƒ.Body(&hof.HTTPBin),
	)
}

// request is a high-order function. It is composed from atomic HTTP I/O into
// the chain of requests.
func request() (*Heap, http.Arrow) {
	var heap Heap

	//
	// HoF combines HTTP requests to
	//  * https://httpbin.org/uuid
	//  * https://httpbin.org/post
	//
	// results of HTTP I/O is persisted in the internal state
	return &heap, http.Join(
		heap.uuid(),
		heap.post(),
	)
}

func main() {
	// instance of http stack
	stack := http.New(http.WithDebugPayload)

	data, lazy := request()

	// executes http I/O
	err := stack.IO(context.Background(), lazy)
	if err != nil {
		panic(err)
	}

	fmt.Printf("==> %+v\n", data)
}
