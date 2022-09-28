package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simplehttpclient/internet"
)

func main() {
	postingInfo()
	gettingInfo()
}

func postingInfo() {
	url := "https://reqres.in/api/users"

	req := map[string]any{
		"name":    "Mirza",
		"address": "Jakarta",
	}

	var res any
	err := internet.NewSimpleHTTPClient(http.MethodPost, url, req).Call(&res)
	if err != nil {
		panic(err.Error())
	}
}

func gettingInfo() {
	url := "https://reqres.in/api/users"

	var res any
	err := internet.NewSimpleHTTPClient(http.MethodGet, url).Call(&res)
	if err != nil {
		panic(err.Error())
	}

	arrByte, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%v\n", string(arrByte))
}
