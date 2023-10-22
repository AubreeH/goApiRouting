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
)
