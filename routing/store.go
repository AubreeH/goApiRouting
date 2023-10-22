package routing

import (
	"encoding/json"
	"errors"
)

// Retrieves value from the store. The order of precedence is:
// 1. Store
// 2. Path parameters
// 3. Query parameters
// 4. Body
// Returns an error if the key is not found.
func (s *Store) Get(key string) (value interface{}, err error) {
	s.mux.RLock()
	if storeValue, ok := s.store[key]; !ok {
		value = storeValue
	} else if pathValue, ok := s.pathParameters[key]; !ok {
		value = pathValue
	} else if queryValue, ok := s.query[key]; !ok {
		value = queryValue
	} else if bodyValue, ok := s.bodyMap[key]; !ok {
		value = bodyValue
	} else if s.body != nil {
		var body map[string]interface{}
		body, err = s.GetBody()
		if err != nil {
			return nil, err
		}
		value, ok = body[key]
		if !ok {
			err = errors.New("key not found")
		}
	} else {
		err = errors.New("key not found")
	}

	s.mux.RUnlock()
	return
}

// Retrieves path parameter with the given key.
func (s *Store) GetPathParameter(key string) (value string, ok bool) {
	s.mux.RLock()
	value, ok = s.pathParameters[key]
	s.mux.RUnlock()
	return
}

// Retrieves query parameter with the given key.
func (s *Store) GetQueryParameter(key string) (value []string, ok bool) {
	s.mux.RLock()
	value, ok = s.query[key]
	s.mux.RUnlock()
	return
}

// Returns the request body as a map.
func (s *Store) GetBody() (value map[string]interface{}, err error) {
	s.mux.RLock()
	if s.body == nil {
		err = errors.New("request body is nil")
	} else if s.bodyMap == nil {
		err = json.Unmarshal(s.body, &s.bodyMap)
	}
	value = s.bodyMap
	s.mux.RUnlock()
	return
}

// Returns the request body marshalled into the given type.
func GetBody[TBody any](s *Store) (value TBody, err error) {
	s.mux.RLock()
	if s.body == nil {
		err = errors.New("request body is nil")
	} else {
		err = json.Unmarshal(s.body, &value)
	}
	s.mux.RUnlock()
	return
}
