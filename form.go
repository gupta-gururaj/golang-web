package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

func main() {
	tmpl := template.Must(template.ParseFiles("forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		details := ContactDetails{
			Email:   r.FormValue("email"),
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message"),
		}

		filedata, err := ioutil.ReadFile("newdata.json")
		if err != nil {
			fmt.Println(err)
		}

		var all []ContactDetails
		err = json.Unmarshal([]byte(filedata), &all)
		if err != nil {
			fmt.Println("Error Unmarshling for user file")
			fmt.Println(err)
		}

		all = append(all, details)
		file, _ := json.MarshalIndent(all, "", " ")
		_ = ioutil.WriteFile("newdata.json", file, 0644)

		tmpl.Execute(w, struct{ Success bool }{true})
	})

	http.ListenAndServe(":8080", nil)
}
