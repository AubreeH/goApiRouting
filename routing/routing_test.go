package routing

import (
	"net/http"
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

		api.Get("*", func(c *Context) Response {
			return Response{Status: 11}
		})

		api.Any("*", func(c *Context) Response {
			return Response{Status: 12}
		})
	})

	// testPrettyPrint(t, r.endpoints)

	testRoute(t, r, http.MethodGet, "/test/099", 1, map[string]interface{}{"value": "099"})
	testRoute(t, r, http.MethodGet, "/test/100/abc", 2, map[string]interface{}{"value": "100"})
	testRoute(t, r, http.MethodPost, "/test/101/ghi", 3, map[string]interface{}{"value": "101"})
	testRoute(t, r, http.MethodPatch, "/test/102/ghi", 4, map[string]interface{}{"value": "102"})
	testRoute(t, r, http.MethodGet, "/test/testing123", 5, map[string]interface{}{})
	testRoute(t, r, http.MethodDelete, "/test/testing123/def", 6, map[string]interface{}{"test": "def"})

	_, _, err := r.getFunc(http.MethodPost, "/test/100/abc")
	if err == nil {
		t.Error("case 7: expected error but didn't get one")
	}

	_, _, err = r.getFunc(http.MethodGet, "/test/testing123/def")
	if err == nil {
		t.Error("case 8: expected error but didn't get one")
	}

	testRoute(t, r, http.MethodGet, "/test/10/abc/de", 9, map[string]interface{}{"value1": "10", "value2": "abc", "value3": "de"})
	testRoute(t, r, http.MethodGet, "/jkl", 10, map[string]interface{}{})
	testRoute(t, r, http.MethodGet, "/mno", 11, map[string]interface{}{})
	testRoute(t, r, http.MethodPut, "/pqr", 12, map[string]interface{}{})
}

func testRoute(t *testing.T, r *Router, method string, path string, expectedStatus int, expectedStoreValues map[string]interface{}) {
	t.Helper()

	function, pathParameters, err := r.getFunc(method, path)
	if err != nil {
		t.Errorf("case %d: %v", expectedStatus, err)
	} else {
		if function == nil {
			t.Errorf("case %d: returned function is nil", expectedStatus)
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
			t.Errorf("case %d: incorrect function returned", expectedStatus)
		}

		for key, value := range expectedStoreValues {
			if pathParameters[key] != value {
				t.Errorf("case %d: path parameter was not retrieved from request", expectedStatus)
			}
		}
	}
}
