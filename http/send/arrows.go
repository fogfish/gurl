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

// Method is base type for HTTP methods
type Method string

// List of supported built-in method constants
const (
	GET    = Method("GET")
	POST   = Method("POST")
	PUT    = Method("PUT")
	DELETE = Method("DELETE")
	PATCH  = Method("PATCH")
)

// Authority is part of URL, use the type to prevent escaping
type Authority string

// Segment is part of URL, use the type to prevent path escaping
type Segment string

// URL defines a mandatory parameters to the request such as
// HTTP method and destination URI
func (method Method) URI(addr string) http.Arrow {
	return func(cat *http.Context) error {
		switch {
		case strings.HasPrefix(addr, "http"):
			req, err := http.NewRequest(string(method), addr)
			if err != nil {
				return err
			}

			cat.Request = req
		default:
			return &gurl.NotSupported{URL: addr}
		}

		return nil
	}
}

// URL defines a mandatory parameters to the request such as
// HTTP method and destination URL, use Params arrow if you
// need to supply URL query params.
func (method Method) URL(uri string, args ...interface{}) http.Arrow {
	return func(cat *http.Context) error {
		addr := mkURL(uri, args...)

		switch {
		case strings.HasPrefix(addr, "http"):
			req, err := http.NewRequest(string(method), addr)
			if err != nil {
				return err
			}

			cat.Request = req
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
	Accept            = HeaderContent("Accept")
	AcceptCharset     = Header("Accept-Charset")
	AcceptEncoding    = Header("Accept-Encoding")
	AcceptLanguage    = Header("Accept-Language")
	Authorization     = Header("Authorization")
	CacheControl      = Header("Cache-Control")
	Connection        = HeaderConnection("Connection")
	ContentEncoding   = Header("Content-Encoding")
	ContentLength     = HeaderContentLength("Content-Length")
	ContentType       = HeaderContent("Content-Type")
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
	TransferEncoding  = HeaderTransferEncoding("Transfer-Encoding")
	UserAgent         = Header("User-Agent")
	Upgrade           = Header("Upgrade")
)

// Is sets value of HTTP header
func (header Header) Is(value string) http.Arrow {
	return func(cat *http.Context) error {
		cat.Request.Header.Add(string(header), value)
		return nil
	}
}

// Content defines headers for content negotiation
type HeaderContent Header

// ApplicationJSON defines header `???: application/json`
func (h HeaderContent) ApplicationJSON(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "application/json")
	return nil
}

// JSON defines header `???: application/json`
func (h HeaderContent) JSON(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "application/json")
	return nil
}

// Form defined Header `???: application/x-www-form-urlencoded`
func (h HeaderContent) Form(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "application/x-www-form-urlencoded")
	return nil
}

// TextPlain defined Header `???: text/plain`
func (h HeaderContent) TextPlain(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "text/plain")
	return nil
}

// Text defined Header `???: text/plain`
func (h HeaderContent) Text(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "text/plain")
	return nil
}

// TextHTML defined Header `???: text/html`
func (h HeaderContent) TextHTML(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "text/html")
	return nil
}

// HTML defined Header `???: text/html`
func (h HeaderContent) HTML(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "text/html")
	return nil
}

// Is sets a literval value of HTTP header
func (h HeaderContent) Is(value string) http.Arrow {
	return Header(h).Is(value)
}

// Lifecycle defines headers for connection management
type HeaderConnection Header

// KeepAlive defines header `???: keep-alive`
func (h HeaderConnection) KeepAlive(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "keep-alive")
	cat.Request.Close = false
	return nil
}

// Close defines header `???: close`
func (h HeaderConnection) Close(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "close")
	cat.Request.Close = true
	return nil
}

// Header TransferEncoding
type HeaderTransferEncoding Header

// Chunked defines header `Transfer-Encoding: chunked`
func (h HeaderTransferEncoding) Chunked(cat *http.Context) error {
	cat.Request.TransferEncoding = []string{"chunked"}
	return nil
}

// Identity defines header `Transfer-Encoding: identity`
func (h HeaderTransferEncoding) Identity(cat *http.Context) error {
	cat.Request.TransferEncoding = []string{"identity"}
	return nil
}

// Is sets a literval value of HTTP header
func (h HeaderTransferEncoding) Is(value string) http.Arrow {
	return func(cat *http.Context) error {
		cat.Request.TransferEncoding = strings.Split(value, ",")
		return nil
	}
}

// Header Content-Length
type HeaderContentLength Header

// Is sets a literval value of HTTP header
func (h HeaderContentLength) Is(value int64) http.Arrow {
	return func(cat *http.Context) error {
		cat.Request.ContentLength = value
		return nil
	}
}

// Params appends query params to request URL. The arrow takes a struct and
// converts it to map[string]string. The function fails if input is not convertable
// to map of strings (e.g. nested struct).
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
		uri := cat.Request.URL

		q := uri.Query()
		for k, v := range req {
			q.Add(k, v)
		}
		uri.RawQuery = q.Encode()
		cat.Request.URL = uri

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
		chunked := cat.Request.Header.Get(string(TransferEncoding)) == "chunked"
		content := cat.Request.Header.Get(string(ContentType))
		if content == "" {
			return fmt.Errorf("unknown Content-Type")
		}

		switch stream := data.(type) {
		case string:
			cat.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(stream)))
			cat.Request.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewBuffer([]byte(stream))), nil
			}
			if !chunked && cat.Request.ContentLength != 0 {
				cat.Request.ContentLength = int64(len(stream))
			}
		case *strings.Reader:
			cat.Request.Body = io.NopCloser(stream)
			snapshot := *stream
			cat.Request.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
			if !chunked && cat.Request.ContentLength != 0 {
				cat.Request.ContentLength = int64(stream.Len())
			}
		case []byte:
			cat.Request.Body = io.NopCloser(bytes.NewBuffer(stream))
			cat.Request.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewBuffer(stream)), nil
			}
			if !chunked && cat.Request.ContentLength != 0 {
				cat.Request.ContentLength = int64(len(stream))
			}
		case *bytes.Buffer:
			cat.Request.Body = io.NopCloser(stream)
			snapshot := stream.Bytes()
			cat.Request.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewBuffer(snapshot)), nil
			}
			if !chunked && cat.Request.ContentLength != 0 {
				cat.Request.ContentLength = int64(stream.Len())
			}
		case *bytes.Reader:
			cat.Request.Body = io.NopCloser(stream)
			snapshot := *stream
			cat.Request.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
			if !chunked && cat.Request.ContentLength != 0 {
				cat.Request.ContentLength = int64(stream.Len())
			}
		case io.Reader:
			rc, ok := stream.(io.ReadCloser)
			if !ok {
				rc = io.NopCloser(stream)
			}
			cat.Request.Body = rc
		default:
			pkt, err := encode(content, data)
			if err != nil {
				return err
			}
			cat.Request.Body = io.NopCloser(pkt)
			snapshot := pkt.Bytes()
			cat.Request.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewBuffer(snapshot)), nil
			}
			if !chunked && cat.Request.ContentLength != 0 {
				cat.Request.ContentLength = int64(pkt.Len())
			}
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
