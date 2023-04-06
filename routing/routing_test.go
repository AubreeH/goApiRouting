package routing

import (
	"errors"
	"net/http"
	"sync"
	"testing"
)

func Test_Router(t *testing.T) {
	config := Config{Port: 8080}
	r := Setup(config)
	r.InitialiseRoutes(func(api BaseApi) {
		api.Group("//test", api.NoMiddleware(), func(api BaseApi) {
			api.Group(`${value="\d\d\d"}`, api.NoMiddleware(), func(api BaseApi) {
				api.Get("abc", func(c *Context) Response {
					return Response{Status: 1}
				})

				api.Post("ghi", func(c *Context) Response {
					return Response{Status: 2}
				})

				api.Patch("ghi", func(c *Context) Response {
					return Response{Status: 3}
				})
			})
			api.Group("testing123", api.NoMiddleware(), func(api BaseApi) {
				api.Delete("${test}", func(c *Context) Response {
					return Response{Status: 4}
				})
			})
		})
	})

	store := make(map[string]interface{})
	function, err := r.getFunc(http.MethodGet, "/test/100/abc", store)
	if err != nil {
		t.Error("case 1:", err)
	} else if function == nil {
		t.Error(errors.New("case 1: returned function is nil"))
	} else if store["value"] != "100" {
		t.Error(errors.New("case 1: path parameter was not added to store"))
	} else {
		var r Response
		wg := sync.WaitGroup{}
		wg.Add(1)

		function(nil, func(response Response) {
			r = response
			wg.Done()
		})

		wg.Wait()
		if r.Status != 1 {
			t.Error("case 1: incorrect function returned")
		}
	}
	store = make(map[string]interface{})
	function, err = r.getFunc(http.MethodPost, "/test/101/ghi", store)
	if err != nil {
		t.Error("case 2:", err)
	} else if function == nil {
		t.Error(errors.New("case 2: returned function is nil"))
	} else if store["value"] != "101" {
		t.Error(errors.New("case 2: path parameter was not added to store"))
	} else {
		var r Response
		wg := sync.WaitGroup{}
		wg.Add(1)

		function(nil, func(response Response) {
			r = response
			wg.Done()
		})

		wg.Wait()
		if r.Status != 2 {
			t.Error("case 2: incorrect function returned")
		}
	}

	store = make(map[string]interface{})
	function, err = r.getFunc(http.MethodPatch, "/test/102/ghi", store)
	if err != nil {
		t.Error("case 3:", err)
	} else if function == nil {
		t.Error(errors.New("case 3: returned function is nil"))
	} else if store["value"] != "102" {
		t.Error(errors.New("case 3: path parameter was not added to store"))
	} else {
		var r Response
		wg := sync.WaitGroup{}
		wg.Add(1)

		function(nil, func(response Response) {
			r = response
			wg.Done()
		})

		wg.Wait()
		if r.Status != 3 {
			t.Error("case 3: incorrect function returned")
		}
	}

	store = make(map[string]interface{})
	function, err = r.getFunc(http.MethodDelete, "/test/testing123/def", store)
	if err != nil {
		t.Error("case 4:", err)
	} else if function == nil {
		t.Error(errors.New("case 4: returned function is nil"))
	} else if store["test"] != "def" {
		t.Error(errors.New("case 3: path parameter was not added to store"))
	} else {
		var r Response
		wg := sync.WaitGroup{}
		wg.Add(1)

		function(nil, func(response Response) {
			r = response
			wg.Done()
		})

		wg.Wait()
		if r.Status != 4 {
			t.Error("case 4: incorrect function returned")
		}
	}

	function, err = r.getFunc(http.MethodPost, "/test/100/abc", store)
	if err == nil {
		t.Error("case 5: expected error but didn't get one")
	}

	function, err = r.getFunc(http.MethodGet, "/test/testing123/def", store)
	if err == nil {
		t.Error("case 6: expected error but didn't get one")
	}
}
