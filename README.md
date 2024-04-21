<p align="center">
  <h3 align="center">ᵍ🆄🆁🅻</h3>
  <p align="center"><strong>Combinator library for network I/O</strong></p>

  <p align="center">
    <!-- Documentation -->
    <a href="http://godoc.org/github.com/fogfish/gurl">
      <img src="https://godoc.org/github.com/fogfish/gurl?status.svg" />
    </a>
    <!-- Build Status  -->
    <a href="https://github.com/fogfish/gurl/actions/">
      <img src="https://github.com/fogfish/gurl/workflows/test/badge.svg" />
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
  </p>
</p>

---

The library implements a **pure functional style** to express communication behavior by hiding the networking complexity using combinators. This construction decorates http i/o pipeline(s) with "programmable commas", allowing to make http requests with few interesting properties such as composition and laziness.

[User Guide](./doc/user-guide.md) |
[Playground](https://goplay.tools/snippet/RLxmdLZ49SC) |
[Examples](./example/) |
[API Specification](http://godoc.org/github.com/fogfish/gurl)

## Inspiration

Microservices have become a design style to evolve system architecture in parallel, implement stable and consistent interfaces. An expressive language is required to design the variety of network communication use-cases. Pure functional languages fit very well to express intent of communication behavior. These languages give rich abstractions to hide the networking complexity and help us to compose a chain of network operations and represent them as pure computation, building new things from small reusable elements. This library is implemented after Erlang's [m_http](https://github.com/fogfish/m_http)

The library attempts to adapt a human-friendly logging syntax of HTTP I/O used by curl and Behavior as a Code paradigm, which connects cause-and-effect (Given/When/Then) with the networking (Input/Process/Output).

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

Given semantic provides an intuitive approach to specify HTTP requests and expected responses. Adoption of this syntax as Go native code provides a rich capabilities for network programming.


## Key features

Standard Golang packages implement a low-level HTTP interface, which requires knowledge about the protocol itself, understanding of Golang implementation aspects, and a bit of boilerplate coding. It also misses standardized chaining (composition) of individual requests. ᵍ🆄🆁🅻 inherits an ability of pure functional languages to express communication behavior by hiding the networking complexity using combinators. Combinators make a chain of network operations as a pure computation. 

* cause-and-effect abstraction of HTTP I/O using Golang naive do-notation
* composition of individual HTTP requests to complex networking computations
* human-friendly, Go native and declarative syntax to depict HTTP operations
* implements a declarative approach for testing of RESTful interfaces
* automatically encodes/decodes Golang native HTTP payload using Content-Type hints
* supports generic transformation to algebraic data types

## Getting started

The library requires **Go 1.18** or later 

The latest version of the library is available at its `main` branch. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.

Use `go get` to retrieve the library and add it as dependency to your application.

```bash
go get -u github.com/fogfish/gurl
```

### Quick Example

The following code snippet demonstrates a typical usage scenario. See runnable [http request example](examples/http-request/main.go).

```go
import (
  "context"

  "github.com/fogfish/gurl/v2/http"
  ø "github.com/fogfish/gurl/v2/http/send"
  ƒ "github.com/fogfish/gurl/v2/http/recv"
)

// Declare the type, used for networking I/O.
type Payload struct {
  Origin string `json:"origin"`
  Url    string `json:"url"`
}

// Define the variable holds results of network I/O
var data Payload

// Declare HTTP I/O specification
lazy := http.GET(
  // specify HTTP request
  ø.GET.URL("http://httpbin.org/get"),
  ø.Accept.JSON,

  // assert HTTP response and "recv" JSON to the variable
  ƒ.Status.OK,
  ƒ.ContentType.JSON,
  ƒ.Recv(&data),
)

// instance of HTTP stack
stack := http.New()

// evaluate HTTP I/O specification
err := stack.IO(context.Background(), lazy)
```

## Next steps

* Study [User Guide](doc/user-guide.md) if defines library concepts and guides about api usage;
* Use [examples](examples) as a reference for further development.

## Extensions

The library supplies extensions
- [x/awsapi](x/awsapi/) enables AWS Signature V4 for HTTP I/O. Allows to use AWS API Gateway with IAM authentication.
- [x/xhtml](x/xhtml/) enables fetching and parsing xHTML content.

## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


The build and testing process requires [Go](https://golang.org) version 1.18 or later.

**Build** and **test** the library in your development console.

```bash
git clone https://github.com/fogfish/gurl
cd gurl
go test ./...
```

### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two questions what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

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

* **Attach** logs, screenshots and exceptions, if possible.

* **Reveal** the steps you took to reproduce the problem, include code snippet or links to your project.


## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/gurl.svg?style=for-the-badge)](LICENSE)
