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
	"image"
	"image/jpeg"
	"os"

	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

type Heap struct {
	image image.Image
}

// declares http I/O
func (h *Heap) request() http.Arrow {
	return http.GET(
		// specify specify the request
		ø.URI("https://avatars.githubusercontent.com/u/716093"),
		ø.Accept.Set("image/*"),

		// specify requirements to the response
		ƒ.Status.OK,
		ƒ.ContentType.Is("image/jpeg"),
		ƒ.Body(&h.image),
	)
}

func main() {
	// instance of http stack
	stack := http.New(http.WithDebugPayload())

	// declares http i/o
	heap := &Heap{}
	lazy := heap.request()

	// executes http I/O
	err := stack.IO(context.Background(), lazy)
	if err != nil {
		panic(err)
	}

	// process image
	jpeg.Encode(os.Stdout, heap.image, &jpeg.Options{Quality: 93})
}
