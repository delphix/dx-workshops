package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v1"
)

func init() {
	log.Infof("Version: %v", version)
	optionStuff()
	configLogging()
}

func configLogging() {
	if len(opts.Debug) > 0 {
		log.SetLevel(log.DebugLevel)
	}
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableLevelTruncation: true,
	})
	logger = log.WithFields(log.Fields{
		"url":      fmt.Sprintf("https://%s", opts.DDPVirtName),
		"username": "delphix_admin",
	})
}

func optionStuff() {
	opts.ConfigFile = func(s string) error {
		ini := flags.NewIniParser(parser)
		return ini.ParseFile(s)
	}

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else if flagsErr.Type == flags.ErrRequired || flagsErr.Type == flags.ErrMarshal {
			var b bytes.Buffer
			parser.WriteHelp(&b)
			log.Fatal(b.String())
		} else {
			log.Warn(flagsErr.Type)
			panic(err)
		}
	}
}

// CreateAPISession returns an APISession object
//v = APIVersion Struct
//l = Locale as an IETF BCP 47 language tag, defaults to 'en-US'.
//c = Client software identification token.
func CreateAPISession(v APIVersionStruct, l string, c string) (APISessionStruct, error) {
	if l == "" {
		l = "en-US"
	}
	if len(c) > 63 {
		err := fmt.Errorf("Client ID specified cannot be longer than 64 characters.\nYou provided %s", c)
		return APISessionStruct{}, err
	}
	apiSession := APISessionStruct{
		Type:    "APISession",
		Version: &v,
		Locale:  l,
		Client:  c,
	}
	return apiSession, nil
}

//CreateAPIVersion returns an APISession object
func CreateAPIVersion(major int, minor int, micro int) (APIVersionStruct, error) {
	maj := new(int)
	min := new(int)
	mic := new(int)
	t := "APIVersion"
	*maj = major
	*min = minor
	*mic = micro

	apiVersion := APIVersionStruct{
		Type:  t,
		Major: maj,
		Minor: min,
		Micro: mic,
	}
	return apiVersion, nil
}

// APISessionStruct - Describes a Delphix web service session and is the result of an
// initial handshake.
// extends TypedObject
type APISessionStruct struct {
	// Client software identification token.
	// required = false
	// maxLength = 64
	Client string `json:"client,omitempty"`
	// Locale as an IETF BCP 47 language tag, defaults to 'en-US'.
	// format = locale
	// required = false
	Locale string `json:"locale,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
	// Version of the API to use.
	// required = true
	Version *APIVersionStruct `json:"version,omitempty"`
}

// APIVersionStruct - Describes an API version.
// extends TypedObject
type APIVersionStruct struct {
	// Major API version number.
	// required = true
	// minimum = 0
	Major *int `json:"major,omitempty"`
	// Micro API version number.
	// minimum = 0
	// required = true
	Micro *int `json:"micro,omitempty"`
	// Minor API version number.
	// required = true
	// minimum = 0
	Minor *int `json:"minor,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// AuthSuccess holds the resty response success message
type AuthSuccess struct {
	ID, Message string
}

// Client the structure of a client request
type Client struct {
	url, username, password string
	version                 *APIVersionStruct
}

// ErrorStruct is the struct of a resty error
type ErrorStruct struct {
	Type          string `json:"type,omitempty"`
	Details       string `json:"details,omitempty"`
	ID            string `json:"id,omitempty"`
	CommandOutput string `json:"commandOutput,omitempty"`
	Diagnosis     string `json:"diagnosis,omitempty"`
}

// LoginRequestStruct - Represents a Delphix user authentication request.
// extends TypedObject
type LoginRequestStruct struct {
	// Whether to keep session alive for all requests or only via
	// 'KeepSessionAlive' request headers. Defaults to ALL_REQUESTS if
	// omitted.
	// enum = [ALL_REQUESTS KEEP_ALIVE_HEADER_ONLY]
	// default = ALL_REQUESTS
	KeepAliveMode string `json:"keepAliveMode,omitempty"`
	// The password of the user to authenticate.
	// format = password
	// required = true
	Password string `json:"password,omitempty"`
	// The authentication domain.
	// enum = [DOMAIN SYSTEM]
	Target string `json:"target,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
	// The username of the user to authenticate.
	// required = true
	Username string `json:"username,omitempty"`
}

// RespError holds the resty response failure message
type RespError struct {
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	ErrorStruct `json:"error,omitempty"`
}

// NewClient creates a new client object
func NewClient(username, password, url string) *Client {
	return &Client{
		url:      url,
		username: username,
		password: password,
	}
}

// LoadAndValidate establishes a new client connection
func (c *Client) LoadAndValidate() (int, error) {

	apiVersion, err := CreateAPIVersionFromString()
	if err != nil {
		return 0, err
	}
	apiStruct, err := CreateAPISession(apiVersion, "", "")
	if err != nil {
		return 0, err
	}

	resp, err := resty.R().
		SetBody(apiStruct).
		Post(c.url + "/session")

	result := resp.Body()
	var resultdat map[string]interface{}
	if err = json.Unmarshal(result, &resultdat); err != nil { //convert the json to go objects
		return resp.StatusCode(), err
	}

	if resultdat["status"].(string) == "ERROR" {
		errorMessage := string(result)
		err = fmt.Errorf(errorMessage)
		if err != nil {
			return resp.StatusCode(), err
		}
	}

	resp, err = resty.R().
		SetResult(AuthSuccess{}).
		SetBody(LoginRequestStruct{
			Type:     "LoginRequest",
			Username: c.username,
			Password: c.password,
		}).
		Post(c.url + "/login")
	if err != nil {
		return resp.StatusCode(), err
	}

	if http.StatusUnauthorized == resp.StatusCode() {
		err = fmt.Errorf("%d", http.StatusUnauthorized)
		if err != nil {
			return resp.StatusCode(), err
		}
	}
	parseHTTPResponseReturnMap(resp)
	return resp.StatusCode(), nil
}

// MaskingLoadAndValidate establishes a new client connection
func (c *Client) MaskingLoadAndValidate() (int, error) {

	resp, err := resty.R().
		SetResult(AuthSuccess{}).
		SetBody(fmt.Sprintf(`{
			"username": "%s",
			"password": "%s"
			}`, c.username, c.password)).
		Post(c.url + "/login")
	if err != nil {
		return resp.StatusCode(), err
	}
	logger.Debug(resp)
	resultdat, _, err := parseHTTPResponseReturnMap(resp)
	if err != nil {
		return resp.StatusCode(), err
	}

	if auth, _ := resultdat["Authorization"].(string); auth != "" {
		resty.SetHeader("Authorization", auth)
	}
	return resp.StatusCode(), nil
}

func bodyToJSON(v interface{}) string {
	switch v := v.(type) {
	case string:
		return string(v)
	case LoginRequestStruct:
		return fmt.Sprintf("{\"type\": \"LoginRequest\",\"username\": \"%s\",\"password\": \"<redacted>\"}", "delphix_admin")
	default:
		if fmt.Sprintf("HERE %v", reflect.ValueOf(v).Kind()) == "struct" {
			tbEnc, _ := json.Marshal(v)
			return (string(tbEnc))
		}
	}
	return ""
}

func (c *Client) initResty() {
	resty.DefaultClient.
		SetTimeout(time.Duration(30 * time.Second)).
		SetRetryCount(5).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second).
		AddRetryCondition(func(r *resty.Response) (bool, error) {
			return r.StatusCode() == http.StatusBadGateway, nil
		})

	resty.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	skipCertValidation(opts.SkipValidate)

	if len(opts.Debug) >= 3 {
		resty.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			var requestdat map[string]interface{}
			strbody := bodyToJSON(req.Body)
			err := json.Unmarshal([]byte(strbody), &requestdat)
			if err == nil {
				requestType, _ := requestdat["type"].(string)
				if requestType == "ERROR" {
					errorMessage := string(req.Body.(string))
					err := fmt.Errorf(errorMessage)
					if err != nil {
						return err
					}
				} else if requestType == "LoginRequest" {
					strbody = fmt.Sprintf("{\"type\": \"LoginRequest\",\"username\": \"%s\",\"password\": \"<redacted>\"}", requestdat["username"])
				}
			}
			log.WithFields(log.Fields{
				"method":  req.Method,
				"url":     req.URL,
				"headers": req.Header,
				"body":    strbody,
			}).Debug("Request:")
			return nil
		})

		resty.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			log.WithFields(log.Fields{
				"statusCode":   resp.StatusCode(),
				"status":       resp.Status(),
				"responseTime": resp.Time(),
				"receivedAt":   resp.ReceivedAt(),
				"body":         resp.String(),
			}).Debug("Response:")
			return nil
		})
	}
}

func skipCertValidation(b bool) {
	if b {
		//Turn off Cert Validation for our testing
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
}

func (c *Client) assembleURL(url string, params []string) string {
	return fmt.Sprintf("%s/%s?%s", c.url, url, strings.Join(params, "&"))
}

// CreateAPIVersionFromString - reads and parses the global var apiVersionString to create the APIVersion object
func CreateAPIVersionFromString() (APIVersionStruct, error) {
	version := strings.Split(apiVersionString, ".")
	var empty APIVersionStruct

	maj, err := strconv.Atoi(version[0])
	if err != nil {
		return empty, err
	}

	min, err := strconv.Atoi(version[1])
	if err != nil {
		return empty, err
	}

	mic, err := strconv.Atoi(version[2])
	if err != nil {
		return empty, err
	}

	return CreateAPIVersion(maj, min, mic)
}

func (c *Client) httpPost(url, body string, params ...string) (map[string]interface{}, int, error) {

	postURL := c.assembleURL(url, params)
	// postURL := fmt.Sprintf("%s/%s", c.url, url)
	logger.WithField("url", postURL)
	resp, err := resty.R().
		SetBody(body).
		Post(postURL)
	if err != nil {
		log.Fatal(err)
	}

	return parseHTTPResponseReturnMap(resp)
}

func (c *Client) httpPut(url, body string, params ...string) (map[string]interface{}, int, error) {

	putURL := c.assembleURL(url, params)
	logger.WithField("url", putURL)
	resp, err := resty.R().
		SetBody(body).
		Put(putURL)
	if err != nil {
		log.Fatal(err)
	}

	return parseHTTPResponseReturnMap(resp)
}

func (c *Client) httpPostBytesReturnSlice(url string, file []byte, params ...string) ([]interface{}, int, error) {

	postURL := c.assembleURL(url, params)
	// postURL := fmt.Sprintf("%s/%s", c.url, url)
	logger.WithField("url", postURL)
	resp, err := resty.R().
		SetBody(file).
		Post(postURL)
	if err != nil {
		log.Fatal(err)
	}

	return parseHTTPResponseReturnSlice(resp)
}

func (c *Client) httpGet(url string, params ...string) (map[string]interface{}, int, error) {
	getURL := c.assembleURL(url, params)
	logger.WithField("url", getURL)
	resp, err := resty.R().
		Get(getURL)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	return parseHTTPResponseReturnMap(resp)
}

func parseHTTPResponse(resp *resty.Response) ([]byte, int, error) {
	var err error

	if http.StatusOK != resp.StatusCode() {
		err = fmt.Errorf("Got an HTTP Status of %v instead of %v\n%v", resp.StatusCode(), http.StatusOK, resp)
		if err != nil {
			return nil, resp.StatusCode(), err
		}
	}
	return resp.Body(), resp.StatusCode(), nil
}

func parseHTTPResponseReturnMap(resp *resty.Response) (resultdat map[string]interface{}, statusCode int, err error) {
	result, statusCode, err := parseHTTPResponse(resp)
	if err != nil {
		return nil, statusCode, err
	}
	if err = json.Unmarshal(result, &resultdat); err != nil { //convert the json to go objects
		return nil, statusCode, err
	}

	if status, _ := resultdat["status"].(string); status == "ERROR" {
		errorMessage := string(result)
		err = fmt.Errorf(errorMessage)
		return nil, statusCode, err
	}
	return resultdat, statusCode, err
}

func parseHTTPResponseReturnSlice(resp *resty.Response) (resultdat []interface{}, statusCode int, err error) {
	result, statusCode, err := parseHTTPResponse(resp)
	if err != nil {
		return nil, statusCode, err
	}
	if err := json.Unmarshal(result, &resultdat); err != nil { //convert the json to go objects
		return nil, statusCode, err
	}

	return resultdat, statusCode, nil
}

func (c *Client) jobWaiter(actionList ...map[string]interface{}) error {
	if actionList == nil {
		logger.Debug("No jobs passed to jobWaiter")
	}
	logger.Debugf("Action list: %v", actionList)
	for _, v := range actionList {
		jobLogger := logger.WithFields(log.Fields{
			"url":    fmt.Sprintf("%s/job/%s", c.url, v["job"]),
			"action": v["action"],
			"job":    v["job"],
		})
		jobLogger.Debugf("acting on %v", v)
		if v["job"] != nil {
			for {

				jobObj, _, err := c.httpGet(fmt.Sprintf("job/%s", v["job"]))
				if err != nil {
					return err
				}
				jobResult := jobObj["result"].(map[string]interface{})
				jobLogger.Infof("%g%s complete", jobResult["percentComplete"], "%")
				if jobState := jobResult["jobState"].(string); jobState == "RUNNING" {
					jobLogger.Debug(c.listObjects("notification", fmt.Sprintf("channel=%s", jobResult["reference"])))
				} else {
					if jobState != "COMPLETED" {
						jobLogger.Fatal(jobState)
					}
					jobLogger.Info(jobState)
					break
				}
			}
		} else {
			for {
				jobLogger := logger.WithFields(log.Fields{
					"url":    fmt.Sprintf("%s/action/%s", c.url, v["action"]),
					"action": v["action"],
				})
				actionObj, _, err := c.httpGet(fmt.Sprintf("action/%s", v["action"]))
				if err != nil {
					return err
				}
				actionResult := actionObj["result"].(map[string]interface{})
				jobLogger.Infof("Waiting for action to complete")
				if actionState := actionResult["state"].(string); actionState == "EXECUTING" || actionState == "WAITING" {
					jobLogger.Debug(c.listObjects("notification", fmt.Sprintf("channel=%s", v["action"])))
				} else {
					if actionState != "COMPLETED" {
						jobLogger.Fatal(actionState)
					}
					jobLogger.Info(actionState)
					break
				}
			}
		}
	}
	return nil
}

// WaitForEngineReady loops until the Client connection is successful or time (t) expires
func (c *Client) waitForEngineReady(p int, t int) error {

	logger.Infof("Waiting up to %v seconds for the DDDP to be ready", t)
	timeOut := 0
	for timeOut < t {
		logger.Info("Waiting for Delphix DDP")
		time.Sleep(time.Duration(p) * time.Second)
		timeOut = timeOut + p
		if _, err := c.LoadAndValidate(); err == nil {
			break
		}
	}
	return nil
}

// WaitForMaskingEngineReady loops until the Client connection is successful or time (t) expires
func (c *Client) waitForMaskingEngineReady(p int, t int) error {
	logger := logger.WithFields(log.Fields{
		"url":      c.url,
		"username": c.username,
	})
	logger.Infof("Waiting up to %v seconds for the DDDP to be ready", t)
	timeOut := 0
	for timeOut < t {
		logger.Info("Waiting for Delphix DDP Masking Engine")
		time.Sleep(time.Duration(p) * time.Second)
		timeOut = timeOut + p
		if _, err := c.MaskingLoadAndValidate(); err == nil {
			break
		}
	}
	return nil
}

func returnObjReference(obj map[string]interface{}, errIn error) (reference interface{}, errOut error) {
	return obj["reference"], errIn
}
