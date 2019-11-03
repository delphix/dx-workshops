package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

type TestSuite struct {
	XMLName    xml.Name `xml:"testsuite"`
	Errors     string   `xml:"errors,attr"`
	Failures   string   `xml:"failures,attr"`
	Name       string   `xml:"name,attr"`
	Tests      int      `xml:"tests,attr"`
	Properties *Properties
	TestCases  []TestCase
}

type Properties struct {
	XMLName       xml.Name `xml:"properties"`
	PropertyArray []*Property
}

type Property struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

type TestCase struct {
	XMLName     xml.Name `xml:"testcase"`
	Classname   string   `xml:"classname,attr"`
	Name        string   `xml:"name,attr"`
	TestError   *TestError
	TestFailure *TestFailure
	TestSuccess *TestSuccess
}

type TestError struct {
	XMLName   xml.Name `xml:"error"`
	Type      string   `xml:"type,attr,omitempty"`
	Timestamp string   `xml:"timestamp,attr,omitempty"`
	Status    int      `xml:"status,attr,omitempty"`
	Message   string   `xml:",chardata"`
}

type TestFailure struct {
	XMLName   xml.Name `xml:"failure"`
	Type      string   `xml:"type,attr,omitempty"`
	Timestamp string   `xml:"timestamp,attr,omitempty"`
	Status    int      `xml:"status,attr,omitempty"`
	Message   string   `xml:",chardata"`
}

type TestSuccess struct {
	XMLName   xml.Name `xml:"success"`
	Type      string   `xml:"type,attr,omitempty"`
	Timestamp string   `xml:"timestamp,attr,omitempty"`
	Status    int      `xml:"status,attr,omitempty"`
	Message   string   `xml:",chardata"`
}

func createReportFile(path, filename string) *os.File {
	err := os.MkdirAll(path, 0744)
	if err != nil {
		log.Fatalf("error creating path: %v", err)
	}
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", path, filename), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return file
}

func writeReport(report *os.File, ts *TestSuite) {
	xmlWriter := io.Writer(report)
	xmlWriter.Write([]byte(xml.Header))
	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	if err := enc.Encode(ts); err != nil {
		log.Fatalf("Unable to write report: %v\n", err)
	}
}
