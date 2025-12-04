package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type LoggedExchange struct {
	Time     time.Time      `json:"time"`
	Request  StoredRequest  `json:"request"`
	Response StoredResponse `json:"response"`
}

type StoredRequest struct {
	Method string              `json:"method"`
	URL    string              `json:"url"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

type StoredResponse struct {
	Status     int                 `json:"status"`
	Header     map[string][]string `json:"header"`
	Body       string              `json:"body"`
	DurationMs int64               `json:"duration_ms"` // time taken by the backend to send its full response to your replayer proxy
}

func DumpReqAndResToFile(r *http.Request, resp *http.Response, duration int64, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644) // 0644 → 0 = octal, 6(owner)=read+write, 4(group)=read, 4(others)=read
	if err != nil {
		return errors.New("Error opening given file in APPEND, WRITE ONLY, CREATE mode\ner: " + err.Error())
	}
	defer file.Close()
	storedReq, err := marshalReqToStoredReq(r)
	if err != nil {
		return err
	}
	storedResp, err := marshallRespToStoredResp(resp, duration)
	if err != nil {
		return err
	}
	loggedExchange := &LoggedExchange{
		Time:     time.Now(),
		Request:  *storedReq,
		Response: *storedResp,
	}
	err = json.NewEncoder(file).Encode(loggedExchange)
	if err != nil {
		return errors.New("Error encoding loggedExchange struct\nerr: " + err.Error())
	}
	return nil
}

func marshalReqToStoredReq(r *http.Request) (*StoredRequest, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("error parsing request body\nerr: " + err.Error())
	}
	return &StoredRequest{
		Method: r.Method,
		URL:    r.URL.String(),
		Header: r.Header,
		Body:   string(bodyBytes),
	}, nil
}

func marshallRespToStoredResp(resp *http.Response, duration int64) (*StoredResponse, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Error reading bytes from response body in marshallRespToStoredResp\nerr: " + err.Error())
	}
	return &StoredResponse{
		Status:     resp.StatusCode,
		Header:     resp.Header,
		Body:       string(bodyBytes),
		DurationMs: duration,
	}, nil
}

func ReplayReqFromFile(fileName string) {
	channel := make(chan *StoredRequest, 50)
	go ReadReqFromFile(fileName, channel)
	go ResendReq(channel)
}

func ReadReqFromFile(fileName string, channel chan *StoredRequest) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("error opening file: ", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() { // reads the next Line
		fmt.Println("read req no: ", count)
		reqLine := scanner.Text() // return the line as string
		var storedReq StoredRequest
		err = json.Unmarshal([]byte(reqLine), &storedReq)
		if err != nil {
			log.Fatal("error decoding reqLine read from file: ", fileName)
		}
		channel <- &storedReq
		count++
	}
	close(channel)
	fmt.Println("READ ALL REQUESTS")
	if scanner.Err() != nil {
		log.Fatal("error while reading file.\nerr: ", err)
	}
}

func ResendReq(channel chan *StoredRequest) {
	count := 0
	for storedReq := range channel {
		fmt.Println("resending req num: ", count)
		body := bytes.NewReader([]byte(storedReq.Body))
		// TODO: ADD THE  HTTP://LOCALHOST:8080 PART TO .ENV so it can be changed for production
		req, err := http.NewRequest(storedReq.Method, "http://localhost:8080"+storedReq.URL, body)
		if err != nil {
			fmt.Println("error created request from stored request\nerr: ", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// TODO: MAKE AN ERROR CHANNEL to handle these errors better
			fmt.Println("error sending request\nerr: ", err)
		}
		respBytes, err := httputil.DumpResponse(resp, true)
		if err != nil {
			// TODO: MAKE AN ERROR CHANNEL to handle these errors better
			fmt.Println("error converting response to bytes\nerr: ", err)
		}
		fmt.Println(string(respBytes))
		count++
	}
}
