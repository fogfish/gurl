package http_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/fogfish/gurl/v2/http"
	µ "github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
	"github.com/fogfish/it/v2"
)

func TestWriteOnceSuccess(t *testing.T) {
	ts := mock()
	defer ts.Close()

	unittest := func() http.Arrow {
		return µ.GET(
			ø.URI("%s/json", ø.Authority(ts.URL)),
			ƒ.Status.OK,
		)
	}

	buf := bytes.Buffer{}
	err := http.WriteOnce(&buf, unittest)
	it.Then(t).Should(it.Nil(err))

	var seq []http.Status
	err = json.Unmarshal(buf.Bytes(), &seq)
	it.Then(t).Should(
		it.Nil(err),
		it.Equal(len(seq), 1),
		it.Equal(seq[0].ID, "github.com/fogfish/gurl/v2/http_test.TestWriteOnceSuccess.func1"),
		it.Equal(seq[0].Status, "success"),
		it.Equal(seq[0].Payload, `{"site": "example.com"}`),
	)
}

func TestWriteOnceNoMatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	unittest := func() http.Arrow {
		return µ.GET(
			ø.URI("%s/json", ø.Authority(ts.URL)),
			ƒ.Status.OK,
			ƒ.ContentType.Form,
		)
	}

	buf := bytes.Buffer{}
	err := http.WriteOnce(&buf, unittest)
	it.Then(t).Should(it.Nil(err))

	var seq []http.Status
	err = json.Unmarshal(buf.Bytes(), &seq)
	it.Then(t).Should(
		it.Nil(err),
		it.Equal(len(seq), 1),
		it.Equal(seq[0].ID, "github.com/fogfish/gurl/v2/http_test.TestWriteOnceNoMatch.func1"),
		it.Equal(seq[0].Status, "failure"),
		it.Equal(seq[0].Payload, `{"site": "example.com"}`),
		it.Equal(seq[0].Reason, "+ Content-Type: application/json\n- Content-Type: application/x-www-form-urlencoded"),
	)
}
