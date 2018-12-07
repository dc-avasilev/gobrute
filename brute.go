package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
	if method == "post" {
		response, err = http.PostForm(target, url.Values{
			loginField:    {login},
			passwordField: {password},
		})
	}

	if method == "get" {
		target = strings.Replace(target, ":password", password, 1)
		response, err = http.Get(target)
	}
	if err != nil {
		return ""
	}
	if response.StatusCode != 200 {
		return fmt.Sprintf("Error status code: %d with password %s", response.StatusCode, password)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	responseText := buf.String()
	if strings.Contains(responseText, stop) {
		return ""
	}
	return password
}
