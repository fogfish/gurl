//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package send_test

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/fogfish/gurl/http"
	ø "github.com/fogfish/gurl/http/send"
	"github.com/fogfish/it/v2"
)

func TestSchema(t *testing.T) {
	cat := http.New()

	t.Run("HTTP", func(t *testing.T) {
		err := cat.IO(context.Background(),
			http.GET(ø.URI("http://example.com")),
		)
		it.Then(t).Should(it.Nil(err))
	})

	t.Run("HTTPS", func(t *testing.T) {
		err := cat.IO(context.Background(),
			http.GET(ø.URI("https://example.com")),
		)
		it.Then(t).Should(it.Nil(err))
	})

	t.Run("Unsupported", func(t *testing.T) {
		err := cat.IO(context.Background(),
			http.GET(ø.URI("other://example.com")),
		)
		it.Then(t).ShouldNot(it.Nil(err))
	})
}

func TestURI(t *testing.T) {
	cat := http.New()

	t.Run("Literal", func(t *testing.T) {
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(ø.URI("https://example.com/a/1")),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(cat.Request.URL.String(), "https://example.com/a/1"),
		)
	})

	t.Run("Format", func(t *testing.T) {
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(ø.URI("https://example.com/%s/%d", "a", 1)),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(cat.Request.URL.String(), "https://example.com/a/1"),
		)
	})

	t.Run("FormatByRef", func(t *testing.T) {
		a, b := "a", 1
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(ø.URI("https://example.com/%s/%d", &a, &b)),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(cat.Request.URL.String(), "https://example.com/a/1"),
		)
	})

	t.Run("Escape", func(t *testing.T) {
		a := "a b"
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(ø.URI("https://example.com/%s/%s", "a b", &a)),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(cat.Request.URL.String(), "https://example.com/a%20b/a%20b"),
		)
	})

	t.Run("NoEscape", func(t *testing.T) {
		a := "a/b"
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(ø.URI("https://example.com/%s/%s", ø.Path("a/b"), (*ø.Path)(&a))),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(cat.Request.URL.String(), "https://example.com/a/b/a/b"),
		)
	})

	t.Run("Authority", func(t *testing.T) {
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(ø.URI("https://%s/%s", ø.Authority("example.com"), ø.Path("a/b"))),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(cat.Request.URL.String(), "https://example.com/a/b"),
		)
	})

	t.Run("url.URL", func(t *testing.T) {
		cat := cat.WithContext(context.Background())
		u, _ := url.Parse("https://example.com/a/b")
		err := cat.IO(
			http.GET(ø.URI("%s/%s", u, "c")),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(cat.Request.URL.String(), "https://example.com/a/b/c"),
		)
	})
}

func TestHeaders(t *testing.T) {
	cat := http.New()

	for val, arr := range map[*[]string]http.Arrow{
		//
		{"accept", "text/plain"}:                        ø.Header("Accept", "text/plain"),
		{"accept", "text/plain"}:                        ø.Accept.Text,
		{"accept", "text/plain"}:                        ø.Accept.TextPlain,
		{"accept", "text/html"}:                         ø.Accept.HTML,
		{"accept", "text/html"}:                         ø.Accept.TextHTML,
		{"accept", "application/json"}:                  ø.Accept.ApplicationJSON,
		{"accept", "application/json"}:                  ø.Accept.JSON,
		{"accept", "application/x-www-form-urlencoded"}: ø.Accept.Form,
		{"accept", "text/plain"}:                        ø.Accept.Set("text/plain"),
		{"connection", "keep-alive"}:                    ø.Connection.KeepAlive,
		{"connection", "close"}:                         ø.Connection.Close,
		{"connection", "close"}:                         ø.Connection.Set("close"),
		{"authorization", "foo bar"}:                    ø.Authorization.Set("foo bar"),
		{"x-value", "1024"}:                             ø.Header("x-value", 1024),
		{"date", "Wed, 01 Feb 2023 10:20:30 UTC"}:       ø.Date.Set(time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)),
	} {
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(
				ø.URI("http://example.com"),
				arr,
			),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(cat.Request.Header.Get((*val)[0]), (*val)[1]),
		)
	}
}

func TestHeaderContentLength(t *testing.T) {
	cat := http.New().WithContext(context.TODO())
	err := cat.IO(
		http.GET(
			ø.URI("http://example.com"),
			ø.ContentLength.Set(1024),
		),
	)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(cat.Request.ContentLength, int64(1024)),
	)
}

func TestHeaderTransferEncoding(t *testing.T) {
	cat := http.New()

	for val, arr := range map[*[]string]http.Arrow{
		{"chunked"}:  ø.TransferEncoding.Chunked,
		{"identity"}: ø.TransferEncoding.Identity,
		{"gzip"}:     ø.TransferEncoding.Set("gzip"),
	} {
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(
				ø.URI("http://example.com"),
				arr,
			),
		)

		it.Then(t).Should(
			it.Nil(err),
			it.Seq(cat.Request.TransferEncoding).Equal(*val...),
		)
	}
}

func TestParams(t *testing.T) {
	cat := http.New()

	t.Run("Struct", func(t *testing.T) {
		type Site struct {
			Site string `json:"site"`
			Host string `json:"host,omitempty"`
		}
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(
				ø.URI("https://example.com"),
				ø.Params(Site{"host", "site"}),
			),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(cat.Request.URL.String(), "https://example.com?host=site&site=host"),
		)
	})

	t.Run("StructInvalid", func(t *testing.T) {
		type Host struct{ Host string }
		type Site struct {
			Host Host `json:"host"`
		}
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(
				ø.URI("https://example.com"),
				ø.Params(Site{Host{"host"}}),
			),
		)
		it.Then(t).ShouldNot(
			it.Nil(err),
		)
	})

	t.Run("KeyVal", func(t *testing.T) {
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(
				ø.URI("https://example.com"),
				ø.Param("host", "site"),
				ø.Param("site", "host"),
			),
		)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(cat.Request.URL.String(), "https://example.com?host=site&site=host"),
		)

	})
}

func TestSend(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host string `json:"host,omitempty"`
	}

	cat := http.New()

	t.Run("Json", func(t *testing.T) {
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(
				ø.URI("https://example.com"),
				ø.ContentType.JSON,
				ø.Send(Site{"host", "site"}),
			),
		)
		buf, _ := io.ReadAll(cat.Request.Body)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(string(buf), "{\"site\":\"host\",\"host\":\"site\"}"),
		)
	})

	t.Run("Form", func(t *testing.T) {
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(
				ø.URI("https://example.com"),
				ø.ContentType.Form,
				ø.Send(Site{"host", "site"}),
			),
		)
		buf, _ := io.ReadAll(cat.Request.Body)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(string(buf), "host=site&site=host"),
		)
	})

	t.Run("Unknown", func(t *testing.T) {
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(
				ø.URI("https://example.com"),
				ø.Send(Site{"host", "site"}),
			),
		)
		it.Then(t).ShouldNot(
			it.Nil(err),
		)
	})

	t.Run("NotSupported", func(t *testing.T) {
		cat := cat.WithContext(context.Background())
		err := cat.IO(
			http.GET(
				ø.URI("https://example.com"),
				ø.ContentType.Set("foo/bar"),
				ø.Send(Site{"host", "site"}),
			),
		)
		it.Then(t).ShouldNot(
			it.Nil(err),
		)
	})
}

func TestSendBytes(t *testing.T) {
	cat := http.New()

	for _, content := range []http.Arrow{
		ø.ContentType.Text,
		ø.ContentType.HTML,
	} {
		for _, val := range []interface{}{
			"host=site",
			strings.NewReader("host=site"),
			[]byte("host=site"),
			bytes.NewBuffer([]byte("host=site")),
			bytes.NewReader([]byte("host=site")),
			io.NopCloser(bytes.NewBuffer([]byte("host=site"))),
		} {
			cat := cat.WithContext(context.Background())
			err := cat.IO(
				http.GET(
					ø.URI("https://example.com"),
					content,
					ø.Send(val),
				),
			)
			buf, _ := io.ReadAll(cat.Request.Body)
			it.Then(t).Should(
				it.Nil(err),
				it.Equal(string(buf), "host=site"),
			)
		}
	}
}

// func TestAliasesURL(t *testing.T) {
// 	for mthd, f := range map[string]func(string, ...interface{}) http.Arrow{
// 		"GET":    ø.GET.URL,
// 		"PUT":    ø.PUT.URL,
// 		"POST":   ø.POST.URL,
// 		"DELETE": ø.DELETE.URL,
// 		"PATCH":  ø.PATCH.URL,
// 	} {
// 		req := f("https://example.com/%s/%v", "a", 1)
// 		cat := http.New().WithContext(context.TODO())

// 		it.Ok(t).
// 			If(cat.IO(req)).Should().Equal(nil).
// 			If(cat.Request.URL.String()).Should().Equal("https://example.com/a/1").
// 			If(cat.Request.Method).Should().Equal(mthd)
// 	}
// }
