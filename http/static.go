//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http

var (
	//
	//
	// StatusCodeContinue ⤳ https://httpstatuses.com/100
	StatusCodeContinue *StatusContinue = NewStatusContinue()
	// StatusCodeSwitchingProtocols ⤳ https://httpstatuses.com/101
	StatusCodeSwitchingProtocols *StatusSwitchingProtocols = NewStatusSwitchingProtocols()
	// StatusCodeProcessing ⤳ https://httpstatuses.com/102
	StatusCodeProcessing *StatusProcessing = NewStatusProcessing()
	// StatusCodeEarlyHints ⤳ https://httpstatuses.com/103
	StatusCodeEarlyHints *StatusEarlyHints = NewStatusEarlyHints()

	//
	//
	// StatusCodeOK ⤳ https://httpstatuses.com/200
	StatusCodeOK *StatusOK = NewStatusOK()
	// StatusCodeCreated ⤳ https://httpstatuses.com/201
	StatusCodeCreated *StatusCreated = NewStatusCreated()
	// StatusCodeAccepted ⤳ https://httpstatuses.com/202
	StatusCodeAccepted *StatusAccepted = NewStatusAccepted()
	// StatusCodeNonAuthoritativeInfo ⤳ https://httpstatuses.com/203
	StatusCodeNonAuthoritativeInfo *StatusNonAuthoritativeInfo = NewStatusNonAuthoritativeInfo()
	// StatusCodeNoContent ⤳ https://httpstatuses.com/204
	StatusCodeNoContent *StatusNoContent = NewStatusNoContent()
	// StatusCodeResetContent ⤳ https://httpstatuses.com/205
	StatusCodeResetContent *StatusResetContent = NewStatusResetContent()
	// StatusCodePartialContent ⤳ https://httpstatuses.com/206
	StatusCodePartialContent *StatusPartialContent = NewStatusPartialContent()
	// StatusCodeMultiStatus ⤳ https://httpstatuses.com/207
	StatusCodeMultiStatus *StatusMultiStatus = NewStatusMultiStatus()
	// StatusCodeAlreadyReported ⤳ https://httpstatuses.com/208
	StatusCodeAlreadyReported *StatusAlreadyReported = NewStatusAlreadyReported()
	// StatusCodeIMUsed ⤳ https://httpstatuses.com/226
	StatusCodeIMUsed *StatusIMUsed = NewStatusIMUsed()

	//
	//
	// StatusCodeMultipleChoices ⤳ https://httpstatuses.com/300
	StatusCodeMultipleChoices *StatusMultipleChoices = NewStatusMultipleChoices()
	// StatusCodeMovedPermanently ⤳ https://httpstatuses.com/301
	StatusCodeMovedPermanently *StatusMovedPermanently = NewStatusMovedPermanently()
	// StatusCodeFound ⤳ https://httpstatuses.com/302
	StatusCodeFound *StatusFound = NewStatusFound()
	// StatusCodeSeeOther ⤳ https://httpstatuses.com/303
	StatusCodeSeeOther *StatusSeeOther = NewStatusSeeOther()
	// StatusCodeNotModified ⤳ https://httpstatuses.com/304
	StatusCodeNotModified *StatusNotModified = NewStatusNotModified()
	// StatusCodeUseProxy ⤳ https://httpstatuses.com/305
	StatusCodeUseProxy *StatusUseProxy = NewStatusUseProxy()
	// StatusCodeTemporaryRedirect ⤳ https://httpstatuses.com/307
	StatusCodeTemporaryRedirect *StatusTemporaryRedirect = NewStatusTemporaryRedirect()
	// StatusCodePermanentRedirect ⤳ https://httpstatuses.com/308
	StatusCodePermanentRedirect *StatusPermanentRedirect = NewStatusPermanentRedirect()

	//
	//
	// StatusCodeBadRequest ⤳ https://httpstatuses.com/400
	StatusCodeBadRequest *StatusBadRequest = NewStatusBadRequest()
	// StatusCodeUnauthorized ⤳ https://httpstatuses.com/401
	StatusCodeUnauthorized *StatusUnauthorized = NewStatusUnauthorized()
	// StatusCodePaymentRequired ⤳ https://httpstatuses.com/402
	StatusCodePaymentRequired *StatusPaymentRequired = NewStatusPaymentRequired()
	// StatusCodeForbidden ⤳ https://httpstatuses.com/403
	StatusCodeForbidden *StatusForbidden = NewStatusForbidden()
	// StatusCodeNotFound ⤳ https://httpstatuses.com/404
	StatusCodeNotFound *StatusNotFound = NewStatusNotFound()
	// StatusCodeMethodNotAllowed ⤳ https://httpstatuses.com/405
	StatusCodeMethodNotAllowed *StatusMethodNotAllowed = NewStatusMethodNotAllowed()
	// StatusCodeNotAcceptable ⤳ https://httpstatuses.com/406
	StatusCodeNotAcceptable *StatusNotAcceptable = NewStatusNotAcceptable()
	// StatusCodeProxyAuthRequired ⤳ https://httpstatuses.com/407
	StatusCodeProxyAuthRequired *StatusProxyAuthRequired = NewStatusProxyAuthRequired()
	// StatusCodeRequestTimeout ⤳ https://httpstatuses.com/408
	StatusCodeRequestTimeout *StatusRequestTimeout = NewStatusRequestTimeout()
	// StatusCodeConflict ⤳ https://httpstatuses.com/409
	StatusCodeConflict *StatusConflict = NewStatusConflict()
	// StatusCodeGone ⤳ https://httpstatuses.com/410
	StatusCodeGone *StatusGone = NewStatusGone()
	// StatusCodeLengthRequired ⤳ https://httpstatuses.com/411
	StatusCodeLengthRequired *StatusLengthRequired = NewStatusLengthRequired()
	// StatusCodePreconditionFailed ⤳ https://httpstatuses.com/412
	StatusCodePreconditionFailed *StatusPreconditionFailed = NewStatusPreconditionFailed()
	// StatusCodeRequestEntityTooLarge ⤳ https://httpstatuses.com/413
	StatusCodeRequestEntityTooLarge *StatusRequestEntityTooLarge = NewStatusRequestEntityTooLarge()
	// StatusCodeRequestURITooLong ⤳ https://httpstatuses.com/414
	StatusCodeRequestURITooLong *StatusRequestURITooLong = NewStatusRequestURITooLong()
	// StatusCodeUnsupportedMediaType ⤳ https://httpstatuses.com/415
	StatusCodeUnsupportedMediaType *StatusUnsupportedMediaType = NewStatusUnsupportedMediaType()
	// StatusCodeRequestedRangeNotSatisfiable ⤳ https://httpstatuses.com/416
	StatusCodeRequestedRangeNotSatisfiable *StatusRequestedRangeNotSatisfiable = NewStatusRequestedRangeNotSatisfiable()
	// StatusCodeExpectationFailed ⤳ https://httpstatuses.com/417
	StatusCodeExpectationFailed *StatusExpectationFailed = NewStatusExpectationFailed()
	// StatusCodeTeapot ⤳ https://httpstatuses.com/418
	StatusCodeTeapot *StatusTeapot = NewStatusTeapot()
	// StatusCodeMisdirectedRequest  ⤳ https://httpstatuses.com/421
	StatusCodeMisdirectedRequest *StatusMisdirectedRequest = NewStatusMisdirectedRequest()
	// StatusCodeUnprocessableEntity  ⤳ https://httpstatuses.com/422
	StatusCodeUnprocessableEntity *StatusUnprocessableEntity = NewStatusUnprocessableEntity()
	// StatusCodeLocked  ⤳ https://httpstatuses.com/423
	StatusCodeLocked *StatusLocked = NewStatusLocked()
	// StatusCodeFailedDependency  ⤳ https://httpstatuses.com/424
	StatusCodeFailedDependency *StatusFailedDependency = NewStatusFailedDependency()
	// StatusCodeTooEarly  ⤳ https://httpstatuses.com/425
	StatusCodeTooEarly *StatusTooEarly = NewStatusTooEarly()
	// StatusCodeUpgradeRequired  ⤳ https://httpstatuses.com/426
	StatusCodeUpgradeRequired *StatusUpgradeRequired = NewStatusUpgradeRequired()
	// StatusCodePreconditionRequired  ⤳ https://httpstatuses.com/428
	StatusCodePreconditionRequired *StatusPreconditionRequired = NewStatusPreconditionRequired()
	// StatusCodeTooManyRequests  ⤳ https://httpstatuses.com/429
	StatusCodeTooManyRequests *StatusTooManyRequests = NewStatusTooManyRequests()
	// StatusCodeRequestHeaderFieldsTooLarge  ⤳ https://httpstatuses.com/431
	StatusCodeRequestHeaderFieldsTooLarge *StatusRequestHeaderFieldsTooLarge = NewStatusRequestHeaderFieldsTooLarge()
	// StatusCodeUnavailableForLegalReasons  ⤳ https://httpstatuses.com/451
	StatusCodeUnavailableForLegalReasons *StatusUnavailableForLegalReasons = NewStatusUnavailableForLegalReasons()

	//
	//
	// StatusCodeInternalServerError ⤳ https://httpstatuses.com/500
	StatusCodeInternalServerError *StatusInternalServerError = NewStatusInternalServerError()
	// StatusCodeNotImplemented ⤳ https://httpstatuses.com/501
	StatusCodeNotImplemented *StatusNotImplemented = NewStatusNotImplemented()
	// StatusCodeBadGateway ⤳ https://httpstatuses.com/502
	StatusCodeBadGateway *StatusBadGateway = NewStatusBadGateway()
	// StatusCodeServiceUnavailable ⤳ https://httpstatuses.com/503
	StatusCodeServiceUnavailable *StatusServiceUnavailable = NewStatusServiceUnavailable()
	// StatusCodeGatewayTimeout ⤳ https://httpstatuses.com/504
	StatusCodeGatewayTimeout *StatusGatewayTimeout = NewStatusGatewayTimeout()
	// StatusCodeHTTPVersionNotSupported ⤳ https://httpstatuses.com/505
	StatusCodeHTTPVersionNotSupported *StatusHTTPVersionNotSupported = NewStatusHTTPVersionNotSupported()
	// StatusCodeVariantAlsoNegotiates ⤳ https://httpstatuses.com/506
	StatusCodeVariantAlsoNegotiates *StatusVariantAlsoNegotiates = NewStatusVariantAlsoNegotiates()
	// StatusCodeInsufficientStorage ⤳ https://httpstatuses.com/507
	StatusCodeInsufficientStorage *StatusInsufficientStorage = NewStatusInsufficientStorage()
	// StatusCodeLoopDetected ⤳ https://httpstatuses.com/508
	StatusCodeLoopDetected *StatusLoopDetected = NewStatusLoopDetected()
	// StatusCodeNotExtended ⤳ https://httpstatuses.com/510
	StatusCodeNotExtended *StatusNotExtended = NewStatusNotExtended()
	// StatusCodeNetworkAuthenticationRequired ⤳ https://httpstatuses.com/511
	StatusCodeNetworkAuthenticationRequired *StatusNetworkAuthenticationRequired = NewStatusNetworkAuthenticationRequired()
)
