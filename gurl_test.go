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
}

func TestSchemaHTTP(t *testing.T) {
	io := gurl.NewIO().URL("GET", "http://example.com")

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestSchemaHTTPS(t *testing.T) {
	io := gurl.NewIO().URL("GET", "https://example.com")

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestSchemaUnsupported(t *testing.T) {
	io := gurl.NewIO().URL("GET", "other://example.com")

	it.Ok(t).
		If(io.Fail).ShouldNot().Equal(nil).
		If(io.Fail).Should().Equal(&gurl.BadSchema{"other"})
}

func TestWith(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(gurl.Accept) == gurl.ApplicationJson {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
	defer ts.Close()

	io := gurl.NewIO().
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

	io := gurl.NewIO().
		POST(ts.URL).
		With(gurl.ContentType, gurl.ApplicationJson).
		Send(Test{"example.com"}).
		Code(200)

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestCodeOk(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{"site": "example.com"}`))
		}),
	)
	defer ts.Close()

	io := gurl.NewIO().
		GET(ts.URL).
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

	io := gurl.NewIO().
		GET(ts.URL).
		Code(200)

	it.Ok(t).
		If(io.Fail).ShouldNot().Equal(nil).
		If(io.Fail).Should().Equal(&gurl.BadMatchCode{[]int{200}, 404})
}

func TestHeadOk(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{"site": "example.com"}`))
		}),
	)
	defer ts.Close()

	io := gurl.NewIO().
		GET(ts.URL).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationJson)

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestHeadAny(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{"site": "example.com"}`))
		}),
	)
	defer ts.Close()

	io := gurl.NewIO().
		GET(ts.URL).
		Code(200).
		Head(gurl.ContentType, "*")

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestHeadFail(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{"site": "example.com"}`))
		}),
	)
	defer ts.Close()

	io := gurl.NewIO().
		GET(ts.URL).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationForm)

	it.Ok(t).
		If(io.Fail).ShouldNot().Equal(nil).
		If(io.Fail).Should().Equal(&gurl.BadMatchHead{gurl.ContentType, gurl.ApplicationForm, gurl.ApplicationJson})
}

func TestRecv(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{"site": "example.com"}`))
		}),
	)
	defer ts.Close()

	var data Test
	io := gurl.NewIO().
		GET(ts.URL).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationJson).
		Recv(&data)

	it.Ok(t).
		If(io.Fail).Should().Equal(nil).
		If(data.Site).Should().Equal("example.com")
}

func TestSeq(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{"site": "example.com"}`))
		}),
	)
	defer ts.Close()

	var data Test
	io := gurl.NewIO().
		GET(ts.URL).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationJson).
		Recv(&data).
		//
		GET(ts.URL).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationJson).
		Recv(&data).
		//
		GET(ts.URL).
		Code(200).
		Head(gurl.ContentType, gurl.ApplicationJson).
		Recv(&data)

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestHoF(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{"site": "example.com"}`))
		}),
	)
	defer ts.Close()

	io := gurl.NewIO()
	val := doThis(io, ts.URL)
	doThat(io, ts.URL, val)

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func doThis(io *gurl.IO, url string) (data Test) {
	io.GET(url).
		Code(200).
		Recv(&data)
	return
}

func doThat(io *gurl.IO, url string, user Test) (data Test) {
	io.PUT(url).
		With(gurl.ContentType, gurl.ApplicationJson).
		Send(user).
		Code(200).
		Recv(&data)
	return
}
