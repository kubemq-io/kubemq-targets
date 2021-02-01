package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/pkg/uuid"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"strconv"
	"time"
)

const mapping = `{
	"settings": {
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings": {
		"properties": {
			"id": {
				"type": "keyword"
			},
			"data": {
				"type": "text"
			}
		}
	}
}`

type logRecord struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

func (l *logRecord) marshal() []byte {
	b, _ := json.Marshal(l)
	return b
}

func newLog(id, data string) *logRecord {
	return &logRecord{Id: id,
		Data: data}
}
func main() {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}

	existRequest := types.NewRequest().
		SetMetadataKeyValue("method", "index.exists").
		SetMetadataKeyValue("index", "log")

	queryExistResponse, err := client.SetQuery(existRequest.ToQuery()).
		SetChannel("query.elastic").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	existResponse, err := types.ParseResponse(queryExistResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("check if log index exist executed, response: %s", existResponse.Metadata.String()))
	exist, err := strconv.ParseBool(existResponse.Metadata["exists"])
	if err != nil {
		log.Fatal(err)
	}

	if exist {
		deleteIndexRequest := types.NewRequest().
			SetMetadataKeyValue("method", "index.delete").
			SetMetadataKeyValue("index", "log")

		queryDeleteIndexResponse, err := client.SetQuery(deleteIndexRequest.ToQuery()).
			SetChannel("query.elastic").
			SetTimeout(10 * time.Second).Send(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		deleteIndexResponse, err := types.ParseResponse(queryDeleteIndexResponse.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(fmt.Sprintf("delete log index executed, response: %s", deleteIndexResponse.Metadata.String()))
	}

	createIndexRequest := types.NewRequest().
		SetMetadataKeyValue("method", "index.create").
		SetMetadataKeyValue("index", "log").
		SetData([]byte(mapping))

	queryCreateIndexResponse, err := client.SetQuery(createIndexRequest.ToQuery()).
		SetChannel("query.elastic").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	createIndexResponse, err := types.ParseResponse(queryCreateIndexResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("create log index executed, response: %s", createIndexResponse.Metadata.String()))

	randomKey := uuid.New().String()
	// set request
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "set").
		SetMetadataKeyValue("id", randomKey).
		SetMetadataKeyValue("index", "log").
		SetData(newLog("some-id", "some-data").marshal())
	querySetResponse, err := client.SetQuery(setRequest.ToQuery()).
		SetChannel("query.elastic").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	setResponse, err := types.ParseResponse(querySetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("set request for key: %s executed, response: %s", randomKey, setResponse.Metadata.String()))

	// get request
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get").
		SetMetadataKeyValue("index", "log").
		SetMetadataKeyValue("id", randomKey)

	queryGetResponse, err := client.SetQuery(getRequest.ToQuery()).
		SetChannel("query.elastic").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getResponse, err := types.ParseResponse(queryGetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get request for key: %s executed, response: %s, data: %s", randomKey, getResponse.Metadata.String(), string(getResponse.Data)))

	// delete request

	delRequest := types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("index", "log").
		SetMetadataKeyValue("id", randomKey)

	queryDelResponse, err := client.SetQuery(delRequest.ToQuery()).
		SetChannel("query.elastic").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	delResponse, err := types.ParseResponse(queryDelResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("delete request for key: %s executed, response: %s", randomKey, delResponse.Metadata.String()))
}
