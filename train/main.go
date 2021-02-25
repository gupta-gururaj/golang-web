package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"net/http"
)

var concurrecyLimit = 10
var done = make(chan bool, concurrecyLimit)

//Train is...
type Train struct {
	ID     string `json:"id"`
	Number string `json:"number"`
	Tname  string `json:"tname"`
	Starts string `json:"starts"`
	Ends   string `json:"ends"`
}

func dbConn() *mongo.Client {
	er := godotenv.Load(".env")
	if er != nil {
		log.Fatalf("Error loading .env file")
	}
	uri := os.Getenv("URI")
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Connected to MongoDB!")
	return client
}

func insert(column []string, collection *mongo.Collection) {
	train := Train{
		ID:     column[0],
		Number: column[1],
		Tname:  column[2],
		Starts: column[3],
		Ends:   column[4],
	}
	_, err := collection.InsertOne(context.TODO(), train)
	if err != nil {
		log.Fatal(err)
	}
	//<-done
}

func storecsv() {
	er := godotenv.Load(".env")
	if er != nil {
		log.Fatalf("Error loading .env file")
	}
	database := os.Getenv("DB")
	col := os.Getenv("COLLECTION")
	client := dbConn()
	csvFile, err := os.Open("All_Indian_Trains.csv")

	defer csvFile.Close()

	if err != nil {
		panic(err)
	}
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	collection := client.Database(database).Collection(col)
	for _, column := range csvLines {
		done <- true
		go insert(column, collection)
	}
	for i := 0; i < concurrecyLimit; i++ {
		done <- true
	}

	fmt.Println("Store Complete")
}

func simpleStore() {
	er := godotenv.Load(".env")
	if er != nil {
		log.Fatalf("Error loading .env file")
	}
	database := os.Getenv("DB")
	col := os.Getenv("COLLECTION")
	client := dbConn()
	csvFile, err := os.Open("All_Indian_Trains.csv")

	defer csvFile.Close()

	if err != nil {
		panic(err)
	}
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	collection := client.Database(database).Collection(col)
	for _, column := range csvLines {
		insert(column, collection)
	}
	fmt.Println("Store Complete")
}

func limitedDisplay(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	x, _ := strconv.Atoi(query.Get("page"))
	x--
	skip := int64(x * 10)
	if skip < 0 {
		skip = 0
	}
	er := godotenv.Load(".env")
	if er != nil {
		log.Fatalf("Error loading .env file")
	}
	uri := os.Getenv("URI")
	database := os.Getenv("DB")
	col := os.Getenv("COLLECTION")

	w.Header().Set("Content-Type", "application/json")
	// Declare host and port options to pass to the Connect() method
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to the MongoDB and return Client instance
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}

	collection := client.Database(database).Collection(col)
	options := options.Find()
	options.SetLimit(10)
	options.SetSkip(skip)
	cursor, err := collection.Find(context.TODO(), bson.D{}, options)
	if err != nil {
		log.Fatal(err)
	}
	var trains []Train
	if err = cursor.All(context.TODO(), &trains); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(trains)
}

func main() {
	//log.Printf("Time taken without concurrecy 285ms")
	csvStore := flag.Bool("storecsv", false, "Store CSV file in mongoDB")
	flag.Parse()

	if *csvStore {
		//start := time.Now()
		//storecsv()
		simpleStore()
		//elapsed := time.Since(start)
		//log.Printf("Time taken with concurrency %s", elapsed)
	}
	fmt.Println("Server - http://localhost:8000/")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/display", display)
	http.HandleFunc("/limit/", limitedDisplay)
	http.ListenAndServe(":8000", nil)
}
