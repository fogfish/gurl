//
// Copyright (C) 2019 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http_test

import (
	µ "github.com/fogfish/gurl/v2/http"
	"github.com/fogfish/it/v2"
	"net/http"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("WithClient", func(t *testing.T) {
		cli := µ.Client()
		cat := µ.New(µ.WithClient(cli)).(*µ.Protocol)
		it.Then(t).Should(it.Equal(cat.Socket.(*http.Client), cli))
	})

	t.Run("WithDebugRequest", func(t *testing.T) {
		cat := µ.New(µ.WithDebugRequest()).(*µ.Protocol)
		it.Then(t).Should(it.Equal(cat.LogLevel, 1))
	})

	t.Run("WithDebugResponse", func(t *testing.T) {
		cat := µ.New(µ.WithDebugResponse()).(*µ.Protocol)
		it.Then(t).Should(it.Equal(cat.LogLevel, 2))
	})

	t.Run("WithDebugPayload", func(t *testing.T) {
		cat := µ.New(µ.WithDebugPayload()).(*µ.Protocol)
		it.Then(t).Should(it.Equal(cat.LogLevel, 3))
	})

	t.Run("WithMemento", func(t *testing.T) {
		cat := µ.New(µ.WithMemento()).(*µ.Protocol)
		it.Then(t).Should(it.True(cat.Memento))
	})

	t.Run("WithDefaultHost", func(t *testing.T) {
		cat := µ.New(µ.WithDefaultHost("https://example.com")).(*µ.Protocol)
		it.Then(t).Should(it.Equal(cat.Host, "https://example.com"))
	})

	t.Run("WithCookieJar", func(t *testing.T) {
		cat := µ.New(µ.WithCookieJar()).(*µ.Protocol)
		it.Then(t).ShouldNot(it.Nil(cat.Socket.(*http.Client).Jar))
	})

	t.Run("WithDefaultRedirectPolicy", func(t *testing.T) {
		cat := µ.New(µ.WithDefaultRedirectPolicy()).(*µ.Protocol)
		it.Then(t).Should(it.Equiv(cat.Socket.(*http.Client).CheckRedirect, nil))
	})
}
