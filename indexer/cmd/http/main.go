package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jmoney8080/example-search-platform/indexer/internal/handle"
)

type LambdaResponse struct {
	message string
}

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func init() {

	infoLogger = log.New(os.Stdout,
		"[INFO]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	warningLogger = log.New(os.Stdout,
		"[WARNING]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	errorLogger = log.New(os.Stderr,
		"[ERROR]: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func respond(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responseBody, _ := json.Marshal(LambdaResponse{message: "cannot read request body"})
		http.Error(w, string(responseBody), http.StatusBadRequest)
		return
	}

	eventRecord := events.DynamoDBEventRecord{}

	err = json.Unmarshal(body, &eventRecord)
	if err != nil {
		responseBody, _ := json.Marshal(LambdaResponse{message: "cannot unmarshal request"})
		http.Error(w, string(responseBody), http.StatusBadRequest)
		return
	}

	records := make([]events.DynamoDBEventRecord, 1)
	records[0] = eventRecord

	wat := handle.Request(
		events.DynamoDBEvent{
			Records: records,
		},
	)

	if wat != nil {
		http.Error(w, "{\"message\": \"I dun goofed\"}", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{\"message\": \"OK\"}"))
}

func main() {
	port := os.Getenv("SERVER_PORT")
	infoLogger.Println("Starting server")
	defer infoLogger.Println("Stopping server")

	http.HandleFunc("/", respond)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
