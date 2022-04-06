package main

import (
	"bytes"
	"context"
	"encoding/json"
	js "encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/buger/jsonparser"
	"gopkg.in/yaml.v2"
)

func HandleRequest(ctx context.Context, event events.SQSEvent) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := s3.New(sess)

	bucket := os.Getenv("S3BUCKET")

	var e map[string]interface{}

	body := event.Records[0].Body
	err := js.Unmarshal([]byte(body), &e)
	if err != nil {
		log.Println(err)
	}
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	jsonparser.ObjectEach(b, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {

		var json map[string]interface{}

		if err := js.Unmarshal(value, &json); err != nil {
			panic(err)
		}

		y, err := yaml.Marshal(json)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		r := bytes.NewReader(y)

		result, err := svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Body:   r,
			Key:    aws.String(string(key) + ".yaml"),
		})

		if err != nil {
			log.Println(err)
		}

		fmt.Println(result)
		return nil
	}, "preprequisities:", "services")

}

func main() {
	lambda.Start(HandleRequest)
}
