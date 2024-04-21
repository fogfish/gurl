//
// Copyright (C) 2019 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package xhtml_test

import (
	"context"
	"testing"

	"github.com/fogfish/gurl/http/x/xhtml"
	"github.com/fogfish/gurl/v2/http/mock"
	"github.com/fogfish/it/v2"
)

func TestFetch(t *testing.T) {
	web := mock.New(
		mock.Header("Content-Type", "text/html"),
		mock.Body([]byte("<html><body><div>Hello World</div></body></html>")),
	)

	cli := xhtml.New(web)
	c, err := cli.Fetch(context.Background(), "http://example.com/")

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(c.Find("div").Text(), "Hello World"),
	)
}
