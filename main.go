package main

import (
	js "encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/buger/jsonparser"
	"gopkg.in/yaml.v2"
)

func main() {

	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	// // Create an uploader with the session and default options
	// uploader := s3manager.NewUploader(sess)

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

		err = os.WriteFile("kubernetes/"+string(key)+".yaml", y, 0644)
		if err != nil {
			panic(err)
		}

		// Upload the file to S3.
		// result, err := uploader.Upload(&s3manager.UploadInput{
		// 	Bucket: aws.String("testbucketportal"),
		// 	Key:    aws.String(string(key) + ".yaml"),
		// 	Body:   y,
		// })

		return nil
	}, "preprequisities:", "services")

	defer jsonFile.Close()
}
