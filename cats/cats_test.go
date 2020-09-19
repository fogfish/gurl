//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package cats_test

import (
	"errors"
	"testing"

	"github.com/fogfish/gurl"
	ç "github.com/fogfish/gurl/cats"
	"github.com/fogfish/it"
)

func identity() gurl.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		return cat
	}
}

func TestFMap(t *testing.T) {
	var s string

	f := gurl.Join(
		identity(),
		ç.FMap(func() error {
			s = "value"
			return nil
		}),
	)

	it.Ok(t).
		If(
			f(gurl.IO()).Fail,
		).Should().Equal(nil).
		If(s).Should().Equal("value")
}

func TestFMapError(t *testing.T) {
	f := gurl.Join(
		identity(),
		ç.FMap(func() error {
			return errors.New("fail")
		}),
	)

	it.Ok(t).
		If(
			f(gurl.IO()).Fail,
		).ShouldNot().Equal(nil)

}

func TestFlatMap(t *testing.T) {
	seq := ""

	f := ç.FMap(func() error {
		seq = seq + "a"
		return nil
	})

	g := gurl.Join(
		f,
		ç.FlatMap(func() gurl.Arrow { return f }),
	)

	it.Ok(t).
		If(
			g(gurl.IO()).Fail,
		).Should().Equal(nil).
		If(seq).Should().Equal("aa")
}

func TestDefined(t *testing.T) {
	type Site struct {
		Site string
		Host string
	}
	var site Site

	f := gurl.Join(
		ç.FMap(func() error {
			site = Site{"site", "host"}
			return nil
		}),
		ç.Defined(&site),
		ç.Defined(&site.Site),
	)

	it.Ok(t).
		If(
			f(gurl.IO()).Fail,
		).Should().Equal(nil)
}

func TestNotDefined(t *testing.T) {
	type Site struct {
		Site string
		Host string
	}
	var site Site

	f := gurl.Join(
		ç.FMap(func() error {
			site = Site{"site", ""}
			return nil
		}),
		ç.Defined(&site),
		ç.Defined(&site.Host),
	)

	it.Ok(t).
		If(
			f(gurl.IO()).Fail,
		).ShouldNot().Equal(nil)
}

func TestValue(t *testing.T) {
	type Site struct {
		Site string
		Host []byte
	}
	var site Site

	f := gurl.Join(
		ç.FMap(func() error {
			site = Site{"site", []byte("abc")}
			return nil
		}),
		ç.Value(&site).Is(&Site{"site", []byte("abc")}),
		ç.Value(&site.Site).String("site"),
		ç.Value(&site.Host).Bytes([]byte("abc")),
	)

	it.Ok(t).
		If(
			f(gurl.IO()).Fail,
		).Should().Equal(nil)
}

func TestValueNoMatch(t *testing.T) {
	type Site struct {
		Site string
		Host []byte
	}
	var site Site

	f := gurl.Join(
		ç.FMap(func() error {
			site = Site{"site", []byte("abc")}
			return nil
		}),
		ç.Value(&site).Is(&Site{"site1", []byte("abc")}),
	)

	it.Ok(t).
		If(
			f(gurl.IO()).Fail,
		).ShouldNot().Equal(nil)
}

type E struct{ Site string }

type Seq []E

func (seq Seq) Len() int                { return len(seq) }
func (seq Seq) Swap(i, j int)           { seq[i], seq[j] = seq[j], seq[i] }
func (seq Seq) Less(i, j int) bool      { return seq[i].Site < seq[j].Site }
func (seq Seq) String(i int) string     { return seq[i].Site }
func (seq Seq) Value(i int) interface{} { return seq[i] }

var seqMock Seq = Seq{
	{Site: "q.example.com"},
	{Site: "a.example.com"},
	{Site: "z.example.com"},
	{Site: "w.example.com"},
	{Site: "s.example.com"},
	{Site: "x.example.com"},
	{Site: "e.example.com"},
	{Site: "d.example.com"},
	{Site: "c.example.com"},
}

func TestSeqHas(t *testing.T) {
	var seq Seq
	expectS := E{Site: "s.example.com"}
	expectZ := E{Site: "z.example.com"}

	f := gurl.Join(
		ç.FMap(func() error {
			seq = seqMock
			return nil
		}),
		ç.Seq(&seq).Has(expectS.Site),
		ç.Seq(&seq).Has(expectS.Site, expectS),
		ç.Seq(&seq).Has(expectZ.Site, expectZ),
	)

	it.Ok(t).
		If(
			f(gurl.IO()).Fail,
		).Should().Equal(nil)
}

func TestSeqHasNotFound(t *testing.T) {
	var seq Seq
	expect0 := E{Site: "0.example.com"}

	f := gurl.Join(
		ç.FMap(func() error {
			seq = seqMock
			return nil
		}),
		ç.Seq(&seq).Has(expect0.Site),
	)

	it.Ok(t).
		If(
			f(gurl.IO()).Fail,
		).ShouldNot().Equal(nil)
}

func TestSeqHasNoMatch(t *testing.T) {
	var seq Seq
	expectS := E{Site: "s.example.com"}
	expectZ := E{Site: "z.example.com"}

	f := gurl.Join(
		ç.FMap(func() error {
			seq = seqMock
			return nil
		}),
		ç.Seq(&seq).Has(expectS.Site, expectZ),
	)

	it.Ok(t).
		If(
			f(gurl.IO()).Fail,
		).ShouldNot().Equal(nil)
}
