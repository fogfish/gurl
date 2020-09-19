//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package gurl_test

import (
	"errors"
	"testing"

	"github.com/fogfish/gurl"
	"github.com/fogfish/it"
)

func identity() gurl.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		return cat
	}
}

func fail() gurl.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		cat.Fail = errors.New("fail")
		return cat
	}
}

func TestJoin(t *testing.T) {
	for _, f := range []gurl.Arrow{
		gurl.Join(identity(), identity()),
		identity().Then(identity()),
	} {
		it.Ok(t).
			If(
				f(gurl.IO()).Fail,
			).Should().Equal(nil)
	}
}

func TestJoinFail(t *testing.T) {
	for _, f := range []gurl.Arrow{
		gurl.Join(identity(), fail()),
		gurl.Join(fail(), identity()),
	} {
		it.Ok(t).
			If(
				f(gurl.IO()).Fail,
			).ShouldNot().Equal(nil)
	}
}
