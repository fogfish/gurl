//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

/*
Package http defines category of HTTP I/O, "do"-notation becomes

	http.Join(
	  ø...,
	  ø...,

	  ƒ...,
	  ƒ...,
	)

Symbol `ø` (option + o) is an convenient alias to module gurl/http/send, which
defines writer morphism that focuses inside and reshapes HTTP protocol request.
The writer morphism is used to declare HTTP method, destination URL, request headers
and payload.

Symbol `ƒ` (option + f) is an convenient alias to module gurl/http/recv, which
defines reader morphism that focuses into side-effect, HTTP protocol response.
The reader morphism is a pattern matcher, is used to match HTTP response code,
headers and response payload. It helps us to declare our expectations on the response.
The evaluation of "program" fails if expectations do not match actual response.

Let's look on step-by-step usage of the category.

**Method and URL** are mandatory. It has to be a first element in the construction.

	http.GET(
	  ø.URI("http://example.com"),
	  ...
	)

Definition of **request headers** is an optional. You can list as many headers as
needed. Either using string literals or variables. Some frequently used headers
implements aliases (e.g. ø.ContentJSON(), ...)

	http.GET(
	  ...
	  ø.Header("Accept", "application/json"),
	  ø.Header("Authorization", &token),
	  ...
	)

The **request payload** is also an optional. You can also use native Golang data types
as egress payload. The library implicitly encodes input structures to binary
using Content-Type as a hint.

	http.GET(
	  ...
	  ø.Send(MyType{Hello: "World"}),
	  ...
	)

The declaration of expected response is always starts with mandatory HTTP **status
code**. The execution fails if peer responds with other than specified value.

	http.GET(
	  ...
	  ƒ.Code(http.StatusCodeOK),
	  ...
	)

It is possible to match presence of header in the response, match its entire
content or lift the header value to a variable. The execution fails if HTTP
response do not match the expectation.

	http.GET(
	  ...
	  ƒ.Header("Content-Type", "application/json"),
	  ...
	)

The library is able to **decode payload** into Golang native data structure
using Content-Type header as a hint.

	var data MyType
	http.GET(
	  ...
	  ƒ.Recv(&data)
	  ...
	)

Please note, the library implements lenses to inline assert of decoded content.
See the documentation of gurl/http/recv module.
*/
package http
