package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/pkg/uuid"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"io/ioutil"
	"log"
	"math/big"
	"time"
)

func main() {

	dat, err := ioutil.ReadFile("./credentials/tableName.txt")
	if err != nil {
		log.Fatal(err)
	}
	tableName := string(dat)
	dat, err = ioutil.ReadFile("./credentials/columnFamily.txt")
	if err != nil {
		log.Fatal(err)
	}
	columnFamily := string(dat)
	if err != nil {
		log.Fatal(err)
	}
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// write_row
	singleRow := map[string]interface{}{"set_row_key": fmt.Sprintf("%d", getRandInt64()), "id": 3, "name": "test1"}
	singleB, err := json.Marshal(singleRow)
	if err != nil {
		log.Fatal(err)
	}
	writeRequest := types.NewRequest().
		SetMetadataKeyValue("method", "write").
		SetMetadataKeyValue("table_name", tableName).
		SetMetadataKeyValue("column_family", columnFamily).
		SetData(singleB)
	queryWriteResponse, err := client.SetQuery(writeRequest.ToQuery()).
		SetChannel("query.gcp.bigtable").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	writeResponse, err := types.ParseResponse(queryWriteResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("set row requests for tableName: %s executed, response: %s", tableName, writeResponse.Data))
	// get rows by column

	getRequest := types.NewRequest().
		SetMetadataKeyValue("column_name", "id").
		SetMetadataKeyValue("table_name", tableName).
		SetMetadataKeyValue("method", "get_all_rows")

	getRows, err := client.SetQuery(getRequest.ToQuery()).
		SetChannel("query.gcp.bigtable").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getRowsResponse, err := types.ParseResponse(getRows.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get rows executed, response: %s", getRowsResponse.Data))
}
func getRandInt64() int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(27))
	if err != nil {
		panic(err)
	}
	return nBig.Int64()
}
