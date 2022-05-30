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
	"io"
	"net/url"
	"reflect"
	"strings"

	"github.com/fogfish/gurl"
	"github.com/fogfish/gurl/http"
)

/*

Method is base type for HTTP methods
*/
type Method string

/*

List of supported built-in method constants
*/
const (
	GET    = Method("GET")
	POST   = Method("POST")
	PUT    = Method("PUT")
	DELETE = Method("DELETE")
	PATCH  = Method("PATCH")
)

/*

Authority is part of URL, use the type to prevent escaping
*/
type Authority string

/*

Segment is part of URL, use the type to prevent path escaping
*/
type Segment string

/*

URL defines a mandatory parameters to the request such as
HTTP method and destination URL, use Params arrow if you
need to supply URL query params.
*/
func (method Method) URL(uri string, args ...interface{}) http.Arrow {
	return func(cat *http.Context) error {
		addr := mkURL(uri, args...)

		switch {
		case strings.HasPrefix(addr, "http"):
			cat.Request = &http.Request{
				Method:  string(method),
				URL:     addr,
				Header:  make(map[string]*string),
				Payload: bytes.NewBuffer(nil),
			}
		default:
			return &gurl.NotSupported{URL: addr}
		}

		return nil
	}
}

func mkURL(uri string, args ...interface{}) string {
	opts := []interface{}{}
	for _, x := range args {
		switch v := x.(type) {
		case *url.URL:
			v.Path = strings.TrimSuffix(v.Path, "/")
			opts = append(opts, v.String())
		case *Segment:
			opts = append(opts, *v)
		case Segment:
			opts = append(opts, v)
		case *Authority:
			opts = append(opts, *v)
		case Authority:
			opts = append(opts, v)
		default:
			opts = append(opts, url.PathEscape(urlSegment(x)))
		}
	}

	return fmt.Sprintf(uri, opts...)
}

func urlSegment(arg interface{}) string {
	val := reflect.ValueOf(arg)

	if val.Kind() == reflect.Ptr {
		return fmt.Sprintf("%v", val.Elem())
	}

	return fmt.Sprintf("%v", val)
}

/*

Header defines HTTP headers to the request, use combinator
to define multiple header values.

  http.Do(
		ø.Header("User-Agent").Is("gurl"),
		ø.Header("Content-Type").Is(content),
	)
*/
type Header string

/*

List of supported HTTP header constants
https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Request_fields
*/
const (
	Accept            = Content("Accept")
	AcceptCharset     = Header("Accept-Charset")
	AcceptEncoding    = Header("Accept-Encoding")
	AcceptLanguage    = Header("Accept-Language")
	Authorization     = Header("Authorization")
	CacheControl      = Header("Cache-Control")
	Connection        = Lifecycle("Connection")
	ContentEncoding   = Header("Content-Encoding")
	ContentLength     = Header("Content-Length")
	ContentType       = Content("Content-Type")
	Cookie            = Header("Cookie")
	Date              = Header("Date")
	Host              = Header("Host")
	IfMatch           = Header("If-Match")
	IfModifiedSince   = Header("If-Modified-Since")
	IfNoneMatch       = Header("If-None-Match")
	IfRange           = Header("If-Range")
	IfUnmodifiedSince = Header("If-Unmodified-Since")
	Origin            = Header("Origin")
	Range             = Header("Range")
	TransferEncoding  = Header("Transfer-Encoding")
	UserAgent         = Header("User-Agent")
	Upgrade           = Header("Upgrade")
)

// Is sets value of HTTP header
func (header Header) Is(value string) http.Arrow {
	return func(cat *http.Context) error {
		h := strings.ToLower(string(header))
		cat.Request.Header[h] = &value
		return nil
	}
}

// Set is combinator to define HTTP header into request
func (header Header) Set(cat *http.Context, value string) error {
	h := strings.ToLower(string(header))
	cat.Request.Header[h] = &value
	return nil
}

// Content defines headers for content negotiation
type Content Header

// ApplicationJSON defines header `???: application/json`
func (h Content) ApplicationJSON(cat *http.Context) error {
	return Header(h).Set(cat, "application/json")
}

// JSON defines header `???: application/json`
func (h Content) JSON(cat *http.Context) error {
	return Header(h).Set(cat, "application/json")
}

// Form defined Header `???: application/x-www-form-urlencoded`
func (h Content) Form(cat *http.Context) error {
	return Header(h).Set(cat, "application/x-www-form-urlencoded")
}

// TextPlain defined Header `???: text/plain`
func (h Content) TextPlain(cat *http.Context) error {
	return Header(h).Set(cat, "text/plain")
}

// Text defined Header `???: text/plain`
func (h Content) Text(cat *http.Context) error {
	return Header(h).Set(cat, "text/plain")
}

// TextHTML defined Header `???: text/html`
func (h Content) TextHTML(cat *http.Context) error {
	return Header(h).Set(cat, "text/html")
}

// HTML defined Header `???: text/html`
func (h Content) HTML(cat *http.Context) error {
	return Header(h).Set(cat, "text/html")
}

// Is sets a literval value of HTTP header
func (h Content) Is(value string) http.Arrow {
	return Header(h).Is(value)
}

// Lifecycle defines headers for connection management
type Lifecycle Header

// KeepAlive defines header `???: keep-alive`
func (h Lifecycle) KeepAlive(cat *http.Context) error {
	return Header(h).Set(cat, "keep-alive")
}

// Close defines header `???: close`
func (h Lifecycle) Close(cat *http.Context) error {
	return Header(h).Set(cat, "close")
}

// Is sets a literval value of HTTP header
func (h Lifecycle) Is(value string) http.Arrow {
	return Header(h).Is(value)
}

/*

Params appends query params to request URL. The arrow takes a struct and
converts it to map[string]string. The function fails if input is not convertable
to map of strings (e.g. nested struct).
*/
func Params[T any](query T) http.Arrow {
	return func(cat *http.Context) error {
		bytes, err := json.Marshal(query)
		if err != nil {
			return err
		}

		var req map[string]string
		err = json.Unmarshal(bytes, &req)
		if err != nil {
			return err
		}

		uri, err := url.Parse(cat.Request.URL)
		if err != nil {
			return err
		}

		q := uri.Query()
		for k, v := range req {
			q.Add(k, v)
		}
		uri.RawQuery = q.Encode()
		cat.Request.URL = uri.String()

		return nil
	}
}

/*

Send payload to destination URL. You can also use native Go data types
(e.g. maps, struct, etc) as egress payload. The library implicitly encodes
input structures to binary using Content-Type as a hint. The function fails
if content type is not supported by the library.

The function accept a "classical" data container such as string, []bytes or
io.Reader interfaces.
*/
func Send(data interface{}) http.Arrow {
	return func(cat *http.Context) error {
		content, ok := cat.Request.Header["content-type"]
		if !ok {
			return fmt.Errorf("unknown Content-Type")
		}

		switch stream := data.(type) {
		case string:
			cat.Request.Payload = bytes.NewBuffer([]byte(stream))
		case []byte:
			cat.Request.Payload = bytes.NewBuffer(stream)
		case io.Reader:
			cat.Request.Payload = stream
		default:
			pkt, err := encode(*content, data)
			if err != nil {
				return err
			}
			cat.Request.Payload = pkt
		}
		return nil
	}
}

func encode(content string, data interface{}) (buf *bytes.Buffer, err error) {
	switch {
	// "application/json" and other variants
	case strings.Contains(content, "json"):
		buf, err = encodeJSON(data)
	// "application/x-www-form-urlencoded"
	case strings.Contains(content, "www-form"):
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
		return nil, fmt.Errorf("encode application/x-www-form-urlencoded: %w", err)
	}

	var payload url.Values = make(map[string][]string)
	for key, val := range req {
		payload[key] = []string{val}
	}
	return bytes.NewBuffer([]byte(payload.Encode())), nil
}
