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
	"net"
	"net/http"
	"time"

	"github.com/fogfish/opts"
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

type Socket interface {
	Do(req *http.Request) (*http.Response, error)
}

// Protocol is an instance of Stack
type Protocol struct {
	Socket
	Host     string
	LogLevel int
	Memento  bool
}

// New instance of HTTP Stack
func New(opt ...Option) Stack {
	cat, err := NewStack(opt...)
	if err != nil {
		panic(err)
	}

	return cat
}

// New instance of HTTP Stack
func NewStack(opt ...Option) (Stack, error) {
	cat := &Protocol{Socket: Client()}

	if err := opts.Apply(cat, opt); err != nil {
		return nil, err
	}

	return cat, nil
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
