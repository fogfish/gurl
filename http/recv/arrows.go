//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

// Package recv defines a pure computations to compose HTTP response receivers
package recv

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ajg/form"
	"github.com/fogfish/gurl"
	"github.com/fogfish/gurl/http"
)

//-------------------------------------------------------------------
//
// core arrows
//
//-------------------------------------------------------------------

/*

Code is a mandatory statement to match expected HTTP Status Code against
received one. The execution fails StatusCode error if service responds
with other value then specified one.
*/
func Code(code ...http.StatusCode) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		if cat = cat.Unsafe(); cat.Fail != nil {
			return cat
		}

		status := cat.HTTP.Recv.StatusCode
		if !hasCode(code, status) {
			cat.Fail = http.NewStatusCode(status, code[0])
		}
		return cat
	}
}

func hasCode(s []http.StatusCode, e int) bool {
	for _, a := range s {
		if a.Value() == e {
			return true
		}
	}
	return false
}

/*

StatusCode is a warpper type over http.StatusCode

  http.Join(
		...
		ƒ.Code(http.StatusOK),
	)

	so that response code is matched using constant
	http.Join(
		...
		ƒ.Status.OK,
	)
*/
type StatusCode int

// HTTP Status Code check constants
const Status = StatusCode(0)

func (StatusCode) eval(code http.StatusCode, cat *gurl.IOCat) *gurl.IOCat {
	if cat = cat.Unsafe(); cat.Fail != nil {
		return cat
	}

	status := cat.HTTP.Recv.StatusCode
	if !hasCode([]http.StatusCode{code}, status) {
		cat.Fail = http.NewStatusCode(status, code)
	}
	return cat
}

// Continue ⟼ http.StatusContinue
func (code StatusCode) Continue(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusContinue, cat)
}

/*
TODO:
	SwitchingProtocols
	Processing
	EarlyHints
*/

// OK ⟼ http.StatusOK
func (code StatusCode) OK(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusOK, cat)
}

// Created ⟼ http.StatusCreated
func (code StatusCode) Created(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusCreated, cat)
}

// Accepted ⟼ http.StatusAccepted
func (code StatusCode) Accepted(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusAccepted, cat)
}

// NonAuthoritativeInfo ⟼ http.StatusNonAuthoritativeInfo
func (code StatusCode) NonAuthoritativeInfo(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusNonAuthoritativeInfo, cat)
}

// NoContent ⟼ http.StatusNoContent
func (code StatusCode) NoContent(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusNoContent, cat)
}

// ResetContent ⟼ http.StatusResetContent
func (code StatusCode) ResetContent(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusResetContent, cat)
}

/*
TODO:
	PartialContent
	MultiStatus
	AlreadyReported
	IMUsed
*/

// MultipleChoices ⟼ http.StatusMultipleChoices
func (code StatusCode) MultipleChoices(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusMultipleChoices, cat)
}

// MovedPermanently ⟼ http.StatusMovedPermanently
func (code StatusCode) MovedPermanently(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusMovedPermanently, cat)
}

// Found ⟼ http.StatusFound
func (code StatusCode) Found(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusFound, cat)
}

// SeeOther ⟼ http.StatusSeeOther
func (code StatusCode) SeeOther(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusSeeOther, cat)
}

// NotModified ⟼ http.StatusNotModified
func (code StatusCode) NotModified(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusNotModified, cat)
}

// UseProxy ⟼ http.StatusUseProxy
func (code StatusCode) UseProxy(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusUseProxy, cat)
}

/*
TODO:
	TemporaryRedirect
	PermanentRedirect
*/

// BadRequest ⟼ http.StatusBadRequest
func (code StatusCode) BadRequest(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusBadRequest, cat)
}

// Unauthorized ⟼ http.StatusUnauthorized
func (code StatusCode) Unauthorized(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusUnauthorized, cat)
}

// PaymentRequired ⟼ http.StatusPaymentRequired
func (code StatusCode) PaymentRequired(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusPaymentRequired, cat)
}

// Forbidden ⟼ http.StatusForbidden
func (code StatusCode) Forbidden(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusForbidden, cat)
}

// NotFound ⟼ http.StatusNotFound
func (code StatusCode) NotFound(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusNotFound, cat)
}

// MethodNotAllowed ⟼ http.StatusMethodNotAllowed
func (code StatusCode) MethodNotAllowed(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusMethodNotAllowed, cat)
}

// NotAcceptable ⟼ http.StatusNotAcceptable
func (code StatusCode) NotAcceptable(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusNotAcceptable, cat)
}

// ProxyAuthRequired ⟼ http.StatusProxyAuthRequired
func (code StatusCode) ProxyAuthRequired(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusProxyAuthRequired, cat)
}

// RequestTimeout ⟼ http.StatusRequestTimeout
func (code StatusCode) RequestTimeout(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusRequestTimeout, cat)
}

// Conflict ⟼ http.StatusConflict
func (code StatusCode) Conflict(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusConflict, cat)
}

// Gone ⟼ http.StatusGone
func (code StatusCode) Gone(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusGone, cat)
}

// LengthRequired ⟼ http.StatusLengthRequired
func (code StatusCode) LengthRequired(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusLengthRequired, cat)
}

// PreconditionFailed ⟼ http.StatusPreconditionFailed
func (code StatusCode) PreconditionFailed(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusPreconditionFailed, cat)
}

// RequestEntityTooLarge ⟼ http.StatusRequestEntityTooLarge
func (code StatusCode) RequestEntityTooLarge(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusRequestEntityTooLarge, cat)
}

// RequestURITooLong ⟼ http.StatusRequestURITooLong
func (code StatusCode) RequestURITooLong(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusRequestURITooLong, cat)
}

// UnsupportedMediaType ⟼ http.StatusUnsupportedMediaType
func (code StatusCode) UnsupportedMediaType(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusUnsupportedMediaType, cat)
}

/*
TODO:
	RequestedRangeNotSatisfiable
	ExpectationFailed
	Teapot
	MisdirectedRequest
	UnprocessableEntity
	Locked
	FailedDependency
	TooEarly
	UpgradeRequired
	PreconditionRequired
	TooManyRequests
	RequestHeaderFieldsTooLarge
	UnavailableForLegalReasons
*/

// InternalServerError ⟼ http.StatusInternalServerError
func (code StatusCode) InternalServerError(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusInternalServerError, cat)
}

// NotImplemented ⟼ http.StatusNotImplemented
func (code StatusCode) NotImplemented(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusNotImplemented, cat)
}

// BadGateway ⟼ http.StatusBadGateway
func (code StatusCode) BadGateway(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusBadGateway, cat)
}

// ServiceUnavailable ⟼ http.StatusServiceUnavailable
func (code StatusCode) ServiceUnavailable(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusServiceUnavailable, cat)
}

// GatewayTimeout ⟼ http.StatusGatewayTimeout
func (code StatusCode) GatewayTimeout(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusGatewayTimeout, cat)
}

// HTTPVersionNotSupported ⟼ http.StatusHTTPVersionNotSupported
func (code StatusCode) HTTPVersionNotSupported(cat *gurl.IOCat) *gurl.IOCat {
	return code.eval(http.StatusHTTPVersionNotSupported, cat)
}

/*
TODO:
	VariantAlsoNegotiates
	InsufficientStorage
	LoopDetected
	NotExtended
	NetworkAuthenticationRequired
*/

/*

Header matches presence of header in the response or match its entire content.
The execution fails with BadMatchHead if the matched value do not meet expectations.

  http.Join(
		...
		ƒ.ContentType.JSON,
		ƒ.ContentEncoding.Is(...),
	)
*/
type Header string

/*

List of supported HTTP header constants
https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Response_fields
*/
const (
	CacheControl     = Header("Cache-Control")
	Connection       = Header("Connection")
	ContentEncoding  = Header("Content-Encoding")
	ContentLanguage  = Header("Content-Language")
	ContentLength    = Header("Content-Length")
	ContentType      = Content("Content-Type")
	Date             = Header("Date")
	ETag             = Header("ETag")
	Expires          = Header("Expires")
	LastModified     = Header("Last-Modified")
	Link             = Header("Link")
	Location         = Header("Location")
	Server           = Header("Server")
	SetCookie        = Header("Set-Cookie")
	TransferEncoding = Header("Transfer-Encoding")
)

// Is matches value of HTTP header, Use wildcard string ("*") to match any header value
func (header Header) Is(value string) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		return header.Match(cat, value)
	}
}

// String matches a header value to closed variable of string type.
func (header Header) String(value *string) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		val := cat.HTTP.Recv.Header.Get(string(header))
		if val == "" {
			cat.Fail = &gurl.Mismatch{
				Diff:    fmt.Sprintf("- %s: *", string(header)),
				Payload: nil,
			}
		} else {
			*value = val
		}

		return cat
	}
}

// Any matches a header value, syntax sugar of Header(...).Is("*")
func (header Header) Any() http.Arrow {
	return header.Is("*")
}

// Match is combinator to check HTTP header value
func (header Header) Match(cat *gurl.IOCat, value string) *gurl.IOCat {
	h := cat.HTTP.Recv.Header.Get(string(header))
	if h == "" {
		cat.Fail = &gurl.Mismatch{
			Diff:    fmt.Sprintf("- %s: %s", string(header), value),
			Payload: nil,
		}
		return cat
	}

	if value != "*" && !strings.HasPrefix(h, value) {
		cat.Fail = &gurl.Mismatch{
			Diff:    fmt.Sprintf("+ %s: %s\n- %s: %s", string(header), h, string(header), value),
			Payload: map[string]string{string(header): h},
		}
		return cat
	}

	return cat
}

// Content defines headers for content negotiation
type Content Header

// JSON matches header `???: application/json`
func (h Content) JSON(cat *gurl.IOCat) *gurl.IOCat {
	return Header(h).Match(cat, "application/json")
}

// Form matches Header `???: application/x-www-form-urlencoded`
func (h Content) Form(cat *gurl.IOCat) *gurl.IOCat {
	return Header(h).Match(cat, "application/x-www-form-urlencoded")
}

// Text matches Header `???: plain/text`
func (h Content) Text(cat *gurl.IOCat) *gurl.IOCat {
	return Header(h).Match(cat, "plain/text")
}

// HTML matches Header `???: plain/html`
func (h Content) HTML(cat *gurl.IOCat) *gurl.IOCat {
	return Header(h).Match(cat, "plain/html")
}

// Is matches value of HTTP header, Use wildcard string ("*") to match any header value
func (h Content) Is(value string) http.Arrow {
	return Header(h).Is(value)
}

// String matches a header value to closed variable of string type.
func (h Content) String(value *string) http.Arrow {
	return Header(h).String(value)
}

/*

Recv applies auto decoders for response and returns either binary or
native Go data structure. The Content-Type header give a hint to decoder.
Supply the pointer to data target data structure.
*/
func Recv(out interface{}) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		cat.Fail = decode(
			cat.HTTP.Recv.Header.Get("Content-Type"),
			cat.HTTP.Recv.Body,
			&out,
		)
		cat.HTTP.Recv.Body.Close()
		cat.HTTP.Recv.Payload = out
		cat.HTTP.Recv.Response = nil
		return cat
	}
}

func decode(content string, stream io.ReadCloser, data interface{}) error {
	switch {
	case strings.Contains(content, "json"):
		return json.NewDecoder(stream).Decode(&data)
	case strings.Contains(content, "www-form"):
		return form.NewDecoder(stream).Decode(&data)
	default:
		return &gurl.Mismatch{
			Diff:    fmt.Sprintf("- Content-Type: application/*\n+ Content-Type: %s", content),
			Payload: map[string]string{"Content-Type": content},
		}
	}
}

/*

Bytes receive raw binary from HTTP response
*/
func Bytes(val *[]byte) http.Arrow {
	return func(cat *gurl.IOCat) *gurl.IOCat {
		*val, cat.Fail = ioutil.ReadAll(cat.HTTP.Recv.Body)
		cat.Fail = cat.HTTP.Recv.Body.Close()
		cat.HTTP.Recv.Response = nil
		cat.HTTP.Recv.Payload = string(*val)
		return cat
	}
}
