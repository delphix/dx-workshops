package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	resty "gopkg.in/resty.v1"
)

func init() {
	log.Printf("Version: %v", version)
}

type Client struct {
	url, username, password string
}

type AppError struct {
	message, timestamp, errType string
	status                      int
}

func newAppError(resp *resty.Response) *AppError {
	var errMap map[string]interface{}
	result := resp.Body()
	if err := json.Unmarshal(result, &errMap); err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", errMap)
	return &AppError{
		message:   errMap["message"].(string),
		timestamp: errMap["timestamp"].(string),
		errType:   errMap["error"].(string),
		status:    int(errMap["status"].(float64)),
	}
}

func (e *AppError) Error() string {
	return e.message
}

// NewClient creates a new client object
func NewClient(username, password, url string) *Client {
	return &Client{
		url:      url,
		username: username,
		password: password,
	}
}

func (c *Client) initResty() {
	log.Printf("Initializing Resty Client")
	resty.DefaultClient.
		SetTimeout(time.Duration(30 * time.Second)).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	resty.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "PostmanRuntime/7.13.0").
		SetHeader("Accept", "*/*").
		SetHeader("Cache-Control", "no-cache").
		SetHeader("accept-encoding", "gzip, deflate").
		SetHeader("Connection", "keep-alive").
		SetHeader("cache-control", "no-cache")
}

func parseHTTPResponse(resp *resty.Response) ([]byte, error) {
	var err error

	if http.StatusOK != resp.StatusCode() {
		err = newAppError(resp)
		if err != nil {
			return nil, err
		}
	}
	return resp.Body(), err
}

func parseHTTPResponseReturnSlice(resp *resty.Response) (resultdat []interface{}, err error) {
	result, err := parseHTTPResponse(resp)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(result, &resultdat); err != nil { //convert the json to go objects
		return nil, err
	}

	return resultdat, nil
}

func parseHTTPResponseReturnMap(resp *resty.Response) (resultdat map[string]interface{}, err error) {
	result, err := parseHTTPResponse(resp)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(result, &resultdat); err != nil { //convert the json to go objects
		return nil, err
	}

	return resultdat, err
}

func getPseudoRandomPatientIDFromPatientsList(patientList map[string]interface{}) (int, error) {
	patientID, ok := patientList["content"].([]interface{})[0].(map[string]interface{})["id"].(float64)
	if !ok {
		err := fmt.Errorf("Did not find a patientID. Something went terribly wrong")
		return 0, err
	}
	return int(patientID), nil
}
