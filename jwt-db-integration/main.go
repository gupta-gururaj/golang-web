package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type users struct {
	Email    string `json:"email"`
	Name     string `gorm:"unique" json:"name"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type response struct {
	Email string `json:"email"`
	Token string
}

var mySigningKey = []byte("captainjacksparrowsayshi")

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Public Page")
}

func admin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "admin work begins from here . . .")
}
func superadmin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "super-user work begins from here . . .")
}

func conn() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	dbname := os.Getenv("DBNAME")
	url := os.Getenv("URL")
	db, err := gorm.Open(dbname, url)
	if err != nil {
		panic(err)
	}
	return db
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := conn()
	defer db.Close()
	var user users
	_ = json.NewDecoder(r.Body).Decode(&user)
	//check if user already exists
	var body []users
	db.Find(&body) //data fetched from database
	exists := false
	for _, v := range body {
		if user.Email == v.Email {
			fmt.Fprintf(w, "User already exists")
			exists = true
		}
	}
	if !exists {
		user.Password, _ = hashPassword(user.Password)
		json.NewEncoder(w).Encode(user)
		db.Create(&user)
		fmt.Fprintf(w, "Data stored successfully")
	}
}

func generateJWT(user users) string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = user.Email
	claims["name"] = user.Name
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Println("Something Went Wrong")
	}
	return tokenString
}

func display(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := conn()
	defer db.Close()
	var body []users
	db.Find(&body)
	json.NewEncoder(w).Encode(body)
}

func logincheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user login
	_ = json.NewDecoder(r.Body).Decode(&user)
	//check with database
	db := conn()
	defer db.Close()
	var body []users
	db.Find(&body) //data fetched
	exists := false
	for _, v := range body {
		if user.Email == v.Email {
			exists = true
			if checkPasswordHash(user.Password, v.Password) {
				//fmt.Fprintf(w, "Welcome!!!")
				var usr users
				db.Where("email = ?", user.Email).Find(&usr)
				var res response
				res.Email = usr.Email
				res.Token = generateJWT(usr)
				fmt.Fprintf(w, "Your credentials are valid, here is your token buddy!!!\n")
				json.NewEncoder(w).Encode(res)
			} else {
				fmt.Fprintf(w, "You're so so dumb, can't even remember a simple password, no token for you :) :) :)\n")
			}
		}
	}
	if !exists {
		fmt.Fprintf(w, "User donot exist, signup now!!!")
	}
}

func newtable() {
	db := conn()
	db.DropTableIfExists(&users{})
	db.AutoMigrate(&users{})
}

func adminMiddleware(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			fmt.Fprintf(w, "Token Empty")
			return
		}
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "Admin" {
				endpoint(w, r)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	}
}

func superAdminMiddleware(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			fmt.Fprintf(w, "Token Empty")
			return
		}
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "Super Admin" {
				endpoint(w, r)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	}
}

func main() {
	fmt.Println("Server - http://localhost:9020/")
	newTable := flag.Bool("new-table", false, "Create new table and drops the old one")
	flag.Parse()

	if *newTable {
		newtable()
	}

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/website", display).Methods("GET")
	r.HandleFunc("/website", signup).Methods("POST")
	r.HandleFunc("/login", logincheck).Methods("GET")
	r.HandleFunc("/admin", adminMiddleware(admin)).Methods("GET")
	r.HandleFunc("/superadmin", superAdminMiddleware(superadmin)).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":9020", r))
}
