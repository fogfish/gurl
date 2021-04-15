//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fogfish/gurl"
	µ "github.com/fogfish/gurl/http"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
	"github.com/fogfish/it"
)

func TestJoin(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET.URL(ts.URL+"/ok"),
		ƒ.Code(µ.StatusOK),
	)
	cat := µ.DefaultIO()

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil)
}

func TestJoinCats(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := gurl.Join(
		µ.Join(
			ø.GET.URL(ts.URL+"/ok"),
			ƒ.Status.OK,
		),
		µ.Join(
			ø.GET.URL(ts.URL),
			ƒ.Code(µ.StatusBadRequest),
		),
	)
	cat := µ.DefaultIO()

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil)
}

//
func mock() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/ok":
				w.WriteHeader(http.StatusOK)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
}
