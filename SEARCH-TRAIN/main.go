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
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Train struct {
	TrainNo                string `json:"trainno"`
	TrainName              string `json:"trainname"`
	SEQ                    string `json:"trainseq"`
	StationCode            string `json:"stationcode"`
	StationName            string `json:"stationname"`
	ArivalTime             string `json:"arrivaltime"`
	DepartureTime          string `json:"departuretime"`
	Distance               string `json:"distance"`
	SourceStation          string `json:"sourcestation"`
	SourceStationname      string `json:"sourcestationname"`
	DestinationStation     string `json:"destinationstation"`
	DestinationStationName string `json:"destinationsationname"`
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

func getLimitTrain(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageint, err := strconv.Atoi(page)
	if err != nil {
		panic(err)
	}
	client := connection()
	defer closedatabase()
	databaseName := os.Getenv("DATABASE_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	findOptions := options.Find() // build a `findOptions`
	findOptions.SetLimit(10)
	if pageint < 0 {
		pageint = 0
	}
	findOptions.SetSkip(int64(pageint) * 10) // set limit for record
	collection := client.Database(databaseName).Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
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

func searchTrain(w http.ResponseWriter, r *http.Request) {

	sourcestation := strings.ToUpper(r.URL.Query().Get("sourcestation"))
	//destinationstation := strings.ToUpper(r.URL.Query().Get("destinatiostation"))
	client := connection()
	defer closedatabase()
	databaseName := os.Getenv("DATABASE_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	collection := client.Database(databaseName).Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), bson.M{"stationcode": sourcestation})
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
	rows = rows[1:]
	channel := make(chan int, limit)
	for _, record := range rows {
		var train Train
		train.TrainNo = record[0]
		train.TrainName = record[1]
		train.SEQ = record[2]
		train.StationCode = record[3]
		train.StationName = record[4]
		train.ArivalTime = record[5]
		train.DepartureTime = record[6]
		train.Distance = record[7]
		train.SourceStation = record[8]
		train.SourceStationname = record[9]
		train.DestinationStation = record[10]
		train.DestinationStationName = record[11]
		channel <- 1
		go func(trainptr *Train) {

			_, err := collection.InsertOne(context.TODO(), trainptr)
			//insertResult, err := collection.InsertOne(context.TODO(), trainptr)

			//fmt.Println("Inserted a single document: ", insertResult.InsertedID)
			if err != nil {
				log.Fatal(err)
			}
			<-channel
		}(&train)

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
	http.HandleFunc("/LimitTrain", getLimitTrain)
	http.HandleFunc("/SearchTrain", searchTrain)

	fmt.Println("server started at http://localhost" + os.Getenv("SERVER_PORT"))
	port := os.Getenv("SERVER_PORT")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
