package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"net/http"
)

//Train is...
type Train struct {
	ID     string `json:"id"`
	Number string `json:"number"`
	Tname  string `json:"tname"`
	Starts string `json:"starts"`
	Ends   string `json:"ends"`
}

func storecsv() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

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

	fmt.Println("Connected to MongoDB!")

	csvFile, err := os.Open("All_Indian_Trains.csv")

	defer csvFile.Close()

	if err != nil {
		panic(err)
	}
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	collection := client.Database("train").Collection("trains")
	for _, column := range csvLines {
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
	}
	fmt.Println("Store Complete")
}

func display(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Declare host and port options to pass to the Connect() method
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to the MongoDB and return Client instance
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}

	collection := client.Database("train").Collection("trains")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
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
	//storecsv()
	fmt.Println("Server - http://localhost:8000/")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/display", display)
	http.ListenAndServe(":8000", nil)
}
