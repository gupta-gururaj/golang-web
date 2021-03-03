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
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"net/http"
	"strings"
)

var (
	concurrecyLimit = 100
	done            = make(chan bool, concurrecyLimit)
)

//Train is...
type Train struct {
	TrainNo   string `json:"trainNo"`
	TrainName string `json:"trainName"`
	SEQ       string `json:"seq"`
	Code      string `json:"code"`
	StName    string `json:"stName"`
	ATime     string `json:"atime"`
	DTime     string `json:"dtime"`
	Distance  string `json:"distance"`
	SS        string `json:"ss"`
	SSname    string `json:"ssname"`
	Ds        string `json:"ds"`
	DsName    string `json:"dsName"`
	timediff  int
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
		TrainNo:   column[0],
		TrainName: column[1],
		SEQ:       column[2],
		Code:      column[3],
		StName:    column[4],
		ATime:     column[5],
		DTime:     column[6],
		Distance:  column[7],
		SS:        column[8],
		SSname:    column[9],
		Ds:        column[10],
		DsName:    column[11],
	}
	_, err := collection.InsertOne(context.TODO(), train)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func concurrentStore() {
	er := godotenv.Load(".env")
	if er != nil {
		log.Fatalf("Error loading .env file")
	}
	database := os.Getenv("DB")
	col := os.Getenv("COLLECTION")
	client := dbConn()
	csvFile, err := os.Open("rail-data.csv")

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
	csvFile, err := os.Open("rail-data.csv")

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

func alltrainsbetween(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ss := query.Get("ss")
	ds := query.Get("ds")
	//fmt.Println(reflect.TypeOf(ds))
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
	cursor, err := collection.Find(context.TODO(), bson.M{"code": ss})
	if err != nil {
		log.Fatal(err)
	}
	var source []Train
	if err = cursor.All(context.TODO(), &source); err != nil {
		log.Fatal(err)
	}
	//fmt.Println(source)
	cursor, err = collection.Find(context.TODO(), bson.M{"code": ds})
	if err != nil {
		log.Fatal(err)
	}
	var destination []Train
	if err = cursor.All(context.TODO(), &destination); err != nil {
		log.Fatal(err)
	}
	//fmt.Println(destination)
	var trains []Train
	count := 0
	for i := 0; i < len(source); i++ {
		for j := 0; j < len(destination); j++ {
			if source[i].TrainNo == destination[j].TrainNo {
				a, _ := strconv.Atoi(source[i].SEQ)
				b, _ := strconv.Atoi(destination[j].SEQ)
				if a < b {
					//fmt.Println(ss, "-", source[i].SEQ, ds, "-", destination[j].SEQ)
					time := subtractTime(source[i].ATime, destination[j].ATime)
					trains = append(trains, source[i])
					trains[count].timediff = time
					count++
				}
			}
		}
	}
	/* for i := 0; i < len(trains); i++ {
		fmt.Println(trains[i].timediff)
	} */
	fmt.Println("-------------")
	sort(trains)
	for i := 0; i < len(trains); i++ {
		fmt.Println(trains[i])
	}
	json.NewEncoder(w).Encode(trains)
	fmt.Println("EOM", count)
}

func sort(trains []Train) []Train { //bubble sort
	n := len(trains)
	swapped := true
	for swapped {
		swapped = false
		for i := 1; i < n; i++ {
			if trains[i-1].timediff > trains[i].timediff {
				trains[i], trains[i-1] = trains[i-1], trains[i]
				swapped = true
			}
		}
	}
	return trains
}

func subtractTime(time1, time2 string) int {
	time1 = strings.Replace(time1, ":", "", -1)
	time2 = strings.Replace(time2, ":", "", -1)
	t1, _ := strconv.Atoi(time1)
	t2, _ := strconv.Atoi(time2)
	if t1 > t2 {
		return t1 - t2
	}
	return t2 - t1
}

func main() {
	csvStore := flag.Bool("storecsv", false, "Store CSV file in mongoDB")
	flag.Parse()

	if *csvStore {
		fmt.Println("Time taken without concurrecy 12s")
		start := time.Now()
		concurrentStore()
		//simpleStore()
		elapsed := time.Since(start)
		fmt.Printf("Time taken with concurrency limit 100 is %s\n", elapsed)
	}
	fmt.Println("Server - http://localhost:9050/")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/limit/", limitedDisplay)
	http.HandleFunc("/alltrains/", alltrainsbetween)
	http.ListenAndServe(":9050", nil)
}
