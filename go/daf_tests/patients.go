package main

import (
	"fmt"
	"log"
	"net/http"

	resty "gopkg.in/resty.v1"
)

func (c *Client) getPatients() (map[string]interface{}, error) {

	url := c.url + "/patients/"
	log.Println("Fetching all patients")
	log.Println(url)
	resp, err := resty.R().
		Get(url)
	if err != nil {
		return nil, err
	}

	if http.StatusUnauthorized == resp.StatusCode() {
		err = fmt.Errorf("%d", http.StatusUnauthorized)
		if err != nil {
			return nil, err
		}
	}
	bodyMap, err := parseHTTPResponseReturnMap(resp)
	if err != nil {
		return nil, err
	}

	return bodyMap, nil
}
