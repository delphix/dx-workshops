package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Takes one param, the namespace you want to list
// i.e database
// i.e selfservice/container
func (c *myClient) listObjects(namespace string, params ...string) (results []interface{}, err error) {
	listResult, _, err := c.httpGet(namespace, params...)
	if err != nil {
		return nil, err
	}
	results = listResult["result"].([]interface{})
	logger.WithField("url", url).Debug()
	return results, nil
}

func (c *myClient) findObjectByAttributeValue(attribute, namespace, value string, params ...string) (result map[string]interface{}, err error) {
	logger := logger.WithFields(log.Fields{
		"attribute":   attribute,
		"value":       value,
		"namespace":   namespace,
		"queryParams": params,
	})
	logger.Debug("Looking for")

	objects, err := c.listObjects(namespace, params...)
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		logger.WithFields(log.Fields{
			"object": object,
		}).Debug("Inspecting:")

		if value == object.(map[string]interface{})[attribute].(string) {
			log.Debug("Found")
			return object.(map[string]interface{}), err
		}
	}
	logger.Debug("Not Found")
	return nil, err
}

func (c *myClient) findObjectByName(namespace, objectName string, params ...string) (result map[string]interface{}, err error) {
	return c.findObjectByAttributeValue("name", namespace, objectName, params...)
}

func (c *myClient) findObjectByNameReturnReference(namespace, objectName string, params ...string) (reference interface{}, err error) {
	return returnObjReference(c.findObjectByAttributeValue("name", namespace, objectName, params...))
}

func (c *myClient) findObjectByReference(namespace, objectReference string, params ...string) (result map[string]interface{}, err error) {
	return c.findObjectByAttributeValue("reference", namespace, objectReference, params...)
}

func (c *myClient) findSourceCongfigByNameAndEnvironmentName(scName, envName string) (results map[string]interface{}, err error) {
	namespace := "environment"
	//Find our Environment of interest
	logger.Info("Searching for Environment by name")
	envObj, err := c.findObjectByName(namespace, envName)
	if err != nil {
		return nil, err
	}
	if envObj == nil {
		logger.Fatalf("Could not find Environment named %s", envName)
	}

	logger.Infof("Found %s: %s", envObj["name"], envObj["reference"])

	namespace = "sourceconfig"
	//Find our SourceConfig of interest
	logger.Info("Searching for SourceConfig by name")
	scObj, err := c.findObjectByName(namespace, scName, fmt.Sprintf("environment=%s", envObj["reference"]))
	if err != nil {
		return nil, err
	}
	if scObj == nil {
		logger.Infof("Could not find SourceConfig named %s on %s", scName, envName)
		return nil, nil
	}

	logger.Infof("Found %s: %s", scObj["name"], scObj["reference"])

	return scObj, err
}

func (c *myClient) findRepoByNameAndEnvironmentName(repoName, envName string) (results map[string]interface{}, err error) {
	namespace := "environment"
	//Find our Environment of interest
	logger.Info("Searching for Environment by name")
	envObj, err := c.findObjectByName(namespace, envName)
	if err != nil {
		return nil, err
	}
	if envObj == nil {
		logger.Fatalf("Could not find Environment named %s", envName)
	}

	logger.Infof("Found %s: %s", envObj["name"], envObj["reference"])

	namespace = "repository"
	//Find our SourceConfig of interest
	logger.Info("Searching for Repository by name")
	repoObj, err := c.findObjectByName(namespace, repoName, fmt.Sprintf("environment=%s", envObj["reference"]))
	if err != nil {
		return nil, err
	}
	if repoObj == nil {
		logger.Infof("Could not find Repo named %s on %s", repoName, envName)
		return nil, nil
	}

	logger.Infof("Found %s: %s", repoObj["name"], repoObj["reference"])

	return repoObj, err
}

func (c *myClient) findActionByTitleAndParentAction(title, parentAction string) (reference map[string]interface{}, err error) {
	return c.findObjectByAttributeValue("title", "action", title, fmt.Sprintf("parentAction=%s", parentAction))
}
