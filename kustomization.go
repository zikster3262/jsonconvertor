package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type resources []string

type Kustomization struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Namespace  string   `json:"namespace"`
	Resources  []string `json:"resources"`
}

func (kustomize Kustomization) Upload(value []byte) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := s3.New(sess)

	bucket := os.Getenv("S3BUCKET")

	r := bytes.NewReader(value)

	result, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Body:   r,
		Key:    aws.String("kustomization.yaml"),
	})

	log.Printf("File: kustomization.yaml was uploaded do S3 %v bucket", bucket)

	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)

}
