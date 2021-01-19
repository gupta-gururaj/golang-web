package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type comment struct {
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
		PostID    int       `json:"post_id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
}

func main() {
	fmt.Println("Comments api fetch using get . . .")
	url := "https://gorest.co.in/public-api/comments"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var all comment
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
	err = ioutil.WriteFile("comments.json", file, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
