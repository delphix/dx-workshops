package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

// Options for the script
type Options struct {
	DDPName      string               `short:"e" long:"ddp_hostname" env:"DELPHIX_DDP_HOSTNAME" description:"The hostname or IP address of the Delphix Dynamic Data Platform" required:"true"`
	UserName     string               `short:"u" long:"username" env:"DELPHIX_USER" description:"The username used to authenticate to the Delphix Engine" required:"true"`
	Password     string               `short:"p" long:"password" env:"DELPHIX_PASS" description:"The password used to authenticate to the Delphix Engine" required:"true"`
	Debug        []bool               `short:"v" long:"debug" env:"DELPHIX_DEBUG" description:"Turn on debugging. -vvv for the most verbose debugging"`
	SkipValidate bool                 `long:"skip-validate" env:"DELPHIX_SKIP_VALIDATE" description:"Don't validate TLS certificate of Delphix Engine"`
	ConfigFile   func(s string) error `short:"c" long:"config" description:"Optional INI config file to pass in for the variables" no-ini:"true"`
	DSourceName  string               `long:"dsource" env:"DELPHIX_DSOURCE" description:"The name of the dSource to snapsync" required:"true"`
	VDBName      string               `long:"vdb" env:"DELPHIX_VDB" description:"The name of the VDB to snapsync" required:"true"`
}

func (c *Client) syncDatabaseByName(dSourceName string) (results map[string]interface{}) {
	databaseNamespace := "database"
	//Find our dSource of interest
	log.Info("Searching for dSource by name")
	databaseObj := c.findObjectByName(databaseNamespace, dSourceName)
	if databaseObj == nil {
		log.Fatalf("Could not find dSource named %s", dSourceName)
	}

	log.Infof("Found %s: %s", databaseObj["name"], databaseObj["reference"])

	//Sync the dSource
	url := fmt.Sprintf("%s/%s/sync", databaseNamespace, databaseObj["reference"])
	action := c.httpPost(url, "")

	c.jobWaiter(action)
	return action
}

func (c *Client) refreshDatabaseByName(vdbName string) (results map[string]interface{}) {
	var refreshParameters string
	databaseNamespace := "database"
	//Find our VDB of interest
	log.Info("Searching for VDB by name")
	databaseObj := c.findObjectByName(databaseNamespace, vdbName)
	if databaseObj == nil {
		log.Fatalf("Could not find VDB named %s", vdbName)
	}

	log.Infof("Found %s: %s", databaseObj["name"], databaseObj["reference"])

	//refresh the VDB
	url := fmt.Sprintf("%s/%s/refresh", databaseNamespace, databaseObj["reference"])

	if objType := databaseObj["type"]; objType != "OracleDatabaseContainer" {
		refreshParameters = `{
		"type": "RefreshParameters", 
		"timeflowPointParameters": {
			"type": "TimeflowPointSemantic",
			"location": "LATEST_SNAPSHOT"
		}
		}`
	} else {
		refreshParameters = `{
		"type": "OracleRefreshParameters", 
		"timeflowPointParameters": {
			"type": "TimeflowPointSemantic",
			"location": "LATEST_SNAPSHOT"
		}
		}`
	}

	action := c.httpPost(url, refreshParameters)

	c.jobWaiter(action)
	return action
}

var (
	opts             Options
	parser           = flags.NewParser(&opts, flags.Default)
	apiVersionString = "1.9.3"
	logger           *log.Entry
	url              string
	Version          = "not set"
)

func main() {
	var err error

	log.Info("Establishing session and logging in")
	client := NewClient(opts.UserName, opts.Password, fmt.Sprintf("https://%s/resources/json/delphix", opts.DDPName))
	client.initResty()
	err = client.LoadAndValidate()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Successfully Logged in")

	client.syncDatabaseByName(opts.DSourceName)
	client.refreshDatabaseByName(opts.VDBName)
	log.Info("Complete")
}
