package config

import (
	"encoding/json"
	"io"
	"os"
	"worker-report-matrix/internal/sqs"
)

var (
	messageStructure *sqs.OutgoingMessage

	SOURCE_REGION = os.Getenv("SOURCE_REGION")
	SOURCE_BUCKET = os.Getenv("SOURCE_BUCKET")
	SOURCE_KEY    = os.Getenv("SOURCE_KEY")

	DESTINATION_REGION = os.Getenv("DESTINATION_REGION")
	DESTINATION_BUCKET = os.Getenv("DESTINATION_BUCKET")
	DESTINATION_KEY    = os.Getenv("DESTINATION_KEY")
)

func GetMessageStructure() (*sqs.OutgoingMessage, error) {
	if messageStructure == nil {
		var err error
		messageStructure, err = LoadMessageStructure()
		if err != nil {
			// App cannot run without message structure
			panic(err)
		}
	}

	ResetMessageStructure()

	return messageStructure, nil
}

func ResetMessageStructure() {
	messageStructure.Source.Path.Region = SOURCE_REGION
	messageStructure.Source.Path.Bucket = SOURCE_BUCKET
	messageStructure.Source.Path.Key    = SOURCE_KEY

	messageStructure.Destination.Path.Region = DESTINATION_REGION
	messageStructure.Destination.Path.Bucket = DESTINATION_BUCKET
	messageStructure.Destination.Path.Key    = DESTINATION_KEY

	messageStructure.DataSource.Content = nil
}

func LoadMessageStructure() (*sqs.OutgoingMessage, error) {
	filePath := "assets/message_structure.json"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var message sqs.OutgoingMessage
	err = json.Unmarshal(bytes, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
