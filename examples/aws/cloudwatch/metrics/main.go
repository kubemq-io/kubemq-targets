package main

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/pkg/uuid"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"time"
)

func main() {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// put Metrics
	namespace := "Logs"

	b := []byte(`[{"Counts":null,"Dimensions":null,"MetricName":"New Metric","StatisticValues":null,"StorageResolution":null,"Timestamp":"2020-08-10T17:09:48.3895822+03:00","Unit":"Count","Value":132,"Values":null}]`)

	putRequest := types.NewRequest().
		SetMetadataKeyValue("method", "put_metrics").
		SetMetadataKeyValue("namespace", namespace).
		SetData(b)

	queryPutResponse, err := client.SetQuery(putRequest.ToQuery()).
		SetChannel("query.aws.cloudwatch.metrics").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	putResponse, err := types.ParseResponse(queryPutResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("put metrics executed, response: %s", putResponse.Data))

	// list metrics
	listRequest := types.NewRequest().
		SetMetadataKeyValue("method", "list_metrics").
		SetMetadataKeyValue("namespace", namespace)

	listReq, err := client.SetQuery(listRequest.ToQuery()).
		SetChannel("query.aws.cloudwatch.metrics").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	listResponse, err := types.ParseResponse(listReq.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("list metrics executed %v, error: %v", listResponse.Data, listResponse.IsError))
}
