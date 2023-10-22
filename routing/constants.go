package routing

type ResponseType string

const (
	// JSONResponse is used in [Response.Type] to instruct the router to use [json.Marshal] to format the body.
	//
	// [Response.Type]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Type
	// [json.Marshal]: https://pkg.go.dev/encoding/json#Marshal
	JSONResponse ResponseType = "application/json"

	// HTMLResponse is used in [Response.Type] to instruct the router to return the body as a plain string.
	//
	// [Response.Type]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Type
	HTMLResponse ResponseType = "text/html"

	// XMLResponse is used in [Response.Type] to instruct the router to use [encoding/xml.Marshal] to format the [Response.Body].
	//
	// [Response.Type]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Type
	// [json.Marshal]: https://pkg.go.dev/encoding/xml#Marshal
	XMLResponse ResponseType = "application/xml"

	// PlainTextResponse is used in [Response.Type] to instruct the router to return the body as plain text.
	//
	// [Response.Type]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Type
	PlainTextResponse ResponseType = "text/plain"

	// FileResponse is used in [Response.Type] to serve a file in the response using [http.ServeFile].
	//
	// Provide the path to the file in the [Response.Body].
	//
	// [Response.Type]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Type
	// [Response.Body]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Body
	// [http.ServeFile]: https://pkg.go.dev/net/http#ServeFile
	FileResponse ResponseType = "file"

	// RedirectResponse is used in [Response.Type] to redirect the user to another page.
	//
	// Provide the URL to redirect to in the [Response.Body].
	//
	// [Response.Type]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Type
	// [Response.Body]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Body
	RedirectResponse ResponseType = "redirect"

	// NoResponse is used in [Response.Type] to instruct the router to not write a response.
	//
	// [Response.Type]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Type
	NoResponse ResponseType = "none"

	// CustomResponse is used in [Response.Type] to instruct the router to use a custom function to write the response.
	//
	// Provide the function in the [Response.Body].
	//
	// [Response.Type]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Type
	// [Response.Body]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Body
	CustomResponse ResponseType = "custom"
)

const (
	StatusContinue           = 100 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/100
	StatusSwitchingProtocols = 101 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/101
	StatusProcessing         = 102 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/102
	StatusEarlyHints         = 103 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/103

	StatusOk                   = 200 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/200
	StatusCreated              = 201 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/201
	StatusAccepted             = 202 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/202
	StatusNonAuthoritativeInfo = 203 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/203
	StatusNoContent            = 204 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/204
	StatusResetContent         = 205 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/205
	StatusPartialContent       = 206 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/206
	StatusMultiStatus          = 207 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/207
	StatusAlreadyReported      = 208 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/208
	StatusIMUsed               = 226 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/226

	StatusMultipleChoices   = 300 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/300
	StatusMovedPermanently  = 301 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/301
	StatusFound             = 302 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/302
	StatusSeeOther          = 303 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/303
	StatusNotModified       = 304 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/304
	StatusTemporaryRedirect = 307 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/307
	StatusPermanentRedirect = 308 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/308

	StatusBadRequest                   = 400 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/400
	StatusUnauthorized                 = 401 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401
	StatusPaymentRequired              = 402 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/402
	StatusForbidden                    = 403 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/403
	StatusNotFound                     = 404 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/404
	StatusMethodNotAllowed             = 405 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/405
	StatusNotAcceptable                = 406 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/406
	StatusProxyAuthRequired            = 407 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/407
	StatusRequestTimeout               = 408 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/408
	StatusConflict                     = 409 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/409
	StatusGone                         = 410 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/410
	StatusLengthRequired               = 411 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/411
	StatusPreconditionFailed           = 412 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/412
	StatusRequestEntityTooLarge        = 413 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/413
	StatusRequestURITooLong            = 414 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/414
	StatusUnsupportedMediaType         = 415 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/415
	StatusRequestedRangeNotSatisfiable = 416 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/416
	StatusExpectationFailed            = 417 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/417
	StatusTeapot                       = 418 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/418
	StatusMisdirectedRequest           = 421 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/421
	StatusUnprocessableEntity          = 422 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/422
	StatusLocked                       = 423 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/423
	StatusFailedDependency             = 424 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/424
	StatusTooEarly                     = 425 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/425
	StatusUpgradeRequired              = 426 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/426
	StatusPreconditionRequired         = 428 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/428
	StatusTooManyRequests              = 429 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429
	StatusRequestHeaderFieldsTooLarge  = 431 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/431
	StatusUnavailableForLegalReasons   = 451 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/451

	StatusInternalServerError           = 500 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/500
	StatusNotImplemented                = 501 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/501
	StatusBadGateway                    = 502 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/502
	StatusServiceUnavailable            = 503 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/503
	StatusGatewayTimeout                = 504 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/504
	StatusHTTPVersionNotSupported       = 505 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/505
	StatusVariantAlsoNegotiates         = 506 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/506
	StatusInsufficientStorage           = 507 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/507
	StatusLoopDetected                  = 508 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/508
	StatusNotExtended                   = 510 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/510
	StatusNetworkAuthenticationRequired = 511 // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/511
)
