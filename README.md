# HTTP High Order Component

A class of High Order Component which can do http requests with few
interesting property such as composition and laziness.

[![Documentation](https://godoc.org/github.com/fogfish/gurl?status.svg)](http://godoc.org/github.com/fogfish/gurl)
[![Build Status](https://secure.travis-ci.org/fogfish/gurl.svg?branch=master)](http://travis-ci.org/fogfish/gurl)
[![Git Hub](https://img.shields.io/github/last-commit/fogfish/gurl.svg)](http://travis-ci.org/fogfish/gurl)
[![Coverage Status](https://coveralls.io/repos/github/fogfish/gurl/badge.svg?branch=master)](https://coveralls.io/github/fogfish/gurl?branch=master)



The library implements rough and naive Haskell's equivalent of 
do-notation, so called monadic binding form. This construction decorates
http i/o pipeline(s) with "programmable commas".

## Inspiration

Microservices have become a design style to evolve system architecture
in parallel, implement stable and consistent interfaces. An expressive
language is required to design the variety of network communication use-cases.
A pure functional languages fits very well to express communication behavior.
The language gives a rich techniques to hide the networking complexity using
monads as abstraction. The IO-monads helps us to compose a chain of network
operations and represent them as pure computation, build a new things from
small reusable elements. The library is implemented after Erlang's [m_http](https://github.com/fogfish/m_http)

The library attempts to adapts a human-friendly syntax of HTTP request/response
logging/definition used by curl with Behavior as a Code paradigm. It tries to
connect cause-and-effect (Given/When/Then) with the networking (Input/Process/Output).

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

This semantic provides an intuitive approach to specify HTTP requests/responses.
Adoption of this syntax as Go native code provides a rich capability to network
programming.

## Key features

* cause-and-effect abstraction of HTTP request/response, naive do-notation
* high-order composition of individual HTTP requests to complex networking computations
* human-friendly, Erlang native and declarative syntax to depict HTTP operations
* implements a declarative approach for testing of RESTful interfaces
* automatically encodes/decodes Go native HTTP payload using Content-Type hints
* supports generic transformation to algebraic data types
* simplify error handling with naive Either implementation

## Getting started

The latest version of the library is available at its `master` branch. All development, including new features and bug fixes, take place on the `master` branch using forking and pull requests as described in contribution guidelines.

Import the library in your code

```go
import (
  "github.com/fogfish/gurl"
)
```

See the [documentation](http://godoc.org/github.com/fogfish/gurl)

## Basics

The following code snippet demonstrates a typical usage scenario. The code
uses dot (.) to compose http primitives and evaluate a "program".

```go
import "github.com/fogfish/gurl"

type Payload struct {
  Origin string `json:"origin"`
  Url    string `json:"url"`
}

var data Payload
io := gurl.NewIO().
  GET("http://httpbin.org/get").
  With("Accept", "application/json").
  Code(200).
  Head("Content-Type", "application/json").
  Recv(&data)

if io.Fail != nil {
  // error handling
}
```

The evaluation of "program" fails if either networking fails or expectations
do not match actual response. There are no needs to check error code after
each operation. The composition is smart enough to terminate "program" execution.

## Composition

The composition of multiple HTTP I/O is an essential part of the library.
The composition is handled in context of IO category. For example,
RESTfull API primitives declared as function, each deals with gurl.IO.

```go
func hof() {
  io := gurl.NewIO()
  token := githubAccessToken(io)
  user := githubUserProfile(io, token)
  orgs := githubUserContribution(io, token)
}

func githubAccessToken(io *gurl.IO) (token AccessToken) {
  io.URL("POST", "...").
    // ...
    Recv(&token)
  return token
}

func githubUserProfile(io *gurl.IO, token AccessToken) (user User) {
  io.URL("POST", "...")
    // ...
    Recv(&user)
  return
}

func githubUserContribution(io *gurl.IO, token string, token AccessToken) (orgs []Org) {
  io.URL("POST", "...")
    // ...
    Recv(&orgs)
  return
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