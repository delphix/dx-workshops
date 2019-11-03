package main

import (
	"fmt"
	"log"

	resty "gopkg.in/resty.v1"
)

func (c *Client) getBillingsFromAnyPatientInList(patientList map[string]interface{}) (map[string]interface{}, error) {
	patientID, err := getPseudoRandomPatientIDFromPatientsList(patientList)
	if err != nil {
		return nil, err
	}
	return c.getBillings(patientID)
}

func (c *Client) getBillings(patient int) (map[string]interface{}, error) {

	url := fmt.Sprintf("%s/patients/%d/billings", c.url, patient)
	log.Printf("Fetching all billings for patient %d", patient)
	log.Println(url)
	resp, err := resty.R().
		Get(url)
	if err != nil {
		return nil, err
	}

	bodyMap, err := parseHTTPResponseReturnMap(resp)
	if err != nil {
		return nil, err
	}

	return bodyMap, nil
}
