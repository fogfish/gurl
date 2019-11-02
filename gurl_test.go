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
)

type Test struct {
	Site string `json:"site"`
}

func TestSchemaHTTP(t *testing.T) {
	io := gurl.NewIO().URL("GET", "http://example.com")

	if io.Fail != nil {
		t.Error(io.Fail)
	}
}

func TestSchemaHTTPS(t *testing.T) {
	io := gurl.NewIO().URL("GET", "https://example.com")

	if io.Fail != nil {
		t.Error(io.Fail)
	}
}

func TestSchemaUnsupported(t *testing.T) {
	io := gurl.NewIO().URL("GET", "other://example.com")

	if io.Fail == nil {
		t.Errorf("accepted unsupported schema.")
	}

	schema := io.Fail.(*gurl.BadSchema).Schema
	if schema != "other" {
		t.Errorf("invalid schema %v at error", schema)
	}
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

	if io.Fail != nil {
		t.Errorf("Unable to define header")
	}
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

	if io.Fail != nil {
		t.Errorf("failed to set a payload")
	}
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

	if io.Fail != nil {
		t.Error(io.Fail)
	}
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

	if io.Fail == nil {
		t.Errorf("failed to match status code")
	}

	code := io.Fail.(*gurl.BadMatchCode).Actual
	if code != http.StatusNotFound {
		t.Errorf("invalid status code %v at error", code)
	}
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

	if io.Fail != nil {
		t.Error(io.Fail)
	}
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

	if io.Fail != nil {
		t.Error(io.Fail)
	}
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

	if io.Fail == nil {
		t.Errorf("failed to match http header")
	}

	head := io.Fail.(*gurl.BadMatchHead).Header
	if head != gurl.ContentType {
		t.Errorf("invalid header %v at error", head)
	}

	code := io.Fail.(*gurl.BadMatchHead).Actual
	if code != gurl.ApplicationJson {
		t.Errorf("invalid header value %v at error", code)
	}
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

	if io.Fail != nil {
		t.Error(io.Fail)
	}
	if data.Site != "example.com" {
		t.Errorf("unable to decode json payload")
	}
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

	if io.Fail != nil {
		t.Error(io.Fail)
	}
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

	if io.Fail != nil {
		t.Error(io.Fail)
	}
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
