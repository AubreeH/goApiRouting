package routing

import (
	"encoding/json"
	"fmt"
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
				return Response{Status: 7}
			})
		})
	})

	e, err := json.MarshalIndent(r.routes, "", "    ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(e))

	testRoute(t, r, http.MethodGet, "/test/099", 1, map[string]interface{}{"value": "099"})
	testRoute(t, r, http.MethodGet, "/test/100/abc", 2, map[string]interface{}{"value": "100"})
	testRoute(t, r, http.MethodPost, "/test/101/ghi", 3, map[string]interface{}{"value": "101"})
	testRoute(t, r, http.MethodPatch, "/test/102/ghi", 4, map[string]interface{}{"value": "102"})
	testRoute(t, r, http.MethodGet, "/test/testing123", 5, map[string]interface{}{})
	testRoute(t, r, http.MethodDelete, "/test/testing123/def", 6, map[string]interface{}{"test": "def"})

	_, err = r.getFunc(http.MethodPost, "/test/100/abc", map[string]interface{}{})
	if err == nil {
		t.Error("case 5: expected error but didn't get one")
	}

	_, err = r.getFunc(http.MethodGet, "/test/testing123/def", map[string]interface{}{})
	if err == nil {
		t.Error("case 6: expected error but didn't get one")
	}

	testRoute(t, r, http.MethodGet, "/test/100/abc/def", 7, map[string]interface{}{"value1": "100", "value2": "abc", "value3": "def"})
}

func testRoute(t *testing.T, r *Router, method string, path string, expectedStatus int, expectedStoreValues map[string]interface{}) {
	t.Helper()

	store := make(map[string]interface{})
	function, err := r.getFunc(method, path, store)
	if err != nil {
		t.Errorf("case %d: %v", expectedStatus, err)
	} else if function == nil {
		t.Errorf("case %d: returned function is nil", expectedStatus)
	} else {
		for key, value := range expectedStoreValues {
			if store[key] != value {
				t.Errorf("case %d: path parameter was not added to store", expectedStatus)
			}
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
	}
}
