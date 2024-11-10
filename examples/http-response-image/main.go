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

type api struct {
	http.Stack
}

func (api api) request(ctx context.Context) (*image.Image, error) {
	return http.IO[image.Image](api.WithContext(ctx),
		http.GET(
			ø.URI("https://avatars.githubusercontent.com/u/716093"),
			ø.Accept.Set("image/*"),

			ƒ.Status.OK,
			ƒ.ContentType.Is("image/jpeg"),
		),
	)
}

func main() {
	api := api{
		Stack: http.New(http.WithDebugPayload),
	}

	img, err := api.request(context.Background())
	if err != nil {
		panic(err)
	}

	jpeg.Encode(os.Stdout, *img, &jpeg.Options{Quality: 93})
}
