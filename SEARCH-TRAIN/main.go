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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	conString      = "mongodb://localhost:27017"
	dbName         = "test"
	collectionName = "railway1"
	limit          = 10
)

//var wg sync.WaitGroup
var (
	ch = make(chan int, limit)
)

type Data struct {
	TrainNo   string `bson:"trainNo   string"`
	TrainName string `bson:"trainName string"`
	SEQ       string `bson:"seq       string"`
	Code      string `bson:"code      string"`
	StName    string `bson:"stName    string"`
	ATime     string `bson:"atime     string"`
	DTime     string `bson:"dtime     string"`
	Distance  string `bson:"distance  string"`
	SS        string `bson:"ss        string"`
	SSname    string `bson:"ssname    string"`
	Ds        string `bson:"ds        string"`
	DsName    string `bson:"dsName    string"`
}

func ReadCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func dbConn() (*mongo.Collection, *mongo.Client) {
	clientOptions := options.Client().ApplyURI(conString)

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

	collection := client.Database(dbName).Collection(collectionName)
	fmt.Println("Connected to MongoDB!")

	return collection, client

}
func getallTrains(w http.ResponseWriter, r *http.Request) {

	collection, client := dbConn()

	defer client.Disconnect(context.TODO())

	cursor, err := collection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}
	var trains []Data
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

func insertData() {

	collection, client := dbConn()

	defer client.Disconnect(context.TODO())

	csvData, err := ReadCsv("Indian_railway1.csv")
	if err != nil {
		panic(err)
	}

	csvData = csvData[1:501]

	for _, line := range csvData {
		ch <- 1
		//wg.Add(1)
		func() {
			//defer wg.Done()
			data := Data{
				TrainNo:   line[0],
				TrainName: line[1],
				SEQ:       line[2],
				Code:      line[3],
				StName:    line[4],
				ATime:     line[5],
				DTime:     line[6],
				Distance:  line[7],
				SS:        line[8],
				SSname:    line[9],
				Ds:        line[10],
				DsName:    line[11],
			}

			_, err := collection.InsertOne(context.TODO(), data)
			if err != nil {
				panic(err)
			}
			<-ch
			//break
		}()
		//wg.Wait()

		//mt.Println("done")
	}
	for i := 0; i < limit; i++ {
		ch <- 1
	}
}
func main() {
	start := time.Now()
	read := flag.Bool("insert", false, "a bool")
	flag.Parse()
	if *read {
		insertData()
	} else {
		fmt.Println("failed")
	}

	elp := time.Since(start)
	fmt.Println(elp)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/Trains", getallTrains)

	http.ListenAndServe(":8080", nil)
}
