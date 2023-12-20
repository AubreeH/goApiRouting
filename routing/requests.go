package routing

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"
)

var (
	methodNotSupportedResponse = Response{
		Body:   "Method not supported",
		Type:   PlainTextResponse,
		Status: 400,
	}
	notFoundResponse = Response{
		Body:   "Not Found",
		Type:   PlainTextResponse,
		Status: 404,
	}
	unexpectedErrorResponse = Response{
		Body:   "Unexpected Error",
		Type:   PlainTextResponse,
		Status: 500,
	}
	invalidRequestBodyResponse = Response{
		Body:   "Invalid Request Body",
		Type:   PlainTextResponse,
		Status: 400,
	}
)

func (r *Router) setupHandler() {
	if r.mux != nil {
		return
	}
	r.mux = http.NewServeMux()
	r.mux.HandleFunc("/", r.handleRequest)
}

func (r *Router) handleRequest(writer http.ResponseWriter, request *http.Request) {
	err := r.vetRequest(request, writer)
	if err != nil {
		return
	}

	var context *Context
	defer cleanup(context)

	request.Body = http.MaxBytesReader(writer, request.Body, r.config.MaxContentLength)

	handler, pathParameters, err := r.getFunc(request.Method, request.URL.Path)
	if err != nil {
		switch err.Error() {
		case "method not supported":
			r.writeResponse(writer, request, methodNotSupportedResponse)
		case "not found":
			r.writeResponse(writer, request, notFoundResponse)
		default:
			r.writeResponse(writer, request, unexpectedErrorResponse)
		}
	} else if requestBody, err := getRequestBody(request); err != nil {
		r.writeResponse(writer, request, invalidRequestBodyResponse)
	} else {
		context = &Context{
			Writer:  writer,
			Request: request,
		}

		context.Store = &Store{
			pathParameters: pathParameters,
			query:          request.URL.Query(),
			body:           requestBody,
			store:          make(map[string]interface{}),
			contentType:    request.Header.Get("Content-Type"),
			context:        context,
		}

		handler(context, func(response Response) {
			r.writeResponse(writer, request, response)
		})
	}
}

func cleanup(c *Context) {
	if c != nil && c.Store != nil && c.Store.files != nil {
		for _, file := range c.Store.files {
			if file.saved && !file.stored {
				os.Remove(file.FilePath)
			}
		}
	}
}

func getRequestBody(request *http.Request) ([]byte, error) {
	if request.Body == nil {
		return nil, nil
	}

	var body []byte
	_, err := request.Body.Read(body)
	if err != nil && err.Error() == "EOF" {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *Router) writeResponse(writer http.ResponseWriter, request *http.Request, response Response) {
	r.writeHeaders(writer, response.Headers)

	// if response.Body == nil {
	// 	response.Body = ""
	// }

	var err error
	switch response.Type {
	case FileResponse:
		handleFileResponse(request, writer, response)
		return
	case RedirectResponse:
		err = handleRedirectResponse(request, writer, response)
	case NoResponse:
		return
	case CustomResponse:
		err = handleCustomResponse(request, writer, response)
	case JSONResponse:
		err = writeJSONResponse(request, writer, response)
	case HTMLResponse:
		err = writeHTMLResponse(request, writer, response)
	case XMLResponse:
		err = writeXMLResponse(request, writer, response)
	default: // case PlainTextResponse:
		err = writePlainTextResponse(request, writer, response)
	}

	if err != nil {
		fmt.Println(err)
	}
}

func (r *Router) writeHeaders(writer http.ResponseWriter, headers map[string]string) {
	for key, value := range r.config.BaseResponseHeaders {
		writer.Header().Set(key, value)
	}

	for key, value := range headers {
		writer.Header().Set(key, value)
	}
}

func write(contentType string, writer http.ResponseWriter, body []byte, response Response) error {
	writer.Header().Set("Content-Type", contentType)
	writer.WriteHeader(response.Status)
	_, err := writer.Write(body)
	return err
}

func writeJSONResponse(request *http.Request, writer http.ResponseWriter, response Response) error {
	var body []byte
	if response.Body == nil {
		body = []byte{}
	} else {
		str, err := json.Marshal(response.Body)
		if err != nil {
			return err
		}
		body = str
	}
	return write("application/json", writer, body, response)
}

func writeHTMLResponse(request *http.Request, writer http.ResponseWriter, response Response) error {
	if response.Template != nil {
		if err := response.Template.Execute(writer, response.Body); err != nil {
			return err
		}
		return nil
	}

	var body string
	if response.Body == nil {
		body = ""
	} else if str, ok := response.Body.(string); ok {
		body = str
	} else {
		return errors.New("body is not a string")
	}
	return write("text/html", writer, []byte(body), response)
}

func writeXMLResponse(request *http.Request, writer http.ResponseWriter, response Response) error {
	var body []byte
	if response.Body == nil {
		body = []byte{}
	} else {
		str, err := xml.Marshal(response.Body)
		if err != nil {
			return err
		}
		body = str
	}
	return write("application/xml", writer, body, response)
}

func writePlainTextResponse(request *http.Request, writer http.ResponseWriter, response Response) error {
	if response.Template != nil {
		if err := response.Template.Execute(writer, response.Body); err != nil {
			return err
		}
		return nil
	}

	var body string
	if response.Body == nil {
		body = ""
	} else if str, ok := response.Body.(string); ok {
		body = str
	} else {
		return errors.New("body is not a string")
	}
	return write("text/plain", writer, []byte(body), response)
}

func handleFileResponse(request *http.Request, writer http.ResponseWriter, response Response) {
	http.ServeFile(writer, request, response.Body.(string))
}

func handleRedirectResponse(request *http.Request, writer http.ResponseWriter, response Response) error {
	if str, ok := response.Body.(string); !ok {
		return fmt.Errorf("redirect body is not a string")
	} else {
		status := response.Status
		if status == 0 {
			status = StatusMovedPermanently
		}
		http.Redirect(writer, request, str, status)
	}
	return nil
}

func handleCustomResponse(request *http.Request, writer http.ResponseWriter, response Response) error {
	if fn, ok := response.Body.(func(*http.Request, http.ResponseWriter, Response)); ok {
		fn(request, writer, response)
		return nil
	}
	return fmt.Errorf("custom response body is not a function")
}
