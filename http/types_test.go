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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	µ "github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
	"github.com/fogfish/it/v2"
)

func TestJoin(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/ok", ø.Authority(ts.URL)),
		ƒ.Code(µ.StatusOK),
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
			case r.URL.Path == "/ok":
				w.WriteHeader(http.StatusOK)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
}
