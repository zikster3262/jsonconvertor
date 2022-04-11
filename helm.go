package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
)

type chartRepo struct {
	ReleaseName string
	URL         string
	ChartName   string
	Namespace   string
	Value       []byte
}

func NewChart(ReleaseName string, URL string, ChartName string, Namespace string, Value []byte) *chartRepo {
	ch := chartRepo{
		ReleaseName: ReleaseName,
		URL:         URL,
		ChartName:   ChartName,
		Namespace:   Namespace,
		Value:       Value,
	}
	return &ch
}

func (chart chartRepo) Template() []byte {

	opt := &helmclient.Options{
		Namespace:        chart.Namespace, // Change this to the namespace you wish the client to operate in.
		RepositoryCache:  "/tmp/.helmcache",
		RepositoryConfig: "/tmp/.helmrepo",
		Debug:            true,
		Linting:          true,
		DebugLog:         func(format string, v ...interface{}) {},
	}

	helmClient, err := helmclient.New(opt)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	_ = helmClient

	chartRepo := repo.Entry{
		Name: chart.ReleaseName,
		URL:  chart.URL,
	}

	if err := helmClient.AddOrUpdateChartRepo(chartRepo); err != nil {
		log.Printf("Error: %v", err)
	}

	chartSpec := helmclient.ChartSpec{
		ReleaseName: chart.ReleaseName,
		ChartName:   chart.ChartName,
		Namespace:   chart.Namespace,
		UpgradeCRDs: true,
		Wait:        true,
		ValuesYaml:  string(chart.Value),
	}

	chartTemp, err := helmClient.TemplateChart(&chartSpec)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	return chartTemp
}

func (chart chartRepo) Upload(name string, value []byte) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := s3.New(sess)

	bucket := os.Getenv("S3BUCKET")

	r := bytes.NewReader(value)

	result, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Body:   r,
		Key:    aws.String(string(name) + ".yaml"),
	})

	log.Printf("File: %v.yaml was uploaded do S3 %v bucket", name, bucket)

	if err != nil {
		log.Printf("Error: %v", err)
	}
	fmt.Println(result)

}
