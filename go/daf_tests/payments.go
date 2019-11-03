package main

import (
	"fmt"
	"log"

	resty "gopkg.in/resty.v1"
)

func (c *Client) getPaymentsFromAnyPatientInList(patientList map[string]interface{}) (map[string]interface{}, error) {
	patientID, err := getPseudoRandomPatientIDFromPatientsList(patientList)
	if err != nil {
		return nil, err
	}

	return c.getPayments(patientID)
}

func (c *Client) getPayments(patient int) (map[string]interface{}, error) {

	url := fmt.Sprintf("%s/patients/%d/payments", c.url, patient)
	log.Printf("Fetching all payments for patient %d", patient)
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

	return bodyMap, err
}
