package routing

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
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
		handler(
			&Context{
				Writer:  writer,
				Request: request,
				Store: &Store{
					pathParameters: pathParameters,
					query:          request.URL.Query(),
					body:           requestBody,
					store:          make(map[string]interface{}),
				},
			},
			func(response Response) {
				r.writeResponse(writer, request, response)
			},
		)
	}
}

func getRequestBody(request *http.Request) ([]byte, error) {
	if request.Body == nil {
		return nil, nil
	}

	var body []byte
	_, err := request.Body.Read(body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *Router) writeResponse(writer http.ResponseWriter, request *http.Request, response Response) {
	r.writeHeaders(writer, response.Headers)

	switch response.Type {
	case JSONResponse:
		writeJSONResponse(request, writer, response)
	case HTMLResponse:
		writeHTMLResponse(request, writer, response)
	case XMLResponse:
		writeXMLResponse(request, writer, response)
	case PlainTextResponse:
		writePlainTextResponse(request, writer, response)
	case FileResponse:
		writeFileResponse(request, writer, response)
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
	body, err := json.Marshal(response.Body)
	if err != nil {
		return err
	}
	return write("application/json", writer, body, response)
}

func writeHTMLResponse(request *http.Request, writer http.ResponseWriter, response Response) error {
	str, ok := response.Body.(string)
	if !ok {
		return errors.New("body is not a string")
	}
	return write("text/html", writer, []byte(str), response)
}

func writeXMLResponse(request *http.Request, writer http.ResponseWriter, response Response) error {
	body, err := xml.Marshal(response.Body)
	if err != nil {
		return err
	}
	return write("application/xml", writer, body, response)
}

func writePlainTextResponse(request *http.Request, writer http.ResponseWriter, response Response) error {
	str, ok := response.Body.(string)
	if !ok {
		return errors.New("body is not a string")
	}
	return write("text/plain", writer, []byte(str), response)
}

func writeFileResponse(request *http.Request, writer http.ResponseWriter, response Response) {
	writer.WriteHeader(response.Status)
	http.ServeFile(writer, request, response.Body.(string))
}
