package main

import (
	"bytes"
	js "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/buger/jsonparser"
	"gopkg.in/yaml.v2"
)

func main() {

	bucket := os.Getenv("S3BUCKET")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := s3.New(sess)

	jsonFile, err := os.Open("input.json")

	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonparser.ObjectEach(byteValue, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {

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

	defer jsonFile.Close()
}
