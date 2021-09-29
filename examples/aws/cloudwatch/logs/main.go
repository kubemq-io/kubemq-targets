package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"io/ioutil"
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
	// put Logs
	dat, err := ioutil.ReadFile("./credentials/aws/cloudwatch/logs/logGroupName.txt")
	if err != nil {
		panic(err)
	}
	logGroupName := string(dat)

	dat, err = ioutil.ReadFile("./credentials/aws/cloudwatch/logs/sequenceToken.txt")
	if err != nil {
		panic(err)
	}
	sequenceToken := string(dat)

	dat, err = ioutil.ReadFile("./credentials/aws/cloudwatch/logs/logStreamName.txt")
	if err != nil {
		panic(err)
	}
	logStreamName := string(dat)
	currentTime := time.Now().UnixNano() / 1000000
	m := make(map[int64]string)
	m[currentTime-15] = "my first message to send"
	m[currentTime] = "my second message to send"
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	listRequest := types.NewRequest().
		SetMetadataKeyValue("method", "put_log_event").
		SetMetadataKeyValue("log_group_name", logGroupName).
		SetMetadataKeyValue("sequence_token", sequenceToken).
		SetMetadataKeyValue("log_stream_name", logStreamName).
		SetData(b)

	queryListResponse, err := client.SetQuery(listRequest.ToQuery()).
		SetChannel("query.aws.cloudwatch.logs").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	listResponse, err := types.ParseResponse(queryListResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("put logs executed, response: %s", listResponse.Data))

	// get Logs
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get_log_event").
		SetMetadataKeyValue("log_group_name", logGroupName).
		SetMetadataKeyValue("log_stream_name", logStreamName)

	getReq, err := client.SetQuery(getRequest.ToQuery()).
		SetChannel("query.aws.cloudwatch.logs").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getResponse, err := types.ParseResponse(getReq.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get logs executed %v, error: %v", getResponse.Data, getResponse.IsError))
}
