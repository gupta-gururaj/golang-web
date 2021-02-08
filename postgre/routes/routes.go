package routes

import (
	"fmt"
	"net/http"
	configuration "postgre/configuration/controller"
)

//StartServer is ...
func StartServer() {
	fmt.Println("Server - http://localhost:9000/")
	http.ListenAndServe(":9000", nil)
}

//Routes is ...
func Routes() {
	http.HandleFunc("/", configuration.Form)
	http.HandleFunc("/display", configuration.Display)
	http.HandleFunc("/delete", configuration.Delete)
	http.HandleFunc("/update", configuration.Update)
}
