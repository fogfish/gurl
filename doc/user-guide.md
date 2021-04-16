<p align="center">
  <h3 align="center">ᵍ🆄🆁🅻</h3>
  <p align="center"><strong>User Guide</strong></p>
</p>

- [Overview](#overview)
- [Composition](#composition)
- [Life-cycle](#life-cycle)
- [Import library](#import-library)
- [Arrow types](#arrow-types)
  - [Writer morphism](#writer-morphism)
    - [Method and URL](#method-and-url)
    - [Query Params](#query-params)
    - [Request Headers](#request-headers)
    - [Request Payload](#request-payload)
  - [Reader morphism](#reader-morphism)


---


## Overview

ᵍ🆄🆁🅻 is a "combinator" library for network I/O. Combinators open up an opportunity to depict computation problems in terms of fundamental elements like physics talks about universe in terms of particles. The only definite purpose of combinators are building blocks for composition of "atomic" functions into computational structures. ᵍ🆄🆁🅻 combinators provide a powerful symbolic expressions in networking domain.

Standard Golang packages implements low-level HTTP interface, which requires knowledge about protocol itself, understanding of Golang implementation aspects, and a bit of boilerplate coding. It also missing standardized chaining (composition) of individual requests. ᵍ🆄🆁🅻 inherits an ability of pure functional languages to express communication behavior by hiding the networking complexity using combinators. The composition becomes a fundamental operation in the library: the codomain of `𝒇` be the domain of `𝒈` so that the composite operation `𝒇 ◦ 𝒈` is defined.

The library uses `Arrow` as a key abstraction of combinators. It is a *pure function* that takes an abstraction of the protocol environments, so called IO category and applies morphism as an "invisible" side-effect of the composition.

```go
/*

Arrow: IO ⟼ IO
*/
type Arrow func(*gurl.IOCat) *gurl.IOCat
```

There are two classes of arrows. The first class is a writer morphism that focuses inside and reshapes HTTP protocol request. The writer morphism is used to declare HTTP method, destination URL, request headers and payload. Second one is a reader morphism that focuses into side-effect of HTTP protocol. The reader morphism is a pattern matcher, is used to match HTTP response code, headers and response payload.

Example of HTTP I/O visualization made by curl give **naive** perspective about arrows.

```
> GET / HTTP/1.1
> Host: example.com
> User-Agent: curl/7.54.0
> Accept: application/json
>
< HTTP/1.1 200 OK
< Content-Type: text/html; charset=UTF-8
< Server: ECS (phd/FD58)
< ...
```


## Composition

`Arrow` can be composed with another `Arrow` into new `Arrow` and so on. The library support only "and-then" style. It builds a strict product Arrow: `A × B × C × ... ⟼ D`. The product type takes a protocol environment and applies "morphism" sequentially unless some step fails. Use variadic function `http.Join` to compose HTTP primitives:

```go
/*

Join composes HTTP arrows to high-order function
(a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
*/
func http.Join(arrows ...http.Arrow) gurl.Arrow

//
var a: http.Arrow = /* ... */
var b: http.Arrow = /* ... */
var c: http.Arrow = /* ... */

d := http.Join(a, b, c)
```

Ease of the composition is one of major intent why ᵍ🆄🆁🅻 library has deviated from standard Golang HTTP interface. `http.Join` produces instance of higher order `gurl.Arrow` type, which is composable “promises” of HTTP I/O and so on. Essentially, the network I/O is just set of `Arrow` functions.

```go
/*

Join composes arrows to high-order function
(a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
*/
func Join(arrows ...gurl.Arrow) gurl.Arrow

//
var a: gurl.Arrow = http.Join(/* ... */)
var b: gurl.Arrow = http.Join(/* ... */)
var c: gurl.Arrow = http.Join(/* ... */)

d := gurl.Join(a, b, c)
```

These rules of `Arrow` composition allow anyone to build a complex HTTP I/O scenario from small re-usable block.



## Life-cycle

```go
lazy := gurl.Join(/* ... */)
```

The instance of `Arrow` produced by one of `Join` functions does not hold a results of HTTP I/O. It only builds a composable "promise" - a pure computation. The computation needs to be evaluated by applying it over the protocol environment. The library provides simple interface to create and customize the environment.  

```go
// HTTP protocol provides default out-of-the-box environment.
// use the default environment with caution in production workload
env := http.DefaultIO() 

// apply the computation over the environment
env = lazy(env)

// handle either networking error or failure of expectation.
// there are not need for error handling after each operation
// in the contrast with classical networking I/O. The combinator
// is smart enough to terminate execution.
if env.Fail != nil {
  // ...
}

// the environment holds failed state until it is recovered 
if err := env.Recover(); err != nil {
  // ...
}
```

Usage of the library for production workload require a careful configuration of HTTP protocol timeouts, TLS policies, etc. Another aspect is thread safeness. The protocol environment is not thread safe. Each golang routine shall create one. Re-use for `http.Client` pointer across environments reduces resources consumption.  

```go
env := gurl.IO(
  http.Stack(&http.Client{/* ... */})
)
```


## Import library

The library consists of multiple packages, import them all 

```go
import (
  // core types
  "github.com/fogfish/gurl"

  // support for http protocol
  "github.com/fogfish/gurl/http"

  // writer morphism is used to declare HTTP method,
  // destination URL, request headers and payload
  // single letter alias (e.g. ø) makes the code less verbose
  ø "github.com/fogfish/gurl/http/send"

  // reader morphism is a pattern matcher for HTTP response code,
  // headers and response payload
  // single letter alias (e.g. ƒ) makes the code less verbose
  ƒ "github.com/fogfish/gurl/http/recv"
)
```


## Arrow types

ᵍ🆄🆁🅻 library delivers set of built-in arrows to deal with HTTP I/O.

### Writer morphism

Writer morphism focuses inside and reshapes HTTP requests. The writer morphism is used to declare HTTP method, destination URL, request headers and payload.


#### Method and URL

Method and URL are only mandatory writer morphism in I/O declaration. Use `type Method string` to declare verb of HTTP request. It's received method `URL` allows to specify a destination endpoint.

```go
http.Join(
  ø.Method("GET").URL("http://example.com"),
)

// The library implements a syntax sugar for mostly used HTTP Verbs
http.Join(
  ø.GET.URL("http://example.com"),
)
```

The `URL` receiver is equivalent to `fmt.Sprintf`. It uses [percent encoding](https://golang.org/pkg/fmt/) to format and escape values.

```go
http.Join(
  ø.GET.URL("http://%s/%s", "example.com", "foo"),
)

// All path segments are escaped by default, use ! symbol to disable it
http.Join(
  // this do not work
  ø.GET.URL("%s/%s", "http://example.com", "foo/bar"),

  // this works
  ø.GET.URL("!%s/%s", "http://example.com", "foo/bar"),
)
```


#### Query Params

It is possible to inline query parameters into URL. However, this is not a type-safe approach.

```go
http.Join(
  ø.GET.URL("http://example.com/?tag=%s", "foo"),
)
```

The `func Params(query interface{}) http.Arrow` combinator lifts any flat structure to query parameters.

```go
type MyParam struct {
  Site string `json:"site,omitempty"`
  Host string `json:"host,omitempty"`
}

http.Join(
  // ...
  ø.Params(MyParam{Site: "site", Host: "host"}),
),
```


#### Request Headers

Use `type Header string` to declare headers and its values.

```go
http.Join(
  // ...
  ø.Header("Content-Type").Is("application/json"),
)

// The library implements a syntax sugar for mostly used HTTP headers
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Request_fields
http.Join(
  // ...
  ø.Authorization.Is("Bearer eyJhbGciOiJIU...adQssw5c"),
)

// The library implements a syntax sugar for content negotiation headers
http.Join(
  // ...
  ø.Accept.JSON,
  ø.ContentType.HTML,
)
```

#### Request payload

The `func Send(data interface{}) http.Arrow` transmits the payload to destination URL. The function takes Go data types (e.g. maps, struct, etc) and encodes its to binary using `Content-Type` header as a hint. The function fails if content type is not defined or not supported by the library.

```go
type MyType struct {
  Site string `json:"site,omitempty"`
  Host string `json:"host,omitempty"`
}

// Encode struct to JSON
http.Join(
  // ...
  ø.ContentType.JSON,
  ø.Send(MyType{Site: "site", Host: "host"}),
)

// Encode map to www-form-urlencoded
http.Join(
  // ...
  ø.ContentType.Form,
  ø.Send(map[string]string{
    "site": "site",
    "host": "host",
  })
)

// Send string, []byte or io.Reader. Just define the right Content-Type
http.Join(
  // ...
  ø.ContentType.Form,
  ø.Send([]byte{"site=site&host=host"}),
)
```

### Reader morphism



The implementation defines an abstraction of the protocol environments and lenses to focus inside it. In other words, the category represents the environment as an "invisible" side-effect of the composition.


calls combinator  using .  

A category is a concept that is defined in abstract terms of objects, arrows together with two functions composition ◦ and identity 𝒊𝒅. These functions shall be compliant with category laws

Associativity : (𝒇 ◦ 𝒈) ◦ 𝒉 = 𝒇 ◦ (𝒈 ◦ 𝒉)
Left identity : 𝒊𝒅 ◦ 𝒇 = 𝒇
Right identity : 𝒇 ◦ 𝒊𝒅 = 𝒇
The category leaves the definition of object, arrows, composition and identity to us, which gives a powerful abstraction!


