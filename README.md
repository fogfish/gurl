<p align="center">
  <h1 align="center">ᵍ🆄🆁🅻</h1>
  <p align="center"><strong>Network I/O combinator library for Golang</strong></p>

  <p align="center">
    <!-- Documentation -->
    <a href="http://godoc.org/github.com/fogfish/gurl">
      <img src="https://godoc.org/github.com/fogfish/gurl?status.svg" />
    </a>
    <!-- Build Status  -->
    <a href="https://github.com/fogfish/gurl/actions/">
      <img src="https://github.com/fogfish/gurl/workflows/Go/badge.svg" />
    </a>
    <!-- GitHub -->
    <a href="http://github.com/fogfish/gurl">
      <img src="https://img.shields.io/github/last-commit/fogfish/gurl.svg" />
    </a>
    <!-- Coverage -->
    <a href="https://coveralls.io/github/fogfish/gurl?branch=master">
      <img src="https://coveralls.io/repos/github/fogfish/gurl/badge.svg?branch=master" />
    </a>
    <!-- Go Card -->
    <a href="https://goreportcard.com/report/github.com/fogfish/gurl">
      <img src="https://goreportcard.com/badge/github.com/fogfish/gurl" />
    </a>
    <!-- Maintainability -->
    <a href="https://codeclimate.com/github/fogfish/gurl/maintainability">
      <img src="https://api.codeclimate.com/v1/badges/b9ff76a1f641ce98cd26/maintainability" />
    </a>
  </p>
</p>


# HTTP High Order Component

A class of High Order Component which can do http requests with few interesting property such as composition and laziness.

[![Documentation](https://godoc.org/github.com/fogfish/gurl?status.svg)](http://godoc.org/github.com/fogfish/gurl)
[![Build Status](https://github.com/fogfish/gurl/workflows/Go/badge.svg)](https://github.com/fogfish/gurl/actions/)
[![Git Hub](https://img.shields.io/github/last-commit/fogfish/gurl.svg)](http://travis-ci.org/fogfish/gurl)
[![Coverage Status](https://coveralls.io/repos/github/fogfish/gurl/badge.svg?branch=master)](https://coveralls.io/github/fogfish/gurl?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/fogfish/gurl)](https://goreportcard.com/report/github.com/fogfish/gurl)
[![Maintainability](https://api.codeclimate.com/v1/badges/b9ff76a1f641ce98cd26/maintainability)](https://codeclimate.com/github/fogfish/gurl/maintainability)


The library implements rough and naive Haskell's equivalent of do-notation, so called monadic binding form. This construction decorates http i/o pipeline(s) with "programmable commas".


## Inspiration

Microservices have become a design style to evolve system architecture in parallel, implement stable and consistent interfaces. An expressive language is required to design the variety of network communication use-cases. A pure functional languages fits very well to express communication behavior. The language gives a rich techniques to hide the networking complexity using monads as abstraction. The IO-monads helps us to compose a chain of network operations and represent them as pure computation, build a new things from small reusable elements. The library is implemented after Erlang's [m_http](https://github.com/fogfish/m_http)

The library attempts to adapts a human-friendly syntax of HTTP request/response logging/definition used by curl with Behavior as a Code paradigm. It tries to connect cause-and-effect (Given/When/Then) with the networking (Input/Process/Output).

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

This semantic provides an intuitive approach to specify HTTP requests/responses. Adoption of this syntax as Go native code provides a rich capability to network programming.


## Key features

* cause-and-effect abstraction of HTTP request/response, naive do-notation
* high-order composition of individual HTTP requests to complex networking computations
* human-friendly, Go native and declarative syntax to depict HTTP operations
* implements a declarative approach for testing of RESTful interfaces
* automatically encodes/decodes Go native HTTP payload using Content-Type hints
* supports generic transformation to algebraic data types
* simplify error handling with naive Either implementation

## Getting started

The library requires **Go 1.13** or later 

The latest version of the library is available at its `master` branch. All development, including new features and bug fixes, take place on the `master` branch using forking and pull requests as described in contribution guidelines.

Import the library in your code

```go
import (
  // core types
  "github.com/fogfish/gurl"

  // support for http protocol
  "github.com/fogfish/gurl/http"

  // module ø (gurl/http/send) - writer morphism is used to declare HTTP method,
  // destination URL, request headers and payload.
  ø "github.com/fogfish/gurl/http/send"

  // module ƒ (gurl/http/recv) - reader morphism is a pattern matcher, is used
  // to match HTTP response code, headers and response payload.
  ƒ "github.com/fogfish/gurl/http/recv"
)
```

See the [documentation](http://godoc.org/github.com/fogfish/gurl)


### IO Category

Standard Golang packages implements low-level HTTP interface, which requires knowledge about protocol itself, understanding of Golang implementation aspects, a bit of boilerplate coding. It also missing standardized chaining (composition) of individual requests.

`gurl` library inherits an ability of pure functional languages to express communication behavior by hiding the networking complexity using category pattern (aka "do"-notation). This pattern helps us to compose a chain of network operations and represent them as pure computation, build a new things from small reusable elements. This library uses the "do"-notation, so called monadic binding form. It is well know in functional programming languages such as Haskell and Scala. The networking becomes a collection of composed "do"-notation in context of a state monad.

A composition of HTTP primitives within the category are written with the following syntax.

```go
  gurl.Join(arrows ...Arrow) Arrow
```

Here, each arrow is a morphism applied to HTTP protocol. The implementation defines an abstraction of the protocol environments and lenses to focus inside it. In other words, the category represents the environment as an "invisible" side-effect of the composition.

The example definition of HTTP I/O within "do"-notation becomes

```go
  gurl.Join(
    http.Join(
      ø...,
      ø...,

      ƒ...,
      ƒ...,
    ),
  )
```

Symbol `ø` (option + o) is an convenient alias to module gurl/http/send, which defines writer morphism that focuses inside and reshapes HTTP protocol request. The writer morphism is used to declare HTTP method, destination URL, request headers and payload.

Symbol `ƒ` (option + f) is an convenient alias to module gurl/http/recv, which defines reader morphism that focuses into side-effect, HTTP protocol response. The reader morphism is a pattern matcher, is used to match HTTP response code, headers and response payload. It helps us to declare our expectations on the response. The evaluation of "program" fails if expectations do not match actual response.

`gurl.Join(arrows ...Arrow) Arrow` and its composition implements lazy I/O. It only returns a "promise", you have to evaluate it in the context of IO instance.

```go
  cat := gurl.IO()
  req := gurl.Join( ... )
  req(cat)
```

Let's look on step-by-step usage of the category.

**Method and URL** are mandatory. It has to be a first element in the construction.

```go
  http.Join(
    ø.GET("http://example.com"),
    ...
  )
```

Definition of **request headers** is an optional. You can list as many headers as needed. Either using string literals or variables. Some frequently used headers implements aliases (e.g. `ø.ContentJSON()`, ...)

```go
  http.Join(
    ...
    ø.Header("Accept").Is("application/json"),
    ø.Header("Authorization").Val(&token),
    ...
  )
```

The **request payload** is also an optional. You can also use native Golang data types as egress payload. The library implicitly encodes input structures to binary using Content-Type as a hint.

```go
  http.Join(
    ...
    ø.Send(MyType{Hello: "World"}),
    ...
  )
```

The declaration of expected response is always starts with mandatory HTTP **status code**. The execution fails if peer responds with other than specified value.

```go
  http.Join(
    ...
    ƒ.Code(http.StatusOK),
    ...
  )
```

It is possible to match presence of header in the response, match its entire content or lift the header value to a variable. The execution fails if HTTP response do not match the expectation.

```go
  http.Join(
    ...
    ƒ.Header("Content-Type").Is("application/json"),
    ...
  )
```

The library is able to **decode payload** into Golang native data structure using Content-Type header as a hint.

```go
  var data MyType
  http.Join(
    ...
    ƒ.Recv(&data)
    ...
  )
```

Please note, the library implements lenses to inline assert of decoded content. See the documentation of gurl/cat module.



### Basics

The following code snippet demonstrates a typical usage scenario. See runnable [example](example/request/main.go).

```go
import (
  "github.com/fogfish/gurl"
  "github.com/fogfish/gurl/http"
  ø "github.com/fogfish/gurl/http/send"
  ƒ "github.com/fogfish/gurl/http/recv"
)

// You can declare any types and use them as part of networking I/O.
type Payload struct {
  Origin string `json:"origin"`
  Url    string `json:"url"`
}

// the variable holds results of network I/O
var data Payload
var reqf := http.Join(
  // declares HTTP method and destination URL
  ø.GET("http://httpbin.org/get"),
  // HTTP content negotiation, declares acceptable types
  ø.Accept("application/json"),

  // requires HTTP Status Code to be 200 OK
  ƒ.Code(gurl.StatusOK),
  // requites HTTP Header to be Content-Type: application/json
  ƒ.Served("application/json"),
  // unmarshal JSON to the variable
  ƒ.Recv(&data),
)

// Note: http do not hold yet, a results of HTTP I/O
//       it is just a composable "promise", you have to
//       evaluate a side-effect of HTTP "computation"
if reqf(http.DefaultIO()).Fail != nil {
  // error handling
}
```

The evaluation of "program" fails if either networking fails or expectations do not match actual response. There are no needs to check error code after each operation. The composition is smart enough to terminate "program" execution.

## Composition

The composition of multiple HTTP I/O is an essential part of the library. The composition is handled in context of IO category. For example, RESTfull API primitives declared as arrow functions, each deals with `gurl.IOCat`. See runnable [high-order function example](example/hof/main.go) and [recursion](example/loop/main.go).


```go
import (
  "github.com/fogfish/gurl"
  "github.com/fogfish/gurl/http"
  ø "github.com/fogfish/gurl/http/send"
  ƒ "github.com/fogfish/gurl/http/recv"
)

func HoF() {
  // Internal state of HoF function
  var (
    token AccessToken
    user  User
    org   Org
  )

  // HoF combines multiple HTTP I/O to chain of execution 
  reqf := gurl.Join(
    AccessToken(&token),
    UserProfile(&token, &user),
    UserContribution(&token, &org)
  )

  if reqf(http.DefaultIO()).Fail != nil {
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
    ø.POST(/* ... */),
    ø.Authorization().Val(token.Bearer),
    // ...
    ƒ.Recv(user),
  )
}

func UserContribution(token *AccessToken, org *Org) {
  return http.Join(
    ø.POST(/* ... */),
    ø.Authorization().Val(token.Bearer),
    // ...
    ƒ.Recv(org),
  )
}
```

## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


The build and testing process requires [Go](https://golang.org) version 1.13 or later.

**Build** and **run** service in your development console. The following command boots Erlang virtual machine and opens Erlang shell.

```bash
git clone https://github.com/fogfish/gurl
cd gurl
go test -cover
```

### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

>
> Short (50 chars or less) summary of changes
>
> More detailed explanatory text, if necessary. Wrap it to about 72 characters or so. In some contexts, the first line is treated as the subject of an email and the rest of the text as the body. The blank line separating the summary from the body is critical (unless you omit the body entirely); tools like rebase can get confused if you run the two together.
> 
> Further paragraphs come after blank lines.
> 
> Bullet points are okay, too
> 
> Typically a hyphen or asterisk is used for the bullet, preceded by a single space, with blank lines in between, but conventions vary here
>
>

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/fogfish/gurl/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 

* **Specify** the configuration of your environment. Include which operating system you use and the versions of runtime environments. 

* **Attach** logs, screenshots and exceptions, in possible.

* **Reveal** the steps you took to reproduce the problem, include code snippet or links to your project.


## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/gurl.svg?style=for-the-badge)](LICENSE)
