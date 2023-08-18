//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ajg/form"
	"github.com/fogfish/gurl/v2"
)

// Arrow is a morphism applied to HTTP protocol stack
type Arrow func(*Context) error

type ReadableHeaderValues interface {
	int | string | time.Time
}

type WriteableHeaderValues interface {
	*int | *string | *time.Time
}

type MatchableHeaderValues interface {
	ReadableHeaderValues | WriteableHeaderValues
}

// Join composes HTTP arrows to high-order function
// (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
func Join(arrows ...Arrow) Arrow {
	return func(cat *Context) error {
		for _, f := range arrows {
			if err := f(cat); err != nil {
				return err
			}
		}

		return nil
	}
}

// Bind composes HTTP arrows to high-order function
// In contrast with Join, input is arrow builders
// (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
func Bind(arrows ...interface{ Arrow() Arrow }) Arrow {
	return func(cat *Context) error {
		for _, arrow := range arrows {
			f := arrow.Arrow()
			if err := f(cat); err != nil {
				return err
			}
		}

		return nil
	}
}

// GET composes HTTP arrows to high-order function for HTTP GET request
// (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
func GET(arrows ...Arrow) Arrow { return method(http.MethodGet, arrows) }

// HEAD composes HTTP arrows to high-order function for HTTP HEAD request
// (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
func HEAD(arrows ...Arrow) Arrow { return method(http.MethodHead, arrows) }

// POST composes HTTP arrows to high-order function for HTTP POST request
// (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
func POST(arrows ...Arrow) Arrow { return method(http.MethodPost, arrows) }

// PUT composes HTTP arrows to high-order function for HTTP PUT request
// (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
func PUT(arrows ...Arrow) Arrow { return method(http.MethodPut, arrows) }

// DELETE composes HTTP arrows to high-order function for HTTP DELETE request
// (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
func DELETE(arrows ...Arrow) Arrow { return method(http.MethodDelete, arrows) }

// PATCH composes HTTP arrows to high-order function for HTTP PATCH request
// (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
func PATCH(arrows ...Arrow) Arrow { return method(http.MethodPatch, arrows) }

func method(verb string, arrows []Arrow) Arrow {
	return func(ctx *Context) error {
		ctx.Method = verb
		for _, f := range arrows {
			if err := f(ctx); err != nil {
				return err
			}
		}

		return nil
	}
}

// Executes protocol operation
func IO[T any](ctx *Context, arrows ...Arrow) (*T, error) {
	for _, f := range arrows {
		if err := f(ctx); err != nil {
			return nil, err
		}
	}

	if ctx.Response == nil {
		return nil, fmt.Errorf("empty response")
	}
	defer ctx.Response.Body.Close()

	var val T
	err := decode(
		ctx.Response.Header.Get("Content-Type"),
		ctx.Response.Body,
		&val,
	)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func decode[T any](content string, stream io.ReadCloser, data *T) error {
	switch {
	case strings.Contains(content, "json"):
		return json.NewDecoder(stream).Decode(data)
	case strings.Contains(content, "www-form"):
		return form.NewDecoder(stream).Decode(data)
	default:
		return &gurl.NoMatch{
			ID:       "http.Recv",
			Diff:     fmt.Sprintf("- Content-Type: application/{json | www-form}\n+ Content-Type: %s", content),
			Protocol: "codec",
			Actual:   content,
		}
	}
}
