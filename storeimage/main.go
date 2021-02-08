// main.go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

var tmpl = template.Must(template.ParseGlob("template/*"))

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "index.html", nil)
		return
	}
	fmt.Fprintf(w, "File Upload Endpoint Hit")

	//S-1 - Fetch file from form
	file, handler, err := r.FormFile("file")
	handleError(err)
	defer file.Close()
	fmt.Println(handler.Size, handler.Filename, handler.Header)

	//S-2 - Create a temporary file
	tempFile, err := ioutil.TempFile("profiles", "uploaded-*.png")
	handleError(err)
	defer tempFile.Close()

	//S-3 - Read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	handleError(err)

	//S-4 - Write this byte array to our temporary file
	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func setupRoutes() {
	http.HandleFunc("/", uploadFile)
	http.ListenAndServe(":9000", nil)
}

func main() {
	fmt.Println("Server - http://localhost:9000/")
	setupRoutes()
}
