//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package gurl

import (
	"fmt"
	"net/http"
)

/*

ProtocolNotSupported is returned if handling of URL schema is not supported by the library
*/
type ProtocolNotSupported string

func (e ProtocolNotSupported) Error() string {
	return fmt.Sprintf("Not supported protocol: %s", string(e))
}

/*

StatusCode is a base type for typesafe HTTP status codes. The library advertises
a usage of "pattern-matching" on HTTP status handling, which helps developers to
catch mismatch of HTTP statuses along with other side-effect failures.

The final values of HTTP statuses embeds StatusCode type. It makes them to look
like a "sum-types" and be compatible with any other error (side effect failures)
within IO category. Use final type instances in the error handling routines.

Use type switch for error handling "branches"

  switch e := io.Fail.(type) {
  case nil:
    // Nothing
  case *gurl.StatusOK:
    // HTTP 200 OK
  case *gurl.StatusNotFound:
    // HTTP 404 NotFound
  default:
    // any other errors
  }

Conditional error handling on expected HTTP Status

  if errors.Is(io.Fail, gurl.NewStatusNotFound()) {
  }

Conditional error handling on any HTTP Status

  if _, ok := io.Fail.(gurl.StatusCodeAny); ok {
  }

*/
type StatusCode int

func mkStatusCode(code int, required int) StatusCode {
	return StatusCode((required << 16) | code)
}

// Error makes StatusCode to be error
func (e StatusCode) Error() string {
	status := e.Value()
	if await := e.Await(); await != 0 {
		return fmt.Sprintf("HTTP Status `%d %s`, required `%d %s`.", status, http.StatusText(status), await, http.StatusText(await))
	}
	return fmt.Sprintf("HTTP %d %s", status, http.StatusText(status))
}

// Is compares wrapped errors
func (e StatusCode) Is(err error) bool {
	if code, ok := err.(StatusCodeAny); ok {
		return e.Value() == code.Value()
	}
	return false
}

// Value transforms StatusCode type to integer value: StatusCode ⟼ int
func (e StatusCode) Value() int {
	return int(e) & 0xffff
}

// Await returns awaited (expected, required) HTTP status
func (e StatusCode) Await() int {
	return int(e) >> 16
}

/*

StatusCodeAny is a type that matches only HTTP status errors.
Use it to conditional handle only HTTP errors.

  if _, ok := io.Fail.(gurl.StatusCodeAny); ok {
  }

  switch e := io.Fail.(type) {
  case nil:
    // Nothing
  case gurl.StatusCodeAny:
    // any HTTP Status
  default:
    // any other errors
	}

*/
type StatusCodeAny interface {
	Error() string
	Value() int
	Await() int
}

// StatusUnknown ⤳ any unmapped type
type StatusUnknown struct{ StatusCode }

//
//
// StatusContinue ⤳ https://httpstatuses.com/100
type StatusContinue struct{ StatusCode }

// NewStatusContinue ⤳ https://httpstatuses.com/100
func NewStatusContinue() *StatusContinue {
	return &StatusContinue{StatusCode(http.StatusContinue)}
}

// StatusSwitchingProtocols ⤳ https://httpstatuses.com/101
type StatusSwitchingProtocols struct{ StatusCode }

// NewStatusSwitchingProtocols ⤳ https://httpstatuses.com/101
func NewStatusSwitchingProtocols() *StatusSwitchingProtocols {
	return &StatusSwitchingProtocols{StatusCode(http.StatusSwitchingProtocols)}
}

// StatusProcessing ⤳ https://httpstatuses.com/102
type StatusProcessing struct{ StatusCode }

// NewStatusProcessing ⤳ https://httpstatuses.com/102
func NewStatusProcessing() *StatusProcessing {
	return &StatusProcessing{StatusCode(http.StatusProcessing)}
}

// StatusEarlyHints ⤳ https://httpstatuses.com/103
type StatusEarlyHints struct{ StatusCode }

// NewStatusEarlyHints ⤳ https://httpstatuses.com/103
func NewStatusEarlyHints() *StatusEarlyHints {
	return &StatusEarlyHints{StatusCode(http.StatusEarlyHints)}
}

//
//
// StatusOK ⤳ https://httpstatuses.com/200
type StatusOK struct{ StatusCode }

// NewStatusOK ⤳ https://httpstatuses.com/200
func NewStatusOK() *StatusOK {
	return &StatusOK{StatusCode(http.StatusOK)}
}

// StatusCreated ⤳ https://httpstatuses.com/201
type StatusCreated struct{ StatusCode }

// NewStatusCreated ⤳ https://httpstatuses.com/201
func NewStatusCreated() *StatusCreated {
	return &StatusCreated{StatusCode(http.StatusCreated)}
}

// StatusAccepted ⤳ https://httpstatuses.com/202
type StatusAccepted struct{ StatusCode }

// NewStatusAccepted ⤳ https://httpstatuses.com/202
func NewStatusAccepted() *StatusAccepted {
	return &StatusAccepted{StatusCode(http.StatusAccepted)}
}

// StatusNonAuthoritativeInfo ⤳ https://httpstatuses.com/203
type StatusNonAuthoritativeInfo struct{ StatusCode }

// NewStatusNonAuthoritativeInfo ⤳ https://httpstatuses.com/203
func NewStatusNonAuthoritativeInfo() *StatusNonAuthoritativeInfo {
	return &StatusNonAuthoritativeInfo{StatusCode(http.StatusNonAuthoritativeInfo)}
}

// StatusNoContent ⤳ https://httpstatuses.com/204
type StatusNoContent struct{ StatusCode }

// NewStatusNoContent ⤳ https://httpstatuses.com/204
func NewStatusNoContent() *StatusNoContent {
	return &StatusNoContent{StatusCode(http.StatusNoContent)}
}

// StatusResetContent ⤳ https://httpstatuses.com/205
type StatusResetContent struct{ StatusCode }

// NewStatusResetContent ⤳ https://httpstatuses.com/205
func NewStatusResetContent() *StatusResetContent {
	return &StatusResetContent{StatusCode(http.StatusResetContent)}
}

// StatusPartialContent ⤳ https://httpstatuses.com/206
type StatusPartialContent struct{ StatusCode }

// NewStatusPartialContent ⤳ https://httpstatuses.com/206
func NewStatusPartialContent() *StatusPartialContent {
	return &StatusPartialContent{StatusCode(http.StatusPartialContent)}
}

// StatusMultiStatus ⤳ https://httpstatuses.com/207
type StatusMultiStatus struct{ StatusCode }

// NewStatusMultiStatus ⤳ https://httpstatuses.com/207
func NewStatusMultiStatus() *StatusMultiStatus {
	return &StatusMultiStatus{StatusCode(http.StatusMultiStatus)}
}

// StatusAlreadyReported ⤳ https://httpstatuses.com/208
type StatusAlreadyReported struct{ StatusCode }

// NewStatusAlreadyReported ⤳ https://httpstatuses.com/208
func NewStatusAlreadyReported() *StatusAlreadyReported {
	return &StatusAlreadyReported{StatusCode(http.StatusAlreadyReported)}
}

// StatusIMUsed ⤳ https://httpstatuses.com/226
type StatusIMUsed struct{ StatusCode }

// NewStatusIMUsed ⤳ https://httpstatuses.com/226
func NewStatusIMUsed() *StatusIMUsed {
	return &StatusIMUsed{StatusCode(http.StatusIMUsed)}
}

//
//
// StatusMultipleChoices ⤳ https://httpstatuses.com/300
type StatusMultipleChoices struct{ StatusCode }

// NewStatusMultipleChoices ⤳ https://httpstatuses.com/300
func NewStatusMultipleChoices() *StatusMultipleChoices {
	return &StatusMultipleChoices{StatusCode(http.StatusMultipleChoices)}
}

// StatusMovedPermanently ⤳ https://httpstatuses.com/301
type StatusMovedPermanently struct{ StatusCode }

// NewStatusMovedPermanently ⤳ https://httpstatuses.com/301
func NewStatusMovedPermanently() *StatusMovedPermanently {
	return &StatusMovedPermanently{StatusCode(http.StatusMovedPermanently)}
}

// StatusFound ⤳ https://httpstatuses.com/302
type StatusFound struct{ StatusCode }

// NewStatusFound ⤳ https://httpstatuses.com/302
func NewStatusFound() *StatusFound {
	return &StatusFound{StatusCode(http.StatusFound)}
}

// StatusSeeOther ⤳ https://httpstatuses.com/303
type StatusSeeOther struct{ StatusCode }

// NewStatusSeeOther ⤳ https://httpstatuses.com/303
func NewStatusSeeOther() *StatusSeeOther {
	return &StatusSeeOther{StatusCode(http.StatusSeeOther)}
}

// StatusNotModified ⤳ https://httpstatuses.com/304
type StatusNotModified struct{ StatusCode }

// NewStatusNotModified ⤳ https://httpstatuses.com/304
func NewStatusNotModified() *StatusNotModified {
	return &StatusNotModified{StatusCode(http.StatusNotModified)}
}

// StatusUseProxy ⤳ https://httpstatuses.com/305
type StatusUseProxy struct{ StatusCode }

// NewStatusUseProxy ⤳ https://httpstatuses.com/305
func NewStatusUseProxy() *StatusUseProxy {
	return &StatusUseProxy{StatusCode(http.StatusUseProxy)}
}

// StatusTemporaryRedirect ⤳ https://httpstatuses.com/307
type StatusTemporaryRedirect struct{ StatusCode }

// NewStatusTemporaryRedirect ⤳ https://httpstatuses.com/307
func NewStatusTemporaryRedirect() *StatusTemporaryRedirect {
	return &StatusTemporaryRedirect{StatusCode(http.StatusTemporaryRedirect)}
}

// StatusPermanentRedirect ⤳ https://httpstatuses.com/308
type StatusPermanentRedirect struct{ StatusCode }

// NewStatusPermanentRedirect ⤳ https://httpstatuses.com/308
func NewStatusPermanentRedirect() *StatusPermanentRedirect {
	return &StatusPermanentRedirect{StatusCode(http.StatusPermanentRedirect)}
}

//
//
// StatusBadRequest ⤳ https://httpstatuses.com/400
type StatusBadRequest struct{ StatusCode }

// NewStatusBadRequest ⤳ https://httpstatuses.com/400
func NewStatusBadRequest() *StatusBadRequest {
	return &StatusBadRequest{StatusCode(http.StatusBadRequest)}
}

// StatusUnauthorized ⤳ https://httpstatuses.com/401
type StatusUnauthorized struct{ StatusCode }

// NewStatusUnauthorized ⤳ https://httpstatuses.com/401
func NewStatusUnauthorized() *StatusUnauthorized {
	return &StatusUnauthorized{StatusCode(http.StatusUnauthorized)}
}

// StatusPaymentRequired ⤳ https://httpstatuses.com/402
type StatusPaymentRequired struct{ StatusCode }

// NewStatusPaymentRequired ⤳ https://httpstatuses.com/402
func NewStatusPaymentRequired() *StatusPaymentRequired {
	return &StatusPaymentRequired{StatusCode(http.StatusPaymentRequired)}
}

// StatusForbidden ⤳ https://httpstatuses.com/403
type StatusForbidden struct{ StatusCode }

// NewStatusForbidden ⤳ https://httpstatuses.com/403
func NewStatusForbidden() *StatusForbidden {
	return &StatusForbidden{StatusCode(http.StatusForbidden)}
}

// StatusNotFound ⤳ https://httpstatuses.com/404
type StatusNotFound struct{ StatusCode }

// NewStatusNotFound ⤳ https://httpstatuses.com/404
func NewStatusNotFound() *StatusNotFound {
	return &StatusNotFound{StatusCode(http.StatusNotFound)}
}

// StatusMethodNotAllowed ⤳ https://httpstatuses.com/405
type StatusMethodNotAllowed struct{ StatusCode }

// NewStatusMethodNotAllowed ⤳ https://httpstatuses.com/405
func NewStatusMethodNotAllowed() *StatusMethodNotAllowed {
	return &StatusMethodNotAllowed{StatusCode(http.StatusMethodNotAllowed)}
}

// StatusNotAcceptable ⤳ https://httpstatuses.com/406
type StatusNotAcceptable struct{ StatusCode }

// NewStatusNotAcceptable ⤳ https://httpstatuses.com/406
func NewStatusNotAcceptable() *StatusNotAcceptable {
	return &StatusNotAcceptable{StatusCode(http.StatusNotAcceptable)}
}

// StatusProxyAuthRequired ⤳ https://httpstatuses.com/407
type StatusProxyAuthRequired struct{ StatusCode }

// NewStatusProxyAuthRequired ⤳ https://httpstatuses.com/407
func NewStatusProxyAuthRequired() *StatusProxyAuthRequired {
	return &StatusProxyAuthRequired{StatusCode(http.StatusProxyAuthRequired)}
}

// StatusRequestTimeout ⤳ https://httpstatuses.com/408
type StatusRequestTimeout struct{ StatusCode }

// NewStatusRequestTimeout ⤳ https://httpstatuses.com/408
func NewStatusRequestTimeout() *StatusRequestTimeout {
	return &StatusRequestTimeout{StatusCode(http.StatusRequestTimeout)}
}

// StatusConflict ⤳ https://httpstatuses.com/409
type StatusConflict struct{ StatusCode }

// NewStatusConflict ⤳ https://httpstatuses.com/409
func NewStatusConflict() *StatusConflict {
	return &StatusConflict{StatusCode(http.StatusConflict)}
}

// StatusGone ⤳ https://httpstatuses.com/410
type StatusGone struct{ StatusCode }

// NewStatusGone ⤳ https://httpstatuses.com/410
func NewStatusGone() *StatusGone {
	return &StatusGone{StatusCode(http.StatusGone)}
}

// StatusLengthRequired ⤳ https://httpstatuses.com/411
type StatusLengthRequired struct{ StatusCode }

// NewStatusLengthRequired ⤳ https://httpstatuses.com/411
func NewStatusLengthRequired() *StatusLengthRequired {
	return &StatusLengthRequired{StatusCode(http.StatusLengthRequired)}
}

// StatusPreconditionFailed ⤳ https://httpstatuses.com/412
type StatusPreconditionFailed struct{ StatusCode }

// NewStatusPreconditionFailed ⤳ https://httpstatuses.com/412
func NewStatusPreconditionFailed() *StatusPreconditionFailed {
	return &StatusPreconditionFailed{StatusCode(http.StatusPreconditionFailed)}
}

// StatusRequestEntityTooLarge ⤳ https://httpstatuses.com/413
type StatusRequestEntityTooLarge struct{ StatusCode }

// NewStatusRequestEntityTooLarge ⤳ https://httpstatuses.com/413
func NewStatusRequestEntityTooLarge() *StatusRequestEntityTooLarge {
	return &StatusRequestEntityTooLarge{StatusCode(http.StatusRequestEntityTooLarge)}
}

// StatusRequestURITooLong ⤳ https://httpstatuses.com/414
type StatusRequestURITooLong struct{ StatusCode }

// NewStatusRequestURITooLong ⤳ https://httpstatuses.com/414
func NewStatusRequestURITooLong() *StatusRequestURITooLong {
	return &StatusRequestURITooLong{StatusCode(http.StatusRequestURITooLong)}
}

// StatusUnsupportedMediaType ⤳ https://httpstatuses.com/415
type StatusUnsupportedMediaType struct{ StatusCode }

// NewStatusUnsupportedMediaType ⤳ https://httpstatuses.com/415
func NewStatusUnsupportedMediaType() *StatusUnsupportedMediaType {
	return &StatusUnsupportedMediaType{StatusCode(http.StatusUnsupportedMediaType)}
}

// StatusRequestedRangeNotSatisfiable ⤳ https://httpstatuses.com/416
type StatusRequestedRangeNotSatisfiable struct{ StatusCode }

// NewStatusRequestedRangeNotSatisfiable ⤳ https://httpstatuses.com/416
func NewStatusRequestedRangeNotSatisfiable() *StatusRequestedRangeNotSatisfiable {
	return &StatusRequestedRangeNotSatisfiable{StatusCode(http.StatusRequestedRangeNotSatisfiable)}
}

// StatusExpectationFailed ⤳ https://httpstatuses.com/417
type StatusExpectationFailed struct{ StatusCode }

// NewStatusExpectationFailed ⤳ https://httpstatuses.com/417
func NewStatusExpectationFailed() *StatusExpectationFailed {
	return &StatusExpectationFailed{StatusCode(http.StatusExpectationFailed)}
}

// StatusTeapot ⤳ https://httpstatuses.com/418
type StatusTeapot struct{ StatusCode }

// NewStatusTeapot ⤳ https://httpstatuses.com/418
func NewStatusTeapot() *StatusTeapot {
	return &StatusTeapot{StatusCode(http.StatusTeapot)}
}

// StatusMisdirectedRequest ⤳ https://httpstatuses.com/421
type StatusMisdirectedRequest struct{ StatusCode }

// NewStatusMisdirectedRequest ⤳ https://httpstatuses.com/421
func NewStatusMisdirectedRequest() *StatusMisdirectedRequest {
	return &StatusMisdirectedRequest{StatusCode(http.StatusMisdirectedRequest)}
}

// StatusUnprocessableEntity ⤳ https://httpstatuses.com/422
type StatusUnprocessableEntity struct{ StatusCode }

// NewStatusUnprocessableEntity ⤳ https://httpstatuses.com/422
func NewStatusUnprocessableEntity() *StatusUnprocessableEntity {
	return &StatusUnprocessableEntity{StatusCode(http.StatusUnprocessableEntity)}
}

// StatusLocked ⤳ https://httpstatuses.com/423
type StatusLocked struct{ StatusCode }

// NewStatusLocked ⤳ https://httpstatuses.com/423
func NewStatusLocked() *StatusLocked {
	return &StatusLocked{StatusCode(http.StatusLocked)}
}

// StatusFailedDependency ⤳ https://httpstatuses.com/424
type StatusFailedDependency struct{ StatusCode }

// NewStatusFailedDependency ⤳ https://httpstatuses.com/424
func NewStatusFailedDependency() *StatusFailedDependency {
	return &StatusFailedDependency{StatusCode(http.StatusFailedDependency)}
}

// StatusTooEarly ⤳ https://httpstatuses.com/425
type StatusTooEarly struct{ StatusCode }

// NewStatusTooEarly ⤳ https://httpstatuses.com/425
func NewStatusTooEarly() *StatusTooEarly {
	return &StatusTooEarly{StatusCode(http.StatusTooEarly)}
}

// StatusUpgradeRequired ⤳ https://httpstatuses.com/426
type StatusUpgradeRequired struct{ StatusCode }

// NewStatusUpgradeRequired ⤳ https://httpstatuses.com/426
func NewStatusUpgradeRequired() *StatusUpgradeRequired {
	return &StatusUpgradeRequired{StatusCode(http.StatusUpgradeRequired)}
}

// StatusPreconditionRequired ⤳ https://httpstatuses.com/428
type StatusPreconditionRequired struct{ StatusCode }

// NewStatusPreconditionRequired ⤳ https://httpstatuses.com/428
func NewStatusPreconditionRequired() *StatusPreconditionRequired {
	return &StatusPreconditionRequired{StatusCode(http.StatusPreconditionRequired)}
}

// StatusTooManyRequests ⤳ https://httpstatuses.com/429
type StatusTooManyRequests struct{ StatusCode }

// NewStatusTooManyRequests ⤳ https://httpstatuses.com/429
func NewStatusTooManyRequests() *StatusTooManyRequests {
	return &StatusTooManyRequests{StatusCode(http.StatusTooManyRequests)}
}

// StatusRequestHeaderFieldsTooLarge ⤳ https://httpstatuses.com/431
type StatusRequestHeaderFieldsTooLarge struct{ StatusCode }

// NewStatusRequestHeaderFieldsTooLarge ⤳ https://httpstatuses.com/431
func NewStatusRequestHeaderFieldsTooLarge() *StatusRequestHeaderFieldsTooLarge {
	return &StatusRequestHeaderFieldsTooLarge{StatusCode(http.StatusRequestHeaderFieldsTooLarge)}
}

// StatusUnavailableForLegalReasons ⤳ https://httpstatuses.com/451
type StatusUnavailableForLegalReasons struct{ StatusCode }

// NewStatusUnavailableForLegalReasons ⤳ https://httpstatuses.com/451
func NewStatusUnavailableForLegalReasons() *StatusUnavailableForLegalReasons {
	return &StatusUnavailableForLegalReasons{StatusCode(http.StatusUnavailableForLegalReasons)}
}

//
//
// StatusInternalServerError ⤳ https://httpstatuses.com/500
type StatusInternalServerError struct{ StatusCode }

// NewStatusInternalServerError ⤳ https://httpstatuses.com/500
func NewStatusInternalServerError() *StatusInternalServerError {
	return &StatusInternalServerError{StatusCode(http.StatusInternalServerError)}
}

// StatusNotImplemented ⤳ https://httpstatuses.com/501
type StatusNotImplemented struct{ StatusCode }

// NewStatusNotImplemented ⤳ https://httpstatuses.com/501
func NewStatusNotImplemented() *StatusNotImplemented {
	return &StatusNotImplemented{StatusCode(http.StatusNotImplemented)}
}

// StatusBadGateway ⤳ https://httpstatuses.com/502
type StatusBadGateway struct{ StatusCode }

// NewStatusBadGateway ⤳ https://httpstatuses.com/502
func NewStatusBadGateway() *StatusBadGateway {
	return &StatusBadGateway{StatusCode(http.StatusBadGateway)}
}

// StatusServiceUnavailable ⤳ https://httpstatuses.com/503
type StatusServiceUnavailable struct{ StatusCode }

// NewStatusServiceUnavailable ⤳ https://httpstatuses.com/503
func NewStatusServiceUnavailable() *StatusServiceUnavailable {
	return &StatusServiceUnavailable{StatusCode(http.StatusServiceUnavailable)}
}

// StatusGatewayTimeout ⤳ https://httpstatuses.com/504
type StatusGatewayTimeout struct{ StatusCode }

// NewStatusGatewayTimeout ⤳ https://httpstatuses.com/504
func NewStatusGatewayTimeout() *StatusGatewayTimeout {
	return &StatusGatewayTimeout{StatusCode(http.StatusGatewayTimeout)}
}

// StatusHTTPVersionNotSupported ⤳ https://httpstatuses.com/505
type StatusHTTPVersionNotSupported struct{ StatusCode }

// NewStatusHTTPVersionNotSupported ⤳ https://httpstatuses.com/505
func NewStatusHTTPVersionNotSupported() *StatusHTTPVersionNotSupported {
	return &StatusHTTPVersionNotSupported{StatusCode(http.StatusHTTPVersionNotSupported)}
}

// StatusVariantAlsoNegotiates ⤳ https://httpstatuses.com/506
type StatusVariantAlsoNegotiates struct{ StatusCode }

// NewStatusVariantAlsoNegotiates ⤳ https://httpstatuses.com/506
func NewStatusVariantAlsoNegotiates() *StatusVariantAlsoNegotiates {
	return &StatusVariantAlsoNegotiates{StatusCode(http.StatusVariantAlsoNegotiates)}
}

// StatusInsufficientStorage ⤳ https://httpstatuses.com/507
type StatusInsufficientStorage struct{ StatusCode }

// NewStatusInsufficientStorage ⤳ https://httpstatuses.com/507
func NewStatusInsufficientStorage() *StatusInsufficientStorage {
	return &StatusInsufficientStorage{StatusCode(http.StatusInsufficientStorage)}
}

// StatusLoopDetected ⤳ https://httpstatuses.com/508
type StatusLoopDetected struct{ StatusCode }

// NewStatusLoopDetected ⤳ https://httpstatuses.com/508
func NewStatusLoopDetected() *StatusLoopDetected {
	return &StatusLoopDetected{StatusCode(http.StatusLoopDetected)}
}

// StatusNotExtended ⤳ https://httpstatuses.com/510
type StatusNotExtended struct{ StatusCode }

// NewStatusNotExtended ⤳ https://httpstatuses.com/510
func NewStatusNotExtended() *StatusNotExtended {
	return &StatusNotExtended{StatusCode(http.StatusNotExtended)}
}

// StatusNetworkAuthenticationRequired ⤳ https://httpstatuses.com/511
type StatusNetworkAuthenticationRequired struct{ StatusCode }

// NewStatusNetworkAuthenticationRequired ⤳ https://httpstatuses.com/511
func NewStatusNetworkAuthenticationRequired() *StatusNetworkAuthenticationRequired {
	return &StatusNetworkAuthenticationRequired{StatusCode(http.StatusNetworkAuthenticationRequired)}
}

// NewStatusCode transforms integer codes to types
func NewStatusCode(code int, required ...StatusCodeAny) StatusCodeAny {
	await := 0
	if len(required) > 0 {
		await = required[0].Value()
	}
	status := mkStatusCode(code, await)

	if codec, ok := decoder[code]; ok {
		return codec(status)
	}

	return &StatusUnknown{status}
}

var decoder = map[int]func(StatusCode) StatusCodeAny{
	// 1xx
	http.StatusContinue:           func(status StatusCode) StatusCodeAny { return &StatusContinue{status} },
	http.StatusSwitchingProtocols: func(status StatusCode) StatusCodeAny { return &StatusSwitchingProtocols{status} },
	http.StatusProcessing:         func(status StatusCode) StatusCodeAny { return &StatusProcessing{status} },
	http.StatusEarlyHints:         func(status StatusCode) StatusCodeAny { return &StatusEarlyHints{status} },

	// 2xx
	http.StatusOK:                   func(status StatusCode) StatusCodeAny { return &StatusOK{status} },
	http.StatusCreated:              func(status StatusCode) StatusCodeAny { return &StatusCreated{status} },
	http.StatusAccepted:             func(status StatusCode) StatusCodeAny { return &StatusAccepted{status} },
	http.StatusNonAuthoritativeInfo: func(status StatusCode) StatusCodeAny { return &StatusNonAuthoritativeInfo{status} },
	http.StatusNoContent:            func(status StatusCode) StatusCodeAny { return &StatusNoContent{status} },
	http.StatusResetContent:         func(status StatusCode) StatusCodeAny { return &StatusResetContent{status} },
	http.StatusPartialContent:       func(status StatusCode) StatusCodeAny { return &StatusPartialContent{status} },
	http.StatusMultiStatus:          func(status StatusCode) StatusCodeAny { return &StatusMultiStatus{status} },
	http.StatusAlreadyReported:      func(status StatusCode) StatusCodeAny { return &StatusAlreadyReported{status} },
	http.StatusIMUsed:               func(status StatusCode) StatusCodeAny { return &StatusIMUsed{status} },

	// 3xx
	http.StatusMultipleChoices:   func(status StatusCode) StatusCodeAny { return &StatusMultipleChoices{status} },
	http.StatusMovedPermanently:  func(status StatusCode) StatusCodeAny { return &StatusMovedPermanently{status} },
	http.StatusFound:             func(status StatusCode) StatusCodeAny { return &StatusFound{status} },
	http.StatusSeeOther:          func(status StatusCode) StatusCodeAny { return &StatusSeeOther{status} },
	http.StatusNotModified:       func(status StatusCode) StatusCodeAny { return &StatusNotModified{status} },
	http.StatusUseProxy:          func(status StatusCode) StatusCodeAny { return &StatusUseProxy{status} },
	http.StatusTemporaryRedirect: func(status StatusCode) StatusCodeAny { return &StatusTemporaryRedirect{status} },
	http.StatusPermanentRedirect: func(status StatusCode) StatusCodeAny { return &StatusPermanentRedirect{status} },

	// 4xx
	http.StatusBadRequest:                   func(status StatusCode) StatusCodeAny { return &StatusBadRequest{status} },
	http.StatusUnauthorized:                 func(status StatusCode) StatusCodeAny { return &StatusUnauthorized{status} },
	http.StatusPaymentRequired:              func(status StatusCode) StatusCodeAny { return &StatusPaymentRequired{status} },
	http.StatusForbidden:                    func(status StatusCode) StatusCodeAny { return &StatusForbidden{status} },
	http.StatusNotFound:                     func(status StatusCode) StatusCodeAny { return &StatusNotFound{status} },
	http.StatusMethodNotAllowed:             func(status StatusCode) StatusCodeAny { return &StatusMethodNotAllowed{status} },
	http.StatusNotAcceptable:                func(status StatusCode) StatusCodeAny { return &StatusNotAcceptable{status} },
	http.StatusProxyAuthRequired:            func(status StatusCode) StatusCodeAny { return &StatusProxyAuthRequired{status} },
	http.StatusRequestTimeout:               func(status StatusCode) StatusCodeAny { return &StatusRequestTimeout{status} },
	http.StatusConflict:                     func(status StatusCode) StatusCodeAny { return &StatusConflict{status} },
	http.StatusGone:                         func(status StatusCode) StatusCodeAny { return &StatusGone{status} },
	http.StatusLengthRequired:               func(status StatusCode) StatusCodeAny { return &StatusLengthRequired{status} },
	http.StatusPreconditionFailed:           func(status StatusCode) StatusCodeAny { return &StatusPreconditionFailed{status} },
	http.StatusRequestEntityTooLarge:        func(status StatusCode) StatusCodeAny { return &StatusRequestEntityTooLarge{status} },
	http.StatusRequestURITooLong:            func(status StatusCode) StatusCodeAny { return &StatusRequestURITooLong{status} },
	http.StatusUnsupportedMediaType:         func(status StatusCode) StatusCodeAny { return &StatusUnsupportedMediaType{status} },
	http.StatusRequestedRangeNotSatisfiable: func(status StatusCode) StatusCodeAny { return &StatusRequestedRangeNotSatisfiable{status} },
	http.StatusExpectationFailed:            func(status StatusCode) StatusCodeAny { return &StatusExpectationFailed{status} },
	http.StatusTeapot:                       func(status StatusCode) StatusCodeAny { return &StatusTeapot{status} },
	http.StatusMisdirectedRequest:           func(status StatusCode) StatusCodeAny { return &StatusMisdirectedRequest{status} },
	http.StatusUnprocessableEntity:          func(status StatusCode) StatusCodeAny { return &StatusUnprocessableEntity{status} },
	http.StatusLocked:                       func(status StatusCode) StatusCodeAny { return &StatusLocked{status} },
	http.StatusFailedDependency:             func(status StatusCode) StatusCodeAny { return &StatusFailedDependency{status} },
	http.StatusTooEarly:                     func(status StatusCode) StatusCodeAny { return &StatusTooEarly{status} },
	http.StatusUpgradeRequired:              func(status StatusCode) StatusCodeAny { return &StatusUpgradeRequired{status} },
	http.StatusPreconditionRequired:         func(status StatusCode) StatusCodeAny { return &StatusPreconditionRequired{status} },
	http.StatusTooManyRequests:              func(status StatusCode) StatusCodeAny { return &StatusTooManyRequests{status} },
	http.StatusRequestHeaderFieldsTooLarge:  func(status StatusCode) StatusCodeAny { return &StatusRequestHeaderFieldsTooLarge{status} },
	http.StatusUnavailableForLegalReasons:   func(status StatusCode) StatusCodeAny { return &StatusUnavailableForLegalReasons{status} },

	// 5xx
	http.StatusInternalServerError:           func(status StatusCode) StatusCodeAny { return &StatusInternalServerError{status} },
	http.StatusNotImplemented:                func(status StatusCode) StatusCodeAny { return &StatusNotImplemented{status} },
	http.StatusBadGateway:                    func(status StatusCode) StatusCodeAny { return &StatusBadGateway{status} },
	http.StatusServiceUnavailable:            func(status StatusCode) StatusCodeAny { return &StatusServiceUnavailable{status} },
	http.StatusGatewayTimeout:                func(status StatusCode) StatusCodeAny { return &StatusGatewayTimeout{status} },
	http.StatusHTTPVersionNotSupported:       func(status StatusCode) StatusCodeAny { return &StatusHTTPVersionNotSupported{status} },
	http.StatusVariantAlsoNegotiates:         func(status StatusCode) StatusCodeAny { return &StatusVariantAlsoNegotiates{status} },
	http.StatusInsufficientStorage:           func(status StatusCode) StatusCodeAny { return &StatusInsufficientStorage{status} },
	http.StatusLoopDetected:                  func(status StatusCode) StatusCodeAny { return &StatusLoopDetected{status} },
	http.StatusNotExtended:                   func(status StatusCode) StatusCodeAny { return &StatusNotExtended{status} },
	http.StatusNetworkAuthenticationRequired: func(status StatusCode) StatusCodeAny { return &StatusNetworkAuthenticationRequired{status} },
}

// Undefined is returned by api if expectation at body value is failed
type Undefined struct {
	Type string
}

func (e *Undefined) Error() string {
	return fmt.Sprintf("Value of type %v is not defined.", e.Type)
}

// Mismatch is returned by api if expectation at body value is failed
type Mismatch struct {
	Diff    string
	Payload interface{}
}

func (e *Mismatch) Error() string {
	return e.Diff
}
