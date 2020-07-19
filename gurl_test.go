//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package gurl_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
	"github.com/fogfish/it"
)

type Test struct {
	Site string `json:"site"`
	Host string `json:"host,omitempty"`
}

func TestJoin(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	http := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.ServedJSON(),
		ƒ.Recv(&data),
	)
	io := gurl.Join(http, http, http)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestStatusSuccess(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	status := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
	)(gurl.IO()).Status("test")

	it.Ok(t).
		If(status.ID).Should().Equal("test").
		If(status.Status).Should().Equal("success").
		If(status.Payload).Should().Equal(&data)
}

func TestStatusFailure(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	status := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeBadRequest),
		ƒ.Recv(&data),
	)(gurl.IO()).Status("test")

	it.Ok(t).
		If(status.ID).Should().Equal("test").
		If(status.Status).Should().Equal("failure").
		If(status.Reason).Should().Equal(gurl.NewStatusCode(200, gurl.StatusCodeBadRequest).Error())
}

func TestStatusFailureMismatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	status := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
		ƒ.Value(&data).Is(Test{Site: "gurl"}),
	)(gurl.IO()).Status("test")

	it.Ok(t).
		If(status.ID).Should().Equal("test").
		If(status.Status).Should().Equal("failure").
		If(status.Payload).Should().Equiv(&Test{Site: "example.com"}).
		If(strings.Contains(status.Reason, `Site: "example.com"`)).Should().Equal(true).
		If(strings.Contains(status.Reason, `Site: "gurl"`)).Should().Equal(true)
}

func TestOnce(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	http := func() gurl.Arrow {
		return gurl.HTTP(
			ø.GET(ts.URL),
			ø.AcceptJSON(),
			ƒ.Code(gurl.StatusCodeOK),
			ƒ.Recv(&data),
			ƒ.Value(&data).Is(&Test{Site: "example.com"}),
		)
	}
	it.Ok(t).
		If(string(gurl.Once(gurl.Tagged{"test", http}))).
		Should().Equal("[{\"id\":\"test\",\"status\":\"success\",\"duration\":0,\"payload\":{\"site\":\"example.com\"}}]")
}

func TestHoF(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.Join(
		doThis(ts.URL, &data),
		doThat(ts.URL, data, &data),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func doThis(url string, data *Test) gurl.Arrow {
	return gurl.HTTP(
		ø.GET(url),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
	)
}

func doThat(url string, user Test, data *Test) gurl.Arrow {
	return gurl.HTTP(
		ø.PUT(url),
		ø.AcceptJSON(),
		ø.ContentJSON(),
		ø.Send(user),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
	)
}

//
func mock() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Header.Get("Accept") == "application/json":
				w.Header().Add("Content-Type", "application/json")
				w.Write([]byte(`{"site": "example.com"}`))
			case r.Header.Get("Accept") == "application/x-www-form-urlencoded":
				w.Header().Add("Content-Type", "application/x-www-form-urlencoded")
				w.Write([]byte("site=example.com"))
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
}
