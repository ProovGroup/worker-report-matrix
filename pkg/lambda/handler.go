package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"worker-report-matrix/internal/config"
	"worker-report-matrix/internal/sqs"

	"github.com/ProovGroup/env"
	"github.com/aws/aws-lambda-go/events"
)

func Handler(ctx context.Context, event events.SQSEvent) error {
	for i := range event.Records {
		e, err := env.GetEnvSqsArnSSM(event.Records[i].EventSourceARN, os.Getenv("AWS_REGION"), env.BDDWrite)
		if err != nil {
			fmt.Println("[ERROR] env.GetEnvSqsArnSSM:", err)
			return err
		}

		var incomingMessage sqs.IncomingMessage
		if err = json.Unmarshal([]byte(event.Records[i].Body), &incomingMessage); err != nil {
			fmt.Println("[ERROR] json.Unmarshal:", err)
			return err
		}

		fmt.Println("[INFO] Generating message for proov_code:", incomingMessage.ProovCode)

		var message *sqs.OutgoingMessage
		message, err = config.GetMessageStructure()
		if err != nil {
			fmt.Println("[ERROR] config.GetMessageStructure:", err)
			return err
		}

		err = message.GetReport(&e, incomingMessage.ProovCode)
		if err != nil {
			fmt.Println("[ERROR] message.GetReport:", err)
			return err
		}

		err = message.Send(&e)
		if err != nil {
			fmt.Println("[ERROR] message.Send:", err)
			return err
		}

		fmt.Println("[INFO] Message sent to SQS successfully")
	}

	return nil
}

