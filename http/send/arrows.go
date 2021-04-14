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

URL defines a mandatory parameters to the request such as
HTTP method and destination URL, use Params arrow if you
need to supply URL query params.
*/
func (method Method) URL(uri string, args ...interface{}) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		var addr *url.URL
		if addr, cat.Fail = mkURL(uri, args...); cat.Fail != nil {
			return cat
		}

		if cat.HTTP == nil {
			cat.HTTP = &gurl.IOCatHTTP{}
		}

		switch addr.Scheme {
		case "http", "https":
			cat.HTTP.Send = &gurl.UpStreamHTTP{
				Method:  string(method),
				URL:     addr,
				Header:  make(map[string]*string),
				Payload: bytes.NewBuffer(nil),
			}
		default:
			cat.Fail = &gurl.NotSupported{URL: addr}
		}
		return cat
	}
}

func mkURL(uri string, args ...interface{}) (*url.URL, error) {
	switch uri[0] {
	case '!':
		return mkEscapedURL(false, uri[1:], args...)
	default:
		return mkEscapedURL(true, uri, args...)
	}
}

func mkEscapedURL(escape bool, uri string, args ...interface{}) (*url.URL, error) {
	opts := []interface{}{}
	for _, x := range args {
		switch v := x.(type) {
		case *url.URL:
			v.Path = strings.TrimSuffix(v.Path, "/")
			opts = append(opts, v.String())
		case func() string:
			opts = append(opts, maybeEscape(escape, v()))
		default:
			opts = append(opts, maybeEscape(escape, urlSegment(x)))
		}
	}

	return url.Parse(fmt.Sprintf(uri, opts...))
}

func urlSegment(arg interface{}) string {
	val := reflect.ValueOf(arg)

	if val.Kind() == reflect.Ptr {
		return fmt.Sprintf("%v", val.Elem())
	}

	return fmt.Sprintf("%v", val)
}

func maybeEscape(escape bool, val string) string {
	if escape {
		return url.PathEscape(val)
	}

	return val
}

/*

Header defines HTTP headers to the request, use combinator
to define multiple header values.

  http.Join(
		ø.Header("Accept").Is(...),
		ø.Header("Content-Type").Is(...),
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

func (header Header) name() string {
	return strings.ToLower(string(header))
}

// Is sets a literval value of HTTP header
func (header Header) Is(value string) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		cat.HTTP.Send.Header[header.name()] = &value
		return cat
	}
}

// Val sets a value of HTTP header from variable
func (header Header) Val(value *string) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		cat.HTTP.Send.Header[header.name()] = value
		return cat
	}
}

// Set is combinator to define HTTP header into request
func (header Header) Set(cat *gurl.IOCat, value string) *gurl.IOCat {
	cat.HTTP.Send.Header[header.name()] = &value
	return cat
}

// Content defines headers for content negotiation
type Content Header

// JSON defines header `???: application/json`
func (h Content) JSON(cat *gurl.IOCat) *gurl.IOCat {
	return Header(h).Set(cat, "application/json")
}

// Form defined Header `???: application/x-www-form-urlencoded`
func (h Content) Form(cat *gurl.IOCat) *gurl.IOCat {
	return Header(h).Set(cat, "application/x-www-form-urlencoded")
}

// Text defined Header `???: plain/text`
func (h Content) Text(cat *gurl.IOCat) *gurl.IOCat {
	return Header(h).Set(cat, "plain/text")
}

// Html defined Header `???: plain/html`
func (h Content) Html(cat *gurl.IOCat) *gurl.IOCat {
	return Header(h).Set(cat, "plain/html")
}

// Is sets a literval value of HTTP header
func (h Content) Is(value string) http.Arrow {
	return Header(h).Is(value)
}

// Val sets a value of HTTP header from variable
func (h Content) Val(value *string) http.Arrow {
	return Header(h).Val(value)
}

// Lifecycle defines headers for connection management
type Lifecycle Header

// KeepAlive defines header `???: keep-alive`
func (h Lifecycle) KeepAlive(cat *gurl.IOCat) *gurl.IOCat {
	return Header(h).Set(cat, "keep-alive")
}

// Close defines header `???: close`
func (h Lifecycle) Close(cat *gurl.IOCat) *gurl.IOCat {
	return Header(h).Set(cat, "close")
}

// Is sets a literval value of HTTP header
func (h Lifecycle) Is(value string) http.Arrow {
	return Header(h).Is(value)
}

// Val sets a value of HTTP header from variable
func (h Lifecycle) Val(value *string) http.Arrow {
	return Header(h).Val(value)
}

/*

Params appends query params to request URL. The arrow takes a struct and
converts it to map[string]string. The function fails if input is not convertable
to map of strings (e.g. nested struct).
*/
func Params(query interface{}) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		bytes, err := json.Marshal(query)
		if err != nil {
			cat.Fail = err
			return cat
		}

		var req map[string]string
		err = json.Unmarshal(bytes, &req)
		if err != nil {
			cat.Fail = err
			return cat
		}

		q := cat.HTTP.Send.URL.Query()
		for k, v := range req {
			q.Add(k, v)
		}
		cat.HTTP.Send.URL.RawQuery = q.Encode()
		return cat
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
	return func(cat *gurl.IOCat) *gurl.IOCat {
		content, ok := cat.HTTP.Send.Header["content-type"]
		if !ok {
			cat.Fail = fmt.Errorf("unknown Content-Type")
			return cat
		}

		switch stream := data.(type) {
		case string:
			cat.HTTP.Send.Payload = bytes.NewBuffer([]byte(stream))
		case []byte:
			cat.HTTP.Send.Payload = bytes.NewBuffer(stream)
		case io.Reader:
			cat.HTTP.Send.Payload = stream
		default:
			cat.HTTP.Send.Payload, cat.Fail = encode(*content, data)
		}
		return cat
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
