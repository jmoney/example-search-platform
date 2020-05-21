package handle

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

type IndexDocument struct {
	AgentID    string `json:agentid`
	DatasetID  string `json:datasetid`
	NumVersion string `json:numversion`
}

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

// Request handling the request
func Request(event events.DynamoDBEvent) error {
	es, err := elasticsearch.NewDefaultClient()

	if err != nil {
		errorLogger.Printf("Error creating the client: %s", err)
		return err
	}

	for _, record := range event.Records {
		agentid := record.Change.Keys["agentid"].String() # FIXME: Check Type
		datasetid := record.Change.Keys["datasetid"].String() # FIXME: Check Type
		numVersions, _ := record.Change.NewImage["numVersions"].Integer() # FIXME: Check Type

		infoLogger.Printf("%s %s", agentid, datasetid)

		// Build the request body.
		indexDocument := IndexDocument{
			AgentID:    agentid,
			DatasetID:  datasetid,
			NumVersion: strconv.FormatInt(numVersions, 10),
		}

		jsonDocument, _ := json.Marshal(indexDocument)

		infoLogger.Printf("Indexable Document: %s", jsonDocument)
		// Set up the request object.
		req := esapi.IndexRequest{
			Index:      "datasets",
			DocumentID: strconv.FormatInt(numVersions, 10),
			Body:       strings.NewReader(string(jsonDocument)),
			Refresh:    "true",
		}

		// Perform the request with the client.
		res, err := req.Do(context.Background(), es)
		if err != nil {
			errorLogger.Printf("Error getting response: %s", err)
			return err
		}
		defer res.Body.Close()

		if res.IsError() {
			errorLogger.Printf("[%s] Error indexing document ID=%d", res.Status(), numVersions)
		} else {
			// Deserialize the response into a map.
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				errorLogger.Printf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and indexed document version.
				infoLogger.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
			}
		}
	}

	return nil
}
