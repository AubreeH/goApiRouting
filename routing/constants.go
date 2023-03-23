package routing

const (
	// JSONResponse is used in Response.Type to instruct the router to use json.Marshal() to format the body.
	JSONResponse = 0
	// HTMLResponse is used in Response.Type to instruct the router to return the body as a plain string. Will panic if body is not a string.
	HTMLResponse = 1
	// XMLResponse is used in Response.Type to instruct the router to use xml.Marshal() to format the body.
	XMLResponse = 2
)
