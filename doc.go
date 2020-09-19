//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

/*

Package gurl is a class of High Order Component which can do http requests
with few interesting property such as composition and laziness. The library
implements rough and naive Haskell's equivalent of do-notation, so called
monadic binding form. This construction decorates http i/o pipeline(s) with
"programmable commas".

Inspiration

Microservices have become a design style to evolve system architecture
in parallel, implement stable and consistent interfaces. An expressive
language is required to design the variety of network communication use-cases.
A pure functional languages fits very well to express communication behavior.
The language gives a rich techniques to hide the networking complexity using
monads as abstraction. The IO-monads helps us to compose a chain of network
operations and represent them as pure computation, build a new things from
small reusable elements. The library is implemented after Erlang's
https://github.com/fogfish/m_http

The library attempts to adapts a human-friendly syntax of HTTP request/response
logging/definition used by curl with Behavior as a Code paradigm. It tries to
connect cause-and-effect (Given/When/Then) with the networking (Input/Process/Output).

	> GET / HTTP/1.1
	> Host: example.com
	> User-Agent: curl/7.54.0
	> Accept: application/json
	>
	< HTTP/1.1 200 OK
	< Content-Type: text/html; charset=UTF-8
	< Server: ECS (phd/FD58)
	< ...

This semantic provides an intuitive approach to specify HTTP requests/responses.
Adoption of this syntax as Go native code provides a rich capability to network
programming.

Key features

↣ cause-and-effect abstraction of HTTP request/response, naive do-notation

↣ high-order composition of individual HTTP requests to complex networking computations

↣ human-friendly, Go native and declarative syntax to depict HTTP operations

↣ implements a declarative approach for testing of RESTful interfaces

↣ automatically encodes/decodes Go native HTTP payload using Content-Type hints

↣ supports generic transformation to algebraic data types

↣ simplify error handling with naive Either implementation


IO Category

Standard Golang packages implements low-level HTTP interface, which
requires knowledge about protocol itself, aspects of Golang implementation,
a bit of boilerplate coding and lack of standardized chaining (composition)
of individual requests.

gurl library inherits an ability of pure functional languages to express
communication behavior by hiding the networking complexity using category pattern
(aka "do"-notation). This pattern helps us to compose a chain of network operations
and represent them as pure computation, build a new things from small reusable
elements. This library uses the "do"-notation, so called monadic binding form.
It is well know in functional programming languages such as Haskell and
Scala. The networking becomes a collection of composed "do"-notation in context
of a state monad.

A composition of HTTP primitives within the category are written with the following
syntax.

  gurl.Join(arrows ...Arrow) Arrow

Here, each arrow is a morphism applied to HTTP protocol. The implementation defines
an abstraction of the protocol environments and lenses to focus inside it. In other
words, the category represents the environment as an "invisible" side-effect of
the composition.

`gurl.Join(arrows ...Arrow) Arrow` and its composition implements lazy I/O. It only
returns a "promise", you have to evaluate it in the context of IO instance.

  io := gurl.IO()
  fn := gurl.Join( ... )
  fn(io)


Basics

The following code snippet demonstrates a typical usage scenario.

  import (
    "github.com/fogfish/gurl"
    "github.com/fogfish/gurl/http"
    ƒ "github.com/fogfish/gurl/http/recv"
    ø "github.com/fogfish/gurl/http/send"
  )

  // You can declare any types and use them as part of networking I/O.
  type Payload struct {
    Origin string `json:"origin"`
    Url    string `json:"url"`
  }

  var data Payload
  var reqf := http.Join(
    // declares HTTP method and destination URL
    ø.GET("http://httpbin.org/get"),
    // HTTP content negotiation, declares acceptable types
    ø.Accept().Is("application/json"),

    // requires HTTP Status Code to be 200 OK
    ƒ.Code(gurl.StatusCodeOK),
    // requites HTTP Header to be Content-Type: application/json
    ƒ.Served().Is("application/json"),
    // unmarshal JSON to the variable
    ƒ.Recv(&data),
  )

  // Note: http do not hold yet, a results of HTTP I/O
  //       it is just a composable "promise", you have to
  //       evaluate a side-effect of HTTP "computation"
  if reqf(gurl.IO()).Fail != nil {
    // error handling
  }

The evaluation of "program" fails if either networking fails or expectations
do not match actual response. There are no needs to check error code after
each operation. The composition is smart enough to terminate "program" execution.

Composition

The composition of multiple HTTP I/O is an essential part of the library.
The composition is handled in context of IO category. For example,
RESTfull API primitives declared as function, each deals with gurl.IOCat.

  import (
    "github.com/fogfish/gurl"
    "github.com/fogfish/gurl/http"
    ƒ "github.com/fogfish/gurl/http/recv"
    ø "github.com/fogfish/gurl/http/send"
  )

  func HoF() {
    var (
      token AccessToken
      user User
      org Org
    )

    reqf := gurl.Join(
      AccessToken(&token),
      UserProfile(&token, &user),
      UserContribution(&token, &org)
    )

    if reqf(gurl.IO()).Fail != nil {
      // error handling
    }
  }

  func AccessToken(token *AccessToken) gurl.Arrow {
    return http.Join(
      // ...
      ƒ.Recv(token),
    )
  }

  func UserProfile(token *AccessToken, user *User) gurl.Arrow {
    return http.Join(
      ø.POST("..."),
      ø.Authorization().Is(token.Bearer),
      // ...
      ƒ.Recv(user),
    )
  }

  func UserContribution(token *AccessToken, org *Org) {
    return http.Join(
      ø.POST("..."),
      ø.Authorization().Is(token.Bearer),
      // ...
      ƒ.Recv(org),
    )
  }
*/
package gurl
