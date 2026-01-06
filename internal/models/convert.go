package models

import (
	"errors"
	"io"
	"net/http"
)

func MarshalReqToStoredReq(r *http.Request) (*StoredRequest, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("Error reading bytes from request body in marshallReqToStoredReq\nerr: " + err.Error())
	}
	return &StoredRequest{
		Method: r.Method,
		URL:    r.URL.String(),
		Header: r.Header,
		Body:   string(bodyBytes),
	}, nil
}
func MarshallRespToStoredResp(resp *http.Response, bodyBytes []byte, duration int64) (*StoredResponse, error) {
	//passing bodyBytes here instead of reading from resp as the resp buffer has been read already and is now empty,instead of filling the buffer again using NopCloser, i am just passing the bodyBytes to the func
	return &StoredResponse{
		Status:     resp.StatusCode,
		Header:     resp.Header,
		Body:       string(bodyBytes),
		DurationMs: duration,
	}, nil
}
