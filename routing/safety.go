package routing

import (
	"errors"
	"net/http"
)

func (r *Router) vetRequest(request *http.Request, writer http.ResponseWriter) error {
	if err := r.checkContentLength(request, writer); err != nil {
		return err
	}

	return nil
}

func (r *Router) checkContentLength(request *http.Request, writer http.ResponseWriter) error {
	if r.config.MaxContentLength == 0 {
		return nil
	}

	if request.ContentLength == -1 {
		r.writeResponse(writer, request, Response{
			Body:   "content length not specified",
			Type:   PlainTextResponse,
			Status: StatusLengthRequired,
		})
	} else if request.ContentLength == 0 && request.Body != nil {
		r.writeResponse(writer, request, Response{
			Body:   "content length not specified",
			Type:   PlainTextResponse,
			Status: StatusLengthRequired,
		})
		return errors.New("content length not specified")
	} else if request.ContentLength > int64(r.config.MaxContentLength) {
		r.writeResponse(writer, request, Response{
			Body:   "content length too large",
			Type:   PlainTextResponse,
			Status: StatusRequestEntityTooLarge,
		})
		return errors.New("content length too large")
	}

	return nil
}
