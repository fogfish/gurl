//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/ajg/form"
	"github.com/fogfish/gurl"
	"github.com/fogfish/logger"
	"golang.org/x/net/publicsuffix"
)

/*

Arrow is a morphism applied to HTTP protocol stack
*/
type Arrow func(*Context) error

/*

Stack is HTTP protocol stack
*/
type Stack interface {
	WithContext(context.Context) *Context
	IO(context.Context, ...Arrow) error
}

/*

Request specify parameters for HTTP requests
*/
type Request struct {
	Method  string
	URL     string
	Header  map[string]*string
	Payload io.Reader
}

/*

Context defines the category of HTTP I/O
*/
type Context struct {
	*Protocol

	// Context of Request / Response
	context.Context
	*Request
	*http.Response
}

// Unsafe evaluates current context of HTTP I/O
func (ctx *Context) Unsafe() error {
	eg, err := http.NewRequest(
		ctx.Request.Method,
		ctx.Request.URL,
		ctx.Request.Payload,
	)
	if err != nil {
		return err
	}

	for head, value := range ctx.Request.Header {
		eg.Header.Set(head, *value)
	}

	if ctx.Context != nil {
		eg = eg.WithContext(ctx.Context)
	}

	logSend(ctx.LogLevel, eg)

	in, err := ctx.Client.Do(eg)
	if err != nil {
		return err
	}

	ctx.Response = in

	logRecv(ctx.LogLevel, in)

	return nil
}

// IO executes protocol operations
func (ctx *Context) IO(arrows ...Arrow) error {
	for _, f := range arrows {
		if err := f(ctx); err != nil {
			return err
		}
	}

	if ctx.Response != nil {
		// Note: due to Golang HTTP pool implementation we need to consume and
		//       discard body. Otherwise, HTTP connection is not returned to
		//       to the pool.
		body := ctx.Response.Body
		ctx.Response = nil

		_, err := io.Copy(ioutil.Discard, body)
		if err != nil {
			return err
		}

		err = body.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

/*

Protocol is an instance of Stack
*/
type Protocol struct {
	*http.Client
	LogLevel int
}

/*

New instantiates category of HTTP I/O
*/
func New(opts ...Config) Stack {
	cat := &Protocol{Client: Client()}

	for _, opt := range opts {
		opt(cat)
	}

	return cat
}

// WithContext create instance of I/O Context
func (cat *Protocol) WithContext(ctx context.Context) *Context {
	return &Context{
		Protocol: cat,
		Context:  ctx,
		Request:  nil,
		Response: nil,
	}
}

// IO executes protocol operations
func (cat *Protocol) IO(ctx context.Context, arrows ...Arrow) error {
	return cat.WithContext(ctx).IO(arrows...)
}

/*

Join composes HTTP arrows to high-order function
(a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
*/
func Join(arrows ...Arrow) Arrow {
	return func(cat *Context) error {
		for _, f := range arrows {
			if err := f(cat); err != nil {
				return err
			}
		}

		return nil
	}
}

// Config for HTTP client
type Config func(*Protocol)

/*

Client Default HTTP client
*/
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

// LogRequest enables debug logging for requests
func LogRequest() Config {
	return func(cat *Protocol) {
		cat.LogLevel = 1
	}
}

// LogResponse enables debug logging for requests
func LogResponse() Config {
	return func(cat *Protocol) {
		cat.LogLevel = 2
	}
}

// LogResponse enables debug logging for requests
func LogDebug() Config {
	return func(cat *Protocol) {
		cat.LogLevel = 3
	}
}

// InsecureTLS disables certificates validation
func InsecureTLS() Config {
	return func(cat *Protocol) {
		switch t := cat.Client.Transport.(type) {
		case *http.Transport:
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = &tls.Config{}
			}
			t.TLSClientConfig.InsecureSkipVerify = true
		default:
			panic(fmt.Errorf("Unsupported transport type %T", t))
		}
	}
}

// CookieJar enables cookie handlings
func CookieJar() Config {
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

/*

IO executes protocol operations
*/
func IO[T any](ctx *Context, arrows ...Arrow) (*T, error) {
	for _, f := range arrows {
		if err := f(ctx); err != nil {
			return nil, err
		}
	}

	if ctx.Response == nil {
		return nil, fmt.Errorf("empty response")
	}
	defer ctx.Response.Body.Close()

	var val T
	err := decode(
		ctx.Response.Header.Get("Content-Type"),
		ctx.Response.Body,
		&val,
	)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func decode[T any](content string, stream io.ReadCloser, data *T) error {
	switch {
	case strings.Contains(content, "json"):
		return json.NewDecoder(stream).Decode(data)
	case strings.Contains(content, "www-form"):
		return form.NewDecoder(stream).Decode(data)
	default:
		return &gurl.NoMatch{
			Diff:    fmt.Sprintf("- Content-Type: application/*\n+ Content-Type: %s", content),
			Payload: map[string]string{"Content-Type": content},
		}
	}
}

//
//
func logSend(level int, eg *http.Request) {
	if level >= 1 {
		if msg, err := httputil.DumpRequest(eg, level == 3); err == nil {
			logger.Debug(">>>>\n%s\n", msg)
		}
	}
}

func logRecv(level int, in *http.Response) {
	if level >= 2 {
		if msg, err := httputil.DumpResponse(in, level == 3); err == nil {
			logger.Debug("<<<<\n%s\n", msg)
		}
	}
}
