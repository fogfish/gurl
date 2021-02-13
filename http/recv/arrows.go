//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

// Package recv defines a pure computations to compose HTTP response receivers
package recv

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ajg/form"
	"github.com/fogfish/gurl"
	"github.com/fogfish/gurl/http"
)

//-------------------------------------------------------------------
//
// core arrows
//
//-------------------------------------------------------------------

/*

Code is a mandatory statement to match expected HTTP Status Code against
received one. The execution fails StatusCode error if service responds
with other value then specified one.
*/
func Code(code ...http.StatusCode) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		if cat = cat.Unsafe(); cat.Fail != nil {
			return cat
		}

		status := cat.HTTP.Recv.StatusCode
		if !hasCode(code, status) {
			cat.Fail = http.NewStatusCode(status, code[0])
		}
		return cat
	}
}

func hasCode(s []http.StatusCode, e int) bool {
	for _, a := range s {
		if a.Value() == e {
			return true
		}
	}
	return false
}

// THeader is tagged string, represents HTTP Header
type THeader struct{ string }

/*

Header matches presence of header in the response or match its entire content.
The execution fails with BadMatchHead if the matched value do not meet expectations.
*/
func Header(header string) THeader {
	return THeader{header}
}

// Is matches value of HTTP header, Use wildcard string ("*") to match any header value
func (header THeader) Is(value string) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		h := cat.HTTP.Recv.Header.Get(header.string)
		if h == "" {
			cat.Fail = &gurl.Mismatch{
				Diff:    fmt.Sprintf("- %s: %s", header.string, value),
				Payload: nil,
			}
			return cat
		}

		if value != "*" && !strings.HasPrefix(h, value) {
			cat.Fail = &gurl.Mismatch{
				Diff:    fmt.Sprintf("+ %s: %s\n- %s: %s", header.string, h, header.string, value),
				Payload: map[string]string{header.string: h},
			}
			return cat
		}

		return cat
	}
}

// String matches a header value to closed variable of string type.
func (header THeader) String(value *string) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		val := cat.HTTP.Recv.Header.Get(header.string)
		if val == "" {
			cat.Fail = &gurl.Mismatch{
				Diff:    fmt.Sprintf("- %s: *", header.string),
				Payload: nil,
			}
		} else {
			*value = val
		}

		return cat
	}
}

// Any matches a header value, syntax sugar of Header(...).Is("*")
func (header THeader) Any() http.Arrow {
	return header.Is("*")
}

/*

Recv applies auto decoders for response and returns either binary or
native Go data structure. The Content-Type header give a hint to decoder.
Supply the pointer to data target data structure.
*/
func Recv(out interface{}) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		cat.Fail = decode(
			cat.HTTP.Recv.Header.Get("Content-Type"),
			cat.HTTP.Recv.Body,
			&out,
		)
		cat.HTTP.Recv.Body.Close()
		cat.HTTP.Recv.Payload = out
		cat.HTTP.Recv.Response = nil
		return cat
	}
}

func decode(content string, stream io.ReadCloser, data interface{}) error {
	switch {
	case strings.Contains(content, "json"):
		return json.NewDecoder(stream).Decode(&data)
	case strings.Contains(content, "www-form"):
		return form.NewDecoder(stream).Decode(&data)
	default:
		return &gurl.Mismatch{
			Diff:    fmt.Sprintf("- Content-Type: application/*\n+ Content-Type: %s", content),
			Payload: map[string]string{"Content-Type": content},
		}
	}
}

/*

Bytes receive raw binary from HTTP response
*/
func Bytes(val *[]byte) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		*val, cat.Fail = ioutil.ReadAll(cat.HTTP.Recv.Body)
		cat.Fail = cat.HTTP.Recv.Body.Close()
		cat.HTTP.Recv.Response = nil
		cat.HTTP.Recv.Payload = string(*val)
		return cat
	}
}

//-------------------------------------------------------------------
//
// alias arrows
//
//-------------------------------------------------------------------

// Served is a syntax sugar of Header("Content-Type")
func Served() THeader {
	return Header("Content-Type")
}

// ServedJSON is a syntax sugar of Header("Content-Type").Is("application/json")
func ServedJSON() http.Arrow {
	return Served().Is("application/json")
}

// ServedForm is a syntax sugar of Header("Content-Type", "application/x-www-form-urlencoded")
func ServedForm() http.Arrow {
	return Served().Is("application/x-www-form-urlencoded")
}
