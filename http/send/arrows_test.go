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
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/fogfish/gurl"
	"github.com/fogfish/gurl/http"
	ø "github.com/fogfish/gurl/http/send"
	"github.com/fogfish/it"
)

func TestSchemaHTTP(t *testing.T) {
	req := ø.GET.URL("http://example.com")
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil)
}

func TestSchemaHTTPS(t *testing.T) {
	req := ø.GET.URL("https://example.com")
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil)
}

func TestSchemaUnsupported(t *testing.T) {
	req := ø.GET.URL("other://example.com")
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).ShouldNot().Equal(nil)
}

func TestURL(t *testing.T) {
	req := ø.GET.URL("https://example.com/%s/%v", "a", 1)
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil).
		If(cat.HTTP.Send.URL.String()).Should().Equal("https://example.com/a/1")
}

func TestURLByRef(t *testing.T) {
	a := "a"
	b := 1
	req := ø.GET.URL("https://example.com/%s/%v", &a, &b)
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil).
		If(cat.HTTP.Send.URL.String()).Should().Equal("https://example.com/a/1")
}

func TestURLEscape(t *testing.T) {
	a := "a b"
	b := 1
	req := ø.GET.URL("https://example.com/%s/%v", &a, &b)
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil).
		If(cat.HTTP.Send.URL.String()).Should().Equal("https://example.com/a%20b/1")
}

func TestURLEscapeSkip(t *testing.T) {
	a := "a/b"
	req := ø.GET.URL("!https://example.com/%s", &a)
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil).
		If(cat.HTTP.Send.URL.String()).Should().Equal("https://example.com/a/b")
}

func TestURLType(t *testing.T) {
	a := "a b"
	b := 1
	p, _ := url.Parse("https://example.com")
	req := ø.GET.URL("%s/%s/%v", p, &a, &b)
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil).
		If(cat.HTTP.Send.URL.String()).Should().Equal("https://example.com/a%20b/1")
}

func TestURLLazyVal(t *testing.T) {
	a := func() string { return "a" }

	req := ø.GET.URL("https://example.com/%s", a)
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil).
		If(cat.HTTP.Send.URL.String()).Should().Equal("https://example.com/a")
}

func TestHeaders(t *testing.T) {
	defAccept := "text/plain"
	defClose := "close"

	for val, arr := range map[*[]string]http.Arrow{
		//
		{"accept", "text/plain"}:                        ø.Header("Accept").Is("text/plain"),
		{"accept", "text/plain"}:                        ø.Header("Accept").Val(&defAccept),
		{"accept", "text/plain"}:                        ø.Accept.Text,
		{"accept", "application/json"}:                  ø.Accept.JSON,
		{"accept", "application/x-www-form-urlencoded"}: ø.Accept.Form,

		{"accept", "text/plain"}: ø.Accept.Is("text/plain"),
		{"accept", "text/plain"}: ø.Accept.Val(&defAccept),
		//
		{"connection", "keep-alive"}: ø.Connection.KeepAlive,
		{"connection", "close"}:      ø.Connection.Is("close"),
		{"connection", "close"}:      ø.Connection.Val(&defClose),
		//
		{"authorization", "foo bar"}: ø.Authorization.Is("foo bar"),
	} {
		req := http.Join(
			ø.GET.URL("http://example.com"),
			arr,
		)
		cat := gurl.IO(http.Default())

		it.Ok(t).
			If(req(cat).Fail).Should().Equal(nil).
			If(*cat.HTTP.Send.Header[(*val)[0]]).Equal((*val)[1])
	}
}

func TestParams(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host string `json:"host,omitempty"`
	}

	req := http.Join(
		ø.GET.URL("https://example.com"),
		ø.Params(Site{"host", "site"}),
	)
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil).
		If(cat.HTTP.Send.URL.String()).Should().Equal("https://example.com?host=site&site=host")
}

func TestParamsInvalidFormat(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host int    `json:"host,omitempty"`
	}

	req := http.Join(
		ø.GET.URL("https://example.com"),
		ø.Params(Site{"host", 100}),
	)
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).ShouldNot().Equal(nil)
}

func TestSendJSON(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host string `json:"host,omitempty"`
	}

	req := http.Join(
		ø.GET.URL("https://example.com"),
		ø.Header("Content-Type").Is("application/json"),
		ø.Send(Site{"host", "site"}),
	)
	cat := req(http.DefaultIO())
	buf, _ := ioutil.ReadAll(cat.HTTP.Send.Payload)

	it.Ok(t).
		If(cat.Fail).Should().Equal(nil).
		If(string(buf)).Should().Equal("{\"site\":\"host\",\"host\":\"site\"}")
}

func TestSendForm(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host string `json:"host,omitempty"`
	}

	req := http.Join(
		ø.GET.URL("https://example.com"),
		ø.Header("Content-Type").Is("application/x-www-form-urlencoded"),
		ø.Send(Site{"host", "site"}),
	)
	cat := req(http.DefaultIO())
	buf, _ := ioutil.ReadAll(cat.HTTP.Send.Payload)

	it.Ok(t).
		If(cat.Fail).Should().Equal(nil).
		If(string(buf)).Should().Equal("host=site&site=host")
}

func TestSendBytes(t *testing.T) {
	for _, content := range []http.Arrow{
		ø.ContentType.Text,
		ø.ContentType.HTML,
	} {
		for _, val := range []interface{}{
			"host=site",
			[]byte("host=site"),
			bytes.NewBuffer([]byte("host=site")),
		} {
			req := http.Join(
				ø.GET.URL("https://example.com"),
				content,
				ø.Send(val),
			)
			cat := req(http.DefaultIO())
			buf, _ := ioutil.ReadAll(cat.HTTP.Send.Payload)

			it.Ok(t).
				If(cat.Fail).Should().Equal(nil).
				If(string(buf)).Should().Equal("host=site")
		}
	}
}

func TestSendUnknown(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host string `json:"host,omitempty"`
	}

	req := http.Join(
		ø.GET.URL("https://example.com"),
		ø.Send(Site{"host", "site"}),
	)
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).ShouldNot().Equal(nil)
}

func TestSendNotSupported(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host string `json:"host,omitempty"`
	}

	req := http.Join(
		ø.GET.URL("https://example.com"),
		ø.Header("Content-Type").Is("foo/bar"),
		ø.Send(Site{"host", "site"}),
	)
	cat := gurl.IO(http.Default())

	it.Ok(t).
		If(req(cat).Fail).ShouldNot().Equal(nil)
}

func TestAliasesURL(t *testing.T) {
	for mthd, f := range map[string]func(string, ...interface{}) http.Arrow{
		"GET":    ø.GET.URL,
		"PUT":    ø.PUT.URL,
		"POST":   ø.POST.URL,
		"DELETE": ø.DELETE.URL,
		"PATCH":  ø.PATCH.URL,
	} {
		req := f("https://example.com/%s/%v", "a", 1)
		cat := gurl.IO(http.Default())

		it.Ok(t).
			If(req(cat).Fail).Should().Equal(nil).
			If(cat.HTTP.Send.URL.String()).Should().Equal("https://example.com/a/1").
			If(cat.HTTP.Send.Method).Should().Equal(mthd)
	}
}
