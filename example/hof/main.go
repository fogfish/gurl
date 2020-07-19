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
	"fmt"

	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

// id implements payload for https://httpbin.org/uuid
type id struct {
	UUID string `json:"uuid,omitempty"`
}

// httpbin implements payload for https://httpbin.org/post
type httpbin struct {
	URL  string `json:"url,omitempty"`
	Data string `json:"data,omitempty"`
}

//
// uuid declares HTTP I/O. Its result is returned via id variable.
func uuid(id *id) gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://httpbin.org/uuid"),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.ServedJSON(),
		ƒ.Recv(id),
	)
}

//
// post declares HTTP I/O. The HTTP request requires uuid.
// Its result is returned via doc variable.
func post(uuid *id, doc *httpbin) gurl.Arrow {
	return gurl.HTTP(
		ø.POST("https://httpbin.org/post"),
		ø.AcceptJSON(),
		ø.ContentJSON(),
		ø.Send(&uuid.UUID),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(doc),
	)
}

//
// hof is a high-order function. It is composed from atomic HTTP I/O into
// the chain of requests. HoF returns results via val variable
func hof(val *string) gurl.Arrow {
	// HoF requires internal state
	var (
		id  id
		doc httpbin
	)
	//
	// HoF combines HTTP requests to
	//  * https://httpbin.org/uuid
	//  * https://httpbin.org/post
	//
	// results of HTTP I/O is persisted in the internal state
	return gurl.Join(
		uuid(&id),
		post(&id, &doc),
		// results of HTTP chain is mapped to return value
		ƒ.FMap(func() error {
			*val = doc.Data
			return nil
		}),
	)
}

func eval() {
	var val string
	http := hof(&val)

	if err := http(gurl.IO()).Fail; err != nil {
		fmt.Printf("fail %v\n", err)
	}
	fmt.Printf("==> %v\n", val)
}

func main() {
	for i := 0; i < 3; i++ {
		eval()
	}
}
