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
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestSchemaHTTP(t *testing.T) {
	io := ø.URL("GET", "http://example.com")(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestSchemaHTTPS(t *testing.T) {
	io := ø.URL("GET", "https://example.com")(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestSchemaUnsupported(t *testing.T) {
	io := ø.URL("GET", "other://example.com")(gurl.IO())

	it.Ok(t).
		If(io.Fail).ShouldNot().Equal(nil).
		If(io.Fail).Should().Equal(&gurl.BadSchema{"other"})
}

func TestMethod(t *testing.T) {
	mthd := []func(string) gurl.Arrow{ø.GET, ø.POST, ø.PUT, ø.DELETE}
	for _, f := range mthd {
		io := f("https://example.com")(gurl.IO())
		it.Ok(t).
			If(io.Fail).Should().Equal(nil)
	}
}

func TestHeader(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.Header("Accept", "application/json"),
		ƒ.Code(200),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).Should().Equal(nil)
}

func TestHeaderAccept(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("https://example.com"),
		ø.Accept("application/json"),
	)(gurl.IO())

	it.Ok(t).
		If(io.HTTP.Header["Accept"]).Should().Equal("application/json")
}

func TestHeaderAcceptJSON(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("https://example.com"),
		ø.AcceptJSON(),
	)(gurl.IO())

	it.Ok(t).
		If(io.HTTP.Header["Accept"]).Should().Equal("application/json")
}

func TestHeaderContent(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("https://example.com"),
		ø.Content("application/json"),
	)(gurl.IO())

	it.Ok(t).
		If(io.HTTP.Header["Content-Type"]).Should().Equal("application/json")
}

func TestHeaderContentJSON(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("https://example.com"),
		ø.ContentJSON(),
	)(gurl.IO())

	it.Ok(t).
		If(io.HTTP.Header["Content-Type"]).Should().Equal("application/json")
}

func TestHeaderKeepAlive(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("https://example.com"),
		ø.KeepAlive(),
	)(gurl.IO())

	it.Ok(t).
		If(io.HTTP.Header["Connection"]).Should().Equal("keep-alive")
}

func TestHeaderAuthorization(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("https://example.com"),
		ø.Authorization("token"),
	)(gurl.IO())

	it.Ok(t).
		If(io.HTTP.Header["Authorization"]).Should().Equal("token")
}

func TestParams(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("https://example.com"),
		ø.Params(Test{"host", "site"}),
	)(gurl.IO())

	it.Ok(t).
		If(io.URL.String()).Should().
		Equal("https://example.com?host=site&site=host")
}

func TestSendJSON(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("https://example.com"),
		ø.ContentJSON(),
		ø.Send(Test{"host", "site"}),
	)(gurl.IO())

	it.Ok(t).
		If(io.HTTP.Payload.String()).Should().
		Equal("{\"site\":\"host\",\"host\":\"site\"}")
}

func TestSendForm(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("https://example.com"),
		ø.ContentForm(),
		ø.Send(Test{"host", "site"}),
	)(gurl.IO())

	it.Ok(t).
		If(io.HTTP.Payload.String()).Should().
		Equal("host=site&site=host")
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

	io := gurl.HTTP(
		ø.POST(ts.URL),
		ø.ContentJSON(),
		ø.Send(Test{"example.com", ""}),
		ƒ.Code(200),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestCodeOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestCodeFail(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer ts.Close()

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ƒ.Code(200),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).ShouldNot().Equal(nil).
		If(io.Fail).Should().Equal(&gurl.BadMatchCode{[]int{200}, 404})
}

func TestHeadOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.ServedJSON(),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestHeadAny(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.Header("Content-Type", "*"),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestHeadFail(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.Served("application/x-www-form-urlencoded"),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).ShouldNot().Equal(nil).
		If(io.Fail).Should().Equal(&gurl.BadMatchHead{"Content-Type", "application/x-www-form-urlencoded", "application/json"})
}

func TestRecv(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&data),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).Should().Equal(nil).
		If(data.Site).Should().Equal("example.com")
}

func TestJoin(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	http := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&data),
	)
	io := gurl.Join(http, http, http)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestDefined(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.Recv(&data),
		ƒ.Defined(&data.Site),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestNotDefined(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.Recv(&data),
		ƒ.Defined(&data.Host),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(&gurl.Undefined{"string"})
}

func TestRequire(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.Recv(&data),
		ƒ.Require(&data.Site, "example.com"),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestRequireFail(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.Recv(&data),
		ƒ.Require(&data.Site, "localhost"),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).Should().Be().Like(&gurl.BadMatch{"localhost", "example.com"})
}

func TestAssert(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.Recv(&data),
		ƒ.FMap(func() (err error) {
			if data.Site != "example.com" {
				err = errors.New("something wrong!")
			}
			return
		}),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestAssertFailure(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.Recv(&data),
		ƒ.FMap(func() (err error) {
			if data.Site == "example.com" {
				err = errors.New("something wrong!")
			}
			return
		}),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).Should().
		Equal(errors.New("something wrong!"))
}

func TestStatusSuccess(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	status := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.Recv(&data),
	)(gurl.IO()).Status("test")

	it.Ok(t).
		If(status.ID).Should().Equal("test").
		If(status.Status).Should().Equal("success").
		If(status.Actual).Should().Equal(&data)
}

func TestStatusFailure(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	status := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(400),
		ƒ.Recv(&data),
	)(gurl.IO()).Status("test")

	it.Ok(t).
		If(status.ID).Should().Equal("test").
		If(status.Status).Should().Equal("failure").
		If(status.Actual).Should().Equal(fmt.Sprint(&gurl.BadMatchCode{[]int{400}, 200}))
}

func TestStatusFailureBadMatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	status := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(200),
		ƒ.Recv(&data),
		ƒ.Require(&data.Site, "gurl"),
	)(gurl.IO()).Status("test")

	it.Ok(t).
		If(status.ID).Should().Equal("test").
		If(status.Status).Should().Equal("failure").
		If(status.Actual).Should().Equiv("example.com").
		If(status.Expect).Should().Equal("gurl")
}

func TestOnce(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	http := func() gurl.Arrow {
		return gurl.HTTP(
			ø.GET(ts.URL),
			ø.AcceptJSON(),
			ƒ.Code(200),
			ƒ.Recv(&data),
			ƒ.Require(&data.Site, "example.com"),
		)
	}
	it.Ok(t).
		If(string(gurl.Once(gurl.Tagged{"test", http}))).
		Should().Equal("[{\"id\":\"test\",\"status\":\"success\",\"duration\":0,\"actual\":{\"site\":\"example.com\"}}]")
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
		ƒ.Code(200),
		ƒ.Recv(&data),
	)
}

func doThat(url string, user Test, data *Test) gurl.Arrow {
	return gurl.HTTP(
		ø.PUT(url),
		ø.AcceptJSON(),
		ø.ContentJSON(),
		ø.Send(user),
		ƒ.Code(200),
		ƒ.Recv(&data),
	)
}

//
func mock() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Accept") == "application/json" {
				w.Header().Add("Content-Type", "application/json")
				w.Write([]byte(`{"site": "example.com"}`))
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
}
