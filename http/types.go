//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"time"

	"github.com/fogfish/gurl"
	"golang.org/x/net/publicsuffix"
)

/*

Arrow is a morphism applied to HTTP
*/
type Arrow func(*gurl.IOCat) *gurl.IOCat

/*

Join composes HTTP arrows to high-order function
(a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
*/
func Join(arrows ...Arrow) gurl.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		if cat.Fail != nil {
			return cat
		}

		for _, f := range arrows {
			if cat = f(cat); cat.Fail != nil {
				return cat
			}
		}

		if cat.HTTP != nil && cat.HTTP.Recv != nil && cat.HTTP.Recv.Response != nil {
			// Note: due to Golang HTTP pool implementation we need to consume and
			//       discard body. Otherwise, HTTP connection is not returned to
			//       to the pool.
			io.Copy(ioutil.Discard, cat.HTTP.Recv.Response.Body)
			cat.Fail = cat.HTTP.Recv.Body.Close()
			cat.HTTP.Recv.Response = nil
		}

		return cat
	}
}

/*

Stack configures custom HTTP stack for the category.
*/
func Stack(client *http.Client) gurl.Config {
	pool := pool{client}
	return gurl.SideEffect(pool.Unsafe)
}

/*

Default configures default HTTP stack for the category.
*/
func Default() gurl.Config {
	pool := pool{Client()}
	return gurl.SideEffect(pool.Unsafe)
}

// DefaultIO creates default HTTP IO category
// use for development only
func DefaultIO(opts ...gurl.Config) *gurl.IOCat {
	args := append(opts, Default())
	return gurl.IO(args...)
}

//
//
type pool struct{ *http.Client }

func (p pool) Unsafe(cat *gurl.IOCat) *gurl.IOCat {
	if cat.Fail != nil {
		return cat
	}

	var eg *http.Request
	eg, cat.Fail = http.NewRequest(
		cat.HTTP.Send.Method,
		cat.HTTP.Send.URL,
		cat.HTTP.Send.Payload,
	)
	if cat.Fail != nil {
		return cat
	}

	for head, value := range cat.HTTP.Send.Header {
		eg.Header.Set(head, *value)
	}

	if cat.Context != nil {
		eg = eg.WithContext(cat.Context)
	}

	var in *http.Response
	in, cat.Fail = p.Client.Do(eg)
	if cat.Fail != nil {
		return cat
	}

	cat.HTTP.Recv = &gurl.DnStreamHTTP{Response: in}

	logSend(cat.LogLevel, eg)
	logRecv(cat.LogLevel, in)

	return cat
}

func logSend(level int, eg *http.Request) {
	if level >= gurl.LogLevelEgress {
		if msg, err := httputil.DumpRequest(eg, level == gurl.LogLevelDebug); err == nil {
			log.Printf(">>>>\n%s\n", msg)
		}
	}
}

func logRecv(level int, in *http.Response) {
	if level >= gurl.LogLevelIngress {
		if msg, err := httputil.DumpResponse(in, level == gurl.LogLevelDebug); err == nil {
			log.Printf("<<<<\n%s\n", msg)
		}
	}
}

// Config for HTTP client
type Config func(*http.Client)

/*

Client creates HTTP client
*/
func Client(opts ...Config) *http.Client {
	cli := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			ReadBufferSize: 128 * 1024,
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			// TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
}

// InsecureTLS disables certificates validation
func InsecureTLS() Config {
	return func(c *http.Client) {
		switch t := c.Transport.(type) {
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
	return func(c *http.Client) {
		jar, err := cookiejar.New(&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		})
		if err != nil {
			panic(err)
		}
		c.Jar = jar
	}
}
