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
	"io"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"

	"github.com/ajg/form"
	"github.com/fogfish/gurl"
	"github.com/google/go-cmp/cmp"
)

/*

Code is a mandatory statement to match expected HTTP Status Code against
received one. The execution fails with BadMatchCode if service responds
with other value then specified one.
*/
func Code(code ...int) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.Unsafe()
		if io.Fail != nil {
			return io
		}

		status := io.HTTP.Ingress.StatusCode
		if !hasCode(code, status) {
			io.Fail = &gurl.BadMatchCode{Expect: code, Actual: status}
		}
		return io
	}
}

func hasCode(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// HtHeader is tagged string, represents HTTP Header
type HtHeader struct{ string }

/*

Header matches presence of header in the response or match its entire content.
The execution fails with BadMatchHead if the matched value do not meet expectations.
*/
func Header(header string) HtHeader {
	return HtHeader{header}
}

// Is matches value of HTTP header, Use wildcard string ("*") to match any header value
func (header HtHeader) Is(value string) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		h := io.HTTP.Ingress.Header.Get(header.string)
		if h == "" {
			io.Fail = &gurl.BadMatchHead{Header: header.string, Expect: value}
		} else if value != "*" && !strings.HasPrefix(h, value) {
			io.Fail = &gurl.BadMatchHead{Header: header.string, Expect: value, Actual: h}
		}

		return io
	}
}

// String matches a header value to closed variable of string type.
func (header HtHeader) String(value *string) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		val := io.HTTP.Ingress.Header.Get(header.string)
		if val == "" {
			io.Fail = &gurl.BadMatchHead{Header: header.string}
		} else {
			*value = val
		}

		return io
	}
}

// Any matches a header value, syntax sugar of Header(...).Is("*")
func (header HtHeader) Any() gurl.Arrow {
	return header.Is("*")
}

// Served is a syntax sugar of Header("Content-Type")
func Served() HtHeader {
	return Header("Content-Type")
}

// ServedJSON is a syntax sugar of Header("Content-Type").Is("application/json")
func ServedJSON() gurl.Arrow {
	return Served().Is("application/json")
}

// ServedForm is a syntax sugar of Header("Content-Type", "application/x-www-form-urlencoded")
func ServedForm() gurl.Arrow {
	return Served().Is("application/x-www-form-urlencoded")
}

/*

Recv applies auto decoders for response and returns either binary or
native Go data structure. The Content-Type header give a hint to decoder.
Supply the pointer to data target data structure.
*/
func Recv(out interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.Fail = decode(
			io.HTTP.Ingress.Header.Get("Content-Type"),
			io.HTTP.Ingress.Body,
			&out,
		)
		io.HTTP.Ingress.Body.Close()
		io.Body = out
		io.HTTP.Ingress = nil
		return io
	}
}

func decode(content string, stream io.ReadCloser, data interface{}) error {
	switch {
	case strings.HasPrefix(content, "application/json"):
		return json.NewDecoder(stream).Decode(&data)
	case strings.HasPrefix(content, "application/x-www-form-urlencoded"):
		return form.NewDecoder(stream).Decode(&data)
	default:
		return &gurl.BadMatchHead{
			Header: "Content-Type",
			Actual: content,
		}
	}
}

/*

Bytes receive raw binary from HTTP response
*/
func Bytes(val *[]byte) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		*val, io.Fail = ioutil.ReadAll(io.HTTP.Ingress.Body)
		io.Fail = io.HTTP.Ingress.Body.Close()
		io.HTTP.Ingress = nil
		return io
	}
}

/*

FMap applies clojure to matched HTTP request.
*/
func FMap(f func() error) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		io.Fail = f()
		return io
	}
}

/*

FlatMap applies closure to matched HTTP request.
It returns an arrow, which continue evaluation.
*/
func FlatMap(f func() gurl.Arrow) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		if g := f(); g != nil {
			return g(io)
		}
		return io
	}
}

/*

Defined checks if the value is defined, use a pointer to the value.
*/
func Defined(value interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		va := reflect.ValueOf(value)
		if va.Kind() == reflect.Ptr {
			va = va.Elem()
		}

		if !va.IsValid() {
			io.Fail = &gurl.Undefined{Type: va.Type().Name()}
		}

		if va.IsValid() && va.IsZero() {
			io.Fail = &gurl.Undefined{Type: va.Type().Name()}
		}
		return io
	}
}

// HtValue is tagged type, represent matchers
type HtValue struct{ actual interface{} }

/*

Value checks if the value equals to defined one.
Supply the pointer to actual value
*/
func Value(val interface{}) HtValue {
	return HtValue{val}
}

// Is matches a value
func (val HtValue) Is(require interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		if diff := cmp.Diff(val.actual, require); diff != "" {
			io.Fail = &gurl.Mismatch{
				Diff:    diff,
				Payload: val.actual,
			}
		}
		return io
	}
}

// String matches a literal value
func (val HtValue) String(require string) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		if diff := cmp.Diff(val.actual, &require); diff != "" {
			io.Fail = &gurl.Mismatch{
				Diff:    diff,
				Payload: val.actual,
			}
		}
		return io
	}
}

// HtSeq is tagged type, represents Sequence of elements
type HtSeq struct{ gurl.Ord }

/*

Seq matches presence of element in the sequence.
*/
func Seq(seq gurl.Ord) HtSeq {
	return HtSeq{seq}
}

/*

Has lookups element using key and matches expected value
*/
func (seq HtSeq) Has(key string, expect ...interface{}) gurl.Arrow {
	return func(io *gurl.IOCat) *gurl.IOCat {
		sort.Sort(seq)
		i := sort.Search(seq.Len(), func(i int) bool { return seq.String(i) >= key })
		if i < seq.Len() && seq.String(i) == key {
			if len(expect) > 0 {
				if diff := cmp.Diff(seq.Value(i), expect[0]); diff != "" {
					io.Fail = &gurl.Mismatch{
						Diff:    diff,
						Payload: seq.Value(i),
					}
				}
			}
			return io
		}
		io.Fail = &gurl.Undefined{Type: key}
		return io
	}
}
