package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Data is ...
type Data struct {
	IP1    string `json:"ip1"`
	IP2    string `json:"ip2"`
	OP     string `json:"op"`
	Result string `json:"result"`
}

func calculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var result Data
	_ = json.NewDecoder(r.Body).Decode(&result)
	x, _ := strconv.ParseFloat(result.IP1, 64)
	y, _ := strconv.ParseFloat(result.IP2, 64)
	var z float64
	if result.OP == "+" {
		z = x + y
	} else if result.OP == "-" {
		z = x - y
	} else if result.OP == "/" {
		z = x / y
	} else if result.OP == "*" {
		z = x * y
	}
	result.Result = fmt.Sprint(z)
	filedata, err := ioutil.ReadFile("newdata.json")
	if err != nil {
		fmt.Println(err)
	}

	var all []Data
	err = json.Unmarshal([]byte(filedata), &all)
	if err != nil {
		fmt.Println("Error Unmarshling for user file")
		fmt.Println(err)
	}

	all = append(all, result)
	file, _ := json.MarshalIndent(all, "", " ")
	_ = ioutil.WriteFile("newdata.json", file, 0644)
	fmt.Println("FOperand-", result.IP1, "SOperand-", result.IP2, "Operator:", result.OP, "Result-", result.Result)
	json.NewEncoder(w).Encode(result)
}

func main() {
	//fmt.Println("GoLang-HTML-CSS-JS integration")
	r := mux.NewRouter()
	fs := http.StripPrefix("/wt/", http.FileServer(http.Dir("./wt")))
	r.PathPrefix("/wt/").Handler(fs)
	http.Handle("/wt/", r)
	/* r.Handle("/", http.FileServer(http.Dir("./wt"))) */
	r.HandleFunc("/cal", calculate).Methods("POST")
	fmt.Println("Server - http://localhost:9050/wt/")
	log.Fatal(http.ListenAndServe(":9050", r))
}
