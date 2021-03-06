package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	DefaultRequestBuilder = func(method string, url string, body io.Reader) (*http.Request, error) {
		req, err := http.NewRequest(http.MethodPost, url, body)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
		return req, err
	}

	DefaultClient = &Client{
		HTTPClient:     http.DefaultClient,
		RequestBuilder: DefaultRequestBuilder,
	}
)

func (c *Client) Call(url string, action string, requestHeaders, request, responseHeaders, response interface{}) error {
	envelope := SOAPEnvelope{
		Header: SOAPHeader{Header: requestHeaders},
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)

	encoder.Encode(envelope)	//Apparently, .Encode() already uses Flush before it returns
	//encoder.Flush()			//so encoder.Flush() is redundant.

	req, _ := c.RequestBuilder(action, url, buffer)

	if action != "" {
		req.Header.Add("SOAPAction", action)
	}

	req.Close = true
	res, _ := c.HTTPClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if len(body) == 0 { //we need this otherwise it does not exit!!!
		return nil
	}

	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Header = SOAPHeader{Header: responseHeaders}
	respEnvelope.Body = SOAPBody{Content: response}
	xml.Unmarshal(body, respEnvelope)

	return nil
}

func Call(url string, action string, requestHeaders, request, responseHeaders, response interface{}) error {
	return DefaultClient.Call(url, action, requestHeaders, request, responseHeaders, response)
}

func main(){
	//read the csv file

	csvFile, _ := os.Open("data.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var numbers []Numbers

	//read the csv file and create Number(FirstNumber, SecondNumber) from each line until we reach the EOF.
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		}
		first, _ := strconv.Atoi(line[0])
		second, _ := strconv.Atoi(line[1])
		numbers = append(numbers, Numbers{
			FirstNumber: first,
			SecondNumber:  second,
		})
	}

	//console output
	fmt.Println("Reading values from CSV:")
	fmt.Println(numbers)

	//Create the output txt file
	f, _ := os.Create("data.txt")

	//make the call for each Number tuple
	for i := 0; i < len(numbers); i++ {
		response := addResponse{}
		err := Call(
			"http://www.dneonline.com/calculator.asmx",
			"http://tempuri.org/Add",
			nil,
			addRequest{IntA: numbers[i].FirstNumber, IntB: numbers[i].SecondNumber},
			nil,
			&response) // soap response pointer

		if err != nil {
			panic(err)
		}

		log.Println("add result is", response.AddResult)
		result := strconv.Itoa(response.AddResult)
		firstNum := strconv.Itoa(numbers[i].FirstNumber)
		secondNum := strconv.Itoa(numbers[i].SecondNumber)
		//add the result into file.
		f.WriteString(firstNum + "+" + secondNum + " = " + result + "\n")
	}
}