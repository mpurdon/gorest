package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/gddo/httputil/header"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Version represents the current version of the rest library
const Version = "0.0.1"

// Method contains the supported HTTP verbs.
type Method string

// Supported HTTP verbs.
const (
	Get    Method = "GET"
	Post   Method = "POST"
	Put    Method = "PUT"
	Patch  Method = "PATCH"
	Delete Method = "DELETE"
)

// Request holds the request to an API Call.
type Request struct {
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
	Method      Method
	BaseURL     string
}

// RestError is a struct for an error handling.
type RestError struct {
	Response *Response
}

// Error is the implementation of the error interface.
func (e *RestError) Error() string {
	return e.Response.Body
}

// DefaultClient is used if no custom HTTP client is defined
var DefaultClient = &Client{
	HTTPClient: http.DefaultClient,
}

// Client allows modification of client properties
// See https://golang.org/pkg/net/http
type Client struct {
	HTTPClient *http.Client
}

// MakeRequest makes the API call.
func (c *Client) MakeRequest(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}

// Send will build your request, make the request, and build your response.
func (c *Client) Send(request Request) (*Response, error) {
	// Build the HTTP request object.
	req, err := BuildRequestObject(request)
	if err != nil {
		return nil, err
	}

	// Build the HTTP client and make the request.
	res, err := c.MakeRequest(req)
	if err != nil {
		return nil, err
	}

	// Build Response object.
	return BuildResponse(res)
}

// Response holds the response from an API call.
type Response struct {
	Headers    map[string][]string // The response headers
	Body       string              // JSON Body
	StatusCode int                 // Http/1.1 Status
}

//
// Json helpers
// See:
// https://www.pixelstech.net/article/1573355516-JSON-unmarshal-in-GoLang
//
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

// Json gets the JSON encoded body as the passed in Entity
func (res *Response) Json(f interface{}) (interface{}, error) {
	b := []byte(res.Body)
	err := json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return f, nil
}

// Json gets the JSON encoded body as a map
func (res *Response) JsonAsMap() (map[string]interface{}, error) {
	var f interface{}

	err := json.Unmarshal([]byte(res.Body), &f)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return f.(map[string]interface{}), nil
}


// AddQueryParameters adds query parameters to the URL.
func AddQueryParameters(baseURL string, queryParams map[string]string) string {
	baseURL += "?"
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	return baseURL + params.Encode()
}

// BuildRequestObject creates the HTTP request object.
func BuildRequestObject(request Request) (*http.Request, error) {
	// Add any query parameters to the URL.
	if len(request.QueryParams) != 0 {
		request.BaseURL = AddQueryParameters(request.BaseURL, request.QueryParams)
	}
	req, err := http.NewRequest(string(request.Method), request.BaseURL, bytes.NewBuffer(request.Body))
	if err != nil {
		return req, err
	}
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	// Add the content-type header if one was not provided
	_, exists := req.Header["Content-Type"]
	if len(request.Body) > 0 && !exists {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, err
}

// MakeRequest makes the API call.
func MakeRequest(req *http.Request) (*http.Response, error) {
	return DefaultClient.HTTPClient.Do(req)
}

// BuildResponse builds the response struct.
func BuildResponse(res *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(res.Body)

	response := Response{
		StatusCode: res.StatusCode,
		Body:       string(body),
		Headers:    res.Header,
	}

	failedToClose := res.Body.Close()
	if failedToClose != nil {
		err = errors.Wrap(failedToClose, "could not close response body")
	}

	return &response, err
}

// Send uses the DefaultClient to send your request
func Send(request Request) (*Response, error) {
	return DefaultClient.Send(request)
}

//
// Experimental JSON Decoding support (WIP)
//

type MalformedRequest struct {
	Status int
	Msg    string
}

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &MalformedRequest{Status: http.StatusUnsupportedMediaType, Msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &MalformedRequest{Status: http.StatusRequestEntityTooLarge, Msg: msg}

		default:
			return err
		}
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	return nil
}
