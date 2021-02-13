package http

import (
	"fmt"
	"net/http"
)

/*

StatusCode is a base type for typesafe HTTP status codes. The library advertises
a usage of "pattern-matching" on HTTP status handling, which helps developers to
catch mismatch of HTTP statuses along with other side-effect failures.

The final values of HTTP statuses embeds StatusCode type. It makes them to look
like a "sum-types" and be compatible with any other error (side effect failures)
within IO category. Use final type instances in the error handling routines.

Use type switch for error handling "branches"

  switch err := cat.Fail.(type) {
  case nil:
    // Nothing
  case StatusCode:
    switch err {
    case http.StatusOK:
      // HTTP 200 OK
    case http.StatusNotFound:
      // HTTP 404 NotFound
    default:
      // any other HTTP errors
    }
  default:
    // any other errors
  }

Conditional error handling on expected HTTP Status

  if errors.Is(cat.Fail, http.StatusNotFound) {
  }

Conditional error handling on any HTTP Status

  if _, ok := cat.Fail.(gurl.StatusCode); ok {
  }

*/
type StatusCode int

// NewStatusCode transforms integer codes to types
func NewStatusCode(code int, required ...StatusCode) StatusCode {
	req := 0
	if len(required) > 0 {
		req = required[0].Value()
	}
	return StatusCode((req << 16) | code)
}

// Error makes StatusCode to be error
func (e StatusCode) Error() string {
	status := e.Value()
	if req := e.Required(); req != 0 {
		return fmt.Sprintf("HTTP Status `%d %s`, required `%d %s`.",
			status, http.StatusText(status), req, http.StatusText(req))
	}
	return fmt.Sprintf("HTTP %d %s", status, http.StatusText(status))
}

// Is compares wrapped errors
func (e StatusCode) Is(err error) bool {
	if code, ok := err.(StatusCode); ok {
		return e.Value() == code.Value()
	}
	return false
}

// Value transforms StatusCode type to integer value: StatusCode ⟼ int
func (e StatusCode) Value() int {
	return int(e) & 0xffff
}

// Required returns (expected, required) HTTP status
func (e StatusCode) Required() int {
	return int(e) >> 16
}

//
const (
	//
	StatusContinue           = StatusCode(http.StatusContinue)
	StatusSwitchingProtocols = StatusCode(http.StatusSwitchingProtocols)
	StatusProcessing         = StatusCode(http.StatusProcessing)
	StatusEarlyHints         = StatusCode(http.StatusEarlyHints)

	//
	StatusOK                   = StatusCode(http.StatusOK)
	StatusCreated              = StatusCode(http.StatusCreated)
	StatusAccepted             = StatusCode(http.StatusAccepted)
	StatusNonAuthoritativeInfo = StatusCode(http.StatusNonAuthoritativeInfo)
	StatusNoContent            = StatusCode(http.StatusNoContent)
	StatusResetContent         = StatusCode(http.StatusResetContent)
	StatusPartialContent       = StatusCode(http.StatusPartialContent)
	StatusMultiStatus          = StatusCode(http.StatusMultiStatus)
	StatusAlreadyReported      = StatusCode(http.StatusAlreadyReported)
	StatusIMUsed               = StatusCode(http.StatusIMUsed)

	//
	StatusMultipleChoices   = StatusCode(http.StatusMultipleChoices)
	StatusMovedPermanently  = StatusCode(http.StatusMovedPermanently)
	StatusFound             = StatusCode(http.StatusFound)
	StatusSeeOther          = StatusCode(http.StatusSeeOther)
	StatusNotModified       = StatusCode(http.StatusNotModified)
	StatusUseProxy          = StatusCode(http.StatusUseProxy)
	StatusTemporaryRedirect = StatusCode(http.StatusTemporaryRedirect)
	StatusPermanentRedirect = StatusCode(http.StatusPermanentRedirect)

	//
	StatusBadRequest                   = StatusCode(http.StatusBadRequest)
	StatusUnauthorized                 = StatusCode(http.StatusUnauthorized)
	StatusPaymentRequired              = StatusCode(http.StatusPaymentRequired)
	StatusForbidden                    = StatusCode(http.StatusForbidden)
	StatusNotFound                     = StatusCode(http.StatusNotFound)
	StatusMethodNotAllowed             = StatusCode(http.StatusMethodNotAllowed)
	StatusNotAcceptable                = StatusCode(http.StatusNotAcceptable)
	StatusProxyAuthRequired            = StatusCode(http.StatusProxyAuthRequired)
	StatusRequestTimeout               = StatusCode(http.StatusRequestTimeout)
	StatusConflict                     = StatusCode(http.StatusConflict)
	StatusGone                         = StatusCode(http.StatusGone)
	StatusLengthRequired               = StatusCode(http.StatusLengthRequired)
	StatusPreconditionFailed           = StatusCode(http.StatusPreconditionFailed)
	StatusRequestEntityTooLarge        = StatusCode(http.StatusRequestEntityTooLarge)
	StatusRequestURITooLong            = StatusCode(http.StatusRequestURITooLong)
	StatusUnsupportedMediaType         = StatusCode(http.StatusUnsupportedMediaType)
	StatusRequestedRangeNotSatisfiable = StatusCode(http.StatusRequestedRangeNotSatisfiable)
	StatusExpectationFailed            = StatusCode(http.StatusExpectationFailed)
	StatusTeapot                       = StatusCode(http.StatusTeapot)
	StatusMisdirectedRequest           = StatusCode(http.StatusMisdirectedRequest)
	StatusUnprocessableEntity          = StatusCode(http.StatusUnprocessableEntity)
	StatusLocked                       = StatusCode(http.StatusLocked)
	StatusFailedDependency             = StatusCode(http.StatusFailedDependency)
	StatusTooEarly                     = StatusCode(http.StatusTooEarly)
	StatusUpgradeRequired              = StatusCode(http.StatusUpgradeRequired)
	StatusPreconditionRequired         = StatusCode(http.StatusPreconditionRequired)
	StatusTooManyRequests              = StatusCode(http.StatusTooManyRequests)
	StatusRequestHeaderFieldsTooLarge  = StatusCode(http.StatusRequestHeaderFieldsTooLarge)
	StatusUnavailableForLegalReasons   = StatusCode(http.StatusUnavailableForLegalReasons)

	//
	StatusInternalServerError           = StatusCode(http.StatusInternalServerError)
	StatusNotImplemented                = StatusCode(http.StatusNotImplemented)
	StatusBadGateway                    = StatusCode(http.StatusBadGateway)
	StatusServiceUnavailable            = StatusCode(http.StatusServiceUnavailable)
	StatusGatewayTimeout                = StatusCode(http.StatusGatewayTimeout)
	StatusHTTPVersionNotSupported       = StatusCode(http.StatusHTTPVersionNotSupported)
	StatusVariantAlsoNegotiates         = StatusCode(http.StatusVariantAlsoNegotiates)
	StatusInsufficientStorage           = StatusCode(http.StatusInsufficientStorage)
	StatusLoopDetected                  = StatusCode(http.StatusLoopDetected)
	StatusNotExtended                   = StatusCode(http.StatusNotExtended)
	StatusNetworkAuthenticationRequired = StatusCode(http.StatusNetworkAuthenticationRequired)
)
