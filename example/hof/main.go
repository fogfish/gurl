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
	"github.com/fogfish/gurl/http"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

// id implements payload for https://httpbin.org/uuid
type tID struct {
	UUID string `json:"uuid,omitempty"`
}

// httpbin implements payload for https://httpbin.org/post
type tHttpBin struct {
	URL  string `json:"url,omitempty"`
	Data string `json:"data,omitempty"`
}

type tHoF struct {
	tID
	tHttpBin
}

//
// uuid declares HTTP I/O. Its result is returned via id variable.
func (hof *tHoF) uuid() gurl.Arrow {
	return http.Join(
		ø.GET.URL("https://httpbin.org/uuid"),
		ø.Accept.JSON,
		ƒ.Status.OK,
		ƒ.ContentType.JSON,
		ƒ.Recv(&hof.tID),
	)
}

//
// post declares HTTP I/O. The HTTP request requires uuid.
// Its result is returned via doc variable.
func (hof *tHoF) post() gurl.Arrow {
	return http.Join(
		ø.POST.URL("https://httpbin.org/post"),
		ø.Accept.JSON,
		ø.ContentType.JSON,
		ø.Send(&hof.tID.UUID),
		ƒ.Status.OK,
		ƒ.Recv(&hof.tHttpBin),
	)
}

//
// hof is a high-order function. It is composed from atomic HTTP I/O into
// the chain of requests. HoF returns results via val variable
func hof(api *tHoF) gurl.Arrow {
	//
	// HoF combines HTTP requests to
	//  * https://httpbin.org/uuid
	//  * https://httpbin.org/post
	//
	// results of HTTP I/O is persisted in the internal state
	return gurl.Join(
		api.uuid(),
		api.post(),
	)
}

func eval(cat *gurl.IOCat) {
	var val tHoF
	req := hof(&val)
	cat = req(cat)

	if err := cat.Recover(); err != nil {
		fmt.Printf("fail %v\n", err)
	}
	fmt.Printf("==> %v\n", val)
}

func main() {
	cat := http.DefaultIO(gurl.Logging(3))

	for i := 0; i < 3; i++ {
		eval(cat)
	}
}
