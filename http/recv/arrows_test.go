package recv_test

import (
	"errors"
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

type Seq []Test

func (c Seq) Len() int                { return len(c) }
func (c Seq) Swap(i, j int)           { c[i], c[j] = c[j], c[i] }
func (c Seq) Less(i, j int) bool      { return c[i].Site < c[j].Site }
func (c Seq) String(i int) string     { return c[i].Site }
func (c Seq) Value(i int) interface{} { return c[i] }

func TestCodeOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).Should().Equal(nil)
}

func TestCodeNoMatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.Accept().Is("text/plain"),
		ƒ.Code(gurl.StatusCodeOK),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).ShouldNot().Equal(nil).
		If(io.Fail).Should().Equal(gurl.NewStatusCode(400, gurl.StatusCodeOK))
}

func TestHeaderOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.ServedJSON(),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestHeaderAny(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Header("Content-Type").Any(),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestHeaderValue(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var mime string
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Header("Content-Type").String(&mime),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).Should().Equal(nil).
		If(mime).Should().Equal("application/json")
}

func TestHeaderFail(t *testing.T) {
	ts := mock()
	defer ts.Close()

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.ServedForm(),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).ShouldNot().Equal(nil).
		If(io.Fail).Should().Equal(
		&gurl.Mismatch{
			Diff:    "+ Content-Type: application/json\n- Content-Type: application/x-www-form-urlencoded",
			Payload: map[string]string{"Content-Type": "application/json"},
		},
	)
}

func TestRecvJSON(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.ServedJSON(),
		ƒ.Recv(&data),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).Should().Equal(nil).
		If(data.Site).Should().Equal("example.com")
}

func TestRecvForm(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptForm(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.ServedForm(),
		ƒ.Recv(&data),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).Should().Equal(nil).
		If(data.Site).Should().Equal("example.com")

}

func TestBytes(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data []byte
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.ServedJSON(),
		ƒ.Bytes(&data),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).Should().Equal(nil).
		If(string(data)).Should().Equal("{\"site\": \"example.com\"}")
}

func TestDefined(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
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
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
		ƒ.Defined(&data.Host),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(&gurl.Undefined{Type: "string"})
}

func TestValue(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
		ƒ.Value(&data).Is(&Test{Site: "example.com"}),
		ƒ.Value(&data.Site).String("example.com"),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestValueBytes(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	var octet []byte
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
		ƒ.FlatMap(func() gurl.Arrow {
			octet = []byte(data.Site)
			return nil
		}),
		ƒ.Value(&octet).Bytes([]byte("example.com")),
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
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
		ƒ.Value(&data.Site).String("localhost"),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).Should().Be().Like(&gurl.Mismatch{})
}

func TestSeqHas(t *testing.T) {
	ts := mockSeq()
	defer ts.Close()

	var data Seq
	expectS := Test{Site: "s.example.com"}
	expectZ := Test{Site: "z.example.com"}

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
		ƒ.Seq(&data).Has(expectS.Site),
		ƒ.Seq(&data).Has(expectS.Site, expectS),
		ƒ.Seq(&data).Has(expectZ.Site, expectZ),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Equal(nil)
}

func TestSeqHasFailure(t *testing.T) {
	ts := mockSeq()
	defer ts.Close()

	var data Seq
	expect0 := Test{Site: "0.example.com"}

	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
		ƒ.Seq(&data).Has(expect0.Site),
	)(gurl.IO())

	it.Ok(t).If(io.Fail).Should().Be().Like(&gurl.Undefined{})
}

func TestSeqHasNoMatch(t *testing.T) {
	ts := mockSeq()
	defer ts.Close()

	var data Seq
	expectS := Test{Site: "s.example.com"}
	expectZ := Test{Site: "z.example.com"}

	err := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
		ƒ.Seq(&data).Has(expectS.Site, expectZ),
	)(gurl.IO()).Fail

	it.Ok(t).
		If(strings.Contains(err.Error(), `Site: "s.example.com"`)).Should().Equal(true).
		If(strings.Contains(err.Error(), `Site: "z.example.com"`)).Should().Equal(true)
}

func TestFMap(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
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

func TestFMapFailure(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	io := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
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

func TestFlatMap(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var data Test
	http := gurl.HTTP(
		ø.GET(ts.URL),
		ø.AcceptJSON(),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&data),
		ƒ.Value(&data).Is(&Test{Site: "example.com"}),
	)

	io := gurl.Join(
		http,
		ƒ.FlatMap(func() gurl.Arrow { return http }),
	)(gurl.IO())

	it.Ok(t).
		If(io.Fail).Should().Equal(nil)
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

func mockSeq() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`[
				{"site": "q.example.com"},
				{"site": "a.example.com"},
				{"site": "z.example.com"},
				{"site": "w.example.com"},
				{"site": "s.example.com"},
				{"site": "x.example.com"},
				{"site": "e.example.com"},
				{"site": "d.example.com"},
				{"site": "c.example.com"}
			]`))
		}),
	)
}
