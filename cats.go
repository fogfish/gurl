//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package gurl

import (
	"bytes"
	"net/http"
	"net/url"
)

/*

IOCatHTTP defines the category of HTTP I/O
*/
type IOCatHTTP struct {
	Send *UpStreamHTTP
	Recv *DnStreamHTTP
}

/*

UpStreamHTTP specify parameters for HTTP requests
*/
type UpStreamHTTP struct {
	Method  string
	URL     *url.URL
	Header  map[string]*string
	Payload *bytes.Buffer
}

/*

DnStreamHTTP specify parameters for HTTP response
*/
type DnStreamHTTP struct {
	*http.Response
	Payload interface{}
}
