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

// Code is a mandatory statement to match expected HTTP Status Code against
// received one. The execution fails StatusCode error if service responds
// with other value then specified one.
func Code(code ...http.StatusCode) http.Arrow {
	return func(cat *http.Context) error {
		if err := cat.Unsafe(); err != nil {
			return err
		}

		status := cat.Response.StatusCode
		if !hasCode(code, status) {
			return http.NewStatusCode(status, code[0])
		}
		return nil
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

// Status is collection of constants for HTTP Status Code checks
const Status = StatusCode(0)

func (StatusCode) eval(code http.StatusCode, cat *http.Context) error {
	if err := cat.Unsafe(); err != nil {
		return err
	}

	status := cat.Response.StatusCode
	if !hasCode([]http.StatusCode{code}, status) {
		return http.NewStatusCode(status, code)
	}

	return nil
}

/*
TODO:
  Continue
	SwitchingProtocols
	Processing
	EarlyHints
*/

// OK ⟼ http.StatusOK
func (code StatusCode) OK(cat *http.Context) error {
	return code.eval(http.StatusOK, cat)
}

// Created ⟼ http.StatusCreated
func (code StatusCode) Created(cat *http.Context) error {
	return code.eval(http.StatusCreated, cat)
}

// Accepted ⟼ http.StatusAccepted
func (code StatusCode) Accepted(cat *http.Context) error {
	return code.eval(http.StatusAccepted, cat)
}

// NonAuthoritativeInfo ⟼ http.StatusNonAuthoritativeInfo
func (code StatusCode) NonAuthoritativeInfo(cat *http.Context) error {
	return code.eval(http.StatusNonAuthoritativeInfo, cat)
}

// NoContent ⟼ http.StatusNoContent
func (code StatusCode) NoContent(cat *http.Context) error {
	return code.eval(http.StatusNoContent, cat)
}

// ResetContent ⟼ http.StatusResetContent
func (code StatusCode) ResetContent(cat *http.Context) error {
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
func (code StatusCode) MultipleChoices(cat *http.Context) error {
	return code.eval(http.StatusMultipleChoices, cat)
}

// MovedPermanently ⟼ http.StatusMovedPermanently
func (code StatusCode) MovedPermanently(cat *http.Context) error {
	return code.eval(http.StatusMovedPermanently, cat)
}

// Found ⟼ http.StatusFound
func (code StatusCode) Found(cat *http.Context) error {
	return code.eval(http.StatusFound, cat)
}

// SeeOther ⟼ http.StatusSeeOther
func (code StatusCode) SeeOther(cat *http.Context) error {
	return code.eval(http.StatusSeeOther, cat)
}

// NotModified ⟼ http.StatusNotModified
func (code StatusCode) NotModified(cat *http.Context) error {
	return code.eval(http.StatusNotModified, cat)
}

// UseProxy ⟼ http.StatusUseProxy
func (code StatusCode) UseProxy(cat *http.Context) error {
	return code.eval(http.StatusUseProxy, cat)
}

/*
TODO:
	TemporaryRedirect
	PermanentRedirect
*/

// BadRequest ⟼ http.StatusBadRequest
func (code StatusCode) BadRequest(cat *http.Context) error {
	return code.eval(http.StatusBadRequest, cat)
}

// Unauthorized ⟼ http.StatusUnauthorized
func (code StatusCode) Unauthorized(cat *http.Context) error {
	return code.eval(http.StatusUnauthorized, cat)
}

// PaymentRequired ⟼ http.StatusPaymentRequired
func (code StatusCode) PaymentRequired(cat *http.Context) error {
	return code.eval(http.StatusPaymentRequired, cat)
}

// Forbidden ⟼ http.StatusForbidden
func (code StatusCode) Forbidden(cat *http.Context) error {
	return code.eval(http.StatusForbidden, cat)
}

// NotFound ⟼ http.StatusNotFound
func (code StatusCode) NotFound(cat *http.Context) error {
	return code.eval(http.StatusNotFound, cat)
}

// MethodNotAllowed ⟼ http.StatusMethodNotAllowed
func (code StatusCode) MethodNotAllowed(cat *http.Context) error {
	return code.eval(http.StatusMethodNotAllowed, cat)
}

// NotAcceptable ⟼ http.StatusNotAcceptable
func (code StatusCode) NotAcceptable(cat *http.Context) error {
	return code.eval(http.StatusNotAcceptable, cat)
}

// ProxyAuthRequired ⟼ http.StatusProxyAuthRequired
func (code StatusCode) ProxyAuthRequired(cat *http.Context) error {
	return code.eval(http.StatusProxyAuthRequired, cat)
}

// RequestTimeout ⟼ http.StatusRequestTimeout
func (code StatusCode) RequestTimeout(cat *http.Context) error {
	return code.eval(http.StatusRequestTimeout, cat)
}

// Conflict ⟼ http.StatusConflict
func (code StatusCode) Conflict(cat *http.Context) error {
	return code.eval(http.StatusConflict, cat)
}

// Gone ⟼ http.StatusGone
func (code StatusCode) Gone(cat *http.Context) error {
	return code.eval(http.StatusGone, cat)
}

// LengthRequired ⟼ http.StatusLengthRequired
func (code StatusCode) LengthRequired(cat *http.Context) error {
	return code.eval(http.StatusLengthRequired, cat)
}

// PreconditionFailed ⟼ http.StatusPreconditionFailed
func (code StatusCode) PreconditionFailed(cat *http.Context) error {
	return code.eval(http.StatusPreconditionFailed, cat)
}

// RequestEntityTooLarge ⟼ http.StatusRequestEntityTooLarge
func (code StatusCode) RequestEntityTooLarge(cat *http.Context) error {
	return code.eval(http.StatusRequestEntityTooLarge, cat)
}

// RequestURITooLong ⟼ http.StatusRequestURITooLong
func (code StatusCode) RequestURITooLong(cat *http.Context) error {
	return code.eval(http.StatusRequestURITooLong, cat)
}

// UnsupportedMediaType ⟼ http.StatusUnsupportedMediaType
func (code StatusCode) UnsupportedMediaType(cat *http.Context) error {
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
func (code StatusCode) InternalServerError(cat *http.Context) error {
	return code.eval(http.StatusInternalServerError, cat)
}

// NotImplemented ⟼ http.StatusNotImplemented
func (code StatusCode) NotImplemented(cat *http.Context) error {
	return code.eval(http.StatusNotImplemented, cat)
}

// BadGateway ⟼ http.StatusBadGateway
func (code StatusCode) BadGateway(cat *http.Context) error {
	return code.eval(http.StatusBadGateway, cat)
}

// ServiceUnavailable ⟼ http.StatusServiceUnavailable
func (code StatusCode) ServiceUnavailable(cat *http.Context) error {
	return code.eval(http.StatusServiceUnavailable, cat)
}

// GatewayTimeout ⟼ http.StatusGatewayTimeout
func (code StatusCode) GatewayTimeout(cat *http.Context) error {
	return code.eval(http.StatusGatewayTimeout, cat)
}

// HTTPVersionNotSupported ⟼ http.StatusHTTPVersionNotSupported
func (code StatusCode) HTTPVersionNotSupported(cat *http.Context) error {
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

// Header matches presence of header in the response or match its entire content.
// The execution fails with BadMatchHead if the matched value do not meet expectations.
//
//	  http.Join(
//			...
//			ƒ.ContentType.JSON,
//			ƒ.ContentEncoding.Is(...),
//		)
type Header string

// List of supported HTTP header constants
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Response_fields
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
	return func(cat *http.Context) error {
		return header.Match(cat, value)
	}
}

// To matches a header value to closed variable of string type.
func (header Header) To(value *string) http.Arrow {
	return func(cat *http.Context) error {
		val := cat.Response.Header.Get(string(header))
		if val == "" {
			return &gurl.NoMatch{
				Diff:    fmt.Sprintf("- %s: *", string(header)),
				Payload: nil,
			}
		}

		*value = val
		return nil
	}
}

// Match is combinator to check HTTP header value
func (header Header) Match(cat *http.Context, value string) error {
	h := cat.Response.Header.Get(string(header))
	if h == "" {
		return &gurl.NoMatch{
			Diff:    fmt.Sprintf("- %s: %s", string(header), value),
			Payload: nil,
		}
	}

	if value != "*" && !strings.HasPrefix(h, value) {
		return &gurl.NoMatch{
			Diff:    fmt.Sprintf("+ %s: %s\n- %s: %s", string(header), h, string(header), value),
			Payload: map[string]string{string(header): h},
		}
	}

	return nil
}

// Any matches a header value, syntax sugar of Header(...).Is("*")
func (header Header) Any(cat *http.Context) error {
	return header.Match(cat, "*")
}

// Content defines headers for content negotiation
type Content Header

// ApplicationJSON matches header `???: application/json`
func (h Content) ApplicationJSON(cat *http.Context) error {
	return Header(h).Match(cat, "application/json")
}

// JSON matches header `???: application/json`
func (h Content) JSON(cat *http.Context) error {
	return Header(h).Match(cat, "application/json")
}

// Form matches Header `???: application/x-www-form-urlencoded`
func (h Content) Form(cat *http.Context) error {
	return Header(h).Match(cat, "application/x-www-form-urlencoded")
}

// TextPlain matches Header `???: text/plain`
func (h Content) TextPlain(cat *http.Context) error {
	return Header(h).Match(cat, "text/plain")
}

// Text matches Header `???: text/plain`
func (h Content) Text(cat *http.Context) error {
	return Header(h).Match(cat, "text/plain")
}

// TextHTML matches Header `???: text/html`
func (h Content) TextHTML(cat *http.Context) error {
	return Header(h).Match(cat, "text/html")
}

// HTML matches Header `???: text/html`
func (h Content) HTML(cat *http.Context) error {
	return Header(h).Match(cat, "text/html")
}

// Any matches a header value `???: *`
func (h Content) Any(cat *http.Context) error {
	return Header(h).Match(cat, "*")
}

// Is matches value of HTTP header, Use wildcard string ("*") to match any header value
func (h Content) Is(value string) http.Arrow {
	return Header(h).Is(value)
}

// String matches a header value to closed variable of string type.
func (h Content) To(value *string) http.Arrow {
	return Header(h).To(value)
}

// Recv applies auto decoders for response and returns either binary or
// native Go data structure. The Content-Type header give a hint to decoder.
// Supply the pointer to data target data structure.
func Recv[T any](out *T) http.Arrow {
	return func(cat *http.Context) (err error) {
		err = decode(
			cat.Response.Header.Get("Content-Type"),
			cat.Response.Body,
			out,
		)
		cat.Response.Body.Close()
		cat.Response = nil
		return
	}
}

func decode[T any](content string, stream io.ReadCloser, data *T) error {
	switch {
	case strings.Contains(content, "json"):
		return json.NewDecoder(stream).Decode(data)
	case strings.Contains(content, "www-form"):
		return form.NewDecoder(stream).Decode(data)
	default:
		return &gurl.NoMatch{
			Diff:    fmt.Sprintf("- Content-Type: application/*\n+ Content-Type: %s", content),
			Payload: map[string]string{"Content-Type": content},
		}
	}
}

// Bytes receive raw binary from HTTP response
func Bytes(val *[]byte) http.Arrow {
	return func(cat *http.Context) (err error) {
		*val, err = io.ReadAll(cat.Response.Body)
		cat.Response.Body.Close()
		cat.Response = nil
		return
	}
}
