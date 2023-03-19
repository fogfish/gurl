//
// Copyright (C) 2019 - 2023 Dmitry Kolesnikov
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
	"strconv"
	"strings"
	"time"

	"github.com/ajg/form"
	"github.com/fogfish/gurl/v2"
	"github.com/fogfish/gurl/v2/http"
	"github.com/google/go-cmp/cmp"
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
			return &gurl.NoMatch{
				ID:       "http.Code",
				Diff:     fmt.Sprintf("+ Status Code: %d\n- Status Code: %d", status, code[0]),
				Protocol: "StatusCode",
				Expect:   code[0],
				Actual:   status,
			}
		}
		return nil
	}
}

func hasCode(s []http.StatusCode, e int) bool {
	for _, a := range s {
		if a.StatusCode() == e {
			return true
		}
	}
	return false
}

// StatusCode is a warpper type over http.StatusCode
//
//	http.Join(
//		...
//		ƒ.Code(http.StatusOK),
//	)
//
// so that response code is matched using constant
//
//	http.Join(
//		...
//		ƒ.Status.OK,
//	)
type StatusCode int

// Status is collection of constants for HTTP Status Code checks
//
//	ƒ.Status.OK
//	ƒ.Status.NotFound
const Status = StatusCode(0)

func (StatusCode) eval(code http.StatusCode, cat *http.Context) error {
	if err := cat.Unsafe(); err != nil {
		return err
	}

	status := cat.Response.StatusCode
	if !hasCode([]http.StatusCode{code}, status) {
		return &gurl.NoMatch{
			ID:       "http.Code",
			Diff:     fmt.Sprintf("+ Status Code: %d\n- Status Code: %d", status, code),
			Protocol: "StatusCode",
			Expect:   code,
			Actual:   status,
		}
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

// helper function to match HTTP header to value
func match(ctx *http.Context, header string, value string) error {
	h := ctx.Response.Header.Get(string(header))
	if h == "" {
		return &gurl.NoMatch{
			ID:       "http.Header",
			Diff:     fmt.Sprintf("- %s: %s", string(header), value),
			Protocol: header,
			Expect:   value,
			Actual:   nil,
		}
	}

	if value != "*" && !strings.HasPrefix(h, value) {
		return &gurl.NoMatch{
			ID:       "http.Header",
			Diff:     fmt.Sprintf("+ %s: %s\n- %s: %s", string(header), h, string(header), value),
			Protocol: header,
			Expect:   value,
			Actual:   h,
		}
	}

	return nil
}

// helper function to lift header value to string
func liftString(ctx *http.Context, header string, value *string) error {
	val := ctx.Response.Header.Get(string(header))
	if val == "" {
		return &gurl.NoMatch{
			ID:       "http.Header",
			Diff:     fmt.Sprintf("- %s: *", string(header)),
			Protocol: header,
		}
	}

	*value = val
	return nil
}

func liftInt(ctx *http.Context, header string, value *int) error {
	val := ctx.Response.Header.Get(string(header))
	if val == "" {
		return &gurl.NoMatch{
			ID:       "http.Header",
			Diff:     fmt.Sprintf("- %s: *", string(header)),
			Protocol: header,
		}
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	*value = num
	return nil
}

func liftTime(ctx *http.Context, header string, value *time.Time) error {
	val := ctx.Response.Header.Get(string(header))
	if val == "" {
		return &gurl.NoMatch{
			ID:       "http.Header",
			Diff:     fmt.Sprintf("- %s: *", string(header)),
			Protocol: header,
		}
	}

	t, err := time.Parse(time.RFC1123, val)
	if err != nil {
		return err
	}

	*value = t
	return nil
}

// Header matches or lifts header value
func Header[T http.MatchableHeaderValues](header string, value T) http.Arrow {
	switch v := any(value).(type) {
	case string:
		return HeaderOf[string](header).Is(v)
	case int:
		return HeaderOf[int](header).Is(v)
	case time.Time:
		return HeaderOf[time.Time](header).Is(v)
	case *string:
		return HeaderOf[string](header).To(v)
	case *int:
		return HeaderOf[int](header).To(v)
	case *time.Time:
		return HeaderOf[time.Time](header).To(v)
	default:
		panic("invalid type")
	}
}

// Header matches presence of header in the response or match its entire content.
// The execution fails with BadMatchHead if the matched value do not meet expectations.
//
//	  http.Join(
//			...
//			ƒ.ContentType.JSON,
//			ƒ.ContentEncoding.Is(...),
//		)
type HeaderOf[T http.ReadableHeaderValues] string

// Matches header to any value
func (h HeaderOf[T]) Any(ctx *http.Context) error {
	return match(ctx, string(h), "*")
}

// Matches value of HTTP header
func (h HeaderOf[T]) Is(value T) http.Arrow {
	switch v := any(value).(type) {
	case string:
		return func(ctx *http.Context) error {
			return match(ctx, string(h), v)
		}
	case int:
		return func(ctx *http.Context) error {
			return match(ctx, string(h), strconv.Itoa(v))
		}
	case time.Time:
		return func(ctx *http.Context) error {
			return match(ctx, string(h), v.UTC().Format(time.RFC1123))
		}
	default:
		panic("invalid type")
	}
}

// Lifts value of HTTP header
func (h HeaderOf[T]) To(value *T) http.Arrow {
	switch v := any(value).(type) {
	case *string:
		return func(ctx *http.Context) error {
			return liftString(ctx, string(h), v)
		}
	case *int:
		return func(ctx *http.Context) error {
			return liftInt(ctx, string(h), v)
		}
	case *time.Time:
		return func(ctx *http.Context) error {
			return liftTime(ctx, string(h), v)
		}
	default:
		panic("invalid type")
	}
}

// Type of HTTP Header, Content-Type enumeration
//
//	const ContentType = HeaderEnumContent("Content-Type")
//	ƒ.ContentType.JSON
type HeaderEnumContent string

// Matches header to any value
func (h HeaderEnumContent) Any(ctx *http.Context) error {
	return match(ctx, string(h), "*")
}

// Matches value of HTTP header
func (h HeaderEnumContent) Is(value string) http.Arrow {
	return func(ctx *http.Context) error {
		return match(ctx, string(h), value)
	}
}

// Matches value of HTTP header
func (h HeaderEnumContent) To(value *string) http.Arrow {
	return func(ctx *http.Context) error {
		return liftString(ctx, string(h), value)
	}
}

// ApplicationJSON defines header `???: application/json`
func (h HeaderEnumContent) ApplicationJSON(ctx *http.Context) error {
	return match(ctx, string(h), "application/json")
}

// JSON defines header `???: application/json`
func (h HeaderEnumContent) JSON(ctx *http.Context) error {
	return match(ctx, string(h), "application/json")
}

// Form defined Header `???: application/x-www-form-urlencoded`
func (h HeaderEnumContent) Form(ctx *http.Context) error {
	return match(ctx, string(h), "application/x-www-form-urlencoded")
}

// TextPlain defined Header `???: text/plain`
func (h HeaderEnumContent) TextPlain(ctx *http.Context) error {
	return match(ctx, string(h), "text/plain")
}

// Text defined Header `???: text/plain`
func (h HeaderEnumContent) Text(ctx *http.Context) error {
	return match(ctx, string(h), "text/plain")
}

// TextHTML defined Header `???: text/html`
func (h HeaderEnumContent) TextHTML(ctx *http.Context) error {
	return match(ctx, string(h), "text/html")
}

// HTML defined Header `???: text/html`
func (h HeaderEnumContent) HTML(ctx *http.Context) error {
	return match(ctx, string(h), "text/html")
}

// Type of HTTP Header, Connection enumeration
//
//	const Connection = HeaderEnumConnection("Connection")
//	ƒ.Connection.KeepAlive
type HeaderEnumConnection string

// Matches header to any value
func (h HeaderEnumConnection) Any(ctx *http.Context) error {
	return match(ctx, string(h), "*")
}

// Matches value of HTTP header
func (h HeaderEnumConnection) Is(value string) http.Arrow {
	return func(ctx *http.Context) error {
		return match(ctx, string(h), value)
	}
}

// Matches value of HTTP header
func (h HeaderEnumConnection) To(value *string) http.Arrow {
	return func(ctx *http.Context) error {
		return liftString(ctx, string(h), value)
	}
}

// KeepAlive defines header `???: keep-alive`
func (h HeaderEnumConnection) KeepAlive(ctx *http.Context) error {
	return match(ctx, string(h), "keep-alive")
}

// Close defines header `???: close`
func (h HeaderEnumConnection) Close(ctx *http.Context) error {
	return match(ctx, string(h), "close")
}

// Type of HTTP Header, Transfer-Encoding enumeration
//
//	const TransferEncoding = HeaderEnumTransferEncoding("Transfer-Encoding")
//	ƒ.TransferEncoding.Chunked
type HeaderEnumTransferEncoding string

// Matches header to any value
func (h HeaderEnumTransferEncoding) Any(ctx *http.Context) error {
	return match(ctx, string(h), "*")
}

// Matches value of HTTP header
func (h HeaderEnumTransferEncoding) Is(value string) http.Arrow {
	return func(ctx *http.Context) error {
		return match(ctx, string(h), value)
	}
}

// Matches value of HTTP header
func (h HeaderEnumTransferEncoding) To(value *string) http.Arrow {
	return func(ctx *http.Context) error {
		return liftString(ctx, string(h), value)
	}
}

// Chunked defines header `Transfer-Encoding: chunked`
func (h HeaderEnumTransferEncoding) Chunked(ctx *http.Context) error {
	return match(ctx, string(h), "chunked")
}

// Identity defines header `Transfer-Encoding: identity`
func (h HeaderEnumTransferEncoding) Identity(ctx *http.Context) error {
	return match(ctx, string(h), "identity")
}

// List of supported HTTP header constants
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Response_fields
const (
	Age              = HeaderOf[int]("Age")
	CacheControl     = HeaderOf[string]("Cache-Control")
	Connection       = HeaderEnumConnection("Connection")
	ContentEncoding  = HeaderOf[string]("Content-Encoding")
	ContentLanguage  = HeaderOf[string]("Content-Language")
	ContentLength    = HeaderOf[int]("Content-Length")
	ContentLocation  = HeaderOf[string]("Content-Location")
	ContentMD5       = HeaderOf[string]("Content-MD5")
	ContentRange     = HeaderOf[string]("Content-Range")
	ContentType      = HeaderEnumContent("Content-Type")
	Date             = HeaderOf[time.Time]("Date")
	ETag             = HeaderOf[string]("ETag")
	Expires          = HeaderOf[time.Time]("Expires")
	LastModified     = HeaderOf[time.Time]("Last-Modified")
	Link             = HeaderOf[string]("Link")
	Location         = HeaderOf[string]("Location")
	RetryAfter       = HeaderOf[time.Time]("Retry-After")
	Server           = HeaderOf[string]("Server")
	SetCookie        = HeaderOf[string]("Set-Cookie")
	TransferEncoding = HeaderEnumTransferEncoding("Transfer-Encoding")
	Via              = HeaderOf[string]("Via")
)

// Body applies auto decoders for response and returns either binary or
// native Go data structure. The Content-Type header give a hint to decoder.
// Supply the pointer to data target data structure.
func Body[T any](out *T) http.Arrow {
	return func(cat *http.Context) error {
		err := decode(
			cat.Response.Header.Get("Content-Type"),
			cat.Response.Body,
			out,
		)
		cat.Response.Body.Close()
		cat.Response = nil
		return err
	}
}

// Recv is alias for Body, maintained only for compatibility
func Recv[T any](out *T) http.Arrow {
	return Body(out)
}

// Match received payload to defined pattern
func Expect[T any](expect T) http.Arrow {
	return func(cat *http.Context) error {
		var actual T
		err := decode(
			cat.Response.Header.Get("Content-Type"),
			cat.Response.Body,
			&actual,
		)
		cat.Response.Body.Close()
		cat.Response = nil

		diff := cmp.Diff(actual, expect)
		if diff != "" {
			return &gurl.NoMatch{
				ID:       "http.Recv",
				Diff:     diff,
				Protocol: "body",
				Expect:   expect,
				Actual:   actual,
			}
		}

		return err
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
			ID:       "http.Recv",
			Diff:     fmt.Sprintf("- Content-Type: application/{json | www-form}\n+ Content-Type: %s", content),
			Protocol: "codec",
			Actual:   content,
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

// Match received payload to defined pattern
func Match(val string) http.Arrow {
	var pat any
	if err := json.Unmarshal([]byte(val), &pat); err != nil {
		panic(err)
	}

	return func(cat *http.Context) (err error) {
		var val any

		err = decode(
			cat.Response.Header.Get("Content-Type"),
			cat.Response.Body,
			&val,
		)
		cat.Response.Body.Close()
		cat.Response = nil

		if !equivVal(pat, val) {
			return &gurl.NoMatch{
				ID:       "http.Match",
				Protocol: "body",
				Expect:   pat,
				Actual:   val,
			}
		}

		return
	}
}

func equivVal(pat, val any) bool {
	if pp, ok := pat.(string); ok && pp == "_" {
		return true
	}

	switch vv := val.(type) {
	case string:
		pp, ok := pat.(string)
		if !ok {
			return false
		}
		return vv == pp
	case float64:
		pp, ok := pat.(float64)
		if !ok {
			return false
		}
		return vv == pp
	case bool:
		pp, ok := pat.(bool)
		if !ok {
			return false
		}
		return vv == pp
	case []any:
		pp, ok := pat.([]any)
		if !ok {
			return false
		}
		if len(pp) != len(vv) {
			return false
		}
		for i, vvx := range vv {
			if !equivVal(pp[i], vvx) {
				return false
			}
		}
		return true
	case map[string]any:
		pp, ok := pat.(map[string]any)
		if !ok {
			return false
		}
		return equivMap(pp, vv)
	}

	return false
}

func equivMap(pat, val map[string]any) bool {
	for k, p := range pat {
		v, has := val[k]
		if !has {
			return false
		}

		if !equivVal(p, v) {
			return false
		}
	}

	return true
}
