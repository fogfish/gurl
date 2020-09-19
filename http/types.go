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
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/fogfish/gurl"
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
	pool := pool{defaultClient()}
	return gurl.SideEffect(pool.Unsafe)
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
		cat.HTTP.Send.URL.String(),
		cat.HTTP.Send.Payload,
	)
	if cat.Fail != nil {
		return cat
	}

	for head, value := range cat.HTTP.Send.Header {
		eg.Header.Set(head, *value)
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
