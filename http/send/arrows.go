//
// Copyright (C) 2019 - 2023 Dmitry Kolesnikov
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
	"strconv"
	"strings"
	"time"

	"github.com/fogfish/gurl/v2"
	"github.com/fogfish/gurl/v2/http"
)

// Method defines HTTP Method/Verb to the request
func Method(verb string) http.Arrow {
	return func(ctx *http.Context) error {
		ctx.Method = verb
		return nil
	}
}

// Authority is part of URL, use the type to prevent escaping
type Authority string

// Path is part of URL, use the type to prevent path escaping
type Path string

// URI defines destination URI
// use Params arrow if you need to supply URL query params.
func URI(url string, args ...any) http.Arrow {
	return func(ctx *http.Context) error {
		if len(args) != 0 {
			url = mkURI(url, args)
		}

		if !strings.HasPrefix(url, "http") {
			url = ctx.Host + url
		}

		if !strings.HasPrefix(url, "http") {
			return &gurl.NotSupported{URL: url}
		}

		req, err := http.NewRequest(ctx.Method, url)
		if err != nil {
			return err
		}

		ctx.Request = req

		return nil
	}
}

func mkURI(uri string, args []any) string {
	opts := []any{}
	for _, x := range args {
		switch v := x.(type) {
		case *url.URL:
			v.Path = strings.TrimSuffix(v.Path, "/")
			opts = append(opts, v.String())
		case *Path:
			opts = append(opts, *v)
		case Path:
			opts = append(opts, v)
		case *Authority:
			opts = append(opts, *v)
		case Authority:
			opts = append(opts, v)
		case string:
			opts = append(opts, url.PathEscape(v))
		case *string:
			opts = append(opts, url.PathEscape(*v))
		case int:
			opts = append(opts, v)
		case *int:
			opts = append(opts, *v)
		default:
			opts = append(opts, url.PathEscape(urlSegment(x)))
		}
	}

	return fmt.Sprintf(uri, opts...)
}

func urlSegment(arg any) string {
	val := reflect.ValueOf(arg)

	if val.Kind() == reflect.Ptr {
		return fmt.Sprintf("%v", val.Elem())
	}

	return fmt.Sprintf("%v", val)
}

// Params appends query params to request URL. The arrow takes a struct and
// converts it to map[string]string. The function fails if input is not convertable
// to map of strings (e.g. contains nested struct).
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

// Param appends query params to request URL.
func Param[T interface{ string | int }](key string, val T) http.Arrow {
	return func(ctx *http.Context) error {
		uri := ctx.Request.URL
		q := uri.Query()
		switch v := any(val).(type) {
		case string:
			q.Add(key, v)
		case int:
			q.Add(key, strconv.Itoa(v))
		}

		uri.RawQuery = q.Encode()
		ctx.Request.URL = uri

		return nil
	}
}

// Header defines HTTP headers to the request
//
//	ø.Header("User-Agent", "gurl"),
func Header[T http.ReadableHeaderValues](header string, value T) http.Arrow {
	return HeaderOf[T](header).Set(value)
}

// Type of HTTP Header
//
//	const Host = HeaderOf[string]("Host")
//	ø.Host.Set("example.com")
type HeaderOf[T http.ReadableHeaderValues] string

// Sets value of HTTP header
func (h HeaderOf[T]) Set(value T) http.Arrow {
	switch v := any(value).(type) {
	case string:
		return func(cat *http.Context) error {
			cat.Request.Header.Add(string(h), v)
			return nil
		}
	case int:
		return func(cat *http.Context) error {
			cat.Request.Header.Add(string(h), strconv.Itoa(v))
			return nil
		}
	case time.Time:
		return func(cat *http.Context) error {
			cat.Request.Header.Add(string(h), v.UTC().Format(time.RFC1123))
			return nil
		}
	default:
		panic("invalid type")
	}
}

// Type of HTTP Header, Content-Type enumeration
//
//	const ContentType = HeaderEnumContent("Content-Type")
//	ø.ContentType.JSON
type HeaderEnumContent string

// Sets value of HTTP header
func (h HeaderEnumContent) Set(value string) http.Arrow {
	return func(cat *http.Context) error {
		cat.Request.Header.Add(string(h), value)
		return nil
	}
}

// ApplicationJSON defines header `???: application/json`
func (h HeaderEnumContent) ApplicationJSON(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "application/json")
	return nil
}

// JSON defines header `???: application/json`
func (h HeaderEnumContent) JSON(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "application/json")
	return nil
}

// Form defined Header `???: application/x-www-form-urlencoded`
func (h HeaderEnumContent) Form(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "application/x-www-form-urlencoded")
	return nil
}

// TextPlain defined Header `???: text/plain`
func (h HeaderEnumContent) TextPlain(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "text/plain")
	return nil
}

// Text defined Header `???: text/plain`
func (h HeaderEnumContent) Text(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "text/plain")
	return nil
}

// TextHTML defined Header `???: text/html`
func (h HeaderEnumContent) TextHTML(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "text/html")
	return nil
}

// HTML defined Header `???: text/html`
func (h HeaderEnumContent) HTML(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "text/html")
	return nil
}

// Type of HTTP Header, Connection enumeration
//
//	const Connection = HeaderEnumConnection("Connection")
//	ø.Connection.KeepAlive
type HeaderEnumConnection string

// Sets value of HTTP header
func (h HeaderEnumConnection) Set(value string) http.Arrow {
	return func(cat *http.Context) error {
		cat.Request.Header.Add(string(h), value)
		return nil
	}
}

// KeepAlive defines header `???: keep-alive`
func (h HeaderEnumConnection) KeepAlive(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "keep-alive")
	cat.Request.Close = false
	return nil
}

// Close defines header `???: close`
func (h HeaderEnumConnection) Close(cat *http.Context) error {
	cat.Request.Header.Add(string(h), "close")
	cat.Request.Close = true
	return nil
}

// Type of HTTP Header, Transfer-Encoding enumeration
//
//	const TransferEncoding = HeaderEnumTransferEncoding("Transfer-Encoding")
//	ø.TransferEncoding.Chunked
type HeaderEnumTransferEncoding string

// Sets value of HTTP header
func (h HeaderEnumTransferEncoding) Set(value string) http.Arrow {
	return func(cat *http.Context) error {
		cat.Request.TransferEncoding = strings.Split(value, ",")
		return nil
	}
}

// Chunked defines header `Transfer-Encoding: chunked`
func (h HeaderEnumTransferEncoding) Chunked(cat *http.Context) error {
	cat.Request.TransferEncoding = []string{"chunked"}
	return nil
}

// Identity defines header `Transfer-Encoding: identity`
func (h HeaderEnumTransferEncoding) Identity(cat *http.Context) error {
	cat.Request.TransferEncoding = []string{"identity"}
	return nil
}

// Header Content-Length
//
//	const ContentLength = HeaderEnumContentLength("Content-Length")
//	ø.ContentLength.Set(1024)
type HeaderEnumContentLength string

// Is sets a literal value of HTTP header
func (h HeaderEnumContentLength) Set(value int64) http.Arrow {
	return func(cat *http.Context) error {
		cat.Request.ContentLength = value
		return nil
	}
}

// List of supported HTTP header constants
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Request_fields
const (
	Accept            = HeaderEnumContent("Accept")
	AcceptCharset     = HeaderOf[string]("Accept-Charset")
	AcceptEncoding    = HeaderOf[string]("Accept-Encoding")
	AcceptLanguage    = HeaderOf[string]("Accept-Language")
	Authorization     = HeaderOf[string]("Authorization")
	CacheControl      = HeaderOf[string]("Cache-Control")
	Connection        = HeaderEnumConnection("Connection")
	ContentEncoding   = HeaderOf[string]("Content-Encoding")
	ContentLength     = HeaderEnumContentLength("Content-Length")
	ContentType       = HeaderEnumContent("Content-Type")
	Cookie            = HeaderOf[string]("Cookie")
	Date              = HeaderOf[time.Time]("Date")
	From              = HeaderOf[string]("From")
	Host              = HeaderOf[string]("Host")
	IfMatch           = HeaderOf[string]("If-Match")
	IfModifiedSince   = HeaderOf[time.Time]("If-Modified-Since")
	IfNoneMatch       = HeaderOf[string]("If-None-Match")
	IfRange           = HeaderOf[string]("If-Range")
	IfUnmodifiedSince = HeaderOf[time.Time]("If-Unmodified-Since")
	Origin            = HeaderOf[string]("Origin")
	Range             = HeaderOf[string]("Range")
	Referer           = HeaderOf[string]("Referer")
	TransferEncoding  = HeaderEnumTransferEncoding("Transfer-Encoding")
	UserAgent         = HeaderOf[string]("User-Agent")
	Upgrade           = HeaderOf[string]("Upgrade")
)

// Send payload to destination URL. You can also use native Go data types
// (e.g. maps, struct, etc) as egress payload. The library implicitly encodes
// input structures to binary using Content-Type as a hint. The function fails
// if content type is not supported by the library.
//
// The function accept a "classical" data container such as string, []bytes or
// io.Reader interfaces.
func Send(data any) http.Arrow {
	return func(cat *http.Context) error {
		chunked := cat.Request.Header.Get(string(TransferEncoding)) == "chunked"
		content := cat.Request.Header.Get(string(ContentType))
		if content == "" {
			return fmt.Errorf("unknown Content-Type")
		}

		switch stream := data.(type) {
		case string:
			cat.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(stream)))
			// cat.Request.GetBody = func() (io.ReadCloser, error) {
			// 	return io.NopCloser(bytes.NewBuffer([]byte(stream))), nil
			// }
			if !chunked && cat.Request.ContentLength == 0 {
				cat.Request.ContentLength = int64(len(stream))
			}
		case *strings.Reader:
			cat.Request.Body = io.NopCloser(stream)
			// snapshot := *stream
			// cat.Request.GetBody = func() (io.ReadCloser, error) {
			// 	r := snapshot
			// 	return io.NopCloser(&r), nil
			// }
			if !chunked && cat.Request.ContentLength == 0 {
				cat.Request.ContentLength = int64(stream.Len())
			}
		case []byte:
			cat.Request.Body = io.NopCloser(bytes.NewBuffer(stream))
			// cat.Request.GetBody = func() (io.ReadCloser, error) {
			// 	return io.NopCloser(bytes.NewBuffer(stream)), nil
			// }
			if !chunked && cat.Request.ContentLength == 0 {
				cat.Request.ContentLength = int64(len(stream))
			}
		case *bytes.Buffer:
			cat.Request.Body = io.NopCloser(stream)
			// snapshot := stream.Bytes()
			// cat.Request.GetBody = func() (io.ReadCloser, error) {
			// 	return io.NopCloser(bytes.NewBuffer(snapshot)), nil
			// }
			if !chunked && cat.Request.ContentLength == 0 {
				cat.Request.ContentLength = int64(stream.Len())
			}
		case *bytes.Reader:
			cat.Request.Body = io.NopCloser(stream)
			// snapshot := *stream
			// cat.Request.GetBody = func() (io.ReadCloser, error) {
			// 	r := snapshot
			// 	return io.NopCloser(&r), nil
			// }
			if !chunked && cat.Request.ContentLength == 0 {
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
			if !chunked && cat.Request.ContentLength == 0 {
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
