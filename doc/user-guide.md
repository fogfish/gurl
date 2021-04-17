<p align="center">
  <h3 align="center">ᵍ🆄🆁🅻</h3>
  <p align="center"><strong>User Guide</strong></p>
</p>

- [Overview](#overview)
- [Compose HoF](#compose-hof)
- [Life-cycle](#life-cycle)
- [Import library](#import-library)
- [Arrow types](#arrow-types)
  - [Writer morphism](#writer-morphism)
    - [Method and URL](#method-and-url)
    - [Query Params](#query-params)
    - [Request Headers](#request-headers)
    - [Request Payload](#request-payload)
  - [Reader morphism](#reader-morphism)
    - [Status Code](#status-code)
    - [Response Headers](#response-headers)
    - [Response Payload](#response-payload)
  - [Using Variables for Dynamic Behavior](#using-variables-for-dynamic-behavior)
- [Assert Protocol Payload](#assert-protocol-payload)
- [Chain Network I/O](#chain-network-io)

---


## Overview

ᵍ🆄🆁🅻 is a "combinator" library for network I/O. Combinators open up an opportunity to depict computation problems in terms of fundamental elements like physics talks about universe in terms of particles. The only definite purpose of combinators are building blocks for composition of "atomic" functions into computational structures. ᵍ🆄🆁🅻 combinators provide a powerful symbolic expressions in networking domain.

Standard Golang packages implements a low-level HTTP interface, which requires knowledge about protocol itself, understanding of Golang implementation aspects, and a bit of boilerplate coding. It also misses standardized chaining (composition) of individual requests. ᵍ🆄🆁🅻 inherits an ability of pure functional languages to express communication behavior by hiding the networking complexity using combinators. The composition becomes a fundamental operation in the library: the codomain of `𝒇` be the domain of `𝒈` so that the composite operation `𝒇 ◦ 𝒈` is defined.

The library uses `Arrow` as a key abstraction of combinators. It is a *pure function* that takes an abstraction of the protocol environments, so called IO category and applies morphism as an "invisible" side-effect of the composition.

```go
/*

Arrow: IO ⟼ IO
*/
type Arrow func(*gurl.IOCat) *gurl.IOCat
```

There are two classes of arrows. The first class is a writer morphism that focuses inside and reshapes HTTP protocol requests. The writer morphism is used to declare HTTP method, destination URL, request headers and payload. Second one is a reader morphism that focuses on the side-effect of HTTP protocol. The reader morphism is a pattern matcher, and is used to match HTTP response code, headers and response payload.

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


## Compose HoF 

`Arrow` can be composed with another `Arrow` into new `Arrow` and so on. The library supports only "and-then" style. It builds a strict product Arrow: `A × B × C × ... ⟼ D`. The product type takes a protocol environment and applies "morphism" sequentially unless some step fails. Use variadic function `http.Join` to compose HTTP primitives:

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

Ease of the composition is one of major intent why ᵍ🆄🆁🅻 library has deviated from standard Golang HTTP interface. `http.Join` produces instances of higher order `gurl.Arrow` type, which is composable “promises” of HTTP I/O and so on. Essentially, the network I/O is just a set of `Arrow` functions.

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

These rules of `Arrow` composition allow anyone to build a complex HTTP I/O scenario from a small reusable block.



## Life-cycle

```go
lazy := gurl.Join(/* ... */)
```

The instance of `Arrow` produced by one of `Join` functions does not hold a result of HTTP I/O. It only builds a composable "promise" ("lazy I/O") - a pure computation. The computation needs to be evaluated by applying it over the protocol environment. The library provides a simple interface to create and customize the environment.  

```go
// HTTP protocol provides a default out-of-the-box environment.
// use the default environment with caution in production workload
env := http.DefaultIO() 

// apply the computation over the environment
env = lazy(env)

// handle either networking error or failure of expectation.
// there are not need for error handling after each operation
// in contrast with classical networking I/O. The combinator
// is smart enough to terminate execution.
if env.Fail != nil {
  // ...
}

// the environment holds failed state until it is recovered 
if err := env.Recover(); err != nil {
  // ...
}
```

Usage of the library for production workload requires a careful configuration of HTTP protocol timeouts, TLS policies, etc. Another aspect is thread safeness. The protocol environment is not thread safe. Each golang routine shall create one. Re-use of `http.Client` pointer across environments reduces resource consumption.  

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

Method and URL are only mandatory writer morphism in I/O declaration. Use `type Method string` to declare the verb of HTTP request. It's received method `URL` allows to specify a destination endpoint.

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
  // this does not work
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

Use `type Header string` to declare headers and its values. Each request might contain declaration of multiple headers.

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

The `func Send(data interface{}) http.Arrow` transmits the payload to the destination URL. The function takes Go data types (e.g. maps, struct, etc) and encodes it to binary using `Content-Type` header as a hint. The function fails if content type is not defined or not supported by the library.

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

Reader morphism focuses on the side-effect of HTTP protocol. It does a pattern matching of HTTP response code, header values and response payload.

#### Status Code

Status code validation is only mandatory reader morphism in I/O declaration. The status code "arrow" checks the code in HTTP response and fails with error if the status code does not match the expected one. The library defines a `type StatusCode int` and constants (e.g. `Status.OK`) for all known HTTP status codes.

```go
http.Join(
  // ...
  ƒ.Status.OK,
)

// Sometime a multiple HTTP status codes has to be accepted
// `ƒ.Code` arrow is variadic function that does it
http.Join(
  // ...
  ƒ.Code(http.StatusOK, http.StatusCreated, http.StatusAccepted),
)
```

#### Response Headers

Use `type Header string` to pattern match presence of HTTP header and its value in the response. The matching fails if the response is missing the header or its value do not equal.

```go
http.Join(
  // ...
  ƒ.Header("Content-Type").Is("application/json"),
)

// The library implements a syntax sugar for mostly used HTTP headers
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Response_fields
http.Join(
  // ...
  ƒ.Authorization.Is("Bearer eyJhbGciOiJIU...adQssw5c"),
)

// The library implements a syntax sugar for content negotiation headers
http.Join(
  // ...
  ƒ.ContentType.JSON,
)

// Any arrow is a syntax sugar of Header("Content-Type").Is("*")
http.Join(
  // ...
  ƒ.Server.Any,
  ƒ.ContentType.Any,
)
```

#### Response Payload

The `func Recv(out interface{}) http.Arrow` decodes the response payload to Golang native data structure using Content-Type header as a hint.

```go
type MyType struct {
  Site string `json:"site,omitempty"`
  Host string `json:"host,omitempty"`
}

var data MyType
http.Join(
  // ...
  ƒ.Recv(&data), // Note: pointer to data structure is required
)
```

The library supports auto decoding of
* `application/json`
* `application/x-www-form-urlencoded`

It also receives raw binaries in case data type is not supported.  

```go
var data []byte
http.Join(
  // ...
  ƒ.Bytes(&data), // Note: pointer to data buffer is required
)
```

### Using Variables for Dynamic Behavior

A pure functional style of development does not have variables or assignment statements. The program is defined by applying type constructors, constants and functions. However, this principle does not closely match current architectures. Programs are implemented using variables such as memory lookups and updates. Any complex real-life networking I/O is not an exception, it requires a global operational state. So far, all examples have used constants and literals but ᵍ🆄🆁🅻 combinators also support dynamic behavior of I/O parameters using pointers to variables.  

```go
func dynamic(host, token, lang *string, req, data *T) gurl.Arrow {
  return http.Join(
    //
    ø.GET.URL("https://%s", host),
    ø.Authorization.Val(token),
    ø.Send(req),
    //
    ƒ.Status.OK,
    ƒ.ContentLanguage.String(lang),
    ƒ.Recv(data),
  )
}
```

## Assert Protocol Payload

ᵍ🆄🆁🅻 library is not only about networking I/O. It also allows to assert the response. It defines a few helper functions that combine assert logic with I/O chain. These functions act as lense that are focused inside the structure, fetching values and asserts them. These helpers abort the evaluation of “program” if expectations do not match actual response. The `func FMap(f func() error) Arrow` lifts any function/closure to composable `Arrow`, allowing to implement assert procedure.  

```go
type T struct {
  ID int
}

// a type receiver to assert the value
func (t *T) CheckValue() error {
  if t.ID == 0 {
    return fmt.Errorf("...")
  }

  return nil
}

func (t *T) SomeIO() gurl.Arrow {
  return http.Join(
    // ...
    ƒ.Recv(t),
  // compose the assertion into I/O chain   
  ).Then(gurl.FMap(t.CheckValue))
}
```


## Chain Network I/O

Ease of the composition is one of major feature in ᵍ🆄🆁🅻 library. It allows chain multiple independent HTTP I/O to the high order computation.

```go
// declare a product type to depict IO context
type HoF struct {
  Token AccessToken
  User  User
  Org   Org
}

// Declare set of independent HTTP I/O.
// Each operation either reads or writes the context
func (hof *HoF) FetchAccessToken() gurl.Arrow {
  return http.Join(
    // ...
    ƒ.Recv(&hof.Token),
  )
}

func (hof *HoF) FetchUser() gurl.Arrow {
  return http.Join(
    ø.POST.URL(/* ... */),
    ø.Authorization().Val(&hof.Token),
    // ...
    ƒ.Recv(&hof.user),
  )
}

func (hof *HoF) FetchContribution() gurl.Arrow {
  return http.Join(
    ø.POST(/* ... */),
    ø.Authorization().Val(&hof.Token),
    // ...
    ƒ.Recv(&hof.Org),
  )
}

// Combine HTTP I/O to sequential chain of execution 
api := &HoF{}
http.Join(
  api.FetchAccessToken(),
  api.FetchUser(),
  api.FetchContribution(),
)
```

See [example](../example) for details about the compositions.
