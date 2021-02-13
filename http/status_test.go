//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	gurl "github.com/fogfish/gurl/http"
	"github.com/fogfish/it"
)

var httpStatusCode = map[int]gurl.StatusCode{
	// 1xx
	http.StatusContinue:           gurl.StatusContinue,
	http.StatusSwitchingProtocols: gurl.StatusSwitchingProtocols,
	http.StatusProcessing:         gurl.StatusProcessing,
	http.StatusEarlyHints:         gurl.StatusEarlyHints,
	// 2xx
	http.StatusOK:                   gurl.StatusOK,
	http.StatusCreated:              gurl.StatusCreated,
	http.StatusAccepted:             gurl.StatusAccepted,
	http.StatusNonAuthoritativeInfo: gurl.StatusNonAuthoritativeInfo,
	http.StatusNoContent:            gurl.StatusNoContent,
	http.StatusResetContent:         gurl.StatusResetContent,
	http.StatusPartialContent:       gurl.StatusPartialContent,
	http.StatusMultiStatus:          gurl.StatusMultiStatus,
	http.StatusAlreadyReported:      gurl.StatusAlreadyReported,
	http.StatusIMUsed:               gurl.StatusIMUsed,
	// 3xx
	http.StatusMultipleChoices:   gurl.StatusMultipleChoices,
	http.StatusMovedPermanently:  gurl.StatusMovedPermanently,
	http.StatusFound:             gurl.StatusFound,
	http.StatusSeeOther:          gurl.StatusSeeOther,
	http.StatusNotModified:       gurl.StatusNotModified,
	http.StatusUseProxy:          gurl.StatusUseProxy,
	http.StatusTemporaryRedirect: gurl.StatusTemporaryRedirect,
	http.StatusPermanentRedirect: gurl.StatusPermanentRedirect,
	// 4xx
	http.StatusBadRequest:                   gurl.StatusBadRequest,
	http.StatusUnauthorized:                 gurl.StatusUnauthorized,
	http.StatusPaymentRequired:              gurl.StatusPaymentRequired,
	http.StatusForbidden:                    gurl.StatusForbidden,
	http.StatusNotFound:                     gurl.StatusNotFound,
	http.StatusMethodNotAllowed:             gurl.StatusMethodNotAllowed,
	http.StatusNotAcceptable:                gurl.StatusNotAcceptable,
	http.StatusProxyAuthRequired:            gurl.StatusProxyAuthRequired,
	http.StatusRequestTimeout:               gurl.StatusRequestTimeout,
	http.StatusConflict:                     gurl.StatusConflict,
	http.StatusGone:                         gurl.StatusGone,
	http.StatusLengthRequired:               gurl.StatusLengthRequired,
	http.StatusPreconditionFailed:           gurl.StatusPreconditionFailed,
	http.StatusRequestEntityTooLarge:        gurl.StatusRequestEntityTooLarge,
	http.StatusRequestURITooLong:            gurl.StatusRequestURITooLong,
	http.StatusUnsupportedMediaType:         gurl.StatusUnsupportedMediaType,
	http.StatusRequestedRangeNotSatisfiable: gurl.StatusRequestedRangeNotSatisfiable,
	http.StatusExpectationFailed:            gurl.StatusExpectationFailed,
	http.StatusTeapot:                       gurl.StatusTeapot,
	http.StatusMisdirectedRequest:           gurl.StatusMisdirectedRequest,
	http.StatusUnprocessableEntity:          gurl.StatusUnprocessableEntity,
	http.StatusLocked:                       gurl.StatusLocked,
	http.StatusFailedDependency:             gurl.StatusFailedDependency,
	http.StatusTooEarly:                     gurl.StatusTooEarly,
	http.StatusUpgradeRequired:              gurl.StatusUpgradeRequired,
	http.StatusPreconditionRequired:         gurl.StatusPreconditionRequired,
	http.StatusTooManyRequests:              gurl.StatusTooManyRequests,
	http.StatusRequestHeaderFieldsTooLarge:  gurl.StatusRequestHeaderFieldsTooLarge,
	http.StatusUnavailableForLegalReasons:   gurl.StatusUnavailableForLegalReasons,
	// 5xx
	http.StatusInternalServerError:           gurl.StatusInternalServerError,
	http.StatusNotImplemented:                gurl.StatusNotImplemented,
	http.StatusBadGateway:                    gurl.StatusBadGateway,
	http.StatusServiceUnavailable:            gurl.StatusServiceUnavailable,
	http.StatusGatewayTimeout:                gurl.StatusGatewayTimeout,
	http.StatusHTTPVersionNotSupported:       gurl.StatusHTTPVersionNotSupported,
	http.StatusVariantAlsoNegotiates:         gurl.StatusVariantAlsoNegotiates,
	http.StatusInsufficientStorage:           gurl.StatusInsufficientStorage,
	http.StatusLoopDetected:                  gurl.StatusLoopDetected,
	http.StatusNotExtended:                   gurl.StatusNotExtended,
	http.StatusNetworkAuthenticationRequired: gurl.StatusNetworkAuthenticationRequired,
}

func TestStatusCodeCodec(t *testing.T) {
	for code, val := range httpStatusCode {
		status := gurl.NewStatusCode(code, gurl.StatusOK)
		it.Ok(t).
			If(code).Should().Equal(val.Value()).
			If(status.Value()).Should().Equal(code).
			If(status.Value()).Should().Equal(val.Value()).
			If(errors.Is(status, val)).Should().Equal(true).
			If(fmt.Sprintf("%T", status)).Should().Equal(fmt.Sprintf("%T", val))
	}
}

func TestStatusCodeRequired(t *testing.T) {
	var code error = gurl.NewStatusCode(200, 201)
	it.Ok(t).
		If(code.Error()).Should().Equal("HTTP Status `200 OK`, required `201 Created`.")
}

func TestStatusCodeText(t *testing.T) {
	var code error = gurl.NewStatusCode(200)
	it.Ok(t).
		If(code.Error()).Should().Equal("HTTP 200 OK").
		If(errors.Is(code, gurl.StatusOK)).Should().Equal(true).
		If(errors.Is(code, gurl.StatusCreated)).Should().Equal(false).
		If(errors.Is(code, fmt.Errorf("some error"))).Should().Equal(false)
}
