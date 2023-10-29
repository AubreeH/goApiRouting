package routing

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"strings"
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

// Sets value in the store.
func (c *Store) Set(key string, value interface{}) {
	c.mux.Lock()
	c.store[key] = value
	c.mux.Unlock()
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
	if s.bodyMap == nil {
		err = s.parseBody()
	}
	value = s.bodyMap
	s.mux.RUnlock()
	return
}

func (s *Store) parseBody() error {
	if s.bodyMap != nil {
		return nil
	} else if s.body == nil {
		return errors.New("request body is nil")
	}

	switch s.contentType {
	case "application/json":
		return json.Unmarshal(s.body, &s.bodyMap)
	case "application/xml":
		return xml.Unmarshal(s.body, &s.bodyMap)
	case "application/x-www-form-urlencoded":
		return s.parseForm()
	}

	if strings.HasPrefix(s.contentType, "multipart/form-data") {
		return s.parseMultipartForm()
	}

	return errors.New("unsupported content type")
}

func (s *Store) parseForm() error {
	req := s.context.Request
	if err := req.ParseForm(); err != nil {
		return err
	}

	err := req.ParseForm()
	if err != nil {
		return err
	}

	for key, value := range req.Form {
		s.bodyMap[key] = value
	}

	return nil
}

func (s *Store) parseMultipartForm() error {
	req := s.context.Request
	reader, err := req.MultipartReader()
	if err != nil {
		return err
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		formName := part.FormName()
		if part.FileName() != "" {
			_, header, err := req.FormFile(formName)
			if err != nil {
				return err
			}
			s.files[formName] = &File{
				formname:       formName,
				formFileHeader: header,
			}
		} else {
			var b []byte
			_, err := part.Read(b)
			if err != nil {
				return err
			}
			s.bodyMap[formName] = string(b)
		}
	}
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
