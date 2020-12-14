package aerospike

type PutRequest struct {
	BinMap    map[string]interface{} `json:"bin_map"`
	KeyName   string                 `json:"key_name"`
	Namespace string                 `json:"namespace"`
	UserKey   interface{}            `json:"user_key"`
}

type GetBatchRequest struct {
	KeyNames  []*string `json:"key_names"`
	BinNames  []string  `json:"bin_names"`
	Namespace string    `json:"namespace"`
}
