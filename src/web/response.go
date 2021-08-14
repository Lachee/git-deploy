package web

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"

	"gopkg.in/yaml.v2"
)

const (
	JSON   = "application/json"
	YAML   = "application/yaml"
	XML    = "application/xml"
	JS     = "application/js"
	TEXT   = "text/plain"
	HTML   = "text/html"
	BINARY = "application/octet-stream"
)

//DefaultContentType defines the default content-type for a new response
var DefaultContentType = TEXT

//ObjectDataType defines the default content-type for a rich object
var ObjectDataType = JSON

type ErrorResponse struct {
	Error string
}

type Response struct {
	StatusCode int
	Headers    map[string]string

	contentType string
	data        interface{}
}

// Creates a new blank response
func EmptyResponse() Response {
	response := Response{
		StatusCode:  http.StatusOK,
		Headers:     make(map[string]string),
		contentType: DefaultContentType,
	}
	response.Headers["Server"] = "Git Deploy (Go)"
	return response
}

// Creates a basic response
func NewResponse(statusCode int, data interface{}) Response {
	response := EmptyResponse()
	response.StatusCode = statusCode
	response.data = data
	response.contentType = GetRecommendedContentType(data)
	return response
}

// Create a new error response
func NewErrorResponse(statusCode int, err error) Response {
	return NewResponse(statusCode, ErrorResponse{Error: err.Error()})
}

//GetRecommendedContentType for the given interface
func GetRecommendedContentType(data interface{}) string {
	// Sets up the default content type
	switch data.(type) {
	default:
		return ObjectDataType
	case string:
		return DefaultContentType
	case int:
	case float64:
	case float32:
	case bool:
	case []byte:
		return BINARY
	}
	return DefaultContentType
}

// SetContentType sets the content type from the request, normally text/plain.
func (response *Response) SetContentType(contentType string) {
	response.contentType = contentType
}

/*
 SetContentTypeFromRequest reads the `content-type` header from the request and uses that for the response.
 Returns the new content type.
*/
func (response *Response) SetContentTypeFromRequest(r *http.Request) string {
	contentType := r.Header.Get("Content-type")
	if contentType != "" {
		response.contentType = contentType
	}
	return response.contentType
}

// Sets the data of the response
func (response *Response) SetData(data interface{}) {
	response.data = data
}

// Writes the response to the response writer
func (response *Response) Write(w http.ResponseWriter) error {

	// Encode the body
	data, encodeErr := response.encodeBody()
	if encodeErr != nil {
		return encodeErr
	}

	// Write the headers
	for key, value := range response.Headers {
		w.Header().Set(key, value)
	}

	// Set content type
	w.Header().Set("content-type", response.contentType)

	// Write the data
	_, writeErr := w.Write(data)
	if writeErr != nil {
		return writeErr
	}

	// Finally write the status code
	w.WriteHeader(response.StatusCode)
	return nil
}

// encodeBody returns the binary representation of data, based of Content-Type
func (response *Response) encodeBody() ([]byte, error) {
	contentType := response.contentType
	switch contentType {
	default:
		switch response.data.(type) {
		case string:
			return []byte(response.data.(string)), nil
		case []byte:
			return response.data.([]byte), nil
		default:
			return nil, errors.New("cannot cast object to plain text")
		}
	case "application/json":
		return json.Marshal(response.data)
	case "application/yaml":
	case "application/yml":
		return yaml.Marshal(response.data)
	case "application/xml":
		return xml.Marshal(response.data)
	}
	return nil, errors.New("invalid content type")
}
