package main

import (
	"encoding/xml"
	"io"
	"net/http"
)

type Client struct {
	// HTTPClient used for send/receive soap request
	HTTPClient     *http.Client
	RequestBuilder func(method string, url string, body io.Reader) (*http.Request, error)
}

type addRequest struct {
	XMLName xml.Name `xml:"http://tempuri.org/ Add"` //requires a blank space before Add

	IntA int `xml:"intA"`
	IntB int `xml:"intB"`
}

type addResponse struct {
	XMLName xml.Name `xml:"http://tempuri.org/ AddResponse"` //requires a blank space before AddResponse

	AddResult int `xml:"AddResult"`
}
//struct to read from the csv
type Numbers struct {
	FirstNumber int
	SecondNumber  int
}
