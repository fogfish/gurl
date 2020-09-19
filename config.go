//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package gurl

/*

Log Level constants, use with Logging config
	- Level 0: disable debug logging (default)
	- Level 1: log only egress traffic
  - Level 2: log only ingress traffic
  - Level 3: log full content of packets
*/
const (
	//
	LogLevelNone    = 0
	LogLevelEgress  = 1
	LogLevelIngress = 2
	LogLevelDebug   = 3
)

/*

Logging enables debug logging of IO traffic
*/
func Logging(level int) Config {
	return func(cat *IOCat) *IOCat {
		cat.LogLevel = level
		return cat
	}
}

/*

SideEffect defines "unsafe" behavior for category
*/
func SideEffect(arrow Arrow) Config {
	return func(cat *IOCat) *IOCat {
		cat.sideEffect = arrow
		return cat
	}
}
