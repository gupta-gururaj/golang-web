package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type ContactDetails struct {
	ID        string
	Age       string
	Email     string
	FirstName string
	LastName  string
}

type pass struct {
	Data []ContactDetails
}

func insert(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@(127.0.0.1:3306)/form")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	tmpl := template.Must(template.ParseFiles("forms.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	details := ContactDetails{
		ID:        r.FormValue(("id")),
		Age:       r.FormValue("age"),
		Email:     r.FormValue("email"),
		FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
	}

	insrt := "INSERT INTO users (ID, age, first_name, last_name, email) VALUES(?,?,?,?,?)"

	stmt, err := db.Prepare(insrt)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(details.ID, details.Age, details.FirstName, details.LastName, details.Email)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

func display(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@(127.0.0.1:3306)/form")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	tmpl := template.Must(template.ParseFiles("display.html"))

	var details []ContactDetails
	rows, err := db.Query("SELECT *FROM users ORDER BY id")
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var detail ContactDetails
		err = rows.Scan(&detail.ID, &detail.Age, &detail.FirstName, &detail.LastName, &detail.Email)
		if err != nil {
			fmt.Println(err)
		}
		details = append(details, detail)
	}
	var datapass pass
	datapass.Data = details
	tmpl.Execute(w, datapass)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println("Delete ID ", id)
	db, err := sql.Open("mysql", "root:password@(127.0.0.1:3306)/form")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	del := "DELETE FROM users WHERE ID = ?"
	stmt, err := db.Prepare(del)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	http.Redirect(w, r, "/display", 301)
}

func update(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("edit.html"))
	id := r.URL.Query().Get("id")
	fmt.Println("Update ID ", id)
	db, err := sql.Open("mysql", "root:password@(127.0.0.1:3306)/form")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	slct := "SELECT * FROM users WHERE id = " + id

	var detail ContactDetails
	rows := db.QueryRow(slct)
	err = rows.Scan(&detail.ID, &detail.Age, &detail.FirstName, &detail.LastName, &detail.Email)
	if err != nil {
		fmt.Println(err)
	}

	if r.Method != http.MethodPost {
		tmpl.Execute(w, detail)
		return
	}

	updt := " UPDATE users SET age = ?, first_name = ?, last_name = ?, Email = ? WHERE ID = ?"

	details := ContactDetails{
		ID:        r.FormValue("id"),
		Age:       r.FormValue("age"),
		Email:     r.FormValue("email"),
		FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
	}
	stmt, err := db.Prepare(updt)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(details.Age, details.FirstName, details.LastName, details.Email, details.ID)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	http.Redirect(w, r, "/display", 301)
}

func main() {
	fmt.Println("Server - http://localhost:9050/")
	http.HandleFunc("/", insert)
	http.HandleFunc("/display", display)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/update", update)
	http.ListenAndServe(":9050", nil)
}
