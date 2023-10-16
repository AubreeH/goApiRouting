package routing

type ResponseType int

const (
	JSONResponse      ResponseType = 0 // JSONResponse is used in [Response.Type] to instruct the router to use [json.Marshal] to format the body.
	HTMLResponse      ResponseType = 1 // HTMLResponse is used in [Response.Type] to instruct the router to return the body as a plain string. Will panic if body is not a string.
	XMLResponse       ResponseType = 2 // XMLResponse is used in [Response.Type] to instruct the router to use [xml.Marshal]() to format the [Response.Body].
	PlainTextResponse ResponseType = 3 // PlainTextResponse is used in [Response.Type] to instruct the router to return the body as plain text. Will panic if [Response.Body] is not a string.
	FileResponse      ResponseType = 4 // FileResponse is used in [Response.Type] to server a file in the response. Provide the path to the file in the [Response.Body]. Will panic if [Response.Body] is not a string.
)
