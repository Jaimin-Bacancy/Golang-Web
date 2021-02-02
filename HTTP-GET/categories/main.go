package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type categories struct {
	Code int `json:"code"`
	Meta struct {
		Pagination struct {
			Total int `json:"total"`
			Pages int `json:"pages"`
			Page  int `json:"page"`
			Limit int `json:"limit"`
		} `json:"pagination"`
	} `json:"meta"`
	Data []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      bool   `json:"status"`
	} `json:"data"`
}

func main() {
	fmt.Println("Categories api")
	url := "https://gorest.co.in/public-api/categories"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var alldata categories
	err = json.Unmarshal([]byte(body), &alldata)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(alldata)

	file, err := json.MarshalIndent(alldata, "", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile("categories.json", file, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
