// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl

package main

// The example shows recursion of HTTP. The recurion is demonstarted as
// sequential retrieval of content until EOF.
//
// In pure functional environment the recursion can be defined as
//
// lookup(Page) ->
//   [m_state ||
//     Head <- request(Token, Url, Page),
//     Tail <- untilEOF(Head, Token, Url, Page),
//     cats:unit(Head ++ Tail)
//   ].

import (
	"context"
	"fmt"

	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

// repo is a payload returned by api
type repo struct {
	Name string `json:"name"`
}

// sequence is a collection accumulated while recursion is evaluated
type seq []repo

// request declares HTTP I/O that fetches a portion (page) from api
func request(page int) (*seq, http.Arrow) {
	var seq seq
	return &seq, http.GET(
		ø.URI("https://api.github.com/users/fogfish/repos"),
		ø.Param("type", "all"),
		ø.Param("page", page),
		ø.Accept.JSON,
		ƒ.Status.OK,
		ƒ.Body(&seq),
	)
}

// HoF recursively composes HTTP I/O until all data is fetched.
// The request is returned via seq variable.
func lookup(cat http.Stack, page int) (seq, error) {
	// internal state to accumulate results of HTTP I/O
	var val seq

	pid := page
	for {
		h, lazy := request(pid)
		err := cat.IO(context.Background(), lazy)
		if err != nil {
			return nil, err
		}

		if len(*h) == 0 {
			return val, nil
		}

		pid = pid + 1
		val = append(val, *h...)
	}
}

func main() {
	cat := http.New()
	val, err := lookup(cat, 1)

	if err != nil {
		fmt.Printf("fail %v\n", err)
	}
	fmt.Printf("==> %v\n", val)
}
