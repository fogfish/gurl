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

// basic declarative request
func request(ctx context.Context, cat http.Stack) error {
	return cat.IO(ctx,
		http.GET(
			ø.URI("https://httpbin.org/get"),
			ø.Accept.JSON,
			ø.UserAgent.Set("gurl"),
			ƒ.Status.OK,
		),
	)
}

func main() {
	cat := http.New()

	a := &tracer{}
	err := request(a.Context(context.Background()), cat)
	if err != nil {
		fmt.Printf("fail %v\n", err)
	}
	println("===> ")

	b := &tracer{}
	err = request(b.Context(context.Background()), cat)
	if err != nil {
		fmt.Printf("fail %v\n", err)
	}
	println("===> ")

	c := &tracer{}
	err = request(c.Context(context.Background()), cat)
	if err != nil {
		fmt.Printf("fail %v\n", err)
	}

	fmt.Printf("==> \n")
}
