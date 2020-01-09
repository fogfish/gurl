//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package gurl_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fogfish/gurl"
	"github.com/fogfish/it"
)

type Test struct {
	Site string `json:"site"`
	Host string `json:"host,omitempty"`
}

func TestSchemaHTTP(t *testing.T) {
	io := gurl.IO().URL("GET", "http://example.com")

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestSchemaHTTPS(t *testing.T) {
	io := gurl.IO().URL("GET", "https://example.com")

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestSchemaUnsupported(t *testing.T) {
	io := gurl.IO().URL("GET", "other://example.com")

	it.Ok(t).
		If(io.Fail).ShouldNot().Equal(nil).
		If(io.Fail).Should().Equal(&gurl.BadSchema{"other"})
}

func TestWith(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.IO().
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200)

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestSend(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var in Test
			defer r.Body.Close()
			err := json.NewDecoder(r.Body).Decode(&in)
			if err == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
	defer ts.Close()

	io := gurl.IO().
		POST(ts.URL).
		With(gurl.ContentType, gurl.ApplicationJson).
		Send(Test{"example.com", ""}).
		Code(200)

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestCodeOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.IO().
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200)

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestCodeFail(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer ts.Close()

	io := gurl.IO().
		GET(ts.URL).
		Code(200)

	it.Ok(t).
		If(io.Fail).ShouldNot().Equal(nil).
		If(io.Fail).Should().Equal(&gurl.BadMatchCode{[]int{200}, 404})
}

func TestHeadOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.IO().
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationJson)

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestHeadAny(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.IO().
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Head(gurl.ContentType, "*")

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestHeadFail(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.IO().
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationForm)

	it.Ok(t).
		If(io.Fail).ShouldNot().Equal(nil).
		If(io.Fail).Should().Equal(&gurl.BadMatchHead{gurl.ContentType, gurl.ApplicationForm, gurl.ApplicationJson})
}

func TestRecv(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.IO().
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationJson).
		Recv(&data)

	it.Ok(t).
		If(io.Fail).Should().Equal(nil).
		If(data.Site).Should().Equal("example.com")
}

func TestSeq(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.IO().
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationJson).
		Recv(&data).
		//
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationJson).
		Recv(&data).
		//
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationJson).
		Recv(&data)

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestDefined(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.IO().
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Recv(&data).
		Defined(data.Site)

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestNotDefined(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.IO().
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Recv(&data).
		Defined(data.Host)

	it.Ok(t).If(io.Fail).Should().Equal(&gurl.Undefined{"string"})
}

func TestRequire(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.IO().
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Recv(&data).
		Require(data.Site, "example.com")

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestRequireFail(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.IO().
		GET(ts.URL).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Recv(&data).
		Require(data.Site, "localhost")

	it.Ok(t).
		If(io.Fail).Should().
		Equal(&gurl.BadMatch{"localhost", "example.com"})
}

type HoF struct {
	*gurl.IOCat
}

func TestHoF(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := &HoF{gurl.IO()}
	val := io.doThis(ts.URL)
	io.doThat(ts.URL, val)

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func (io *HoF) doThis(url string) (data Test) {
	io.GET(url).
		With(gurl.Accept, gurl.ApplicationJson).
		Code(200).
		Recv(&data)
	return
}

func (io *HoF) doThat(url string, user Test) (data Test) {
	io.PUT(url).
		With(gurl.Accept, gurl.ApplicationJson).
		With(gurl.ContentType, gurl.ApplicationJson).
		Send(user).
		Code(200).
		Recv(&data)
	return
}

//
func mock() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(gurl.Accept) == gurl.ApplicationJson {
				w.Header().Add("Content-Type", "application/json")
				w.Write([]byte(`{"site": "example.com"}`))
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
}
