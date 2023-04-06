package routing

type ResponseType int

const (
	// JSONResponse is used in Response.Type to instruct the router to use json.Marshal() to format the body.
	JSONResponse ResponseType = 0
	// HTMLResponse is used in Response.Type to instruct the router to return the body as a plain string. Will panic if body is not a string.
	HTMLResponse ResponseType = 1
	// XMLResponse is used in Response.Type to instruct the router to use xml.Marshal() to format the body.
	XMLResponse ResponseType = 2
	// PlainTextResponse is used in Response.Type to instruct the router to return the body as plain text. Will panic if body is not a string.
	PlainTextResponse ResponseType = 3
	// FileResponse is used in Response.Type to server a file in the response. Provide the path to the file in the Response.Body. Will panic if Response.Body is not a string.
	FileResponse ResponseType = 4
)
