package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemakerruntime"
)

var endpoint1 = "linear-learner-2023-12-02-18-11-46-657"
var endpoint2 = "linear-learner-2023-12-02-17-57-03-454"

type InputData struct {
	Month   int    `json:"month"`
	Carrier string `json:"carrier"`
	Airport string `json:"airport"`
}

var carrierClasses map[string]int
var airportClasses map[string]int

func init() {
	carrierClassesBytes, err := ioutil.ReadFile("data/carrier_classes.json")
	if err != nil {
		log.Fatalf("Failed to read carrier_classes.json: %v", err)
	}
	err = json.Unmarshal(carrierClassesBytes, &carrierClasses)
	if err != nil {
		log.Fatalf("Failed to parse carrier_classes.json: %v", err)
	}

	airportClassesBytes, err := ioutil.ReadFile("data/airport_classes.json")
	if err != nil {
		log.Fatalf("Failed to read airport_classes.json: %v", err)
	}
	err = json.Unmarshal(airportClassesBytes, &airportClasses)
	if err != nil {
		log.Fatalf("Failed to parse airport_classes.json: %v", err)
	}
}

func main() {
	http.HandleFunc("/predict", handleRequest)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request")
    var data InputData
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    dataCopy := data

    carrierClass, ok := carrierClasses[data.Carrier]
    if !ok {
        http.Error(w, "Invalid carrier code", http.StatusBadRequest)
        return
    }
    data.Carrier = strconv.Itoa(carrierClass)

    airportClass, ok := airportClasses[data.Airport]
    if !ok {
        http.Error(w, "Invalid airport code", http.StatusBadRequest)
        return
    }
    data.Airport = strconv.Itoa(airportClass)

    prediction1, err := getPrediction(dataCopy, endpoint1)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    prediction2, err := getPrediction(dataCopy, endpoint2)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := fmt.Sprintf("Predicted Delay Time: %s\nOdds of a Delay: %s", prediction1, prediction2)

    w.Write([]byte(response))
}

func inputDataToCSV(data InputData) string {
	b := &strings.Builder{}
	csvWriter := csv.NewWriter(b)
	csvWriter.Write([]string{
		strconv.Itoa(data.Month),
		data.Carrier,
		data.Airport,
	})
	csvWriter.Flush()

	return b.String()
}

func getPrediction(data InputData, endpointName string) (string, error) {
    carrierCode, ok := carrierClasses[data.Carrier]
    if !ok {
        return "", fmt.Errorf("carrier not found in carrier_classes.json")
    }
    airportCode, ok := airportClasses[data.Airport]
    if !ok {
        return "", fmt.Errorf("airport not found in airport_classes.json")
    }
    csvData := fmt.Sprintf("%d,%d,%d", data.Month, carrierCode, airportCode)
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
        Config: aws.Config{
            Region: aws.String("us-east-2"),
        },
    }))
    sagemakerClient := sagemakerruntime.New(sess)
    input := &sagemakerruntime.InvokeEndpointInput{
        Body:         []byte(csvData),
        ContentType:  aws.String("text/csv"),
        EndpointName: aws.String(endpointName),
    }
    output, err := sagemakerClient.InvokeEndpoint(input)
    if err != nil {
        return "", err
    }
    prediction := string(output.Body)
    return prediction, nil
}
