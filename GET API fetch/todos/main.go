package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type todo struct {
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
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
}

func main() {
	fmt.Println("Todo api fetch using get . . .")
	url := "https://gorest.co.in/public-api/todos"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var all todo
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
	err = ioutil.WriteFile("todos.json", file, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
