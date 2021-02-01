package main

//erase data

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "form"
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

func form(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	tmpl := template.Must(template.ParseFiles("forms.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := ContactDetails{
		Age:       r.FormValue("age"),
		Email:     r.FormValue("email"),
		FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
	}

	sqlStatement := `
	INSERT INTO users (age, email, first_name, last_name)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, details.Age, details.Email, details.FirstName, details.LastName).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)
	http.Redirect(w, r, "/", 301)
}

func display(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	sqlStatement := "DELETE FROM users where id = $1"
	_, err = db.Exec(sqlStatement, id)
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	sqlStatement := `SELECT * FROM users WHERE id = $1`

	var detail ContactDetails
	err = db.QueryRow(sqlStatement, id).Scan(&detail.ID, &detail.Age, &detail.FirstName, &detail.LastName, &detail.Email)
	if err != nil {
		fmt.Println(err)
	}

	if r.Method != http.MethodPost {
		tmpl.Execute(w, detail)
		return
	}

	sqlStatement = `
	UPDATE users
	SET age = $2, first_name = $3, last_name = $4, email = $5
	WHERE id = $1;`

	details := ContactDetails{
		Age:       r.FormValue("age"),
		Email:     r.FormValue("email"),
		FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
	}
	_, err = db.Exec(sqlStatement, id, details.Age, details.FirstName, details.LastName, details.Email)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	http.Redirect(w, r, "/display", 301)
}

func main() {
	fmt.Println("Server - http://localhost:9000/")
	http.HandleFunc("/", form)
	http.HandleFunc("/display", display)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/update", update)
	http.ListenAndServe(":9000", nil)
}
