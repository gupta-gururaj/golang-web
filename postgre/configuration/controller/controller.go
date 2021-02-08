package controller

import (
	//Databse Package
	_ "database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"postgre/common"
	"postgre/configuration/structures"

	//GitHub repository if needed in case of error
	_ "github.com/lib/pq"
)

var tmpl = common.GetTemplate()

//Form is ...
func Form(w http.ResponseWriter, r *http.Request) {
	db := common.Conn()
	defer db.Close()
	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "forms.html", nil)
		return
	}

	file, _, err1 := r.FormFile("image")
	common.HandleError(err1)

	//file stores in byte array
	x, err2 := ioutil.ReadAll(file)
	common.HandleError(err2)

	details := structures.ContactDetails{
		Age:       r.FormValue("age"),
		Email:     r.FormValue("email"),
		FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
		Img:       x,
	}

	sqlStatement := `
	INSERT INTO users (age, email, first_name, last_name, image)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`
	id := 0
	var err error
	common.HandleError(err)
	err = db.QueryRow(sqlStatement, details.Age, details.Email, details.FirstName, details.LastName, details.Img).Scan(&id)
	common.HandleError(err)
	fmt.Println("New record ID is:", id)
	http.Redirect(w, r, "/", 301)
}

//Display is ...
func Display(w http.ResponseWriter, r *http.Request) {
	db := common.Conn()
	defer db.Close()

	var details []structures.ContactDetails
	rows, err := db.Query("SELECT *FROM users ORDER BY id")
	common.HandleError(err)

	for rows.Next() {
		var detail structures.ContactDetails
		err = rows.Scan(&detail.ID, &detail.Age, &detail.FirstName, &detail.LastName, &detail.Email, &detail.Img)
		if err != nil {
			fmt.Println(err)
		}
		detail.Image = base64.StdEncoding.EncodeToString(detail.Img)
		details = append(details, detail)
	}
	var datapass structures.Pass
	datapass.Data = details
	tmpl.ExecuteTemplate(w, "display.html", datapass)
}

//Delete is ...
func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println("Delete ID ", id)
	db := common.Conn()
	sqlStatement := "DELETE FROM users where id = $1"
	_, err := db.Exec(sqlStatement, id)
	common.HandleError(err)
	defer db.Close()
	http.Redirect(w, r, "/display", 301)
}

//Update is ...
func Update(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println("Update ID ", id)
	db := common.Conn()

	sqlStatement := `SELECT * FROM users WHERE id = $1`

	var detail structures.ContactDetails
	err := db.QueryRow(sqlStatement, id).Scan(&detail.ID, &detail.Age, &detail.FirstName, &detail.LastName, &detail.Email, &detail.Img)
	common.HandleError(err)

	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "edit.html", detail)
		return
	}

	sqlStatement = `
	UPDATE users
	SET age = $2, first_name = $3, last_name = $4, email = $5, image = $6
	WHERE id = $1;`

	file, _, err1 := r.FormFile("image")
	common.HandleError(err1)
	x, err2 := ioutil.ReadAll(file)
	common.HandleError(err2)

	details := structures.ContactDetails{
		Age:       r.FormValue("age"),
		Email:     r.FormValue("email"),
		FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
		Img:       x,
	}
	_, err = db.Exec(sqlStatement, id, details.Age, details.FirstName, details.LastName, details.Email, detail.Img)
	common.HandleError(err)
	defer db.Close()
	http.Redirect(w, r, "/display", 301)
}
