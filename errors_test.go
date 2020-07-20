//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https: gurl.StatusCode,
//

package gurl_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/fogfish/gurl"
	"github.com/fogfish/it"
)

var httpStatusCode = map[int]gurl.StatusCodeAny{
	// 1xx
	http.StatusContinue:           gurl.StatusCodeContinue,
	http.StatusSwitchingProtocols: gurl.StatusCodeSwitchingProtocols,
	http.StatusProcessing:         gurl.StatusCodeProcessing,
	http.StatusEarlyHints:         gurl.StatusCodeEarlyHints,
	// 2xx
	http.StatusOK:                   gurl.StatusCodeOK,
	http.StatusCreated:              gurl.StatusCodeCreated,
	http.StatusAccepted:             gurl.StatusCodeAccepted,
	http.StatusNonAuthoritativeInfo: gurl.StatusCodeNonAuthoritativeInfo,
	http.StatusNoContent:            gurl.StatusCodeNoContent,
	http.StatusResetContent:         gurl.StatusCodeResetContent,
	http.StatusPartialContent:       gurl.StatusCodePartialContent,
	http.StatusMultiStatus:          gurl.StatusCodeMultiStatus,
	http.StatusAlreadyReported:      gurl.StatusCodeAlreadyReported,
	http.StatusIMUsed:               gurl.StatusCodeIMUsed,
	// 3xx
	http.StatusMultipleChoices:   gurl.StatusCodeMultipleChoices,
	http.StatusMovedPermanently:  gurl.StatusCodeMovedPermanently,
	http.StatusFound:             gurl.StatusCodeFound,
	http.StatusSeeOther:          gurl.StatusCodeSeeOther,
	http.StatusNotModified:       gurl.StatusCodeNotModified,
	http.StatusUseProxy:          gurl.StatusCodeUseProxy,
	http.StatusTemporaryRedirect: gurl.StatusCodeTemporaryRedirect,
	http.StatusPermanentRedirect: gurl.StatusCodePermanentRedirect,
	// 4xx
	http.StatusBadRequest:                   gurl.StatusCodeBadRequest,
	http.StatusUnauthorized:                 gurl.StatusCodeUnauthorized,
	http.StatusPaymentRequired:              gurl.StatusCodePaymentRequired,
	http.StatusForbidden:                    gurl.StatusCodeForbidden,
	http.StatusNotFound:                     gurl.StatusCodeNotFound,
	http.StatusMethodNotAllowed:             gurl.StatusCodeMethodNotAllowed,
	http.StatusNotAcceptable:                gurl.StatusCodeNotAcceptable,
	http.StatusProxyAuthRequired:            gurl.StatusCodeProxyAuthRequired,
	http.StatusRequestTimeout:               gurl.StatusCodeRequestTimeout,
	http.StatusConflict:                     gurl.StatusCodeConflict,
	http.StatusGone:                         gurl.StatusCodeGone,
	http.StatusLengthRequired:               gurl.StatusCodeLengthRequired,
	http.StatusPreconditionFailed:           gurl.StatusCodePreconditionFailed,
	http.StatusRequestEntityTooLarge:        gurl.StatusCodeRequestEntityTooLarge,
	http.StatusRequestURITooLong:            gurl.StatusCodeRequestURITooLong,
	http.StatusUnsupportedMediaType:         gurl.StatusCodeUnsupportedMediaType,
	http.StatusRequestedRangeNotSatisfiable: gurl.StatusCodeRequestedRangeNotSatisfiable,
	http.StatusExpectationFailed:            gurl.StatusCodeExpectationFailed,
	http.StatusTeapot:                       gurl.StatusCodeTeapot,
	http.StatusMisdirectedRequest:           gurl.StatusCodeMisdirectedRequest,
	http.StatusUnprocessableEntity:          gurl.StatusCodeUnprocessableEntity,
	http.StatusLocked:                       gurl.StatusCodeLocked,
	http.StatusFailedDependency:             gurl.StatusCodeFailedDependency,
	http.StatusTooEarly:                     gurl.StatusCodeTooEarly,
	http.StatusUpgradeRequired:              gurl.StatusCodeUpgradeRequired,
	http.StatusPreconditionRequired:         gurl.StatusCodePreconditionRequired,
	http.StatusTooManyRequests:              gurl.StatusCodeTooManyRequests,
	http.StatusRequestHeaderFieldsTooLarge:  gurl.StatusCodeRequestHeaderFieldsTooLarge,
	http.StatusUnavailableForLegalReasons:   gurl.StatusCodeUnavailableForLegalReasons,
	// 5xx
	http.StatusInternalServerError:           gurl.StatusCodeInternalServerError,
	http.StatusNotImplemented:                gurl.StatusCodeNotImplemented,
	http.StatusBadGateway:                    gurl.StatusCodeBadGateway,
	http.StatusServiceUnavailable:            gurl.StatusCodeServiceUnavailable,
	http.StatusGatewayTimeout:                gurl.StatusCodeGatewayTimeout,
	http.StatusHTTPVersionNotSupported:       gurl.StatusCodeHTTPVersionNotSupported,
	http.StatusVariantAlsoNegotiates:         gurl.StatusCodeVariantAlsoNegotiates,
	http.StatusInsufficientStorage:           gurl.StatusCodeInsufficientStorage,
	http.StatusLoopDetected:                  gurl.StatusCodeLoopDetected,
	http.StatusNotExtended:                   gurl.StatusCodeNotExtended,
	http.StatusNetworkAuthenticationRequired: gurl.StatusCodeNetworkAuthenticationRequired,
}

func TestStatusCodeCodec(t *testing.T) {
	for code, val := range httpStatusCode {
		status := gurl.NewStatusCode(code, gurl.StatusCodeOK)
		it.Ok(t).
			If(status.Value()).Should().Equal(code).
			If(status.Value()).Should().Equal(val.Value()).
			If(errors.Is(status, val)).Should().Equal(true).
			If(status.Await()).Should().Equal(http.StatusOK).
			If(code).Should().Equal(val.Value())
	}
}
