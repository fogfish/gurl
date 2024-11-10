//
// Copyright (C) 2019 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/fogfish/opts"
	"golang.org/x/net/publicsuffix"
)

// HTTP Stack config option
type Option = opts.Option[Protocol]

var (
	// Set custom implementation of HTTP client.
	// It requires anything that implements Socket interface (aka http.Client)
	//
	//	type Socket interface {
	//	  Do(req *http.Request) (*http.Response, error)
	//	}
	WithClient = opts.ForType[Protocol, Socket]()

	// Set the default host for http stack.
	// The host is used when request URI does not contain any host.
	WithHost = opts.ForName[Protocol, string]("Host")

	// Enables HTTP Response buffering
	WithMemento = opts.ForName[Protocol, bool]("Memento")

	// Buffers HTTP Response Payload into context.
	WithMementoPayload = WithMemento(true)

	// Disables TLS certificate validation for HTTP(S) sessions.
	WithInsecureTLS = opts.From(withInsecureTLS)

	// Enables automated cookie handling across requests originated from the session.
	WithCookieJar = opts.From(withCookieJar)

	// Disables default [gurl] redirect policy to Golang's one.
	// It enables the HTTP stack automatically follows redirects
	WithRedirects = opts.From(withRedirects)

	// Enable log level
	WithLogLevel = opts.ForName[Protocol, int]("LogLevel")

	// Enables debug logging.
	// The logger outputs HTTP requests only.
	WithDebugRequest = WithLogLevel(1)

	// Enables debug logging.
	// The logger outputs HTTP requests and responses.
	WithDebugResponse = WithLogLevel(2)

	// Enable debug logging.
	WithDebugPayload = WithLogLevel(3)
)

func withInsecureTLS(cat *Protocol) error {
	if cli, ok := cat.Socket.(*http.Client); ok {
		switch t := cli.Transport.(type) {
		case *http.Transport:
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = &tls.Config{}
			}
			t.TLSClientConfig.InsecureSkipVerify = true
		default:
			return fmt.Errorf("unsupported transport type %T", t)
		}
	}
	return nil
}

func withCookieJar(cat *Protocol) error {
	if cli, ok := cat.Socket.(*http.Client); ok {
		jar, err := cookiejar.New(&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		})
		if err != nil {
			return err
		}
		cli.Jar = jar
	}
	return nil
}

func withRedirects(cat *Protocol) error {
	if cli, ok := cat.Socket.(*http.Client); ok {
		cli.CheckRedirect = nil
	}
	return nil
}
