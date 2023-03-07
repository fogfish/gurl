//
// Copyright (C) 2019 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"
)

//
// The file implements the protocol stack, type owning HTTP client
//

// Creates instance of HTTP Request
func NewRequest(method, url string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}

// Stack is HTTP protocol stack
type Stack interface {
	WithContext(context.Context) *Context
	IO(context.Context, ...Arrow) error
}

// Protocol is an instance of Stack
type Protocol struct {
	*http.Client
	Host     string
	LogLevel int
	Memento  bool
}

// Allocate instance of HTTP Stack
func New(opts ...Config) Stack {
	cat := &Protocol{Client: Client()}

	for _, opt := range opts {
		opt(cat)
	}

	return cat
}

// WithContext create instance of I/O Context
func (stack *Protocol) WithContext(ctx context.Context) *Context {
	return &Context{
		Context:  ctx,
		Host:     stack.Host,
		Method:   http.MethodGet,
		Request:  nil,
		Response: nil,
		stack:    stack,
	}
}

func (stack *Protocol) IO(ctx context.Context, arrows ...Arrow) error {
	c := stack.WithContext(ctx)

	for _, f := range arrows {
		if err := f(c); err != nil {
			c.discardBody()
			return err
		}
		if err := c.discardBody(); err != nil {
			return err
		}
	}

	return nil
}

// Config option for HTTP client
type Config func(*Protocol)

// Creates default HTTP client
func Client() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			ReadBufferSize: 128 * 1024,
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
			// TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// WithClient replaces default client with custom instance
func WithClient(client *http.Client) Config {
	return func(cat *Protocol) {
		cat.Client = client
	}
}

// Enables debug logging for requests
func WithDebugRequest() Config {
	return func(cat *Protocol) {
		cat.LogLevel = 1
	}
}

// Enables debug logging for requests & response
func WithDebugResponse() Config {
	return func(cat *Protocol) {
		cat.LogLevel = 2
	}
}

// Enables debug logging for requests, responses and payload
func WithDebugPayload() Config {
	return func(cat *Protocol) {
		cat.LogLevel = 3
	}
}

// WithMemorize the response payload
func WithMemento() Config {
	return func(cat *Protocol) {
		cat.Memento = true
	}
}

func WithDefaultHost(host string) Config {
	return func(cat *Protocol) {
		cat.Host = host
	}
}

// WithInsecureTLS disables certificates validation
func WithInsecureTLS() Config {
	return func(cat *Protocol) {
		switch t := cat.Client.Transport.(type) {
		case *http.Transport:
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = &tls.Config{}
			}
			t.TLSClientConfig.InsecureSkipVerify = true
		default:
			panic(fmt.Errorf("unsupported transport type %T", t))
		}
	}
}

// WithCookieJar enables cookie handlings
func WithCookieJar() Config {
	return func(cat *Protocol) {
		jar, err := cookiejar.New(&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		})
		if err != nil {
			panic(err)
		}
		cat.Client.Jar = jar
	}
}
