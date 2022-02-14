package main

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"io/ioutil"
	"log"
	"time"
)

type testStructure struct {
	region string

	json     string
	domain   string
	index    string
	endpoint string
	service  string
	id       string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./credentials/aws/region.txt")
	if err != nil {
		return nil, err
	}
	t.region = string(dat)

	dat, err = ioutil.ReadFile("./credentials/aws/elasticsearch/signer/json.txt")
	if err != nil {
		return nil, err
	}
	t.json = string(dat)
	dat, err = ioutil.ReadFile("./credentials/aws/elasticsearch/signer/domain.txt")
	if err != nil {
		return nil, err
	}
	t.domain = string(dat)
	dat, err = ioutil.ReadFile("./credentials/aws/elasticsearch/signer/index.txt")
	if err != nil {
		return nil, err
	}
	t.index = string(dat)
	dat, err = ioutil.ReadFile("./credentials/aws/elasticsearch/signer/endpoint.txt")
	if err != nil {
		return nil, err
	}
	t.endpoint = string(dat)
	dat, err = ioutil.ReadFile("./credentials/aws/elasticsearch/signer/service.txt")
	if err != nil {
		return nil, err
	}
	t.service = string(dat)
	dat, err = ioutil.ReadFile("./credentials/aws/elasticsearch/signer/id.txt")
	if err != nil {
		return nil, err
	}
	t.id = string(dat)

	return t, nil
}

func main() {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	dat, err := getTestStructure()
	if err != nil {
		log.Fatal(err)
	}
	// POST
	request := types.NewRequest().
		SetMetadataKeyValue("method", "POST").
		SetMetadataKeyValue("region", dat.region).
		SetMetadataKeyValue("json", dat.json).
		SetMetadataKeyValue("domain", dat.domain).
		SetMetadataKeyValue("endpoint", dat.endpoint).
		SetMetadataKeyValue("index", dat.index).
		SetMetadataKeyValue("service", dat.service).
		SetMetadataKeyValue("id", dat.id)
	queryResponse, err := client.SetQuery(request.ToQuery()).
		SetChannel("query.aws.elasticsearch").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	response, err := types.ParseResponse(queryResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("executed, response: %s", response.Data))
}
