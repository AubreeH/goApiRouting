package routing

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
)

func Test_Router(t *testing.T) {
	config := Config{Port: 8080}

	r := Setup(config)
	r.InitialiseRoutes(func(api BaseApi) {
		api.Group("test", api.NoMiddleware(), func(api BaseApi) {
			api.Group(`${value="\d\d\d"}`, api.NoMiddleware(), func(api BaseApi) {
				api.Get("", func(c *Context) Response {
					return Response{Status: 1}
				})

				api.Get("abc", func(c *Context) Response {
					return Response{Status: 2}
				})

				api.Post("ghi", func(c *Context) Response {
					return Response{Status: 3}
				})

				api.Patch("ghi", func(c *Context) Response {
					return Response{Status: 4}
				})
			})

			api.Get("testing123", func(c *Context) Response {
				return Response{Status: 5}
			})

			api.Group("testing123", api.NoMiddleware(), func(api BaseApi) {
				api.Delete("${test}", func(c *Context) Response {
					return Response{Status: 6}
				})
			})

			api.Get(`${value1="\d\d"}/${value2}/${value3="[a-z][a-z]"}`, func(c *Context) Response {
				return Response{Status: 9}
			})
		})

		api.Any("jkl", func(c *Context) Response {
			return Response{Status: 10}
		})

		api.Get("", func(c *Context) Response {
			return Response{Status: 11}
		})

		api.Any("*", func(c *Context) Response {
			return Response{Status: 12}
		})

		api.Get("test/*", func(c *Context) Response {
			return Response{Status: 13}
		})

		api.Any("test/*", func(c *Context) Response {
			return Response{Status: 14}
		})
	})

	// testPrettyPrint(t, r.endpoints)

	testRouteSuccess(1, t, r, http.MethodGet, "/test/099", 1, map[string]interface{}{"value": "099"})
	testRouteSuccess(2, t, r, http.MethodGet, "/test/100/abc", 2, map[string]interface{}{"value": "100"})
	testRouteSuccess(3, t, r, http.MethodPost, "/test/101/ghi", 3, map[string]interface{}{"value": "101"})
	testRouteSuccess(4, t, r, http.MethodPatch, "/test/102/ghi", 4, map[string]interface{}{"value": "102"})
	testRouteSuccess(5, t, r, http.MethodGet, "/test/testing123", 5, map[string]interface{}{})
	testRouteSuccess(6, t, r, http.MethodDelete, "/test/testing123/def", 6, map[string]interface{}{"test": "def"})
	testRouteError(7, t, r, http.MethodPost, "/test/100/abc", errors.New("method not supported"))
	testRouteError(8, t, r, http.MethodPost, "/test/testing123/def", errors.New("method not supported"))
	testRouteSuccess(9, t, r, http.MethodGet, "/test/10/abc/de", 9, map[string]interface{}{"value1": "10", "value2": "abc", "value3": "de"})
	testRouteSuccess(10, t, r, http.MethodGet, "/jkl", 10, map[string]interface{}{})
	testRouteSuccess(11, t, r, http.MethodGet, "", 11, map[string]interface{}{})
	testRouteSuccess(12, t, r, http.MethodGet, "/", 11, map[string]interface{}{})
	testRouteSuccess(13, t, r, http.MethodGet, "/mno", 12, map[string]interface{}{})
	testRouteSuccess(14, t, r, http.MethodPut, "/pqr", 12, map[string]interface{}{})
	testRouteSuccess(15, t, r, http.MethodPost, "/", 12, map[string]interface{}{})
	testRouteSuccess(16, t, r, http.MethodDelete, "", 12, map[string]interface{}{})
	testRouteSuccess(18, t, r, http.MethodGet, "/test/abc", 13, map[string]interface{}{})
}

func testRouteSuccess(testId int, t *testing.T, r *Router, method string, path string, expectedStatus int, expectedStoreValues map[string]interface{}) {
	t.Helper()

	function, pathParameters, err := r.getFunc(method, path)
	if err != nil {
		testError(t, "case %d: function return error '%v'", testId, err)
	} else {
		if function == nil {
			testError(t, "case %d: returned function is nil", testId)
		}

		var r Response
		wg := sync.WaitGroup{}
		wg.Add(1)

		function(nil, func(response Response) {
			r = response
			wg.Done()
		})

		wg.Wait()

		if r.Status != expectedStatus {
			testError(t, "case %d: incorrect function returned. expected %d but recieved %d", testId, expectedStatus, r.Status)
		}

		for key, value := range expectedStoreValues {
			if v, ok := pathParameters[key]; ok && v != value {
				testError(t, "case %d: path parameter '%s' did not match. expected '%v' but received '%v'", testId, key, value, v)
			} else if !ok {
				testError(t, "case %d: path parameter '%s' not found", testId, key)
			}
		}
	}
}

func testRouteError(testId int, t *testing.T, r *Router, method string, path string, expectedError error) {
	t.Helper()

	_, _, err := r.getFunc(method, path)
	if err == nil {
		testError(t, "case %d: expected error but didn't get one", testId)
	} else if err.Error() != expectedError.Error() {
		testError(t, "case %d: incorrect error returned. expected '%v' but received '%v'", testId, expectedError, err)
	}
}

func testError(t *testing.T, format string, args ...any) {
	t.Helper()
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Printf(format, args...)
	t.Errorf(format, args...)
}
