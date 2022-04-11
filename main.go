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
	services := resources{}

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
		// create resources yaml for kustomization.yaml file
		services = append(services, string(key)+".yaml")

		m[string(key)] = y

		return nil
	}, "preprequisities:", "services")

	// Create Kustomization yaml with resources from services resource
	kustomization := Kustomization{
		APIVersion: "kustomize.config.k8s.io/v1beta1",
		Kind:       "Kustomization",
		Namespace:  "default",
		Resources:  services,
	}

	kyaml, err := yaml.Marshal(kustomization)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	kustomization.Upload(kyaml)

	if m["fluent-bit"] != nil {
		fluentbit := NewChart("fluent", "https://fluent.github.io/helm-charts", "fluent/fluent-bit", "default", m["fluent-bit"])
		fluentbit.Upload("fluent-bit", fluentbit.Template())
	}

	if m["kube-prometheus-stack"] != nil {
		prometheus := NewChart("helm-charts", "https://prometheus-community.github.io/helm-charts", "helm-charts/kube-prometheus-stack", "default", m["kube-prometheus-stack"])
		prometheus.Upload("kube-prometheus-stack", prometheus.Template())
	}

	if m["ingress-nginx"] != nil {
		nginx := NewChart("ingress-nginx", "https://kubernetes.github.io/ingress-nginx", "ingress-nginx/ingress-nginx", "default", m["ingress-nginx"])
		nginx.Upload("ingress-nginx", nginx.Template())
	}

	if m["external-dns"] != nil {
		externaldns := NewChart("bitnami", "https://charts.bitnami.com/bitnami", "bitnami/external-dns", "default", m["external-dns"])
		externaldns.Upload("external-dns", externaldns.Template())
	}

	log.Printf("Servces: %v", services)

}

func main() {
	lambda.Start(HandleRequest)
}
