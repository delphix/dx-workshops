package main

import (
	"fmt"
	"log"
)

var (
	version = "not set"
)

func main() {
	testFail := false
	v := &TestSuite{Name: "daf_tests"}
	v.Properties = &Properties{}
	v.Properties.PropertyArray = append(v.Properties.PropertyArray, &Property{Name: "tests-group", Value: "daf_tests"})
	report := createReportFile("build/reports", "test.xml")
	defer report.Close()

	server := "testweb"
	port := "8080"
	username := "patients_admin"
	password := "delphix"
	client := NewClient(username, password, fmt.Sprintf("http://%s:%s", server, port))
	err := client.login()
	if err != nil {
		log.Fatal(err)
	}
	testCase := TestCase{Classname: "Read All Patients"}
	patientList, err := client.getPatients()
	if err != nil {
		appError := err.(*AppError)
		testCase.TestFailure = &TestFailure{
			Type:      appError.errType,
			Timestamp: appError.timestamp,
			Status:    appError.status,
			Message:   appError.message,
		}
		testFail = true
	} else {
		testCase.TestSuccess = &TestSuccess{Message: fmt.Sprintf("%v", patientList)}
	}
	v.TestCases = append(v.TestCases, testCase)

	testCase = TestCase{Classname: "Read All Records"}
	records, err := client.getRecordsFromAnyPatientInList(patientList)
	if err != nil {
		appError := err.(*AppError)
		testCase.TestFailure = &TestFailure{
			Type:      appError.errType,
			Timestamp: appError.timestamp,
			Status:    appError.status,
			Message:   appError.message,
		}
		testFail = true
	} else {
		testCase.TestSuccess = &TestSuccess{Message: fmt.Sprintf("%v", records)}
	}
	v.TestCases = append(v.TestCases, testCase)

	testCase = TestCase{Classname: "Read All Billings"}
	billings, err := client.getBillingsFromAnyPatientInList(patientList)
	if err != nil {
		appError := err.(*AppError)
		testCase.TestFailure = &TestFailure{
			Type:      appError.errType,
			Timestamp: appError.timestamp,
			Status:    appError.status,
			Message:   appError.message,
		}
		testFail = true
	} else {
		testCase.TestSuccess = &TestSuccess{Message: fmt.Sprintf("%v", billings)}
	}
	v.TestCases = append(v.TestCases, testCase)

	testCase = TestCase{Classname: "Read All Payments"}
	payments, err := client.getPaymentsFromAnyPatientInList(patientList)
	if err != nil {
		appError := err.(*AppError)
		testCase.TestFailure = &TestFailure{
			Type:      appError.errType,
			Timestamp: appError.timestamp,
			Status:    appError.status,
			Message:   appError.message,
		}
		testFail = true
	} else {
		testCase.TestSuccess = &TestSuccess{Message: fmt.Sprintf("%v", payments)}
	}
	v.TestCases = append(v.TestCases, testCase)

	testCase = TestCase{Classname: "Add/Delete User"}
	bob := &User{"bobbyb", "bob", "bee", "beebee"}
	userID, err := client.signUpUserAndReturnID(bob)

	log.Printf("User %s assigned ID %d", bob.username, userID)
	testClient := NewClient(bob.username, bob.password, client.url)
	err = testClient.loginTest()
	if err != nil {
		appError := err.(*AppError)
		testCase.TestFailure = &TestFailure{
			Type:      appError.errType,
			Timestamp: appError.timestamp,
			Status:    appError.status,
			Message:   appError.message,
		}
		testFail = true
	}
	client.login()
	err = client.deleteUser(userID)
	if err != nil {
		appError := err.(*AppError)
		testCase.TestFailure = &TestFailure{
			Type:      appError.errType,
			Timestamp: appError.timestamp,
			Status:    appError.status,
			Message:   appError.message,
		}
		testFail = true
	}
	v.TestCases = append(v.TestCases, testCase)

	testCase = TestCase{Classname: "Add Duplicate User"}
	michael := &User{"mcred", "Michael", "Credible", "Ibreakthings"}
	userID, err = client.signUpUserAndReturnID(michael)
	if err != nil {
		testCase.TestSuccess = &TestSuccess{Message: fmt.Sprintf("%v", err)}
	} else {
		testCase.TestFailure = &TestFailure{
			Type:    "Unique Constraint",
			Status:  200,
			Message: fmt.Sprintf("User %s assigned ID %d", michael.username, userID),
		}
		testFail = true
	}

	v.TestCases = append(v.TestCases, testCase)
	writeReport(report, v)
	if testFail == true {
		// os.Exit(2)
		log.Print("tests failed")
	}
}
