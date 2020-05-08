//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package gurl

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sort"
	"time"
)

// IOCat defines the category or type for HTTP I/O. A composition of
// HTTP primitives within the category are written with the following syntax:
//
//   gurl.HTTP(Arrow1, ..., ArrowN)
//
// Here, each Arrow is a morphism applied to HTTP protocol, they composition
// is defined using "programmable comma". Effectively the implementation resembles the
// state monad. It defines an abstraction of environments and lenses to focus
// inside it. In other words, the category represents the environment as an
// "invisible" side-effect of the composition.
type IOCat struct {
	URL  *url.URL
	HTTP *IOSpec
	Body interface{}
	Fail error
	pool *http.Client
	dur  time.Duration
}

// Arrow is a morphism applied to IO category
type Arrow func(*IOCat) *IOCat

// IOSpec defines parameters of IO transactor
type IOSpec struct {
	Method  string
	Header  map[string]*string
	Payload *bytes.Buffer
	Ingress *http.Response
}

// Join composes arrows to high-order function
// (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
func Join(arrows ...Arrow) Arrow {
	return func(io *IOCat) *IOCat {
		if io.Fail != nil {
			return io
		}
		for _, f := range arrows {
			if io = f(io); io.Fail != nil {
				return io
			}
		}
		return io
	}
}

// HTTP builds high-order protocol closure
func HTTP(arrows ...Arrow) Arrow {
	return func(io *IOCat) *IOCat {
		if io.Fail != nil {
			return io
		}
		for _, f := range arrows {
			if io = f(io); io.Fail != nil {
				return io
			}
		}
		if io.HTTP != nil && io.HTTP.Ingress != nil {
			ioutil.ReadAll(io.HTTP.Ingress.Body)
			io.Fail = io.HTTP.Ingress.Body.Close()
			io.HTTP.Ingress = nil
		}
		return io
	}
}

// IO creates the instance of HTTP I/O category with default HTTP client.
// Please note that default client disables TLS verification.
// Use this only for testing.
func IO(client ...*http.Client) *IOCat {
	if len(client) < 1 {
		return &IOCat{pool: defaultClient()}
	}

	return &IOCat{pool: client[0]}
}

func defaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			ReadBufferSize: 128 * 1024,
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// Unsafe performs networking side-effect
func (io *IOCat) Unsafe() *IOCat {
	if io.Fail != nil {
		return io
	}

	var eg *http.Request
	eg, io.Fail = http.NewRequest(io.HTTP.Method, io.URL.String(), io.HTTP.Payload)
	if io.Fail != nil {
		return io
	}

	for head, value := range io.HTTP.Header {
		eg.Header.Set(head, *value)
	}

	t := time.Now()
	io.HTTP.Ingress, io.Fail = io.pool.Do(eg)
	io.dur = time.Now().Sub(t)
	return io
}

// Ord extends sort.Interface with ability to lookup element by string
type Ord interface {
	sort.Interface
	// String return primary key as string type
	String(int) string
	// Value return value at index
	Value(int) interface{}
}
