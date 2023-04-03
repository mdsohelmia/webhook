package webhook

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	response *http.Response
	body     []byte
}

// StatusCode returns the status code of the response
func (receiver *Response) StatusCode() int {
	return receiver.response.StatusCode
}

// Header returns the header of the response
func (receiver *Response) Headers() http.Header {
	return receiver.response.Header
}

// Body returns the body of the response
func (receiver *Response) Unmarshal(v any) error {
	if err := json.NewDecoder(receiver.response.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

// Status returns the status of the response
func (receiver *Response) Status() string {
	return receiver.response.Status
}

// Close closes the response body
func (receiver *Response) Close() error {
	return receiver.response.Body.Close()
}

// Ok returns true if the status code is between 200 and 300
func (receiver *Response) Ok() bool {
	return receiver.StatusCode() >= 200 && receiver.StatusCode() < 300
}

// StandardHttpResponse returns the standard http response
func (receiver *Response) StandardHttpResponse() *http.Response {
	return receiver.response
}

// Body returns the body of the response
func (receiver *Response) Body() []byte {
	return receiver.body
}
