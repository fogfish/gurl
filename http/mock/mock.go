//
// Copyright (C) 2019 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package mock

import (
	"bytes"
	µ "github.com/fogfish/gurl/v2/http"
	"io"
	"net/http"
)

// Mocks HTTP client
type Mock struct {
	r   *http.Response
	err error
}

func (mock *Mock) Do(req *http.Request) (*http.Response, error) {
	if mock.err != nil {
		return nil, mock.err
	}

	return mock.r, nil
}

// Mocks Option
type Option func(m *Mock)

func Preset(opts ...Option) Option {
	return func(m *Mock) {
		for _, opt := range opts {
			opt(m)
		}
	}
}

// Mock failure of HTTP client
func Fail(err error) Option {
	return func(m *Mock) {
		m.err = err
	}
}

// Mock response with status code (default 200)
func Status(code int) Option {
	return func(m *Mock) {
		m.r.StatusCode = code
	}
}

// Mock response with HTTP header (default none)
func Header(h, v string) Option {
	return func(m *Mock) {
		m.r.Header.Set(h, v)
	}
}

// Mock response body (default empty)
func Body(body []byte) Option {
	return func(m *Mock) {
		m.r.Body = io.NopCloser(bytes.NewBuffer(body))
	}
}

// Mock response body I/O error (default nil)
func IOError(err error) Option {
	return func(m *Mock) {
		m.r.Body = io.NopCloser(errReader{err})
	}
}

// Mock HTTP Client
func New(opts ...Option) µ.Option {
	m := &Mock{
		r: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{},
			Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
		},
	}

	for _, opt := range opts {
		opt(m)
	}

	return µ.WithClient(m)
}

type errReader struct{ err error }

func (r errReader) Read(p []byte) (n int, err error) { return 0, r.err }
