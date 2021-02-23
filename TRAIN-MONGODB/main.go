package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Train struct {
	NO     string `json:"no"`
	NAME   string `json:"name"`
	STARTS string `json:"starts"`
	ENDS   string `json:"ends"`
}

func connection() *mongo.Client {
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
	return client
}

func closedatabase() {
	client := connection()
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to MongoDB closed.")
}

func getallTrains(w http.ResponseWriter, r *http.Request) {

	client := connection()
	defer closedatabase()

	collection := client.Database("traindb").Collection("trains")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}
	var trains []Train
	if err = cursor.All(context.TODO(), &trains); err != nil {
		log.Fatal(err)
	}
	bytedata, err := json.MarshalIndent(trains, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytedata)
}

func readcsv(w http.ResponseWriter, r *http.Request) {
	csvFile, err := os.Open("All_Indian_Trains.csv")
	client := connection()
	defer closedatabase()
	collection := client.Database("traindb").Collection("trains")
	defer csvFile.Close()

	if err != nil {
		panic(err)
	}

	// load file then skip Header
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Read()
	var trains []Train
	// loop through each record create a vehicle object and import
	for {
		var train Train
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		train.NO = record[1]
		train.NAME = record[2]
		train.STARTS = record[3]
		train.ENDS = record[4]
		trains = append(trains, train)
		insertResult, err := collection.InsertOne(context.TODO(), train)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	}

	fmt.Println(trains)
	fmt.Fprintf(w, "reading csv done")
}

func main() {
	fs := http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates")))
	http.Handle("/templates/", fs)
	http.HandleFunc("/readcsv", readcsv)
	http.HandleFunc("/Trains", getallTrains)
	fmt.Println("server started at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
