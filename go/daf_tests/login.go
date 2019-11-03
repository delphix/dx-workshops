package main

import (
	"fmt"
	"log"
	"net/http"

	resty "gopkg.in/resty.v1"
)

func (c *Client) loginTest() error {
	log.Printf("Testing login")
	err := c.login()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) login() error {
	c.initResty()
	url := c.url + "/auth/login"
	log.Printf("Logging in to %s", url)

	payload := fmt.Sprintf("{\n\t\"username\": \"%s\",\n\t\"password\": \"%s\"\n}", c.username, c.password)

	resp, err := resty.R().
		SetBody(payload).
		Post(url)
	if err != nil {
		return err
	}

	if http.StatusUnauthorized == resp.StatusCode() {
		err = fmt.Errorf("%d", http.StatusUnauthorized)
		if err != nil {
			return err
		}
	}

	_, err = parseHTTPResponse(resp)
	if err != nil {
		return err
	}

	resty.
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", resp.Body()))

	return nil
}
