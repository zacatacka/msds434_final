package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemakerruntime"
)

var endpoint1 = "linear-learner-2023-11-30-23-03-02-270"
var endpoint2 = "linear-learner-2023-11-30-23-18-22-710"

type InputData struct {
	Month   int    `json:"month"`
	Carrier string `json:"carrier"`
	Airport string `json:"airport"`
}

func main() {
	http.HandleFunc("/predict", handleRequest)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var data InputData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	csvData := inputDataToCSV(data)

	prediction1, err := getPrediction(csvData, endpoint1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prediction2, err := getPrediction(csvData, endpoint2)
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

func getPrediction(csvData string, endpointName string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
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
