package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/lib/pq"
)

func createTable() {
	db := conn()
	db.DropTableIfExists(&users{})
	db.AutoMigrate(&users{})
}

func conn() *gorm.DB {
	db, err := gorm.Open("postgres", "user=postgres password=root dbname=gorm sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}

type users struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func insert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := conn()
	defer db.Close()
	var user users
	_ = json.NewDecoder(r.Body).Decode(&user)
	json.NewEncoder(w).Encode(user)
	db.Create(&user)
}

func display(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := conn()
	defer db.Close()
	var body []users
	db.Find(&body)
	jsonData, err := json.MarshalIndent(body, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	w.Write(jsonData)
}

func particularDisplay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := conn()
	defer db.Close()
	id := mux.Vars(r)["id"]
	var body []users
	db.Find(&body, id)
	jsonData, err := json.MarshalIndent(body, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	w.Write(jsonData)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := conn()
	defer db.Close()
	id := mux.Vars(r)["id"]
	db.Table("users").Where("id= ?", id).Delete(&users{})
	var user users
	_ = json.NewDecoder(r.Body).Decode(&user)
	db.Create(&user)
	json.NewEncoder(w).Encode(user)
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := conn()
	defer db.Close()
	id := mux.Vars(r)["id"]
	db.Table("users").Where("id= ?", id).Delete(&users{})
	var body []users
	db.Find(&body)
	jsonData, err := json.MarshalIndent(body, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	w.Write(jsonData)
}

func main() {
	fmt.Println("Server - http://localhost:9030/")
	r := mux.NewRouter()
	r.HandleFunc("/database", insert).Methods("POST")
	r.HandleFunc("/database", display).Methods("GET")
	r.HandleFunc("/database/{id}", particularDisplay).Methods("GET")
	r.HandleFunc("/database/{id}", update).Methods("PUT")
	r.HandleFunc("/database/{id}", delete).Methods("DELETE")
	http.ListenAndServe(":9030", r)
}
