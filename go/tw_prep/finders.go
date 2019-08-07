package main

import (
	log "github.com/sirupsen/logrus"
)

// Takes one param, the namespace you want to list
// i.e database
// i.e selfservice/container
func (c *Client) listObjects(namespace string, params ...string) (results []interface{}) {
	listResult := c.httpGet(namespace, params...)
	results = listResult["result"].([]interface{})
	logger.WithField("url", url).Debug()
	return results
}

func (c *Client) findObjectByAttributeValue(attribute, namespace, value string, params ...string) (result map[string]interface{}) {
	logger.WithFields(log.Fields{
		"attribute":   attribute,
		"value":       value,
		"namespace":   namespace,
		"queryParams": params,
	})
	logger.Debug("Looking for")

	objects := c.listObjects(namespace, params...)

	for _, object := range objects {
		logger.WithFields(log.Fields{
			"object": object,
		}).Debug("Inspecting:")

		if value == object.(map[string]interface{})[attribute].(string) {
			log.Debug("Found")
			return object.(map[string]interface{})
		}
	}
	logger.Warn("Not Found")
	return nil
}

func (c *Client) findObjectByName(namespace, objectName string, params ...string) (result map[string]interface{}) {
	return c.findObjectByAttributeValue("name", namespace, objectName, params...)
}

func (c *Client) findObjectByReference(namespace, objectReference string, params ...string) (result map[string]interface{}) {
	return c.findObjectByAttributeValue("reference", namespace, objectReference, params...)
}
