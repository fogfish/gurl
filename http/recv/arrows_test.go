//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package recv_test

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

func TestCodeOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
	)
	cat := gurl.IO(µ.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil)
}

func TestCodeNoMatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET(ts.URL+"/other"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
	)
	cat := gurl.IO(µ.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Be().Like(µ.StatusCodeBadRequest)
}

func TestHeaderOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.Header("content-type").Is("application/json"),
	)
	cat := gurl.IO(µ.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil)
}

func TestHeaderAny(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.Header("content-type").Any(),
	)
	cat := gurl.IO(µ.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil)
}

func TestHeaderVal(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var content string
	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.Header("content-type").String(&content),
	)
	cat := gurl.IO(µ.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil).
		If(content).Should().Equal("application/json")
}

func TestHeaderMismatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.Header("content-type").Is("foo/bar"),
	)
	cat := gurl.IO(µ.Default())

	it.Ok(t).
		If(req(cat).Fail).ShouldNot().Equal(nil)
}

func TestHeaderUndefinedWithLit(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.Header("x-content-type").Is("foo/bar"),
	)
	cat := gurl.IO(µ.Default())

	it.Ok(t).
		If(req(cat).Fail).ShouldNot().Equal(nil)
}

func TestHeaderUndefinedWithVal(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var val string
	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.Header("x-content-type").String(&val),
	)
	cat := gurl.IO(µ.Default())

	it.Ok(t).
		If(req(cat).Fail).ShouldNot().Equal(nil)
}

func TestRecvJSON(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
	}

	ts := mock()
	defer ts.Close()

	var site Site
	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.ServedJSON(),
		ƒ.Recv(&site),
	)
	cat := gurl.IO(µ.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil).
		If(site.Site).Should().Equal("example.com")
}

func TestRecvForm(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
	}

	ts := mock()
	defer ts.Close()

	var site Site
	req := µ.Join(
		ø.GET(ts.URL+"/form"),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.ServedForm(),
		ƒ.Recv(&site),
	)
	cat := gurl.IO(µ.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil).
		If(site.Site).Should().Equal("example.com")
}

func TestRecvBytes(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
	}

	ts := mock()
	defer ts.Close()

	var data []byte
	req := µ.Join(
		ø.GET(ts.URL+"/form"),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.Served().Any(),
		ƒ.Bytes(&data),
	)
	cat := gurl.IO(µ.Default())

	it.Ok(t).
		If(req(cat).Fail).Should().Equal(nil).
		If(string(data)).Should().Equal("site=example.com")
}

//
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
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
}
