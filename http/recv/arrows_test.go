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
	"time"

	µ "github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
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

	for _, arrow := range []µ.Arrow{
		ƒ.Status.OK,
		ƒ.Code(200),
	} {

		req := µ.GET(
			ø.URI("%s/other", ø.Authority(ts.URL)),
			ø.Accept.JSON,
			arrow,
		)
		cat := µ.New()
		err := cat.IO(context.Background(), req)

		it.Then(t).Should(
			it.Equal(err.Error(), "+ Status Code: 400\n- Status Code: 200"),
		)
	}
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
		ƒ.Header("Date", time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)),
		ƒ.Header("X-Value", 1024),
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
		ƒ.ContentType.Any,
		ƒ.Header("Content-Type", "*"),
		ƒ.HeaderOf[string]("Content-Type").Any,
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

	var (
		date    time.Time
		content string
		value   int
	)
	req := µ.GET(
		ø.URI("%s/json", ø.Authority(ts.URL)),
		ø.Accept.JSON,
		ƒ.Status.OK,
		ƒ.ContentType.To(&content),
		ƒ.Header("Date", &date),
		ƒ.Header("X-Value", &value),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(content, "application/json"),
		it.Equal(date.Format(time.RFC1123), "Wed, 01 Feb 2023 10:20:30 UTC"),
		it.Equal(value, 1024),
	)
}

func TestHeaderMismatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var (
		date  time.Time
		value int
	)

	for _, header := range []µ.Arrow{
		ƒ.ContentType.Is("foo/bar"),
		ƒ.Header("X-FOO", &value),
		ƒ.Header("X-FOO", &date),
	} {
		req := µ.GET(
			ø.URI("%s/json", ø.Authority(ts.URL)),
			ø.Accept.JSON,
			ƒ.Status.OK,
			header,
		)
		cat := µ.New()
		err := cat.IO(context.Background(), req)

		it.Then(t).ShouldNot(
			it.Nil(err),
		)
	}
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
		ƒ.ContentType.ApplicationJSON,
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

func TestExpectJSON(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
	}

	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/json", ø.Authority(ts.URL)),
		ƒ.Status.OK,
		ƒ.ContentType.ApplicationJSON,
		ƒ.ContentType.JSON,
		ƒ.Expect(Site{"example.com"}),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).Should(
		it.Nil(err),
	)
}

func TestExpectJSONFailed(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
	}

	ts := mock()
	defer ts.Close()

	req := µ.GET(
		ø.URI("%s/json", ø.Authority(ts.URL)),
		ƒ.Status.OK,
		ƒ.ContentType.ApplicationJSON,
		ƒ.ContentType.JSON,
		ƒ.Expect(Site{"some.com"}),
	)
	cat := µ.New()
	err := cat.IO(context.Background(), req)

	it.Then(t).ShouldNot(
		it.Equal(err.Error(), ""),
	)
}

func TestRecvBytes(t *testing.T) {
	ts := mock()
	defer ts.Close()

	for path, content := range map[string]µ.Arrow{
		"/text":   ƒ.ContentType.Text,
		"/text/1": ƒ.ContentType.TextPlain,
		"/html":   ƒ.ContentType.HTML,
		"/html/2": ƒ.ContentType.TextHTML,
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

func TestMatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	t.Run("Match", func(t *testing.T) {
		for _, pat := range []string{
			`{"a":"_"}`,
			`{"b":"_"}`,
			`{"c":"_"}`,
			`{"d":"_"}`,
			`{"e":"_"}`,
			`{"f":"_"}`,
			`{"a":"a"}`,
			`{"b":101}`,
			`{"c":1.1}`,
			`{"f":true}`,
			`{"a":"a", "b":101, "c":1.1}`,
			`{"d":["a", "b", "c"]}`,
			`{"e":{"a":"_"}}`,
			`{"e":{"a":"a"}}`,
		} {
			req := µ.GET(
				ø.URI("%s/match", ø.Authority(ts.URL)),
				ƒ.Status.OK,
				ƒ.Match(pat),
			)

			cat := µ.New()
			err := cat.IO(context.Background(), req)

			it.Then(t).Should(
				it.Nil(err),
			)
		}
	})

	t.Run("NoMatch", func(t *testing.T) {
		for _, pat := range []string{
			`{"g":"_"}`,
			`{"a":"b"}`,
			`{"a":101}`,
			`{"b":111}`,
			`{"b":"a"}`,
			`{"c":2.1}`,
			`{"c":"a"}`,
			`{"f":false}`,
			`{"f":"false"}`,
			`{"a":"a", "b":111, "c":1.1}`,
			`{"d":["a", "b"]}`,
			`{"d":["a", "d", "c"]}`,
			`{"d":"abc"}`,
			`{"e":{"f":"_"}}`,
			`{"e":{"a":"b"}}`,
			`{"e":"ab"}`,
		} {
			req := µ.GET(
				ø.URI("%s/match", ø.Authority(ts.URL)),
				ƒ.Status.OK,
				ƒ.Match(pat),
			)

			cat := µ.New()
			err := cat.IO(context.Background(), req)

			it.Then(t).ShouldNot(
				it.Nil(err),
			)
		}
	})
}

func mock() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case strings.HasPrefix(r.URL.Path, "/json"):
				w.Header().Add("Content-Type", "application/json")
				w.Header().Add("Date", "Wed, 01 Feb 2023 10:20:30 UTC")
				w.Header().Add("X-Value", "1024")
				w.Write([]byte(`{"site": "example.com"}`))
			case strings.HasPrefix(r.URL.Path, "/form"):
				w.Header().Add("Content-Type", "application/x-www-form-urlencoded")
				w.Write([]byte("site=example.com"))
			case strings.HasPrefix(r.URL.Path, "/text"):
				w.Header().Add("Content-Type", "text/plain")
				w.Write([]byte("site=example.com"))
			case strings.HasPrefix(r.URL.Path, "/html"):
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
			case strings.HasPrefix(r.URL.Path, "/match"):
				w.Header().Add("Content-Type", "application/json")
				w.Write([]byte(`{"a":"a", "b":101, "c":1.1, "d":["a", "b", "c"], "e": {"a":"a", "b":101, "c":1.1}, "f": true}`))
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
