//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

// Package send defines a pure computations to compose HTTP request senders
package send

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/fogfish/gurl"
)

/*

URL defines a mandatory parameters to the request such as
HTTP method and destination URL, use Params arrow if you
need to supply URL query params.
*/
func URL(method, uri string) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.HTTP = nil
		io.URL, io.Fail = url.Parse(uri)

		switch io.URL.Scheme {
		case "http", "https":
			io.HTTP = &gurl.IOSpec{
				Method:  method,
				Header:  make(map[string]string),
				Payload: bytes.NewBuffer(nil),
				Ingress: nil,
			}
		default:
			io.Fail = &gurl.BadSchema{io.URL.Scheme}
		}
		return io
	}
}

// GET is syntax sugar of URL("GET", ...)
func GET(uri string) gurl.Arrow {
	return URL("GET", uri)
}

// POST is syntax sugar of URL("POST", ...)
func POST(uri string) gurl.Arrow {
	return URL("POST", uri)
}

// PUT is syntax sugar of URL("PUT", ...)
func PUT(uri string) gurl.Arrow {
	return URL("PUT", uri)
}

// DELETE is syntax sugar of URL("DELETE", ...)
func DELETE(uri string) gurl.Arrow {
	return URL("DELETE", uri)
}

/*

Header defines HTTP headers to the request, use combinator
to define multiple header values.

  gurl.HTTP(
		ø.Header("Accept", ...),
		ø.Header("Content-Type", ...),
	)
*/
func Header(header, value string) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.HTTP.Header[header] = value
		return io
	}
}

// Accept is syntax sugar of Header("Accept", ...)
func Accept(mime string) gurl.Arrow {
	return Header("Accept", mime)
}

// AcceptJSON is syntax sugar of Header("Accept", "application/json")
func AcceptJSON() gurl.Arrow {
	return Header("Accept", "application/json")
}

// Content is syntax sugar of Header("Content-Type", ...)
func Content(mime string) gurl.Arrow {
	return Header("Content-Type", mime)
}

// ContentJSON is syntax sugar of Header("Content-Type", "application/json")
func ContentJSON() gurl.Arrow {
	return Header("Content-Type", "application/json")
}

// KeepAlive is a syntax sugar of Header("Connection", "keep-alive")
func KeepAlive() gurl.Arrow {
	return Header("Connection", "keep-alive")
}

// Authorization is syntax sugar of Header("Authorization", ...)
func Authorization(token string) gurl.Arrow {
	return Header("Authorization", token)
}

/*

Params appends query params to request URL. The arrow takes a struct and
converts it to map[string]string. The function fails if input is not convertable
to map of strings (e.g. nested struct).
*/
func Params(query interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		bytes, err := json.Marshal(query)
		if err != nil {
			io.Fail = err
			return io
		}

		var req map[string]string
		err = json.Unmarshal(bytes, &req)
		if err != nil {
			io.Fail = err
			return io
		}

		q := io.URL.Query()
		for k, v := range req {
			q.Add(k, v)
		}
		io.URL.RawQuery = q.Encode()
		return io
	}
}

/*

Send payload to destination URL. You can also use native Go data types
(e.g. maps, struct, etc) as egress payload. The library implicitly encodes
input structures to binary using Content-Type as a hint. The function fails
if content type is not supported by the library.
*/
func Send(data interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.HTTP.Payload, io.Fail = encode(io.HTTP.Header["Content-Type"], data)
		return io
	}
}

func encode(content string, data interface{}) (buf *bytes.Buffer, err error) {
	switch content {
	case "application/json":
		buf, err = encodeJSON(data)
	case "application/x-www-form-urlencoded":
		buf, err = encodeForm(data)
	default:
		err = fmt.Errorf("unsupported Content-Type %v", content)
	}

	return
}

func encodeJSON(data interface{}) (*bytes.Buffer, error) {
	json, err := json.Marshal(data)
	return bytes.NewBuffer(json), err
}

func encodeForm(data interface{}) (*bytes.Buffer, error) {
	bin, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var req map[string]string
	err = json.Unmarshal(bin, &req)
	if err != nil {
		return nil, err
	}

	var payload url.Values = make(map[string][]string)
	for key, val := range req {
		payload[key] = []string{val}
	}

	return bytes.NewBuffer([]byte(payload.Encode())), nil
}
