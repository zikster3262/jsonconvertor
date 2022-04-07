package main

import (
	"context"
	"encoding/json"
	js "encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/buger/jsonparser"
	"gopkg.in/yaml.v2"
)

func HandleRequest(ctx context.Context, event events.SQSEvent) {

	m := make(map[string][]byte)

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

		m[string(key)] = y

		return nil
	}, "preprequisities:", "services")

	prometheus := NewChart("helm-charts", "https://prometheus-community.github.io/helm-charts", "helm-charts/kube-prometheus-stack", "default", m["kube-prometheus-stack"])
	nginx := NewChart("ingress-nginx", "https://kubernetes.github.io/ingress-nginx", "ingress-nginx/ingress-nginx", "default", m["ingress-nginx"])
	fluentbit := NewChart("fluent", "https://fluent.github.io/helm-charts", "fluent/fluent-bit", "default", m["fluent-bit"])

	prometheus.Upload("prometheus", prometheus.Template())
	nginx.Upload("nginx", nginx.Template())
	fluentbit.Upload("fluent-bit", nginx.Template())
}

func main() {
	lambda.Start(HandleRequest)
}
