package send_test

import (
	"testing"

	"github.com/fogfish/gurl"
	ø "github.com/fogfish/gurl/http/send"
	"github.com/fogfish/it"
)

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
	mthd := []func(string, ...interface{}) gurl.Arrow{ø.GET, ø.POST, ø.PUT, ø.DELETE}
	for _, f := range mthd {
		io := f("https://example.com")(gurl.IO())
		it.Ok(t).
			If(io.Fail).Should().Equal(nil).
			If(io.URL.String()).Should().Equal("https://example.com")
	}
}

func TestURL(t *testing.T) {
	mthd := []func(string, ...interface{}) gurl.Arrow{ø.GET, ø.POST, ø.PUT, ø.DELETE}
	for _, f := range mthd {
		io := f("https://example.com/%v", 1)(gurl.IO())
		it.Ok(t).
			If(io.Fail).Should().Equal(nil).
			If(io.URL.String()).Should().Equal("https://example.com/1")
	}
}

func TestHeaderIs(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("http://example.com"),
		ø.Header("Accept").Is("text/plain"),
	)(gurl.IO())

	it.Ok(t).
		If(*io.HTTP.Header["Accept"]).Should().Equal("text/plain")
}

func TestHeaderVal(t *testing.T) {
	val := "text/plain"

	io := gurl.HTTP(
		ø.GET("http://example.com"),
		ø.Header("Accept").Val(&val),
	)(gurl.IO())

	it.Ok(t).
		If(*io.HTTP.Header["Accept"]).Should().Equal("text/plain")
}

func TestHeaderAccept(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("http://example.com"),
		ø.Accept().Is("text/plain"),
	)(gurl.IO())

	it.Ok(t).
		If(*io.HTTP.Header["Accept"]).Should().Equal("text/plain")
}

func TestHeaderAcceptJSON(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("http://example.com"),
		ø.AcceptJSON(),
	)(gurl.IO())

	it.Ok(t).
		If(*io.HTTP.Header["Accept"]).Should().Equal("application/json")
}
func TestHeaderAcceptForm(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("http://example.com"),
		ø.AcceptForm(),
	)(gurl.IO())

	it.Ok(t).
		If(*io.HTTP.Header["Accept"]).Should().Equal("application/x-www-form-urlencoded")
}

func TestHeaderContent(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("http://example.com"),
		ø.Content().Is("text/plain"),
	)(gurl.IO())

	it.Ok(t).
		If(*io.HTTP.Header["Content-Type"]).Should().Equal("text/plain")
}

func TestHeaderContentJSON(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("http://example.com"),
		ø.ContentJSON(),
	)(gurl.IO())

	it.Ok(t).
		If(*io.HTTP.Header["Content-Type"]).Should().Equal("application/json")
}

func TestHeaderContentForm(t *testing.T) {
	io := gurl.HTTP(
		ø.GET("http://example.com"),
		ø.ContentForm(),
	)(gurl.IO())

	it.Ok(t).
		If(*io.HTTP.Header["Content-Type"]).Should().Equal("application/x-www-form-urlencoded")
}
