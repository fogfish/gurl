package http

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

//
// The file implements the context for Arrow
//

// Context of HTTP I/O
type Context struct {
	context.Context

	Method   string
	Request  *http.Request
	Response *http.Response
	stack    *Protocol
}

// IO executes protocol operations
func (ctx *Context) IO(arrows ...Arrow) error {
	for _, f := range arrows {
		if err := f(ctx); err != nil {
			return err
		}
	}

	if ctx.Response != nil {
		// Note: due to Golang HTTP pool implementation we need to consume and
		//       discard body. Otherwise, HTTP connection is not returned to
		//       to the pool.
		body := ctx.Response.Body
		ctx.Response = nil

		_, err := io.Copy(io.Discard, body)
		if err != nil {
			return err
		}

		err = body.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// Unsafe evaluates current context of HTTP I/O
func (ctx *Context) Unsafe() error {
	eg := ctx.Request

	if ctx.Context != nil {
		eg = eg.WithContext(ctx.Context)
	}

	ctx.logSend(ctx.stack.LogLevel, eg)

	in, err := ctx.stack.Do(eg)
	if err != nil {
		return err
	}

	ctx.Response = in

	ctx.logRecv(ctx.stack.LogLevel, in)

	return nil
}

func (ctx *Context) discardBody() error {
	if ctx.Response != nil {
		// Note: due to Golang HTTP pool implementation we need to consume and
		//       discard body. Otherwise, HTTP connection is not returned to
		//       to the pool.
		body := ctx.Response.Body
		ctx.Response = nil

		if _, err := io.Copy(io.Discard, body); err != nil {
			return err
		}

		if err := body.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *Context) logSend(level int, eg *http.Request) {
	if level >= 1 {
		if msg, err := httputil.DumpRequest(eg, level == 3); err == nil {
			log.Printf(">>>>\n%s\n", msg)
		}
	}
}

func (ctx *Context) logRecv(level int, in *http.Response) {
	if level >= 2 {
		if msg, err := httputil.DumpResponse(in, level == 3); err == nil {
			log.Printf("<<<<\n%s\n", msg)
		}
	}
}
