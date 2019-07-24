package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v1"
)

// Options for the script
type Options struct {
	DDPVirtName  string               `short:"e" long:"ddp_virt_hostname" env:"DELPHIX_DDP_VIRT_HOSTNAME" description:"The hostname or IP address of the Delphix Dynamic Data Platform - Virtualization Engine" required:"true"`
	DDPMaskName  string               `short:"m" long:"ddp_mask_hostname" env:"DELPHIX_DDP_MASK_HOSTNAME" description:"The hostname or IP address of the Delphix Dynamic Data Platform - Masking Engine" required:"true"`
	Password     string               `short:"p" long:"password" env:"DELPHIX_PASS" description:"The password used to authenticate to the Delphix Engine as the delphix_admin user" required:"true"`
	Debug        []bool               `short:"v" long:"debug" env:"DELPHIX_DEBUG" description:"Turn on debugging. -vvv for the most verbose debugging"`
	SkipValidate bool                 `long:"skip-validate" env:"DELPHIX_SKIP_VALIDATE" description:"Don't validate TLS certificate of Delphix Engine"`
	ConfigFile   func(s string) error `short:"c" long:"config" description:"Optional INI config file to pass in for the variables" no-ini:"true"`
	SysPass      string               `long:"sysadmin-password" env:"SYSADMIN_PASS" description:"The password used to authenticate to the Delphix Engine as the sysadmin user" required:"true"`
}

// VDBParams - parameters for constructing a VDB
type VDBParams struct {
	vdbName    string
	dbName     string
	groupName  string
	pdbName    string
	cdbName    string
	envName    string
	maskingJob string
}

// PatientContainer - parameters for constructing the Patients Data Pods
type PatientContainer struct {
	name       string
	vdbName    string
	sourceName string
	owners     []string
}

// User - parameters for constructing the Users
type User struct {
	name                  string
	password              string
	passwordUpdateRequest string
	roles                 []string
}

// createUser takes two parameters
// userName: the name of the user to create
// userPassword: the password to set for the user
func createUser(name, password, passwordUpdateRequest string) string {
	return fmt.Sprintf(`{
    "type": "User",
    "name": "%s",
    "credential": {
        "type": "PasswordCredential",
		"password": "%s",
		"passwordUpdateRequest": "%s"
    }
}`, name, password, passwordUpdateRequest)
}

// createAuthorization takes two parameters>
// roleRef: the reference of the role to assign to the user
// userRef: the reference of the user to assign to the role
func createAuthorization(roleRef, userRef string) string {
	return fmt.Sprintf(`{
    "type": "Authorization",
    "user": "%s",
    "role": "%s",
    "target": "%s"
}`, userRef, roleRef, userRef)
}

// createUserAndAuthorization takes four parameters
// user: a User struct
// wait: Wait for the action to complete before completing
func (c *Client) createUserAndAuthorizations(user User, wait bool) (results map[string]interface{}, err error) {
	userObj, err := c.createUser(user, true)
	if err != nil {
		return nil, err
	}

	userObjRef := userObj["reference"]
	if userObjRef == nil {
		return nil, fmt.Errorf("user %s was not found", user.name)
	}

	for _, roleName := range user.roles {

		roleObj, err := c.findObjectByName("role", roleName)
		if err != nil {
			return nil, err
		}

		roleObjRef := roleObj["reference"]

		c.createAuthorization(roleObjRef.(string), userObjRef.(string), true)
	}
	return userObj, err
}

// batchCreateUserAndAuthorizations takes one parameter:
// userList: a slice of JSUsers to create
func (c *Client) batchCreateUserAndAuthorizations(userList ...User) (resultsList []map[string]interface{}, err error) {
	for _, v := range userList {
		if action, err := c.createUserAndAuthorizations(v, true); action != nil && err == nil {
			resultsList = append(resultsList, action)
		} else if err != nil {
			return nil, err
		}
	}

	return resultsList, err
}

// createAuthorization takes three parameters>
// roleRef: the reference of the role to assign to the user
// userRef: the reference of the user to assign to the role
// wait: Wait for the action to complete before completing
func (c *Client) createAuthorization(roleRef, userRef string, wait bool) (results map[string]interface{}, err error) {
	namespace := "authorization"

	url := fmt.Sprintf("%s", namespace)
	postBody := createAuthorization(roleRef, userRef)
	action, _, err := c.httpPost(url, postBody)
	if err != nil {
		return nil, err
	}

	if wait {
		c.jobWaiter(action)
	}
	return action, err
}

// createUser takes two parameters>
// user: a User struct
// wait: Wait for the action to complete before completing
func (c *Client) createUser(user User, wait bool) (results map[string]interface{}, err error) {
	namespace := "user"
	userObj, err := c.findObjectByName(namespace, user.name)
	if err != nil {
		return userObj, err
	}

	userObjRef := userObj["reference"]
	if userObjRef == nil {
		url := fmt.Sprintf("%s", namespace)
		postBody := createUser(user.name, user.password, user.passwordUpdateRequest)
		action, _, err := c.httpPost(url, postBody)
		if err != nil {
			return nil, err
		}

		if wait {
			c.jobWaiter(action)
			userObj, err = c.findObjectByName(namespace, user.name)
		} else {
			return action, err
		}
	}
	log.Debug(userObjRef)
	log.Infof("%s already exists", user.name)
	return userObj, err

}

// discoverCDB takes five parameters:
// envName: The name of the Environment in Delphix that contains the CDB
// cdbName: The name of the database name of the CDB
// username: database username of the CDB (ie. c##delphixdb)
// password: the password of the CDB user
// wait: Wait for the action to complete before completing
func (c *Client) discoverCDB(envName, cdbName, username, password string, wait bool) (results map[string]interface{}, err error) {
	namespace := "sourceconfig"
	scObj, err := c.findSourceCongfigByDBNameAndEnvironmentName(cdbName, envName)
	scObjRef := scObj["reference"]
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s", namespace, scObjRef)
	postBody := fmt.Sprintf(`{
		"type": "OracleSIConfig", 
		"user": "%s",
		"credentials": {
			"type": "PasswordCredential",
			"password": "%s"
			}
		}`, username, password)
	action, _, err := c.httpPost(url, postBody)
	if err != nil {
		return nil, err
	}

	if wait {
		c.jobWaiter(action)
	}
	return action, err

}

// createDatasetGroup takes two parameters>
// groupName: the name of the dataset group to create
// wait: Wait for the action to complete before completing
func (c *Client) createDatasetGroup(groupName string, wait bool) (results map[string]interface{}, err error) {
	namespace := "group"
	groupObj, err := c.findObjectByName(namespace, groupName)
	if err != nil {
		return nil, err
	}

	groupObjRef := groupObj["reference"]
	if groupObjRef == nil {
		url := fmt.Sprintf("%s", namespace)
		postBody := fmt.Sprintf(`{
		"type": "Group", 
		"name": "%s"
		}`, groupName)
		action, _, err := c.httpPost(url, postBody)
		if err != nil {
			return nil, err
		}

		if wait {
			c.jobWaiter(action)
		}
		return action, err
	}
	log.Debug(groupObjRef)
	log.Infof("%s already exists", groupName)
	return nil, err

}

// batchcreateDatasetGroup takes one parameter:
// groupNameList: a slice of groupNames to create
func (c *Client) batchCreateDatasetGroup(groupNameList ...string) (resultsList []map[string]interface{}, err error) {
	for _, v := range groupNameList {
		if action, err := c.createDatasetGroup(v, false); err != nil {
			return nil, err
		} else if action != nil {
			resultsList = append(resultsList, action)
		}
	}
	c.jobWaiter(resultsList...)

	return resultsList, err
}

func (c *Client) linkDatabase(wait bool) (results map[string]interface{}, err error) {
	namespace := "database"
	dbName := "Patients Prod"
	groupName := "Prod"
	environmentUserName := "delphix"
	dbUser := "delphixdb"
	dbPass := "delphixdb"
	pdbName := "PATPDB"
	envName := "proddb"

	if dbObjRef, err := c.findObjectByNameReturnReference(namespace, dbName); dbObjRef == nil && err == nil {
		envRef, err := c.findObjectByNameReturnReference("environment", envName)
		if err != nil {
			return nil, err
		}
		if envRef == nil {
			logger.Fatalf("Environment %s not found", envName)
		}
		groupRef, err := c.findObjectByNameReturnReference("group", groupName)
		if err != nil {
			return nil, err
		}
		if groupRef == nil {
			logger.Fatalf("Group %s not found", groupName)
		}
		environmentUserRef, err := c.findObjectByNameReturnReference("environment/user", environmentUserName, fmt.Sprintf("environment=%s", envRef))
		if err != nil {
			return nil, err
		}
		if environmentUserRef == nil {
			logger.Fatalf("User %s not found", environmentUserName)
		}
		scObjRef, err := returnObjReference(c.findSourceCongfigByDBNameAndEnvironmentName(pdbName, envName))
		if err != nil {
			return nil, err
		}
		if scObjRef == nil {
			logger.Fatalf("PDB %s not found on %s", pdbName, envName)
		}
		postBody := fmt.Sprintf(`{
			"type": "LinkParameters",
			"name": "%s",
			"group": "%s",
			"linkData": {
				"type": "OraclePDBLinkData",
				"config": "%s",
				"linkNow": false,
				"environmentUser": "%s",
				"dbUser": "%s",
				"dbCredentials": {
					"type": "PasswordCredential",
					"password": "%s"
				}
			}
		}`, dbName, groupRef, scObjRef, environmentUserRef, dbUser, dbPass)
		url := fmt.Sprintf("%s/link", namespace)
		action, _, err := c.httpPost(url, postBody)
		if err != nil {
			return nil, err
		}
		logger.Debugf("DBLINK ACTION: %v", action)
		err = c.jobWaiter(action)
		if err != nil {
			return nil, err
		}
		url = fmt.Sprintf("%s/%s/sync", namespace, action["result"])
		syncAction, _, err := c.httpPost(url, "")
		if err != nil {
			return nil, err
		}
		if wait {
			c.jobWaiter(syncAction)
		}
		return action, err
	} else if err != nil {
		return nil, err
	} else {
		log.Debug(dbObjRef)
		log.Infof("%s already exists", dbName)
		return nil, err
	}
}

// turnOnMasking deprecated with >=5.3.3
// func turnOnMasking() {
// 	logger := log.WithFields(log.Fields{
// 		"connection": fmt.Sprintf("ssh://%s:22", opts.DDPVirtName),
// 		"username":   "sysadmin",
// 	})
// 	logger.Info("Attempting to start masking")
// 	config := &ssh.ClientConfig{
// 		User: "sysadmin",
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(opts.SysPass),
// 		},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}

// 	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", opts.DDPVirtName), config)
// 	if err != nil {
// 		logger.Fatal("Unable to create client: ", err)
// 	}
// 	defer client.Close()

// 	session, err := client.NewSession()
// 	if err != nil {
// 		logger.Fatal("Unable to create session: ", err)
// 	}
// 	defer session.Close()

// 	err = session.Run("cd system; startMasking; commit")
// 	if err != nil {
// 		logger.Fatal("Unable to start masking: ", err)
// 	}
// 	logger.Info("Masking started")
// }

func (c *Client) linkMaskingJob() (results map[string]interface{}, err error) {
	namespace := "maskingjob"
	dbName := "Patients Prod"
	maskingName := "Patients Mask"
	maskingRef, err := c.findObjectByNameReturnReference("maskingjob", maskingName)
	if err != nil {
		return nil, err
	}
	if maskingRef == nil {
		logger.Fatalf("Masking Job %s not found", maskingName)
	}
	dbRef, err := c.findObjectByNameReturnReference("database", dbName)
	if err != nil {
		return nil, err
	}
	if dbRef == nil {
		logger.Fatalf("Database %s not found", dbName)
	}
	postBody := fmt.Sprintf(`{
        "type": "MaskingJob",
        "associatedContainer": "%s"
	}`, dbRef)
	url := fmt.Sprintf("%s/%s", namespace, maskingRef)
	action, _, err := c.httpPost(url, postBody)
	if err != nil {
		return nil, err
	}
	c.jobWaiter(action)
	return action, err
}

func (c *Client) provisionVDB(vdbParams VDBParams, wait bool) (results map[string]interface{}, err error) {
	namespace := "database"
	maskingJSON := ""
	var maskingJobRef interface{}

	if vdbObjRef, err := c.findObjectByNameReturnReference(namespace, vdbParams.vdbName); vdbObjRef == nil && err == nil {
		pdbRef, err := c.findObjectByNameReturnReference(namespace, vdbParams.pdbName)
		if err != nil {
			return nil, err
		}
		if pdbRef == nil {
			logger.Fatalf("PDB %s not found", vdbParams.pdbName)
		}
		envRef, err := c.findObjectByNameReturnReference("environment", vdbParams.envName)
		if err != nil {
			return nil, err
		}
		if envRef == nil {
			logger.Fatalf("Environment %s not found", vdbParams.envName)
		}
		groupRef, err := c.findObjectByNameReturnReference("group", vdbParams.groupName)
		if err != nil {
			return nil, err
		}
		if groupRef == nil {
			logger.Fatalf("Group %s not found", vdbParams.groupName)
		}
		scObjRef, err := returnObjReference(c.findSourceCongfigByDBNameAndEnvironmentName(vdbParams.cdbName, vdbParams.envName))
		if err != nil {
			return nil, err
		}
		if scObjRef == nil {
			logger.Fatalf("CDB %s not found on %s", vdbParams.cdbName, vdbParams.envName)
		}
		if vdbParams.maskingJob != "" {
			maskingJobRef, err = c.findObjectByNameReturnReference("maskingjob", vdbParams.maskingJob)
			if err != nil {
				return nil, err
			}
			if maskingJobRef == nil {
				logger.Fatalf("Masking Job %s not found", vdbParams.maskingJob)
			}
			maskingJSON = fmt.Sprintf(`,
			"maskingJob": "%s"
			`, maskingJobRef)
		} else {
			maskingJobRef = nil
		}
		logger.Debugf("MASKING JOB NAME: %v", vdbParams.maskingJob)
		logger.Debugf("MASKING JOB REF: %v", maskingJobRef)
		postBody := fmt.Sprintf(`{
			"type": "OracleMultitenantProvisionParameters",
			"container": {
				"type": "OracleDatabaseContainer",
				"name": "%s",
				"group": "%s"
			},
			"source": {
				"type": "OracleVirtualPdbSource",
				"mountBase": "/mnt/provision",
				"allowAutoVDBRestartOnHostReboot": true
			},
			"sourceConfig": {
				"type": "OraclePDBConfig",
				"databaseName": "%s",
				"cdbConfig": "%s"
			},
			"timeflowPointParameters": {
				"type": "TimeflowPointSemantic",
				"location": "LATEST_SNAPSHOT",
				"container": "%s"
			}%s
		}`, vdbParams.vdbName, groupRef, vdbParams.dbName, scObjRef, pdbRef, maskingJSON)
		// logger.Fatal(postBody)
		url := fmt.Sprintf("%s/provision", namespace)
		action, _, err := c.httpPost(url, postBody)
		if err != nil {
			return nil, err
		}
		if wait {
			c.jobWaiter(action)
		}
		return action, err
	} else if err != nil {
		return nil, err
	} else {
		log.Debug(vdbObjRef)
		log.Infof("%s already exists", vdbParams.vdbName)
		return nil, err
	}
}

// batchProvisionVDB takes one parameter:
// vdbParamsList: a slice of VDBParams to create
func (c *Client) batchProvisionVDB(vdbParamsList ...VDBParams) (resultsList []map[string]interface{}, err error) {
	for _, v := range vdbParamsList {
		if action, err := c.provisionVDB(v, false); action != nil && err == nil {
			resultsList = append(resultsList, action)
		} else if err != nil {
			return nil, err
		}
	}
	err = c.jobWaiter(resultsList...)

	return resultsList, err
}

func (c *Client) createSelfServiceTemplate(wait bool) (resultsList map[string]interface{}, err error) {
	templateName := "Patients"
	namespace := "selfservice/template"
	vdbName := "Patients Masked Master"
	sourceName := "Masked Master"
	if templateRef, err := c.findObjectByNameReturnReference(namespace, templateName); templateRef == nil && err == nil {
		vdbRef, err := c.findObjectByNameReturnReference("database", vdbName)
		if err != nil {
			return nil, err
		}
		if vdbRef == nil {
			logger.Fatalf("VDB %s not found", vdbName)
		}
		url := fmt.Sprintf("%s", namespace)
		postBody := fmt.Sprintf(`{
		"type": "JSDataTemplateCreateParameters", 
		"name": "%s",
		"dataSources": [{
			"type": "JSDataSourceCreateParameters",
			"container": "%s",
			"source": {
				"type": "JSDataSource",
				"name": "%s"
			}
		}]
		}`, templateName, vdbRef, sourceName)
		action, _, err := c.httpPost(url, postBody)
		if wait {
			err = c.jobWaiter(action)
		}
		return action, err
	} else if err != nil {
		return nil, err
	} else {
		log.Debug(templateRef)
		log.Infof("%s already exists", templateName)
		return nil, err
	}
}

func (c *Client) createSelfServiceContainer(container PatientContainer, wait bool) (resultsList map[string]interface{}, err error) {
	templateName := "Patients"
	namespace := "selfservice/container"
	var ownerRefs []string
	if container.owners == nil {
		logger.Fatal("This container did not have any owners specified")
	}
	if containerRef, err := c.findObjectByNameReturnReference(namespace, container.name); containerRef == nil && err == nil {
		templateRef, err := c.findObjectByNameReturnReference("selfservice/template", templateName)
		if err != nil {
			return nil, err
		}
		if templateRef == nil {
			logger.Fatalf("Template %s not found", templateName)
		}
		vdbRef, err := c.findObjectByNameReturnReference("database", container.vdbName)
		if err != nil {
			return nil, err
		}
		if vdbRef == nil {
			logger.Fatalf("VDB %s not found", container.vdbName)
		}
		for _, v := range container.owners {
			userRef, err := c.findObjectByNameReturnReference("user", v)
			if err != nil {
				return nil, err
			}
			if userRef == nil {
				logger.Fatalf("User %s not found", v)
			}
			ownerRefs = append(ownerRefs, strconv.Quote(userRef.(string)))
		}
		url := fmt.Sprintf("%s", namespace)
		postBody := fmt.Sprintf(`{
		"type": "JSDataContainerCreateWithoutRefreshParameters", 
		"name": "%s",
		"template": "%s",
		"dataSources": [{
			"type": "JSDataSourceCreateParameters",
			"container": "%s",
			"source": {
				"type": "JSDataSource",
				"name": "%s"
			}
		}],
		"owners": [%s]
		}`, container.name, templateRef, vdbRef, container.sourceName, strings.Join(ownerRefs, ","))
		action, _, err := c.httpPost(url, postBody)
		if err != nil {
			return nil, err
		}
		if wait {
			c.jobWaiter(action)
		}
		return action, err
	} else if err != nil {
		return nil, err
	} else {
		log.Debug(containerRef)
		log.Infof("%s already exists", container.name)
		return nil, err
	}
}

func (c *Client) populateMasking() (err error) {
	logger := logger.WithFields(log.Fields{
		"url":      c.url,
		"username": c.username,
	})

	appName := "Patients Application"
	envName := "Patients Environment"
	appExists := false
	envExists := false
	globFileName := "global_objects.json"
	mjFileName := "masking_job.json"
	glObjFile, err := ioutil.ReadFile(globFileName)
	if err != nil {
		logger.Fatal("Unable to open: ", globFileName, err)
	}
	logger.Info("Uploading global objects")
	resty.DefaultClient.
		SetTimeout(time.Duration(60 * time.Second))
	c.httpPostBytesReturnSlice("import", glObjFile, "force_overwrite=true")
	envQuery, _, err := c.httpGet("environments")
	if err != nil {
		return err
	}
	envList := envQuery["responseList"].([]interface{})
	for _, v := range envList {
		if v.(map[string]interface{})["environmentName"] == envName {
			envExists = true
			logger.Infof("%s and %s already exists", appName, envName)
		}
	}
	if envExists == false {
		appQuery, _, err := c.httpGet("applications")
		if err != nil {
			return err
		}
		appList := appQuery["responseList"].([]interface{})
		for _, v := range appList {
			if v.(map[string]interface{})["applicationName"] == appName {
				appExists = true
				logger.Infof("%s already exists", appName)
			}
		}
		if appExists == false {
			c.httpPost("applications", fmt.Sprintf("{\"applicationName\":\"%s\"}", appName))
		}
		c.httpPost("environments", fmt.Sprintf(`{
			"environmentName": "%s",
			"application": "%s",
			"purpose": "MASK"
		}`, envName, appName))
	}
	mjFile, err := ioutil.ReadFile(mjFileName)
	if err != nil {
		logger.Fatal("Unable to open: ", mjFileName, err)
	}
	logger.Info("Uploading masking job")
	c.httpPostBytesReturnSlice("import", mjFile, "force_overwrite=true", "environment_id=1")
	logger.Info("updating connector")
	c.httpPut("database-connectors/1", fmt.Sprintf(`
	{
		"connectorName": "Patients Prod - Do Not Mask",
		"databaseType": "ORACLE",
		"environmentId": 1,
		"jdbc": "jdbc:oracle:thin:@10.0.1.20:1521/patpdb",
		"schemaName": "DELPHIXDB",
		"username": "DELPHIXDB",
		"password": "delphixdb"
	}`))
	return err
}

// batchCreateSelfServiceContainer takes one parameter:
// containerList: a slice of PatientContainer to create
func (c *Client) batchCreateSelfServiceContainer(containerList ...PatientContainer) (resultsList []map[string]interface{}, err error) {
	for _, v := range containerList {
		if action, err := c.createSelfServiceContainer(v, false); action != nil && err == nil {
			resultsList = append(resultsList, action)
		} else if err != nil {
			return nil, err
		}
	}
	c.jobWaiter(resultsList...)

	return resultsList, err
}

// updateUserPasswordByName updates the specified username's(u) password(p)
func (c *Client) updateUserPasswordByName(u string, p string) (err error) {
	userRef, err := c.findObjectByNameReturnReference("user", u)
	if err != nil {
		return err
	}

	if userRef == nil {
		logger.Fatalf("%s not found. Exiting", u)
	}
	logger.Infof("Changing %s's password\n", u)
	logger.Infof("Using " + c.username)
	postBody := fmt.Sprintf(`
			{
				"type": "CredentialUpdateParameters",
				"newCredential": {
					"type": "PasswordCredential",
					"password": "%s"
				}
			}
		`, p)
	_, _, err = c.httpPost(fmt.Sprintf("/user/%s/updateCredential", userRef), postBody)
	if err != nil {
		return err
	}
	if c.username == u {
		c.password = p
		c.LoadAndValidate()
	}
	return err
}

// getStorageDevices returns the list of unassigned storage devices on the Delphix engine
func (c *Client) getStorageDevices() (devices []string, err error) {

	results, err := c.listObjects("storage/device") //grab the query results
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		err = fmt.Errorf("No devices available for assignment into the domain")
		return nil, err
	}
	logger.Info("The following devices are available:")
	for _, result := range results { //loop through the results
		if result.(map[string]interface{})["configured"] != true { //if the device is not already configured
			logger.Info(result.(map[string]interface{})["reference"]) //grab the device reference
			devices = append(devices, result.(map[string]interface{})["reference"].(string))
		}
	}
	return devices, nil
}

// initializeSystem  adds all available storage devices into the pool, and complete initial setup Wizard
func (c *Client) initializeSystem(d string, p string) (b bool, err error) {
	devices, err := c.getStorageDevices()
	if err != nil {
		return
	}
	if len(devices) == 0 {
		logger.Info("All devices are configured already")
		return
	}
	var deviceString string
	for i, device := range devices {
		device = fmt.Sprintf("%s", strconv.Quote(device))
		if i == 0 {
			deviceString = device
			continue
		}
		deviceString = fmt.Sprintf("%s,%s", deviceString, device)
	}
	postBody := fmt.Sprintf(`
		{
					"type": "SystemInitializationParameters",
					"defaultUser": "%s",
					"defaultPassword": "%s",
					"devices": [
						%s
					]
				}
		`, d, p, deviceString)
	logger.Info("Assigning devices to storage domain")
	action, _, err := c.httpPost("/domain/initializeSystem", postBody)
	if err != nil {
		return
	}
	c.jobWaiter(action)
	c.waitForEngineReady(10, 300)
	return true, err
}

func (c *Client) checkForEnvironments(envList []string) (err error) {
	namespace := "environment"
	for _, envName := range envList {
		envRef, err := c.findObjectByNameReturnReference(namespace, envName)
		if err != nil {
			return err
		}
		if envRef == nil {
			err = fmt.Errorf("Environment %s not found", envName)
		}
	}
	return err
}

func (c *Client) checkForBusyEnvironments(envList []string) (err error) {
	namespace := "environment"
	for _, envName := range envList {
		envRef, err := c.findObjectByNameReturnReference(namespace, envName)
		if err != nil {
			return err
		}
		if envRef == nil {
			err := fmt.Errorf("Environment %s not found", envName)
			logger.Info(err)
			return err
		}
		runningJobs, err := c.listObjects("job", "jobState=RUNNING", fmt.Sprintf("target=%s", envRef))
		if err != nil {
			return err
		}
		if len(runningJobs) > 0 {
			err = fmt.Errorf("Environment %s is busy", envName)
			logger.Info(err)
			return err
		}
	}

	return err
}

// waitForEnvironmentsReady loops until the environments are ready or time (t) expires
// p: period (interval) of time between attempts
// t: total time to try
func (c *Client) waitForEnvironmentsReady(p, t int, envList ...string) (err error) {

	logger.Infof("Waiting up to %v seconds for the environments to be ready", t)
	timeOut := 0
	for timeOut < t {
		logger.Info("Waiting for the environments")
		time.Sleep(time.Duration(p) * time.Second)
		timeOut = timeOut + p
		if err = c.checkForBusyEnvironments(envList); err == nil {
			break
		}
	}
	if timeOut >= t {
		err = fmt.Errorf("waitForEnvironmentsReady timed out")
	}
	return err
}

func shutdownHTTPServer(srv *http.Server) {
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

func initializePlatform() (err error) {
	logger.Info("Establishing session and logging in to Virtualization Engine")
	sysClient := NewClient("sysadmin", "sysadmin", fmt.Sprintf("https://%s/resources/json/delphix", opts.DDPVirtName))
	sysClient.initResty()

	err = sysClient.waitForEngineReady(10, 300)
	if err != nil && err.Error() == "401" {
		logger.Info("Looks like the engine may have previously been configured. Moving on to next phase.")
		err = nil
	} else {
		err = sysClient.updateUserPasswordByName("sysadmin", opts.SysPass)
		if err != nil {
			return err
		}
		didInitialize, err := sysClient.initializeSystem("delphix_admin", "delphix")
		if err != nil {
			return err
		}
		if didInitialize == true {
			logger.Info("Configuring for Virtualization")
			result, _, err := sysClient.httpPost("system", `{
			"type": "SystemInfo",
			"engineType": "VIRTUALIZATION"
		}`)
			if err != nil {
				return err
			}
			logger.Debug(result)
		}
	}

	logger = log.WithFields(log.Fields{
		"url":      fmt.Sprintf("https://%s", opts.DDPMaskName),
		"username": "admin",
	})
	logger.Info("Establishing session and logging in to Masking Engine")
	sysClient = NewClient("sysadmin", "sysadmin", fmt.Sprintf("https://%s/resources/json/delphix", opts.DDPMaskName))
	sysClient.initResty()

	err = sysClient.waitForEngineReady(10, 300)
	if err != nil && err.Error() == "401" {
		logger.Info("Looks like the engine may have previously been configured. Moving on to next phase.")
		err = nil
	} else {
		err = sysClient.updateUserPasswordByName("sysadmin", opts.SysPass)
		if err != nil {
			return err
		}
		didInitialize, err := sysClient.initializeSystem("admin", opts.Password)
		if err != nil {
			return err
		}
		if didInitialize == true {
			logger.Info("Configuring for Masking")
			resty.DefaultClient.
				SetTimeout(time.Duration(180 * time.Second))
			result, _, err := sysClient.httpPost("system", `{
		"type": "SystemInfo",
		"engineType": "MASKING"
	}`)
			if err != nil {
				return err
			}
			logger.Debug(result)
		}
	}

	return nil
}

func (c *Client) setInitialPlatformPassword() {
	logger = log.WithFields(log.Fields{
		"url":      c.url,
		"username": c.username,
	})
	statusCode, err := c.LoadAndValidate()
	if err != nil && statusCode == 401 {
		logger.Info("Looks like the engine may have previously been configured. Moving on to next phase.")
		c.password = opts.Password
		c.initResty()
		_, err = c.LoadAndValidate()
		if err != nil {
			logger.Fatal(err)
		}

	} else {
		err := c.updateUserPasswordByName(c.username, opts.Password)
		if err != nil {
			logger.Fatal(err)
		}
	}
}

func (c *Client) setInitialMaskingPassword() {
	logger := logger.WithFields(log.Fields{
		"url":      c.url,
		"username": c.username,
	})
	statusCode, err := c.MaskingLoadAndValidate()
	if err != nil && statusCode == 401 {
		logger.Info("Looks like the engine may have previously been configured. Moving on to next phase.")
		c.password = opts.Password
		c.initResty()
		_, err = c.MaskingLoadAndValidate()
		if err != nil {
			logger.Fatal(err)
		}

	} else {
		logger.Infof("Updating Masking user: %s", c.username)
		result, _, err := c.httpPut("users/5", fmt.Sprintf(`{
			"userName": "%s",
			"password": "%s",
			"firstName": "First",
			"lastName": "Last",
			"email": "user@delphix.com",
			"isAdmin": true,
			"showWelcome": false,
			"isLocked": false
		}`, c.username, opts.Password))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Debug(result)
	}
}

var (
	opts             Options
	parser           = flags.NewParser(&opts, flags.Default)
	apiVersionString = "1.9.3"
	logger           *log.Entry
	url              string
	version          = "undefined"
)

func main() {
	var err error

	err = initializePlatform()
	if err != nil {
		logger.Fatal(err)
	}
	virtualizationClient := NewClient("delphix_admin", "delphix", fmt.Sprintf("https://%s/resources/json/delphix", opts.DDPVirtName))
	virtualizationClient.initResty()
	virtualizationClient.setInitialPlatformPassword()

	maskingClient := NewClient("admin", "Admin-12", fmt.Sprintf("https://%s/resources/json/delphix", opts.DDPMaskName))
	maskingClient.initResty()
	maskingClient.setInitialPlatformPassword()

	maskingClient = NewClient("admin", "Admin-12", fmt.Sprintf("http://%s/masking/api", opts.DDPMaskName))
	maskingClient.initResty()
	maskingClient.setInitialMaskingPassword()

	logger.Info("Configuring Masking Service")
	result, _, err := virtualizationClient.httpPost("maskingjob/serviceconfig/MASKING_SERVICE_CONFIG-1", fmt.Sprintf(`{
			"type": "MaskingServiceConfig",
			"server": "10.0.1.11",
			"port": 80,
			"credentials": {
				"type": "PasswordCredential",
				"password": "%s"
			}
		}`, opts.Password))
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug(result)

	srv := startHTTPServer()
	defer shutdownHTTPServer(srv)
	err = virtualizationClient.waitForEnvironmentsReady(10, 300, "proddb", "devdb")

	if err != nil {
		logger.Fatal(err)
	}

	devUser := User{
		name:                  "dev",
		password:              "delphix",
		passwordUpdateRequest: "NONE",
		roles:                 []string{"Self-Service User"},
	}
	qaUser := User{
		name:                  "qa",
		password:              "delphix",
		passwordUpdateRequest: "NONE",
		roles:                 []string{"Self-Service User"},
	}
	_, err = virtualizationClient.batchCreateUserAndAuthorizations(devUser, qaUser)

	_, err = virtualizationClient.discoverCDB("proddb", "patcdb", "c##delphixdb", "delphixdb", true)
	if err != nil {
		logger.Fatal(err)
	}
	_, err = virtualizationClient.discoverCDB("devdb", "tcdb", "c##delphixdb", "delphixdb", true)
	if err != nil {
		logger.Fatal(err)
	}
	_, err = virtualizationClient.batchCreateDatasetGroup("Prod", "Masked Masters", "Non Prod")
	if err != nil {
		logger.Fatal(err)
	}
	_, err = virtualizationClient.linkDatabase(true)
	if err != nil {
		logger.Fatal(err)
	}

	err = maskingClient.populateMasking()
	if err != nil {
		logger.Fatal(err)
	}

	returnVal, _, err := virtualizationClient.httpGet("maskingjob/fetch")
	if err != nil {
		logger.Fatal(err)
	}
	err = virtualizationClient.jobWaiter(returnVal)
	if err != nil {
		logger.Fatal(err)
	}
	// logger.Fatal(returnVal)
	_, err = virtualizationClient.linkMaskingJob()
	if err != nil {
		logger.Fatal(err)
	}

	patmm := VDBParams{
		vdbName:    "Patients Masked Master",
		dbName:     "patmm",
		groupName:  "Masked Masters",
		pdbName:    "Patients Prod",
		cdbName:    "patcdb",
		envName:    "proddb",
		maskingJob: "Patients Mask",
	}

	_, err = virtualizationClient.provisionVDB(patmm, true)
	if err != nil {
		logger.Fatal(err)
	}

	devdb := VDBParams{
		vdbName:   "Patients Dev",
		dbName:    "devdb",
		groupName: "Non Prod",
		pdbName:   "Patients Masked Master",
		cdbName:   "tcdb",
		envName:   "devdb",
	}

	testdb := VDBParams{
		vdbName:   "Patients Test",
		dbName:    "testdb",
		groupName: "Non Prod",
		pdbName:   "Patients Masked Master",
		cdbName:   "tcdb",
		envName:   "devdb",
	}

	prep1 := VDBParams{
		vdbName:   "prep1",
		dbName:    "prep1",
		groupName: "Non Prod",
		pdbName:   "Patients Masked Master",
		cdbName:   "tcdb",
		envName:   "devdb",
	}

	prep2 := VDBParams{
		vdbName:   "prep2",
		dbName:    "prep2",
		groupName: "Non Prod",
		pdbName:   "Patients Masked Master",
		cdbName:   "tcdb",
		envName:   "devdb",
	}

	prep3 := VDBParams{
		vdbName:   "prep3",
		dbName:    "prep3",
		groupName: "Non Prod",
		pdbName:   "Patients Masked Master",
		cdbName:   "tcdb",
		envName:   "devdb",
	}

	prep4 := VDBParams{
		vdbName:   "prep4",
		dbName:    "prep4",
		groupName: "Non Prod",
		pdbName:   "Patients Masked Master",
		cdbName:   "tcdb",
		envName:   "devdb",
	}

	prep5 := VDBParams{
		vdbName:   "prep5",
		dbName:    "prep5",
		groupName: "Non Prod",
		pdbName:   "Patients Masked Master",
		cdbName:   "tcdb",
		envName:   "devdb",
	}

	prep6 := VDBParams{
		vdbName:   "prep6",
		dbName:    "prep6",
		groupName: "Non Prod",
		pdbName:   "Patients Masked Master",
		cdbName:   "tcdb",
		envName:   "devdb",
	}

	_, err = virtualizationClient.batchProvisionVDB(devdb, testdb, prep1, prep2, prep3, prep4, prep5, prep6)
	if err != nil {
		logger.Fatal(err)
	}

	_, err = virtualizationClient.createSelfServiceTemplate(true)
	if err != nil {
		logger.Fatal(err)
	}

	developContainer := PatientContainer{
		vdbName:    devdb.vdbName,
		sourceName: "Dev",
		name:       "Develop",
		owners:     []string{"delphix_admin", "dev"},
	}

	testContainer := PatientContainer{
		vdbName:    testdb.vdbName,
		sourceName: "Test",
		name:       "Test",
		owners:     []string{"delphix_admin", "qa"},
	}

	prep1Container := PatientContainer{
		vdbName:    prep1.vdbName,
		sourceName: prep1.vdbName,
		name:       prep1.vdbName,
		owners:     []string{"delphix_admin"},
	}

	prep2Container := PatientContainer{
		vdbName:    prep2.vdbName,
		sourceName: prep2.vdbName,
		name:       prep2.vdbName,
		owners:     []string{"delphix_admin"},
	}

	prep3Container := PatientContainer{
		vdbName:    prep3.vdbName,
		sourceName: prep3.vdbName,
		name:       prep3.vdbName,
		owners:     []string{"delphix_admin"},
	}

	prep4Container := PatientContainer{
		vdbName:    prep4.vdbName,
		sourceName: prep4.vdbName,
		name:       prep4.vdbName,
		owners:     []string{"delphix_admin"},
	}

	prep5Container := PatientContainer{
		vdbName:    prep5.vdbName,
		sourceName: prep5.vdbName,
		name:       prep5.vdbName,
		owners:     []string{"delphix_admin"},
	}

	prep6Container := PatientContainer{
		vdbName:    prep6.vdbName,
		sourceName: prep6.vdbName,
		name:       prep6.vdbName,
		owners:     []string{"delphix_admin"},
	}

	_, err = virtualizationClient.batchCreateSelfServiceContainer(developContainer, testContainer, prep1Container, prep2Container, prep3Container, prep4Container, prep5Container, prep6Container)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Complete")
}
