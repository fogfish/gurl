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
	sysio "io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"sort"
	"time"
)

// IOSpec defines parameters of IO transactor.
type IOSpec struct {
	Method  string
	Header  map[string]*string
	Payload *bytes.Buffer
	Ingress *http.Response
}

/*

HTTP builds high-order protocol closure in the context of IOCat.
A composition of HTTP primitives within the category are written
with the following syntax:

  gurl.HTTP(Arrow1, ..., ArrowN)

Here, each Arrow is a morphism applied to HTTP protocol, they composition
is defined using "programmable comma". Effectively the implementation resembles the
state monad. It defines an abstraction of environments and lenses to focus
inside it. In other words, the category represents the environment as an
"invisible" side-effect of the composition.
*/
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
			// Note: due to Golang HTTP pool implementation we need to consume and
			//       discard body. Otherwise, HTTP connection is not returned to
			//       to the pool.
			sysio.Copy(ioutil.Discard, io.HTTP.Ingress.Body)
			io.Fail = io.HTTP.Ingress.Body.Close()
			io.HTTP.Ingress = nil
		}
		return io
	}
}

/*

Join composes arrows to high-order function
(a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
*/
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

/*

Verbose enables debug logging of IO traffic
 - Level 0: disable debug logging (default)
 - Level 1: log only egress traffic
 - Level 2: log only ingress traffic
 - Level 3: log full content of packets
*/
func Verbose(level int) Config {
	return func(io *IOCat) *IOCat {
		io.verbose = level
		return io
	}
}

/*

Protocol defines a custom protocol infrastructure for the category.
*/
func Protocol(client *http.Client) Config {
	return func(io *IOCat) *IOCat {
		io.pool = client
		return io
	}
}

/*

IO creates the instance of HTTP I/O category use Config type to parametrize
the behavior. Note, the default client disables TLS verification.

Use returned instance of IO Category to evaluate IO "promise" returned by `gurl.HTTP`,
`gurl.Join`, etc.

  io := gurl.IO()
  fn := gurl.HTTP( ... )
  fn(io)

*/
func IO(opts ...Config) *IOCat {
	io := &IOCat{}
	for _, opt := range opts {
		io = opt(io)
	}

	if io.pool == nil {
		io.pool = defaultClient()
	}
	return io
}

func defaultClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			ReadBufferSize: 128 * 1024,
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

/*

Unsafe performs networking side-effect
*/
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

	logbody := io.verbose > 2
	if io.verbose > 0 {
		if msg, err := httputil.DumpRequest(eg, logbody); err == nil {
			log.Printf(">>>>\n%s\n", msg)
		}
	}
	if io.verbose > 1 {
		if msg, err := httputil.DumpResponse(io.HTTP.Ingress, logbody); err == nil {
			log.Printf("<<<<\n%s\n", msg)
		}
	}

	return io
}

/*

Ord extends sort.Interface with ability to lookup element by string.
This interface is a helper abstraction to evaluate presence of element in the sequence.

  gurl.HTTP(
    ...
    ƒ.Recv(&seq),
    ƒ.Seq(&seq).Has("example"),
    ...
  )

The example above shows a typical usage of Ord interface. The remote peer returns sequence
of elements. The lens Seq and Has focuses on the required element. A reference
implementation of the interface is

  type Seq []MyType

  func (c Seq) Len() int                { return len(c) }
  func (c Seq) Swap(i, j int)           { c[i], c[j] = c[j], c[i] }
  func (c Seq) Less(i, j int) bool      { return c[i].MyKey < c[j].MyKey }
  func (c Seq) String(i int) string     { return c[i].MyKey }
  func (c Seq) Value(i int) interface{} { return c[i] }

*/
type Ord interface {
	sort.Interface
	// String return primary key as string type
	String(int) string
	// Value return value at index
	Value(int) interface{}
}
