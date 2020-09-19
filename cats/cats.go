//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package cats

import (
	"reflect"
	"sort"

	"github.com/fogfish/gurl"
	"github.com/google/go-cmp/cmp"
)

/*

FMap applies clojure to category.
The function lifts any computation to the category and make it composable
with the "program".
*/
func FMap(f func() error) gurl.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		cat.Fail = f()
		return cat
	}
}

/*

FlatMap applies closure to matched HTTP request.
It returns an arrow, which continue evaluation.
*/
func FlatMap(f func() gurl.Arrow) gurl.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		if g := f(); g != nil {
			return g(cat)
		}
		return cat
	}
}

/*

Defined checks if the value is defined, use a pointer to the value.
*/
func Defined(value interface{}) gurl.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		va := reflect.ValueOf(value)
		if va.Kind() == reflect.Ptr {
			va = va.Elem()
		}

		if !va.IsValid() {
			cat.Fail = &gurl.Undefined{Type: va.Type().Name()}
		}

		if va.IsValid() && va.IsZero() {
			cat.Fail = &gurl.Undefined{Type: va.Type().Name()}
		}
		return cat
	}
}

// TValue is tagged type, represent matchers
type TValue struct{ actual interface{} }

/*
Value checks if the value equals to defined one.
Supply the pointer to actual value
*/
func Value(val interface{}) TValue {
	return TValue{val}
}

// Is matches a value
func (val TValue) Is(require interface{}) gurl.Arrow {
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
func (val TValue) String(require string) gurl.Arrow {
	return val.Is(&require)
}

// Bytes matches a literal value of bytes
func (val TValue) Bytes(require []byte) gurl.Arrow {
	return val.Is(&require)
}

// TSeq is tagged type, represents Sequence of elements
type TSeq struct{ gurl.Ord }

/*

Seq matches presence of element in the sequence.
*/
func Seq(seq gurl.Ord) TSeq {
	return TSeq{seq}
}

/*

Has lookups element using key and matches expected value
*/
func (seq TSeq) Has(key string, expect ...interface{}) gurl.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		sort.Sort(seq)
		i := sort.Search(seq.Len(), func(i int) bool { return seq.String(i) >= key })
		if i < seq.Len() && seq.String(i) == key {
			if len(expect) > 0 {
				if diff := cmp.Diff(seq.Value(i), expect[0]); diff != "" {
					cat.Fail = &gurl.Mismatch{
						Diff:    diff,
						Payload: seq.Value(i),
					}
				}
			}
			return cat
		}
		cat.Fail = &gurl.Undefined{Type: key}
		return cat
	}
}
