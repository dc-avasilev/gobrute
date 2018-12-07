package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	mode := flag.String("mode", "none", "")
	port := flag.Int("p", -1, "")
	target := flag.String("t", "none", "")
	passwordPath := flag.String("P", "none", "")
	serverAddress := flag.String("address", "none", "")
	requestType := flag.String("type", "none", "")
	loginField := flag.String("LF", "none", "")
	passwordField := flag.String("PF", "none", "")
	login := flag.String("login", "none", "")
	stop := flag.String("stop", "none", "")
	flag.Parse()

	if *mode == "none" {
		fmt.Println("Mode parameter not found")
		return
	}
	if *mode == "server" {
		if *passwordPath == "none" {
			fmt.Println("Password path (-P) not found")
			return
		}
		if strings.ToLower(*requestType) != "get" && strings.ToLower(*requestType) != "post" {
			fmt.Println("Request type (-type) not found")
			return
		}
		server_mode(*port, *target, *login, *stop, *passwordPath, *requestType, *loginField, *passwordField)
		return
	}
	if *mode == "client" {
		client_mode(*serverAddress)
		return
	}
	fmt.Println("Mode not correct")
	return
}
