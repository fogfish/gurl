//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//
package main

/*

The example shows recursion of HTTP. The recurion is demonstarted as
sequential retrival of content until EOF.

In pure functional environment the recursion can be defined as

lookup(Page) ->
  [m_state ||
    Head <- request(Token, Url, Page),
    Tail <- untilEOF(Head, Token, Url, Page),
    cats:unit(Head ++ Tail)
  ].

*/

import (
	"fmt"
	"strconv"

	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

// repo is a payload returned by api
type repo struct {
	Name string `json:"name"`
}

// sequence is a collection accumulated while recursion is evaluated
type seq []repo

// request declares HTTP I/O that fetches a portion (page) from api
func (s *seq) request(page int) gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://api.github.com/users/fogfish/repos"),
		ø.Params(map[string]string{"type": "all", "page": strconv.Itoa(page)}),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.Recv(s),
	)
}

// untilEOF declares continuation of HTTP I/O until EOF is reached
func (s *seq) untilEOF(head seq, page int) gurl.Arrow {
	if len(head) == 0 {
		return nil
	}

	// internal state is accumulated and execution of next page is scheduled
	*s = append(*s, head...)
	return s.lookup(page + 1)
}

// HoF recursively composes HTTP I/O until all data is fetched.
// The request is returned via seq variable.
func (s *seq) lookup(page int) gurl.Arrow {
	// internal state to accumulate results of HTTP I/O
	var head seq

	//
	// HoF combines HTTP requests with a logic that continues evaluation.
	return gurl.Join(
		head.request(page),
		ƒ.FlatMap(func() gurl.Arrow { return s.untilEOF(head, page) }),
	)
}

func main() {
	var val seq
	http := val.lookup(1)

	if err := http(gurl.IO()).Fail; err != nil {
		fmt.Printf("fail %v\n", err)
	}
	fmt.Printf("==> %v\n", val)
}
