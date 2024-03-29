//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/png"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	µ "github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
	"github.com/fogfish/it/v2"
)

func TestMethod(t *testing.T) {
	cat := µ.New()

	for expect, method := range map[string]func(arrows ...µ.Arrow) µ.Arrow{
		"GET":     µ.GET,
		"HEAD":    µ.HEAD,
		"POST":    µ.POST,
		"PUT":     µ.PUT,
		"DELETE":  µ.DELETE,
		"PATCH":   µ.PATCH,
		"OPTIONS": func(arrows ...µ.Arrow) µ.Arrow { return µ.Join(ø.Method("OPTIONS"), µ.Join(arrows...)) },
	} {
		cat := cat.WithContext(context.Background())
		err := cat.IO(method(ø.URI("https://example.com")))
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(cat.Method, expect),
			it.Equal(cat.Request.Method, expect),
		)
	}
}

func TestStackOptions(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/ok", ø.Authority(ts.URL)),
		ƒ.Code(µ.StatusOK),
	)
	cat := µ.New(
		µ.WithCookieJar(),
		µ.WithInsecureTLS(),
	)
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
	)
}

func TestJoin(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/opts", ø.Authority(ts.URL)),
		µ.Join(
			ø.Param("a", "1"),
			ø.Param("b", "2"),
		),
		ƒ.Code(µ.StatusOK),
		ƒ.Match(`{"opts": "a=1&b=2"}`),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
	)
}

func TestJoinCats(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		µ.GET(
			ø.URI("%s/ok", ø.Authority(ts.URL)),
			ƒ.Status.OK,
		),
		µ.GET(
			ø.URI(ts.URL),
			ƒ.Code(µ.StatusBadRequest),
		),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
	)
}

type opt struct{ key, val string }

func (opt opt) Arrow() µ.Arrow { return ø.Param(opt.key, opt.val) }

type err string

func (err err) Arrow() µ.Arrow { return func(ctx *µ.Context) error { return fmt.Errorf("%s", err) } }

func TestBind(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/opts", ø.Authority(ts.URL)),
		µ.Bind(
			opt{"a", "1"},
			opt{"b", "2"},
		),
		ƒ.Code(µ.StatusOK),
		ƒ.Match(`{"opts": "a=1&b=2"}`),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
	)
}

func TestBindFailed(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/opts", ø.Authority(ts.URL)),
		µ.Bind(
			opt{"a", "1"},
			err("failed"),
		),
		ƒ.Code(µ.StatusOK),
		ƒ.Match(`{"opts": "a=1&b=2"}`),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).ShouldNot(
		it.Nil(err),
	)
}

func TestIOWithContext(t *testing.T) {
	ts := mock()
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	req := µ.GET(
		ø.URI("%s/ok", ø.Authority(ts.URL)),
		ƒ.Status.OK,
	)

	cat := µ.New()
	err := cat.IO(ctx, req)

	it.Then(t).ShouldNot(
		it.Nil(err),
	)

	cancel()
}

func TestIOWithContextCancel(t *testing.T) {
	ts := mock()
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Nanosecond)
	req := µ.GET(
		ø.URI("%s/ok", ø.Authority(ts.URL)),
		ƒ.Status.OK,
	)

	cat := µ.New()
	cancel()
	err := cat.IO(ctx, req)

	it.Then(t).ShouldNot(
		it.Nil(err),
	)
}

func TestIO(t *testing.T) {
	ts := mock()
	defer ts.Close()

	type Site struct {
		Site string `json:"site"`
	}

	cat := µ.New()

	t.Run("JSON", func(t *testing.T) {
		val, err := µ.IO[Site](cat.WithContext(context.Background()),
			µ.GET(
				ø.URI("%s/json", ø.Authority(ts.URL)),
				ƒ.Status.OK,
			),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(val.Site, "example.com"),
		)
	})

	t.Run("Form", func(t *testing.T) {
		val, err := µ.IO[Site](cat.WithContext(context.Background()),
			µ.GET(
				ø.URI("%s/form", ø.Authority(ts.URL)),
				ƒ.Status.OK,
			),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(val.Site, "example.com"),
		)
	})

	t.Run("Image", func(t *testing.T) {
		val, err := µ.IO[image.Image](cat.WithContext(context.Background()),
			µ.GET(
				ø.URI("%s/image", ø.Authority(ts.URL)),
				ƒ.Status.OK,
				ƒ.ContentType.Is("image/png"),
			),
		)

		it.Then(t).Should(
			it.Nil(err),
		).ShouldNot(
			it.Nil(val),
		)
	})
}

func mock() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/json":
				w.Header().Add("Content-Type", "application/json")
				w.Write([]byte(`{"site": "example.com"}`))
			case r.URL.Path == "/form":
				w.Header().Add("Content-Type", "application/x-www-form-urlencoded")
				w.Write([]byte("site=example.com"))
			case strings.HasPrefix(r.URL.Path, "/image"):
				w.Header().Add("Content-Type", "image/png")
				dst, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEUAAACnej3aAAAAAXRSTlMAQObYZgAAAApJREFUCNdjYAAAAAIAAeIhvDMAAAAASUVORK5CYII=")
				if err != nil {
					panic(err)
				}
				w.Write(dst)
			case r.URL.Path == "/ok":
				w.WriteHeader(http.StatusOK)
			case r.URL.Path == "/opts":
				w.Header().Add("Content-Type", "application/json")
				w.Write([]byte(`{"opts": "` + r.URL.RawQuery + `"}`))
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
}
