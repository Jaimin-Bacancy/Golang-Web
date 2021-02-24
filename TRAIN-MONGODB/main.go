package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
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

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}

func connection() *mongo.Client {
	// Set client options
	dburl := os.Getenv("DBURL")
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
	databaseName := os.Getenv("DATABASE_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	collection := client.Database(databaseName).Collection(collectionName)
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

func readCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return rows, nil
}

func insertToDatabase() {
	filename := os.Getenv("CSV_FILENAME")
	rows, err := readCsv(filename)
	if err != nil {
		panic(err)
	}
	client := connection()
	defer closedatabase()
	collection := getColletion(client, "traindb", "trains")

	limit := 10
	channel := make(chan int, limit)
	for _, record := range rows {
		var train Train
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

	start := time.Now()
	useCsvread := flag.Bool("readcsv", false, "")
	flag.Parse()
	fmt.Println(*useCsvread)
	if *useCsvread {
		insertToDatabase()
	}
	fmt.Println(time.Since(start))
	fs := http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates")))
	http.Handle("/templates/", fs)
	http.HandleFunc("/Trains", getallTrains)
	fmt.Println("server started at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
