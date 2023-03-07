//
// Copyright (C) 2019 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
	"github.com/fogfish/it/v2"
)

func TestWriteOnceSuccess(t *testing.T) {
	ts := mock()
	defer ts.Close()

	unittest := func() http.Arrow {
		return http.GET(
			ø.URI("/json"),
			ƒ.Status.OK,
		)
	}

	buf := bytes.Buffer{}
	hts := http.New(http.WithMemento(), http.WithDefaultHost(ts.URL))
	err := http.WriteOnce(&buf, hts, unittest)
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
		return http.GET(
			ø.URI("/json"),
			ƒ.Status.OK,
			ƒ.ContentType.Form,
		)
	}

	buf := bytes.Buffer{}
	hts := http.New(http.WithMemento(), http.WithDefaultHost(ts.URL))
	err := http.WriteOnce(&buf, hts, unittest)
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
