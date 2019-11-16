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

// PGVDBParams - parameters for constructing an AppData VDB
type PGVDBParams struct {
	configClone         string
	dbName              string
	environmentUserName string
	groupName           string
	masked              bool
	port                int
	sourceDBName        string
	sourceConfig        AppDataStagedSourceConfig
	vdbName             string
}

// PVDBParams - parameters for constructing a PDB VDB
type PVDBParams struct {
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

// LinkParams - parameters for constructing a dSource
type LinkParams struct {
	dbName              string
	groupName           string
	environmentUserName string
	sourceConfig        AppDataStagedSourceConfig
	dbUser              string
	dbPass              string
	preSync             string
}

// AppDataStagedSourceConfig - parameters for constructing a Source Config object
type AppDataStagedSourceConfig struct {
	linkingEnabled bool
	name           string
	parameters     string
	dbPath         string
	repoName       string
	envName        string
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
func (c *myClient) createUserAndAuthorizations(user User, wait bool) (results map[string]interface{}, err error) {
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
func (c *myClient) batchCreateUserAndAuthorizations(userList ...User) (resultsList []map[string]interface{}, err error) {
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
func (c *myClient) createAuthorization(roleRef, userRef string, wait bool) (results map[string]interface{}, err error) {
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
func (c *myClient) createUser(user User, wait bool) (results map[string]interface{}, err error) {
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
func (c *myClient) discoverCDB(envName, cdbName, username, password string, wait bool) (results map[string]interface{}, err error) {
	namespace := "sourceconfig"
	scObj, err := c.findSourceCongfigByNameAndEnvironmentName(cdbName, envName)
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
func (c *myClient) createDatasetGroup(groupName string, wait bool) (results map[string]interface{}, err error) {
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
func (c *myClient) batchCreateDatasetGroup(groupNameList ...string) (resultsList []map[string]interface{}, err error) {
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

// Oracle PDB
func (c *myClient) linkOracle12PDB(wait bool) (results map[string]interface{}, err error) {
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
		scObjRef, err := returnObjReference(c.findSourceCongfigByNameAndEnvironmentName(pdbName, envName))
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

// Postgres Database
func (c *myClient) linkPostgresDatabase(linkParams LinkParams, wait bool) (results map[string]interface{}, err error) {
	namespace := "database"
	if dbObjRef, err := c.findObjectByNameReturnReference(namespace, linkParams.dbName); dbObjRef == nil && err == nil {
		envRef, err := c.findObjectByNameReturnReference("environment", linkParams.sourceConfig.envName)
		if err != nil {
			return nil, err
		}
		if envRef == nil {
			logger.Fatalf("Environment %s not found", linkParams.sourceConfig.envName)
		}
		groupRef, err := c.findObjectByNameReturnReference("group", linkParams.groupName)
		if err != nil {
			return nil, err
		}
		if groupRef == nil {
			logger.Fatalf("Group %s not found", linkParams.groupName)
		}
		environmentUserRef, err := c.findObjectByNameReturnReference("environment/user", linkParams.environmentUserName, fmt.Sprintf("environment=%s", envRef))
		if err != nil {
			return nil, err
		}
		if environmentUserRef == nil {
			logger.Fatalf("User %s not found", linkParams.environmentUserName)
		}
		scObj, err := c.createSourceConfig(linkParams.sourceConfig, true)
		if err != nil {
			return nil, err
		}
		scObjRef, ok := scObj["result"].(string)
		if !ok {
			scObjRef = scObj["reference"].(string)
		}

		log.Debug(scObj)
		log.Debug(scObjRef)
		prodDBIP, err := getIP("proddb")
		if err != nil {
			return nil, err
		}

		postBody := fmt.Sprintf(`{
		"name": "%s",
		"group": "%s",
		"description": "",
		"linkData": {
			"config": "%s",
			"stagingMountBase": "/var/lib/pgsql/staging",
			"stagingEnvironment": "%s",
			"stagingEnvironmentUser": "%s",
			"environmentUser": "%s",
			"operations": {
			"preSync": [%s],
			"postSync": [],
			"type": "LinkedSourceOperations"
			},
			"parameters": {
			"delphixInitiatedBackupFlag": true,
			"keepStagingInSync": true,
			"postgresPort": 5433,
			"externalBackup": [],
			"delphixInitiatedBackup": [
				{
				"userName": "%s",
				"postgresSourcePort": 5432,
				"userPass": "%s",
				"sourceHostAddress": "%s"
				}
			],
			"configSettingsStg": []
			},
			"sourcingPolicy": { "logsyncEnabled": false, "type": "SourcingPolicy" },
			"type": "AppDataStagedLinkData"
		},
		"type": "LinkParameters"
		}`, linkParams.dbName, groupRef, scObjRef, envRef, environmentUserRef, environmentUserRef, linkParams.preSync, linkParams.dbUser, linkParams.dbPass, prodDBIP)
		log.Debug(postBody)
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

		syncAction, err := c.findActionByTitleAndParentAction("DB_SYNC", action["action"].(string))
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
		log.Infof("%s already exists", linkParams.dbName)
		results := make(map[string]interface{})
		results["result"] = dbObjRef.(string)
		return results, err
	}
}

func (c *myClient) linkMaskingJob() (results map[string]interface{}, err error) {
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

// provisionPVDB provisions virtual PDB into existing CDB's
func (c *myClient) provisionPVDB(vdbParams PVDBParams, wait bool) (results map[string]interface{}, err error) {
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
		scObjRef, err := returnObjReference(c.findSourceCongfigByNameAndEnvironmentName(vdbParams.cdbName, vdbParams.envName))
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

// provisionVDB provisions VDBs
func (c *myClient) provisionVDB(vdbParams PGVDBParams, wait bool) (results map[string]interface{}, err error) {
	namespace := "database"

	if vdbObjRef, err := c.findObjectByNameReturnReference(namespace, vdbParams.vdbName); vdbObjRef == nil && err == nil {
		sourceDBRef, err := c.findObjectByNameReturnReference(namespace, vdbParams.sourceDBName)
		if err != nil {
			return nil, err
		}
		if sourceDBRef == nil {
			logger.Fatalf("Source Database %s not found", vdbParams.sourceDBName)
		}
		envRef, err := c.findObjectByNameReturnReference("environment", vdbParams.sourceConfig.envName)
		if err != nil {
			return nil, err
		}
		if envRef == nil {
			logger.Fatalf("Environment %s not found", vdbParams.sourceConfig.envName)
		}
		groupRef, err := c.findObjectByNameReturnReference("group", vdbParams.groupName)
		if err != nil {
			return nil, err
		}
		if groupRef == nil {
			logger.Fatalf("Group %s not found", vdbParams.groupName)
		}
		scObj, err := c.createSourceConfig(vdbParams.sourceConfig, true)
		if err != nil {
			return nil, err
		}
		scRepoRef := scObj["repository"]
		if scRepoRef == nil {
			logger.Fatalf("Repository for SourceConfig %s not found on %s", vdbParams.sourceConfig.name, vdbParams.sourceConfig.envName)
		}
		environmentUserRef, err := c.findObjectByNameReturnReference("environment/user", vdbParams.environmentUserName, fmt.Sprintf("environment=%s", envRef))
		if err != nil {
			return nil, err
		}
		if environmentUserRef == nil {
			logger.Fatalf("User %s not found", vdbParams.environmentUserName)
		}

		postBody := fmt.Sprintf(`{
			"container": {
				"sourcingPolicy": { "logsyncEnabled": false, "type": "SourcingPolicy" },
				"group": "%s",
				"name": "%s",
				"type": "AppDataContainer"
			},
			"source": {
				"operations": {
				"configureClone": [%s],
				"preRefresh": [],
				"postRefresh": [],
				"preRollback": [],
				"postRollback": [],
				"preSnapshot": [],
				"postSnapshot": [],
				"preStart": [],
				"postStart": [],
				"preStop": [],
				"postStop": [],
				"type": "VirtualSourceOperations"
				},
				"parameters": {
					"postgresPort": %d,
					"configSettingsStg": [
						{ "propertyName": "listen_addresses", "value": "*" }
					]
				},
				"additionalMountPoints": [],
				"allowAutoVDBRestartOnHostReboot": true,
				"logCollectionEnabled": false,
				"name": "%s",
				"type": "AppDataVirtualSource"
			},
			"sourceConfig": {
				"path": "/mnt/provision/%s",
				"name": "%s",
				"repository": "%s",
				"linkingEnabled": true,
				"environmentUser": "%s",
				"type": "AppDataDirectSourceConfig"
			},
			"timeflowPointParameters": {
				"type": "TimeflowPointSemantic",
				"location": "LATEST_SNAPSHOT",
				"container": "%s"
			},
			"masked": %t,
			"type": "AppDataProvisionParameters"
			}`, groupRef, vdbParams.vdbName, vdbParams.configClone, vdbParams.port, vdbParams.vdbName, vdbParams.dbName, vdbParams.vdbName, scRepoRef, environmentUserRef, sourceDBRef, vdbParams.masked)
		logger.Debug(postBody)
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

func (c *myClient) createSourceConfig(scParams AppDataStagedSourceConfig, wait bool) (results map[string]interface{}, err error) {

	scObj, err := c.findSourceCongfigByNameAndEnvironmentName(scParams.name, scParams.envName)
	if err != nil {
		return nil, err
	}
	if scObj != nil {
		logger.Infof("SourceConfig %s already exists on %s", scParams.name, scParams.envName)
		logger.Info(scObj)
		return scObj, err
	}

	repoObj, err := c.findRepoByNameAndEnvironmentName(scParams.repoName, scParams.envName)
	if err != nil {
		return nil, err
	}
	if repoObj == nil {
		log.Fatalf("Repo %s on %s not found", scParams.repoName, scParams.envName)
	}

	postBody := fmt.Sprintf(`{
		"name":"%s",
		"repository":"%s",
		"parameters":{
			"dbPath":"%s"
			},
		"linkingEnabled":%t,
		"type":"AppDataStagedSourceConfig"
	}`, scParams.name, repoObj["reference"], scParams.parameters, scParams.linkingEnabled)
	log.Print(postBody)
	action, _, err := c.httpPost("sourceconfig", postBody)
	if err != nil {
		return nil, err
	}
	if wait {
		c.jobWaiter(action)
	}
	return action, err
}

// batchProvisionVDB takes one parameter:
// vdbParamsList: a slice of VDBParams to create
func (c *myClient) batchProvisionVDB(vdbParamsList ...PGVDBParams) (resultsList []map[string]interface{}, err error) {
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

func (c *myClient) createSelfServiceTemplate(wait bool) (resultsList map[string]interface{}, err error) {
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

func (c *myClient) createSelfServiceContainer(container PatientContainer, wait bool) (resultsList map[string]interface{}, err error) {
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

func (c *myClient) populateMasking(user, pass string, port int) (err error) {
	logger := logger.WithFields(log.Fields{
		"url":      c.url,
		"username": c.username,
	})

	appName := "Patients Application"
	envName := "Patients Environment"
	appExists := false
	envExists := false
	globFileName := "global_objects.json"
	mjFileName := "masking_job_postgres.json"
	glObjFile, err := ioutil.ReadFile(globFileName)
	if err != nil {
		logger.Fatal("Unable to open: ", globFileName, err)
	}
	logger.Info("Uploading global objects")
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
	prodDBIP, err := getIP("proddb")
	if err != nil {
		return err
	}
	logger.Infof("updating masking connector host to %s", prodDBIP)
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

// batchCreateSelfServiceContainer takes one parameter:
// containerList: a slice of PatientContainer to create
func (c *myClient) batchCreateSelfServiceContainer(containerList ...PatientContainer) (resultsList []map[string]interface{}, err error) {
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
func (c *myClient) updateUserPasswordByName(u string, p string) (err error) {
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
	c.LoadAndValidate()
	_, _, err = c.httpPost(fmt.Sprintf("user/%s/updateCredential", userRef), postBody)
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
func (c *myClient) getStorageDevices() (devices []string, err error) {

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
func (c *myClient) initializeSystem(d string, p string) (b bool, err error) {
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
	action, _, err := c.httpPost("domain/initializeSystem", postBody)
	if err != nil {
		return
	}
	c.jobWaiter(action)
	c.waitForEngineReady(10, 300)
	return true, err
}

func (c *myClient) checkForEnvironments(envList []string) (err error) {
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

func (c *myClient) checkForBusyEnvironments(envList []string) (err error) {
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
func (c *myClient) waitForEnvironmentsReady(p, t int, envList ...string) (err error) {

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
	sysCR := NewClientRequest("sysadmin", "sysadmin", fmt.Sprintf("https://%s/resources/json/delphix", opts.DDPVirtName))
	sysClient := sysCR.initResty()

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
	sysCR = NewClientRequest("sysadmin", "sysadmin", fmt.Sprintf("https://%s/resources/json/delphix", opts.DDPMaskName))
	sysClient = sysCR.initResty()

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

func (c *myClient) setInitialPlatformPassword() {
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

func (c *myClient) setInitialMaskingPassword() {
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
		}`, maskingIP, opts.Password))
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug(result)
	return result, err
}

func (c *myClient) createDemoRetentionPolicy() (string, error) {
	namespace := "policy"
	policyName := "Demo"
	if policyRef, err := c.findObjectByNameReturnReference(namespace, policyName); policyRef == nil && err == nil {
		postBody := fmt.Sprintf(`{
		"dataDuration": 2,
		"dataUnit": "YEAR",
		"logDuration": 1,
		"logUnit": "YEAR",
		"customized": false,
		"name": "%s",
		"type": "RetentionPolicy"
		}`, policyName)

		action, _, err := c.httpPost(namespace, postBody)
		if err != nil {
			return "", err
		}

		c.jobWaiter(action)
		return action["result"].(string), err
	} else {
		return policyRef.(string), err
	}
}

func (c *myClient) applyDemoRetentionPolicy(policyRef, containerRef string) error {
	namespace := "policy"
	postBody := fmt.Sprintf(`{
			"target":"%s",
			"type":"PolicyApplyTargetParameters"
		}`, containerRef)

	action, _, err := c.httpPost(fmt.Sprintf("%s/%s/apply", namespace, policyRef), postBody)
	if err != nil {
		return err
	}

	c.jobWaiter(action)
	return err

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

	prodSC := AppDataStagedSourceConfig{
		linkingEnabled: true,
		name:           "Patients Prod",
		parameters:     "{dbPath: \\\"Patients Prod\\\"}",
		dbPath:         "Patients Prod",
		repoName:       "Postgres vFiles (11.5)",
		envName:        "proddb",
	}

	devSC := AppDataStagedSourceConfig{
		linkingEnabled: true,
		name:           "Patients Non-Prod",
		parameters:     "{dbPath: \\\"Patients Non-Prod\\\"}",
		dbPath:         "Patients Non-Prod",
		repoName:       "Postgres vFiles (11.5)",
		envName:        "devdb",
	}

	proddb := LinkParams{
		dbName:              "Patients Prod",
		groupName:           "Prod",
		environmentUserName: "postgres",
		sourceConfig:        prodSC,
		dbUser:              "delphixdb",
		dbPass:              "delphixdb",
		preSync: fmt.Sprintf(`{
          "command": "psql -c 'create table test ( did integer NOT NULL);drop table test;select pg_switch_wal();'",
          "name": "Switch WAL Log",
          "type": "RunBashOnSourceOperation"
        }`),
	}

	maskingIP, err := getIP("maskingengine")
	if err != nil {
		logger.Fatal(err)
	}

	patmm := PGVDBParams{
		configClone: fmt.Sprintf(`{
							"command": "#!/usr/bin/env bash\n#\n# Copyright (c) 2019 by Delphix. All rights reserved.\n#\n#v1.0\n#2019 - Adam Bowen\n#requires curl and jq to be installed on the host machine\nDMHOST=\"%s\"\nDMPORT=\"80\"\nURL=\"http://${DMHOST}:${DMPORT}/masking/api\"\nENVIRONMENT=\"Patients Environment\"\nMASKINGJOB=\"Patients Mask - PG\"\nCONNECTOR=\"PatientsMM - PG\"\nDMUSER=Admin\nDMPASS=%s\n\necho \"Connecting to ${URL}\"\n\necho \"Authenticating\"\n\nAUTH=$(curl -sX POST --header 'Content-Type: application/json' --header 'Accept: application/json'\\\n  -d \"{ \\\"username\\\": \\\"${DMUSER}\\\", \\\"password\\\": \\\"${DMPASS}\\\"}\" \"${URL}/login\" | jq -r .Authorization)\n\n[[ -z $AUTH || $AUTH == \"null\" ]] && echo \"Was unable to get authenticate. Please try again\" && exit 1\n\nENVID=$(curl -sX GET --header 'Accept: application/json' --header \"Authorization: ${AUTH}\" \"${URL}/environments\" | \\\n  jq -r \".responseList[]| select(.environmentName==\\\"${ENVIRONMENT}\\\").environmentId\")\n\n[[ -z $ENVID || $ENVID == \"null\" ]] && echo \"Was unable to find Job ${ENVIRONMENT}. Please try again\" && echo ${EXECID} && exit 1\n\nJOBID=$(curl -sX GET --header 'Accept: application/json' --header \"Authorization: ${AUTH}\" \\\n  \"${URL}/masking-jobs?environment_id=${ENVID}\" | jq -r \".responseList[] | select(.jobName==\\\"${MASKINGJOB}\\\").maskingJobId\")\n\n[[ -z $JOBID || $JOBID == \"null\" ]] && echo \"Was unable to find Job ${MASKINGJOB}. Please try again\" && echo ${JOBID} && exit 1\n\nCONID=$(curl -sX GET --header 'Accept: application/json' --header \"Authorization: ${AUTH}\" \\\n  \"${URL}/database-connectors?environment_id=${ENVID}\" | jq -r \".responseList[] | select(.connectorName==\\\"${CONNECTOR}\\\").databaseConnectorId\")\n\n[[ -z $CONID || $CONID == \"null\" ]] && echo \"Was unable to find Job ${CONNECTOR}. Please try again\" && echo ${CONID} && exit 1\n\necho \"Executing Job \\\"${ENVIRONMENT}/${MASKINGJOB}\\\" with \\\"${CONNECTOR}\\\": ${JOBID}\"\n\nEXECID=$(curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json'\\\n  --header \"Authorization: ${AUTH}\" -d \"{\\\"jobId\\\": ${JOBID}, \\\"targetConnectorId\\\": ${CONID}}}\"\\\n  \"${URL}/executions\"|jq -r .executionId)\n\n[[ -z $EXECID || $EXECID == \"null\" ]] && echo \"Was unable to start Job ${JOBID}. Please try again\" && echo ${EXECID} && exit 1\n\necho \"Waiting for execution ${EXECID} to finish\"\n\nwhile [[ \"$STATUS\" != \"SUCCEEDED\" && \"$STATUS\" != \"FAILED\" ]]\n    do\n    sleep 3\n    STATUS=$(curl -sX GET --header 'Accept: application/json' --header \"Authorization: ${AUTH}\"\\\n    \"${URL}/executions/${EXECID}\" | jq -r .status)\n    [[ -z $STATUS || $STATUS == \"null\" ]] && echo \"Was unable to get status of execution ${EXECID}. Please try again\" && exit 1\ndone\n\n[[ \"$STATUS\" == \"FAILED\" ]] && echo \"\\\"${ENVIRONMENT}/${MASKINGJOB}\\\" failed execution ${EXECID}.\\nCheck logs for details.\" && exit 1\n\necho \"\\\"${ENVIRONMENT}/${MASKINGJOB}\\\" : $JOBID successfully ran\"",
							"name": "Masking Job",
							"type": "RunBashOnSourceOperation"
						}`, maskingIP, opts.Password),
		dbName:              "patmm",
		environmentUserName: "postgres",
		groupName:           "Masked Masters",
		masked:              true,
		port:                5435,
		sourceDBName:        "Patients Prod",
		sourceConfig:        prodSC,
		vdbName:             "Patients Masked Master",
	}

	devdb := PGVDBParams{
		dbName:              "devdb",
		environmentUserName: "postgres",
		groupName:           "Non Prod",
		port:                5454,
		sourceConfig:        devSC,
		sourceDBName:        "Patients Masked Master",
		vdbName:             "Patients Dev",
	}

	testdb := PGVDBParams{
		dbName:              "testdb",
		environmentUserName: "postgres",
		groupName:           "Non Prod",
		port:                5455,
		sourceConfig:        devSC,
		sourceDBName:        "Patients Masked Master",
		vdbName:             "Patients Test",
	}

	prep1 := PGVDBParams{
		dbName:              "prep1",
		environmentUserName: "postgres",
		groupName:           "Non Prod",
		port:                5461,
		sourceConfig:        devSC,
		sourceDBName:        "Patients Masked Master",
		vdbName:             "prep1",
	}

	prep2 := PGVDBParams{
		dbName:              "prep2",
		environmentUserName: "postgres",
		groupName:           "Non Prod",
		port:                5462,
		sourceConfig:        devSC,
		sourceDBName:        "Patients Masked Master",
		vdbName:             "prep2",
	}

	prep3 := PGVDBParams{
		dbName:              "prep3",
		environmentUserName: "postgres",
		groupName:           "Non Prod",
		port:                5463,
		sourceConfig:        devSC,
		sourceDBName:        "Patients Masked Master",
		vdbName:             "prep3",
	}

	prep4 := PGVDBParams{
		dbName:              "prep4",
		environmentUserName: "postgres",
		groupName:           "Non Prod",
		port:                5464,
		sourceConfig:        devSC,
		sourceDBName:        "Patients Masked Master",
		vdbName:             "prep4",
	}

	prep5 := PGVDBParams{
		dbName:              "prep5",
		environmentUserName: "postgres",
		groupName:           "Non Prod",
		port:                5465,
		sourceConfig:        devSC,
		sourceDBName:        "Patients Masked Master",
		vdbName:             "prep5",
	}

	prep6 := PGVDBParams{
		dbName:              "prep6",
		environmentUserName: "postgres",
		groupName:           "Non Prod",
		port:                5466,
		sourceConfig:        devSC,
		sourceDBName:        "Patients Masked Master",
		vdbName:             "prep6",
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

	err = initializePlatform()
	if err != nil {
		logger.Fatal(err)
	}
	virtualizationCR := NewClientRequest("delphix_admin", "delphix", fmt.Sprintf("https://%s/resources/json/delphix", opts.DDPVirtName))
	virtualizationClient := virtualizationCR.initResty()
	virtualizationClient.setInitialPlatformPassword()

	err = virtualizationClient.uploadPlugin("postgres.json")
	if err != nil {
		log.Fatal(err)
	}

	retPolicy, err := virtualizationClient.createDemoRetentionPolicy()
	if err != nil {
		log.Fatal(err)
	}

	maskingCR := NewClientRequest("admin", "Admin-12", fmt.Sprintf("https://%s/resources/json/delphix", opts.DDPMaskName))
	maskingClient := maskingCR.initResty()
	maskingClient.setInitialPlatformPassword()

	maskingCR = NewClientRequest("admin", "Admin-12", fmt.Sprintf("http://%s/masking/api", opts.DDPMaskName))
	maskingClient = maskingCR.initResty()
	maskingClient.setInitialMaskingPassword()

	result, err := virtualizationClient.updateMaskingService()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(result)

	srv := startHTTPServer()
	defer shutdownHTTPServer(srv)
	err = virtualizationClient.waitForEnvironmentsReady(10, 300, prodSC.envName, devSC.envName)

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

	_, err = virtualizationClient.batchCreateDatasetGroup("Prod", "Masked Masters", "Non Prod")
	if err != nil {
		logger.Fatal(err)
	}

	dSourceRef, err := virtualizationClient.linkPostgresDatabase(proddb, true)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info(dSourceRef)
	err = virtualizationClient.applyDemoRetentionPolicy(retPolicy, dSourceRef["result"].(string))
	err = maskingClient.populateMasking(proddb.dbUser, proddb.dbPass, patmm.port)
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

	// _, err = virtualizationClient.linkMaskingJob()
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	_, err = virtualizationClient.provisionVDB(patmm, true)
	if err != nil {
		logger.Fatal(err)
	}

	// //create the devSource config upfront, so that we can create other VDB's in parallel
	_, err = virtualizationClient.createSourceConfig(devSC, true)
	if err != nil {
		logger.Fatal(err)
	}

	_, err = virtualizationClient.batchProvisionVDB(devdb, testdb, prep1, prep2, prep3, prep4, prep5, prep6)
	if err != nil {
		logger.Fatal(err)
	}

	_, err = virtualizationClient.createSelfServiceTemplate(true)
	if err != nil {
		logger.Fatal(err)
	}

	_, err = virtualizationClient.batchCreateSelfServiceContainer(developContainer, testContainer, prep1Container, prep2Container, prep3Container, prep4Container, prep5Container, prep6Container)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Complete")
}
