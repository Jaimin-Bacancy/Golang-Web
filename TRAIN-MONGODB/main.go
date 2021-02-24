package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
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
	dburl := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(dburl)

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

func getColletion(client *mongo.Client, dbname string, colletionname string) *mongo.Collection {
	collection := client.Database(dbname).Collection(colletionname)
	return collection
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

func readcsv() {
	csvFile, err := os.Open("All_Indian_Trains.csv")
	defer csvFile.Close()
	client := connection()
	defer closedatabase()
	collection := getColletion(client, "traindb", "trains")

	if err != nil {
		panic(err)
	}

	// load file then skip Header
	reader := csv.NewReader(bufio.NewReader(csvFile))

	// loop through each record create a train object and import
	limit := 5
	channel := make(chan int, limit)
	for {

		var train Train
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		channel <- 1
		go func(train *Train, record *[]string) {
			train.NO = (*record)[1]
			train.NAME = (*record)[2]
			train.STARTS = (*record)[3]
			train.ENDS = (*record)[4]
			insertResult, err := collection.InsertOne(context.TODO(), train)
			fmt.Println("Inserted a single document: ", insertResult.InsertedID)
			if err != nil {
				log.Fatal(err)
			}
			<-channel
		}(&train, &record)

	}

	for i := 1; i <= limit; i++ {
		channel <- 1
	}
	fmt.Println("reading csv done")
}

func main() {
	useCsvread := flag.Bool("readcsv", false, "")
	flag.Parse()
	if *useCsvread {
		readcsv()
	}
	fs := http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates")))
	http.Handle("/templates/", fs)
	http.HandleFunc("/Trains", getallTrains)
	fmt.Println("server started at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
