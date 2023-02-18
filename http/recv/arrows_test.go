//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package recv_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	µ "github.com/fogfish/gurl/http"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
	"github.com/fogfish/it/v2"
)

func TestCodeOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/json", ø.Authority(ts.URL)),
		ø.Accept.JSON,
		ƒ.Code(µ.StatusOK),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
	)
}

func TestCodeNoMatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/other", ø.Authority(ts.URL)),
		ø.Accept.JSON,
		ƒ.Status.OK,
	)
	cat := µ.New()
	var err interface{ StatusCode() int }
	f := func() error { return cat.IO(context.Background(), req) }

	it.Then(t).Should(
		it.Fail(f).With(&err),
		it.Equal(err.(µ.StatusCode).StatusCode(), µ.StatusBadRequest.StatusCode()),
	)
}

func TestStatusCodes(t *testing.T) {
	ts := mock()
	defer ts.Close()

	for code, check := range map[µ.StatusCode]µ.Arrow{
		//
		µ.StatusOK:                   ƒ.Status.OK,
		µ.StatusCreated:              ƒ.Status.Created,
		µ.StatusAccepted:             ƒ.Status.Accepted,
		µ.StatusNonAuthoritativeInfo: ƒ.Status.NonAuthoritativeInfo,
		µ.StatusNoContent:            ƒ.Status.NoContent,
		µ.StatusResetContent:         ƒ.Status.ResetContent,
		//
		µ.StatusMultipleChoices:  ƒ.Status.MultipleChoices,
		µ.StatusMovedPermanently: ƒ.Status.MovedPermanently,
		µ.StatusFound:            ƒ.Status.Found,
		µ.StatusSeeOther:         ƒ.Status.SeeOther,
		µ.StatusNotModified:      ƒ.Status.NotModified,
		µ.StatusUseProxy:         ƒ.Status.UseProxy,
		//
		µ.StatusBadRequest:            ƒ.Status.BadRequest,
		µ.StatusUnauthorized:          ƒ.Status.Unauthorized,
		µ.StatusPaymentRequired:       ƒ.Status.PaymentRequired,
		µ.StatusForbidden:             ƒ.Status.Forbidden,
		µ.StatusNotFound:              ƒ.Status.NotFound,
		µ.StatusMethodNotAllowed:      ƒ.Status.MethodNotAllowed,
		µ.StatusNotAcceptable:         ƒ.Status.NotAcceptable,
		µ.StatusProxyAuthRequired:     ƒ.Status.ProxyAuthRequired,
		µ.StatusRequestTimeout:        ƒ.Status.RequestTimeout,
		µ.StatusConflict:              ƒ.Status.Conflict,
		µ.StatusGone:                  ƒ.Status.Gone,
		µ.StatusLengthRequired:        ƒ.Status.LengthRequired,
		µ.StatusPreconditionFailed:    ƒ.Status.PreconditionFailed,
		µ.StatusRequestEntityTooLarge: ƒ.Status.RequestEntityTooLarge,
		µ.StatusRequestURITooLong:     ƒ.Status.RequestURITooLong,
		µ.StatusUnsupportedMediaType:  ƒ.Status.UnsupportedMediaType,
		//
		µ.StatusInternalServerError:     ƒ.Status.InternalServerError,
		µ.StatusNotImplemented:          ƒ.Status.NotImplemented,
		µ.StatusBadGateway:              ƒ.Status.BadGateway,
		µ.StatusServiceUnavailable:      ƒ.Status.ServiceUnavailable,
		µ.StatusGatewayTimeout:          ƒ.Status.GatewayTimeout,
		µ.StatusHTTPVersionNotSupported: ƒ.Status.HTTPVersionNotSupported,
	} {
		req := µ.GET(
			ø.URI("%s/code/%d", ø.Authority(ts.URL), code.StatusCode()),
			check,
		)
		cat := µ.New()
		err := cat.IO(context.Background(), req)

		it.Then(t).Should(
			it.Nil(err),
		)
	}
}

func TestHeaderOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/json", ø.Authority(ts.URL)),
		ø.Accept.JSON,
		ƒ.Status.OK,
		ƒ.ContentType.JSON,
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
	)
}

func TestHeaderAny(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/json", ø.Authority(ts.URL)),
		ø.Accept.JSON,
		ƒ.Status.OK,
		ƒ.ContentType.Is("*"),
		ƒ.Header("Content-Type", "*"),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
	)
}

func TestHeaderVal(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var content string
	req := µ.GET(
		ø.URI("%s/json", ø.Authority(ts.URL)),
		ø.Accept.JSON,
		ƒ.Status.OK,
		ƒ.ContentType.To(&content),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(content, "application/json"),
	)
}

func TestHeaderMismatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/json", ø.Authority(ts.URL)),
		ø.Accept.JSON,
		ƒ.Status.OK,
		ƒ.ContentType.Is("foo/bar"),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).ShouldNot(
		it.Nil(err),
	)
}

func TestHeaderUndefinedWithLit(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/json", ø.Authority(ts.URL)),
		ø.Accept.JSON,
		ƒ.Status.OK,
		ƒ.Header("x-content-type", "foo/bar"),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).ShouldNot(
		it.Nil(err),
	)
}

func TestHeaderUndefinedWithVal(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var val string
	req := µ.GET(
		ø.URI("%s/json", ø.Authority(ts.URL)),
		ø.Accept.JSON,
		ƒ.Status.OK,
		ƒ.Header("x-content-type", &val),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).ShouldNot(
		it.Nil(err),
	)
}

func TestRecvJSON(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
	}

	ts := mock()
	defer ts.Close()

	var site Site
	req := µ.GET(
		ø.URI("%s/json", ø.Authority(ts.URL)),
		ƒ.Status.OK,
		ƒ.ContentType.JSON,
		ƒ.Recv(&site),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(site.Site, "example.com"),
	)
}

func TestRecvForm(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
	}

	ts := mock()
	defer ts.Close()

	var site Site
	req := µ.GET(
		ø.URI("%s/form", ø.Authority(ts.URL)),
		ƒ.Status.OK,
		ƒ.ContentType.Form,
		ƒ.Recv(&site),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(site.Site, "example.com"),
	)
}

func TestRecvBytes(t *testing.T) {
	ts := mock()
	defer ts.Close()

	for path, content := range map[string]µ.Arrow{
		"/text": ƒ.ContentType.Text,
		"/html": ƒ.ContentType.HTML,
	} {

		var data []byte
		req := µ.GET(
			ø.URI(ts.URL+path),
			ƒ.Status.OK,
			content,
			ƒ.Bytes(&data),
		)
		cat := µ.New()
		err := cat.IO(context.Background(), req)

		it.Then(t).Should(
			it.Nil(err),
			it.Equal(string(data), "site=example.com"),
		)
	}
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
			case r.URL.Path == "/text":
				w.Header().Add("Content-Type", "text/plain")
				w.Write([]byte("site=example.com"))
			case r.URL.Path == "/html":
				w.Header().Add("Content-Type", "text/html")
				w.Write([]byte("site=example.com"))
			case r.URL.Path == "/code/301":
				w.Header().Add("Location", "http://127.1")
				w.WriteHeader(301)
			case r.URL.Path == "/code/302":
				w.Header().Add("Location", "http://127.1")
				w.WriteHeader(302)
			case r.URL.Path == "/code/303":
				w.Header().Add("Location", "http://127.1")
				w.WriteHeader(303)
			case strings.HasPrefix(r.URL.Path, "/code"):
				seq := strings.Split(r.URL.Path, "/")
				code, _ := strconv.Atoi(seq[2])
				w.WriteHeader(code)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
}
