package main

import (
	"fmt"
	"net"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

// Options for the script
type Options struct {
	DDPVirtName     string               `short:"e" long:"ddp_virt_hostname" env:"DDP_VIRT_HOSTNAME" description:"The hostname or IP address of the Delphix Dynamic Data Platform - Virtualization Engine" required:"true"`
	DDPMaskName     string               `short:"m" long:"ddp_mask_hostname" env:"DDP_MASK_HOSTNAME" description:"The hostname or IP address of the Delphix Dynamic Data Platform - Masking Engine" required:"true"`
	VirtUserName    string               `long:"virt_username" env:"DDP_VIRT_USERNAME" description:"The username used to authenticate to the Virtualization Engine" required:"true"`
	VirtPassword    string               `long:"virt_password" env:"DDP_VIRT_PASSWORD" description:"The password used to authenticate to the Virtualization Engine" required:"true"`
	MaskUserName    string               `long:"mask_username" env:"DDP_MASK_USER" description:"The username used to authenticate to the Masking Engine" required:"true"`
	MaskPassword    string               `long:"mask_password" env:"DDP_MASK_PASS" description:"The password used to authenticate to the Mas Engine" required:"true"`
	Debug           []bool               `short:"v" long:"debug" env:"DELPHIX_DEBUG" description:"Turn on debugging. -vvv for the most verbose debugging"`
	SkipValidate    bool                 `long:"skip-validate" env:"DELPHIX_SKIP_VALIDATE" description:"Don't validate TLS certificate of Delphix Engine"`
	ConfigFile      func(s string) error `short:"c" long:"config" description:"Optional INI config file to pass in for the variables" no-ini:"true"`
	DSourceList     []string             `long:"dsource" env:"DELPHIX_DSOURCE" description:"The name of the dSource to enable" required:"true"`
	VDBList         []string             `long:"vdb" env:"DELPHIX_VDB" description:"The name of the VDB to enable and start" required:"true"`
	EnvironmentList []string             `long:"environment" env:"DELPHIX_ENV" description:"The name of the environment to refresh" required:"true"`
}

func (c *myClient) syncDatabaseByName(dSourceName string, reSync bool) (results map[string]interface{}, err error) {
	var postBody string
	databaseNamespace := "database"
	//Find our dSource of interest
	log.Info("Searching for dSource by name")
	databaseObj, err := c.findObjectByName(databaseNamespace, dSourceName)
	if err != nil {
		return nil, err
	}
	if databaseObj == nil {
		log.Fatalf("Could not find dSource named %s", dSourceName)
	}

	log.Infof("Found %s: %s", databaseObj["name"], databaseObj["reference"])

	if reSync == true {
		postBody = fmt.Sprintf(`{
		"resync":true,
		"type":"AppDataSyncParameters"
		}`)
	}
	log.Debugf(postBody)
	//Sync the dSource
	url := fmt.Sprintf("%s/%s/sync", databaseNamespace, databaseObj["reference"])
	action, _, err := c.httpPost(url, postBody)
	if err != nil {
		return nil, err
	}

	c.jobWaiter(action)
	return action, err
}

func (c *myClient) refreshDatabaseByName(vdbName string) (results map[string]interface{}, err error) {
	namespace := "database"
	//Find our VDB of interest
	log.Info("Searching for VDB by name")
	databaseObj, err := c.findObjectByName(namespace, vdbName)
	if err != nil {
		return nil, err
	}
	if databaseObj == nil {
		log.Fatalf("Could not find VDB named %s", vdbName)
	}

	log.Infof("Found %s: %s", databaseObj["name"], databaseObj["reference"])

	//refresh the VDB
	url := fmt.Sprintf("%s/%s/refresh", namespace, databaseObj["reference"])
	refreshParameters := `{
		"type": "RefreshParameters", 
		"timeflowPointParameters": {
			"type": "TimeflowPointSemantic",
			"location": "LATEST_SNAPSHOT"
		}
		}`
	action, _, err := c.httpPost(url, refreshParameters)
	if err != nil {
		return nil, err
	}

	c.jobWaiter(action)
	return action, err
}

func (c *myClient) refreshEnvironmentByName(envName string, wait bool) (results map[string]interface{}, err error) {
	namespace := "environment"
	//Find our Environment of interest
	log.Info("Searching for Environment by name")
	envObj, err := c.findObjectByName(namespace, envName)
	if err != nil {
		return nil, err
	}
	if envObj == nil {
		log.Fatalf("Could not find Environment named %s", envName)
	}

	log.Infof("Found %s: %s", envObj["name"], envObj["reference"])

	//refresh the Environment
	url := fmt.Sprintf("%s/%s/refresh", namespace, envObj["reference"])
	action, _, err := c.httpPost(url, "")
	if err != nil {
		return nil, err
	}

	if wait {
		c.jobWaiter(action)
	}
	return action, err
}

func (c *myClient) batchRefreshEnvironmentByName(envList []string) (resultsList []map[string]interface{}, err error) {
	for _, v := range envList {
		action, err := c.refreshEnvironmentByName(v, false)
		if err != nil {
			return nil, err
		}
		resultsList = append(resultsList, action)
	}

	c.jobWaiter(resultsList...)

	return resultsList, err
}

func (c *myClient) updateEnvironmentHostByHostName(envName string, wait bool) (results map[string]interface{}, err error) {
	namespace := "environment"
	//Find our Environment of interest
	log.Info("Searching for Environment by name")
	envObj, err := c.findObjectByName(namespace, envName)
	if err != nil {
		return nil, err
	}
	if envObj == nil {
		log.Fatalf("Could not find Environment named %s", envName)
	}

	log.Infof("Found %s: %s", envName, envObj["reference"])

	namespace = "host"
	//Find our Host of interest
	log.Info("Searching for Host by Environment")
	hostObj, err := c.listObjects(namespace, fmt.Sprintf("environment=%s", envObj["reference"]))
	if err != nil {
		return nil, err
	}
	if hostObj == nil {
		log.Fatalf("Could not find Host linked to Environment %s", envName)
	}

	log.Infof("Found %s: %s", envName, envObj["reference"])
	log.Debug(hostObj)

	//Get IP address of the Environment Name
	ips, err := net.LookupIP(envName)
	if err != nil {
		log.Error(fmt.Errorf("Could not get IPs: %v", err))
	}

	log.Info(ips[0])

	//Update the Host
	url := fmt.Sprintf("%s/%s", namespace, hostObj[0].(map[string]interface{})["reference"])
	hostUpdate := fmt.Sprintf(`{
		"type": "UnixHost", 
		"address": "%s"
		}
	`, ips[0])
	action, _, err := c.httpPost(url, hostUpdate)
	if err != nil {
		return nil, err
	}

	if wait {
		c.jobWaiter(action)
	}
	return action, err
}

func (c *myClient) batchUpdateEnvironmentHostByHostName(envList []string) (resultsList []map[string]interface{}, err error) {
	for _, v := range envList {
		action, err := c.updateEnvironmentHostByHostName(v, false)
		if err != nil {
			return nil, err
		}
		resultsList = append(resultsList, action)
	}

	c.jobWaiter(resultsList...)

	return resultsList, nil
}

func (c *myClient) findSourceByDatabaseRef(databaseRef string) (map[string]interface{}, error) {
	namespace := "source"
	//Find our VDB of interest
	log.Info("Searching for source by reference")
	obj, err := c.listObjects(namespace, fmt.Sprintf("database=%s", databaseRef))
	if err != nil {
		return nil, err
	}
	log.Debug(obj)
	if obj == nil {
		log.Fatalf("Could not find sourceconfig for VDB reference %s", databaseRef)
	} else if len(obj) > 1 {
		log.Fatalf("More than one result was returned. Exiting.\n%v", obj)
	}
	return obj[0].(map[string]interface{}), err
}

func (c *myClient) updatePGSource(sourceRef string) error {
	var err error
	prodDBIP, err := getIP("proddb")
	if err != nil {
		return err
	}
	logger.Infof("updating staging host to %s", prodDBIP)

	postBody := fmt.Sprintf(`{
		"parameters": {
			"delphixInitiatedBackupFlag": true,
			"keepStagingInSync": true,
			"postgresPort": 5433,
			"externalBackup": [],
			"delphixInitiatedBackup": [
			{
				"userName": "delphixdb",
				"postgresSourcePort": 5432,
				"userPass": "********",
				"sourceHostAddress": "%s"
			}
			],
			"configSettingsStg": []
		},
		"type": "AppDataLinkedStagedSource"
		}`, prodDBIP)

	namespace := "source"
	url := fmt.Sprintf("%s/%s", namespace, sourceRef)
	action, _, err := c.httpPost(url, postBody)
	if err != nil {
		return err
	}

	c.jobWaiter(action)
	return err
}

func (c *myClient) updateMMSource(sourceRef string) error {
	var err error
	maskingIP, err := getIP("maskingengine")
	if err != nil {
		return err
	}
	logger.Infof("updating masking engine to %s", maskingIP)

	postBody := fmt.Sprintf(`{
		"operations": {
			"configureClone": [
			{
				"command": "#!/usr/bin/env bash\n#\n# Copyright (c) 2019 by Delphix. All rights reserved.\n#\n#v1.0\n#2019 - Adam Bowen\n#requires curl and jq to be installed on the host machine\nDMHOST=\"%s\"\nDMPORT=\"80\"\nURL=\"http://${DMHOST}:${DMPORT}/masking/api\"\nENVIRONMENT=\"Patients Environment\"\nMASKINGJOB=\"Patients Mask - PG\"\nCONNECTOR=\"PatientsMM - PG\"\nDMUSER=Admin\nDMPASS=Landshark-12\n\necho \"Connecting to ${URL}\"\n\necho \"Authenticating\"\n\nAUTH=$(curl -sX POST --header 'Content-Type: application/json' --header 'Accept: application/json'\\\n  -d \"{ \\\"username\\\": \\\"${DMUSER}\\\", \\\"password\\\": \\\"${DMPASS}\\\"}\" \"${URL}/login\" | jq -r .Authorization)\n\n[[ -z $AUTH || $AUTH == \"null\" ]] && echo \"Was unable to get authenticate. Please try again\" && exit 1\n\nENVID=$(curl -sX GET --header 'Accept: application/json' --header \"Authorization: ${AUTH}\" \"${URL}/environments\" | \\\n  jq -r \".responseList[]| select(.environmentName==\\\"${ENVIRONMENT}\\\").environmentId\")\n\n[[ -z $ENVID || $ENVID == \"null\" ]] && echo \"Was unable to find Job ${ENVIRONMENT}. Please try again\" && echo ${EXECID} && exit 1\n\nJOBID=$(curl -sX GET --header 'Accept: application/json' --header \"Authorization: ${AUTH}\" \\\n  \"${URL}/masking-jobs?environment_id=${ENVID}\" | jq -r \".responseList[] | select(.jobName==\\\"${MASKINGJOB}\\\").maskingJobId\")\n\n[[ -z $JOBID || $JOBID == \"null\" ]] && echo \"Was unable to find Job ${MASKINGJOB}. Please try again\" && echo ${JOBID} && exit 1\n\nCONID=$(curl -sX GET --header 'Accept: application/json' --header \"Authorization: ${AUTH}\" \\\n  \"${URL}/database-connectors?environment_id=${ENVID}\" | jq -r \".responseList[] | select(.connectorName==\\\"${CONNECTOR}\\\").databaseConnectorId\")\n\n[[ -z $CONID || $CONID == \"null\" ]] && echo \"Was unable to find Job ${CONNECTOR}. Please try again\" && echo ${CONID} && exit 1\n\necho \"Executing Job \\\"${ENVIRONMENT}/${MASKINGJOB}\\\" with \\\"${CONNECTOR}\\\": ${JOBID}\"\n\nEXECID=$(curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json'\\\n  --header \"Authorization: ${AUTH}\" -d \"{\\\"jobId\\\": ${JOBID}, \\\"targetConnectorId\\\": ${CONID}}}\"\\\n  \"${URL}/executions\"|jq -r .executionId)\n\n[[ -z $EXECID || $EXECID == \"null\" ]] && echo \"Was unable to start Job ${JOBID}. Please try again\" && echo ${EXECID} && exit 1\n\necho \"Waiting for execution ${EXECID} to finish\"\n\nwhile [[ \"$STATUS\" != \"SUCCEEDED\" && \"$STATUS\" != \"FAILED\" ]]\n    do\n    sleep 3\n    STATUS=$(curl -sX GET --header 'Accept: application/json' --header \"Authorization: ${AUTH}\"\\\n    \"${URL}/executions/${EXECID}\" | jq -r .status)\n    [[ -z $STATUS || $STATUS == \"null\" ]] && echo \"Was unable to get status of execution ${EXECID}. Please try again\" && exit 1\ndone\n\n[[ \"$STATUS\" == \"FAILED\" ]] && echo \"\\\"${ENVIRONMENT}/${MASKINGJOB}\\\" failed execution ${EXECID}.\\nCheck logs for details.\" && exit 1\n\necho \"\\\"${ENVIRONMENT}/${MASKINGJOB}\\\" : $JOBID successfully ran\"",
				"name": "Masking Job",
				"type": "RunBashOnSourceOperation"
			}
			],
			"type": "VirtualSourceOperations"
		},
		"type": "AppDataVirtualSource"
		}`, maskingIP)

	namespace := "source"
	url := fmt.Sprintf("%s/%s", namespace, sourceRef)
	action, _, err := c.httpPost(url, postBody)
	if err != nil {
		return err
	}

	c.jobWaiter(action)
	return err
}

func (c *myClient) updatePGdSource(dName string) error {
	namespace := "database"
	//Find our VDB of interest
	log.Info("Searching for dSource by name")
	obj, err := c.findObjectByName(namespace, dName)
	if err != nil {
		return err
	}
	log.Debug(obj)
	if obj == nil {
		log.Fatalf("Could not find dSource named %s", dName)
	}
	log.Infof("Found %s: %s", obj["name"], obj["reference"])
	sourceObj, err := c.findSourceByDatabaseRef(obj["reference"].(string))
	if err != nil {
		return err
	}
	err = c.updatePGSource(sourceObj["reference"].(string))
	if err != nil {
		return err
	}
	_, err = c.syncDatabaseByName(dName, true)
	return err
}

func (c *myClient) updateMMConfigClone(dName string) error {
	namespace := "database"
	//Find our VDB of interest
	log.Info("Searching for VDB by name")
	obj, err := c.findObjectByName(namespace, dName)
	if err != nil {
		return err
	}
	log.Debug(obj)
	if obj == nil {
		log.Fatalf("Could not find dSource named %s", dName)
	}
	log.Infof("Found %s: %s", obj["name"], obj["reference"])
	sourceObj, err := c.findSourceByDatabaseRef(obj["reference"].(string))
	if err != nil {
		return err
	}
	err = c.updateMMSource(sourceObj["reference"].(string))
	return err
}

func (c *myClient) startVDBByName(vdbName string, wait bool) (results map[string]interface{}, err error) {
	namespace := "database"
	//Find our VDB of interest
	log.Info("Searching for VDB by name")
	obj, err := c.findObjectByName(namespace, vdbName)
	if err != nil {
		return nil, err
	}
	log.Debug(obj)
	if obj == nil {
		log.Fatalf("Could not find VDB named %s", vdbName)
	}
	log.Infof("Found %s: %s", obj["name"], obj["reference"])
	sourceObj, err := c.findSourceByDatabaseRef(obj["reference"].(string))
	if err != nil {
		return nil, err
	}
	//start
	namespace = "source"
	url := fmt.Sprintf("%s/%s/start", namespace, sourceObj["reference"])
	action, _, err := c.httpPost(url, "")
	if err != nil {
		return nil, err
	}

	if wait {
		c.jobWaiter(action)
	}
	return action, err
}

func (c *myClient) batchStartVDBByName(vdbList ...string) (resultsList []map[string]interface{}, err error) {
	for _, v := range vdbList {
		action, err := c.startVDBByName(v, false)
		if err != nil {
			return nil, err
		}
		resultsList = append(resultsList, action)
	}

	c.jobWaiter(resultsList...)

	return resultsList, err
}

func (c *myClient) updateMaskingConnector(user, pass string, port int) (err error) {
	logger := logger.WithFields(log.Fields{
		"url":      c.url,
		"username": c.username,
	})

	prodDBIP, err := getIP("proddb")
	if err != nil {
		return err
	}
	logger.Infof("updating connector host to %s", prodDBIP)
	putBody := fmt.Sprintf(`{
      "connectorName": "PatientsMM - PG",
      "databaseType": "POSTGRES",
      "environmentId": 1,
      "databaseName": "dafdb",
      "host": "%s",
      "port": %d,
      "schemaName": "public",
      "username": "%s",
      "kerberosAuth": false,
	  "password": "%s"
	}`, prodDBIP, port, user, pass)
	log.Debug(putBody)
	result, _, err := c.httpPut("database-connectors/1", putBody)
	if err != nil {
		return err
	}
	log.Debug(result)
	return err
}

func (c *myClient) updateMaskingService() (map[string]interface{}, error) {

	maskingIP, err := getIP("maskingengine")
	if err != nil {
		return nil, err
	}
	logger.Infof("updating masking service server to %s", maskingIP)

	result, _, err := c.httpPost("maskingjob/serviceconfig/MASKING_SERVICE_CONFIG-1", fmt.Sprintf(`{
			"type": "MaskingServiceConfig",
			"server": "%s",
			"port": 80,
			"credentials": {
				"type": "PasswordCredential",
				"password": "%s"
			}
		}`, maskingIP, opts.MaskPassword))
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug(result)
	return result, err
}

var (
	opts             Options
	parser           = flags.NewParser(&opts, flags.Default)
	apiVersionString = "1.9.3"
	logger           *log.Entry
	url              string
	version          = "not set"
)

func main() {
	var err error

	log.Info("Establishing session and logging in")

	virtualizationCR := NewClientRequest(opts.VirtUserName, opts.VirtPassword, fmt.Sprintf("https://%s/resources/json/delphix", opts.DDPVirtName))
	virtualizationClient := virtualizationCR.initResty()
	err = virtualizationClient.waitForEngineReady(10, 600)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Successfully Logged in")
	_, err = virtualizationClient.updateMaskingService()
	if err != nil {
		log.Fatal(err)
	}
	maskingCR := NewClientRequest(opts.MaskUserName, opts.MaskPassword, fmt.Sprintf("http://%s/masking/api", opts.DDPMaskName))
	maskingClient := maskingCR.initResty()
	maskingClient.waitForMaskingEngineReady(10, 600)
	err = maskingClient.updateMaskingConnector("delphixdb", "delphixdb", 5435)
	if err != nil {
		log.Fatal(err)
	}
	_, err = virtualizationClient.batchUpdateEnvironmentHostByHostName(opts.EnvironmentList)
	if err != nil {
		log.Fatal(err)
	}
	_, err = virtualizationClient.batchRefreshEnvironmentByName(opts.EnvironmentList)
	if err != nil {
		log.Fatal(err)
	}
	err = virtualizationClient.updatePGdSource(opts.DSourceList[0])
	if err != nil {
		log.Fatal(err)
	}
	err = virtualizationClient.updateMMConfigClone("Patients Masked Master")
	if err != nil {
		log.Fatal(err)
	}
	_, err = virtualizationClient.batchStartVDBByName(opts.VDBList...)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Complete")
}
