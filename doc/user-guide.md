<p align="center">
  <h3 align="center">ᵍ🆄🆁🅻</h3>
  <p align="center"><strong>User Guide</strong></p>
</p>

The combinators fit very well to express intent of communication behavior. It gives rich abstractions to hide the networking complexity and help us to compose a chain of network operations and represent them as pure computation, building new things from small reusable elements.

This is a "combinator" library for network I/O. Combinators open up an opportunity to depict computation problems in terms of fundamental elements like physics talks about universe in terms of particles. The only definite purpose of combinators are building blocks for composition of "atomic" functions into computational structures. It uses a powerful symbolic expressions of combinators to implement declarative language for testing suite development.

- [Background](#background)
  - [Combinators](#combinators)
  - [High-Order Functions](#high-order-functions)
  - [Lifecycle](#lifecycle)
- [Import library](#import-library)
- [Writer combinators](#writer-combinators)
  - [Method](#method)
  - [Target URI](#target-uri)
  - [Query Params](#query-params)
  - [Request Headers](#request-headers)
  - [Request Payload](#request-payload)
- [Reader combinators](#reader-combinators)
  - [Status Code](#status-code)
  - [Response Headers](#response-headers)
  - [Response Payload](#response-payload)
  - [Assert Payload](#assert-payload)
  - [Using Variables for Dynamic Behavior](#using-variables-for-dynamic-behavior)
- [Chain Networking I/O](#chain-networking-io)

---


## Background

### Combinators

Let’s formalize principles that help us to define our own abstraction applicable in functional programming through composition. The composition becomes a fundamental operation: the codomain of 𝒇 be the domain of 𝒈 so that the composite operation 𝒇 ◦ 𝒈 is defined. Our formalism uses `Arrow: IO ⟼ IO` as a key abstraction of networking combinators.

```go
// Arrow: IO ⟼ IO
type Arrow func(*Context) error
```

It is a pure function that takes an abstraction of the protocol context and applies morphism as an "invisible" side-effect of the composition.

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

Following the Input/Process/Output protocols paradigm, the two classes of combinators are defined:
* The first class is **writer** (emitter) morphism combinators, denoted by the symbol `ø` across this guide and example code. It focuses inside the protocol stack and reshapes requests. In the context of HTTP protocol, the writer morphism is used to declare HTTP method, destination URL, request headers and payload. 
* Second one is **reader** (matcher) morphism combinators, denoted by the symbol `ƒ`. It focuses on the side-effects of the protocol stack. The reader morphism is a pattern matcher, and is used to match response code, headers and response payload, etc. Its major property is “fail fast” with error if the received value does not match the expected pattern.

### High-order functions 

`Arrow` can be composed with another `Arrow` into new `Arrow` and so on. Only product "and-then" composition style is supported. It builds a strict product `Arrow: A ◦ B ◦ C ◦ ... ⟼ D`. The product type takes a protocol context and applies "morphism" sequentially unless some step fails. Use variadic function `http.Join`, `http.GET`, `http.POST`, `http.PUT` and so on to compose HTTP primitives:

```go
// Join composes HTTP arrows to high-order function
// (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
func http.Join(arrows ...http.Arrow) http.Arrow

//
var a: http.Arrow = /* ... */
var b: http.Arrow = /* ... */
var c: http.Arrow = /* ... */

d := http.Join(a, b, c)
```

Ease of the composition is one of major intent why syntax deviates from standard Golang HTTP interface. `http.Join` produces instances of higher order `http.Arrow` type, which is composable “promises” of HTTP I/O and so on. Essentially, the network I/O is just a set of `Arrow` functions. These rules of Arrow composition allow anyone to build a complex HTTP I/O scenario from a small reusable block.

### Lifecycle

```go
lazy := http.GET(/* ... */)
```

The instance of `Arrow` produced by one of `Join` functions does not hold a result of HTTP I/O. It only builds a composable "promise" ("lazy I/O") - a pure computation. The computation needs to be evaluated by applying it over the protocol context. The library provides a simple interface to create and customize the environment.  

```go
// HTTP protocol provides a default out-of-the-box environment.
// use the default environment with caution in production workload
cat := http.New() 

// apply the computation over the environment
err := cat.IO(context.TODO(), lazy)
```

Usage of the library for production workload requires a careful configuration of HTTP protocol timeouts, TLS policies, etc.

```go
cat := http.New(
  http.WithClient(/* ... */),
  http.InsecureTLS(),
  http.CookieJar(),
  http.LogRequest(),
  http.LogResponse(),
  http.LogPayload(),
)
```

## Import library

The combinator domain specific language consists of multiple packages, import them all into Golang module

```go
import (
  // context of http protocol stack
  "github.com/fogfish/gurl/http"

  // Writer (emitter) morphism combinators. It focuses inside the protocol stack
  // and reshapes requests. In the context of HTTP protocol, the writer morphism
  // is used to declare HTTP method, destination URL, request headers and payload.
  // single letter symbol (e.g. ø) makes the code less verbose
  ø "github.com/fogfish/gurl/http/send"

  // Reader (matcher) morphism combinators. It focuses on the side-effects of
  // the protocol stack. The reader morphism is a pattern matcher, and is used
  // to match response code, headers and response payload, etc. Its major
  // property is “fail fast” with error if the received value does not match
  // the expected pattern. 
  // single letter alias (e.g. ƒ) makes the code less verbose
  ƒ "github.com/fogfish/gurl/http/recv"
)
```

## Writer combinators

Writer (emitter) morphism combinators. It focuses inside the protocol stack and reshapes requests. In the context of HTTP protocol, the writer morphism is used to declare HTTP method, destination URL, request headers and payload.

### Method

Use `http.GET(/* ... */)` combinator to declare the verb of HTTP request. The language declares a combinator for most of HTTP verbs: `http.GET`, `http.HEAD`, `http.POST`, `http.PUT`, `http.DELETE` and `http.PATCH`.

```go
func SomeGetXxx() http.Arrow {
  return http.GET(/* ... */)
}

func SomePutXxx() http.Arrow {
  return http.PUT(/* ... */)
}
```

Use `ø.Method` combinator to declare other verbs

```go
func SomeXxx() http.Arrow {
  return http.Join(
    ø.Method("OPTIONS"),
    /* ... */)
}
```

### Target URI

Use `ø.URI(string)` combinator to specifies target URI for HTTP request. The combinator uses absolute URI to specify protocol, target host and path of the endpoint.

```go
func SomeXxx() http.Arrow {
  return http.GET(
    ø.URI("http://example.com"),
    /* ... */)
}
```

The `ø.URI` combinator is equivalent to `fmt.Sprintf`. It uses [percent encoding](https://golang.org/pkg/fmt/) to format and escape values.

```go
http.GET(
  ø.URI("http://example.com/%s", "foo bar"),
)

// All path segments are escaped by default, use ø.Authority or ø.Path
// types to disable escaping

// BAD, DOES NOT WORK
http.GET(
  ø.URI("%s/%s", "http://example.com", "foo/bar"),
)

// GOOD, IT WORKS
http.GET(
  ø.URI("%s/%s", ø.Authority("http://example.com"), ø.PATH("foo/bar")),
)
```

### Query Params

Use `ø.Params(any)` combinator to lifts the flat structure or individual values into query parameters of specified URI. 

```go
type MyParam struct {
  Site string `json:"site,omitempty"`
  Host string `json:"host,omitempty"`
}

func SomeXxx() http.Arrow {
  return http.GET(
    /* ... */
    ø.Params(MyParam{Site: "example.com", Host: "127.1"}),
    /* ... */
  )
}
```

Use `ø.Param` to declare individual query parameters, this combinator is suitable for simple queries, where definition of dedicated type seen as an overhead 

```go
func SomeXxx() http.Arrow {
  return http.GET(
    /* ... */
    ø.Param("site", "example.com"),
    ø.Param("host", "127.1"),
    /* ... */
  )
}
```

### Request Headers

Use `ø.Header[T any](string, T)` to declares headers and its values into HTTP requests. The [standard HTTP headers](https://en.wikipedia.org/wiki/List_of_HTTP_header_fields) are accomplished by a dedicated combinator making it type safe and easy to use e.g. `ø.ContentType.ApplicationJSON`.

```go
func SomeXxx() http.Arrow {
  return http.GET(
    /* ... */
    ø.Header("Client", "curl/7.64.1"),
    ø.Authorization.Set("Bearer eyJhbGciOiJIU...adQssw5c"),
    ø.ContentType.ApplicationJSON,
    ø.Accept.JSON,
    /* ... */
  )
}
```

### Request payload

Use `ø.Send` to transmits the payload to the destination URI. The combinator takes standard data types (e.g. maps, struct, etc) and encodes it to binary using Content-Type header as a hint. It fails if content type header is not defined or not supported by the library.

```go
type MyType struct {
  Site string `json:"site,omitempty"`
  Host string `json:"host,omitempty"`
}

func SomeSendJSON() http.Arrow {
  return http.GET(
    // ...
    ø.ContentType.JSON,
    ø.Send(MyType{Site: "example.com", Host: "127.1"}),
  )
}

func SomeSendForm() http.Arrow {
  return http.GET(
    // ...
    ø.ContentType.Form,
    ø.Send(map[string]string{
      "site": "example.com",
      "host": "127.1",
    })
  )
}

func SomeSendOctetStream() http.Arrow {
  return http.GET(
    // ...
    ø.ContentType.Form,
    ø.Send([]byte{"site=example.com&host=127.1"}),
  )
}
```

On top of the shown type, it also support a raw octet-stream payload presented after one of the following Golang types: `string`, `*strings.Reader`, `[]byte`, `*bytes.Buffer`, `*bytes.Reader`, `io.Reader` and any arbitrary `struct`.


## Reader combinators

Reader (matcher) morphism combinators. It focuses on the side-effects of the protocol stack. The reader morphism is a pattern matcher, and is used to match response code, headers and response payload, etc. Its major property is “fail fast” with error if the received value does not match the expected pattern.


### Status Code

Use `ƒ.Status.OK` checks the code in HTTP response and fails with error if the status code does not match the expected one. Status code is only mandatory reader combinator to be declared. The all well-known HTTP status codes are accomplished by a dedicated combinator making it type safe (e.g. `ƒ.Status` is constant with all known HTTP status codes as combinators).

```go
func SomeXxx() http.Arrow {
  return http.GET(
    // ...
    ƒ.Status.OK,
  )
}
```

Sometime a multiple HTTP status codes has to be accepted `ƒ.Code` arrow is variadic function that does it

```go
func SomeXxx() http.Arrow {
  return http.GET(
    // ...
    ƒ.Code(http.StatusOK, http.StatusCreated, http.StatusAccepted),
  )
}
```

### Response Headers

Use `ƒ.Header` combinator to matches the presence of HTTP header and its value in the response. The matching fails if the response is missing the header or its value does not correspond to the expected one. The [standard HTTP headers](https://en.wikipedia.org/wiki/List_of_HTTP_header_fields) are accomplished by a dedicated combinator making it type safe and easy to use e.g. `ƒ.ContentType.ApplicationJSON`.

```go
func SomeXxx() http.Arrow {
  return http.GET(
    // ...
    ƒ.Header("Content-Type", "application/json"),
    ƒ.Authorization.Is("Bearer eyJhbGciOiJIU...adQssw5c"),
    ƒ.ContentType.JSON,
    ƒ.Server.Any,
  )
}
```

The combinator support "lifting" of header value into the variable for the further usage in the application.

```go
func SomeXxx() http.Arrow {
  var (
    date time.Time
    mime string
    some string
  )

  return http.GET(
    // ...
    ƒ.Date.To(&date),
    ƒ.ContentType.To(&mime),
    ƒ.Header("X-Some", &some),
  )
}
```

### Response Payload

Use `ƒ.Body` consumes payload from HTTP requests and decodes the value into the type associated with the lens using Content-Type header as a hint. It fails if the body cannot be consumed.


```go
type MyType struct {
  Site string `json:"site,omitempty"`
  Host string `json:"host,omitempty"`
}

func SomeXxx() http.Arrow {
  var data MyType

  return http.GET(
    // ...
    ƒ.Body(&data), // Note: pointer to data structure is required
  )
}
```

So far, utility support auto decoding of the following `Content-Types` into structs
* `application/json`
* `application/x-www-form-urlencoded`
* `image/*`

The library automatically decodes images into `image.Image` data type.   

```go
import (
  _ "image/jpeg"
)

func SomeXxx() http.Arrow {
  var data image.Image

  return http.GET(
    // ...
    ƒ.Body(&data), // Note: pointer to data structure is required
  )
}
```


For all other cases, there is `ƒ.Bytes` combinator that receives raw binaries.  

```go
func SomeXxx() http.Arrow {
  var data []byte

  return http.GET(
    // ...
    ƒ.Bytes(&data), // Note: pointer to buffer is required
  )
}
```

### Assert Payload

Combinators is not only about pure networking but also supports assertion of responses. Assert combinator aborts the evaluation of computation if expected value do not match the response. There are three type of asserts: type safe `ƒ.Expect`, loosely typed `ƒ.Match` and customer combinator.

**Type safe**: Use `ƒ.Expect` to define expected value as Golang struct. The combinator fails if received value do not strictly equals to expected one.

```go
func TestXxx() http.Arrow {
  return http.GET(
    // ...
    ƒ.Expect(MyType{Site: "example.com", Host: "127.1"}),
  )
}
```

**Loosely typed**: Use `ƒ.Match` to define expected value as string pattern. In the contrast to type safe combinator, the combinator takes a valid JSON object as string.
It matches only defined values and supports wildcard matching. For example: 

```go
// matches anything
`"_"`

// matches any object with key "site"
`{"site": "_"}`

// matches array of length 1 
`["_"]`

// matches any object with key "site" equal to  "example.com"
`{"site": "example.com"}`

// matches any array of length 2 with first object having the key 
`[{"site": "_"}, "_"]`

// matches nested objects
`{"site": {"host": "_"}}`

// and so on ...

func TestXxx() http.Arrow {
  return http.GET(
    // ...
    ƒ.Match(`{"site": "example.com", "host": "127.1"}`),
  )
}
```

**Custom combinator**: The `type Arrow func(*http.Context) error` is "open" interface to combine assert logic with networking I/O. These functions act as lense -- focuses inside the structure, fetching values and asserts them. These helpers can do anything with the computation including its termination: 

```go
type MyType struct {
  Site string `json:"site,omitempty"`
  Host string `json:"host,omitempty"`
}

// a type receiver to assert the value
func (t *MyType) CheckValue(*http.Context) error {
  if t.Host != "127.1" {
    return fmt.Errorf("...")
  }

  return nil
}

func TestXxx() http.Arrow {
  var data MyType

  return http.GET(
    // ...
    ƒ.Recv(data),
    t.CheckValue,
  )
}
```

### Using Variables for Dynamic Behavior

A pure functional style of development does not have variables or assignment statements. The program is defined by applying type constructors, constants and functions. However, this principle does not closely match current architectures. Programs are implemented using variables such as memory lookups and updates. Any complex real-life networking I/O is not an exception, it requires a global operational state. So far, all examples have used constants and literals but ᵍ🆄🆁🅻 combinators also support dynamic behavior of I/O parameters using pointers to variables.  

```go
type MyClient http.Stack

func (cli MyClient) Request(host, token string, req T) (*T, error) {
  return http.IO[T](cat.WithContext(context.Background()),
    http.GET(
      ø.GET.URL("https://%s", host),
      ø.Authorization.Set(token),
      ø.Send(req),
      ƒ.Status.OK,
    )
  )
}
```


## Chain networking I/O

Ease of the composition is one of major intent why combinators has been defined. `http.Join` produces instances of higher order combinator, which is composable into higher order constructs. Let's consider an example where sequence of requests needs to be executed one after another (e.g. interaction with GitHub API):   

```go
// 1. declare a product type that depict the context of networking I/O. 
type State struct {
  Token AccessToken
  User  User
  Org   Org
}

// 2. declare collection of independent requests, each either reads or writes
// the context
func (s *State) FetchAccessToken() http.Arrow {
  return http.GET(
    // ...
    ƒ.Recv(&s.Token),              // writes access token to context
  )
}

func (s *State) FetchUser() error {
  return http.POST(
    ø.URI(/* ... */),
    ø.Authorization.Set(&s.Token), // reads access token from context
    // ...
    ƒ.Recv(&hof.User),               // writes user object to context
  )
}

func (s *State) FetchContribution() error {
  return http.POST(
    ø.URI(&s.User.Repos),          // reads user object from context
    ø.Authorization.Set(&s.Token), // reads access token from context
    // ...
    ƒ.Recv(&s.Org),                // writes user's contribution to context
  )
}

// 3. Composed sequence of requests into the chained sequence
func HighOrderFunction() (*State, http.Arrow) {
	var state State

	//
	// HoF combines HTTP requests to
	//  * https://httpbin.org/uuid
	//  * https://httpbin.org/post
	//
	// results of HTTP I/O is persisted in the internal state
	return &state, http.Join(
    state.FetchAccessToken(),
    state.FetchUser(),
    state.FetchContribution(),
	)
}
```

Hopefully you find it useful, and the docs easy to follow.

Feel free to [create an issue](https://github.com/fogfish/gurl/issues) if you find something that's not clear.

See [example](../examples) for details about the compositions.
