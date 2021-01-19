package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type products struct {
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
		ID             int    `json:"id"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		Image          string `json:"image"`
		Price          string `json:"price"`
		DiscountAmount string `json:"discount_amount"`
		Status         bool   `json:"status"`
		Categories     []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"categories"`
	} `json:"data"`
}

func main() {
	fmt.Println("Products api fetch using get . . .")
	url := "https://gorest.co.in/public-api/products"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var all products
	err = json.Unmarshal([]byte(body), &all)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(all)

	file, err := json.MarshalIndent(all, "", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile("products.json", file, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
