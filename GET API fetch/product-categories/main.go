package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type product_categories struct {
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
		ID         int `json:"id"`
		ProductID  int `json:"product_id"`
		CategoryID int `json:"category_id"`
	} `json:"data"`
}

func main() {
	fmt.Println("Product-Categories api fetch using get . . .")
	url := "https://gorest.co.in/public-api/product-categories"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var all product_categories
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
	err = ioutil.WriteFile("product-categories.json", file, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
