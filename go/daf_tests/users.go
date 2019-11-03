package main

import (
	"fmt"
	"log"
	"net/http"

	resty "gopkg.in/resty.v1"
)

type User struct {
	username, firstname, lastname, password string
}

func (c *Client) getUsers() map[string]interface{} {

	url := c.url + "/users/"
	log.Println("Fetching all users")
	log.Println(url)
	resp, err := resty.R().
		Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if http.StatusUnauthorized == resp.StatusCode() {
		err = fmt.Errorf("%d", http.StatusUnauthorized)
		if err != nil {
			log.Fatal(err)
		}
	}
	bodyMap, err := parseHTTPResponseReturnMap(resp)
	if err != nil {
		log.Fatal(err)
	}

	return bodyMap
}

func (c *Client) signUpUser(user *User) (map[string]interface{}, error) {

	url := c.url + "/auth/sign-up"
	log.Println(url)
	log.Printf("Creating %v", user)
	payload := fmt.Sprintf("{\n\t\"username\": \"%s\",\n\t\"firstname\": \"%s\",\n\t\"lastname\": \"%s\"\n,\n\t\"password\": \"%s\"\n}", user.username, user.firstname, user.lastname, user.password)
	log.Printf(payload)
	resp, err := resty.R().
		SetBody(payload).
		Post(url)
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

func (c *Client) deleteUser(userId int) error {

	url := fmt.Sprintf("%s/users/%d", c.url, userId)
	log.Println(url)
	log.Printf("Deleting %d", userId)
	resp, err := resty.R().
		Delete(url)
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
	log.Printf("UserID %d deleted", userId)
	return nil
}

func (c *Client) signUpUserAndReturnID(user *User) (int, error) {
	id, err := c.signUpUser(user)
	if err != nil {
		return 0, err
	}
	return int(id["id"].(float64)), nil
}
