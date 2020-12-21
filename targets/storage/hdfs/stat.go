package hdfs

import (
	"encoding/json"
	"os"
	"time"
)

type Stat struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	IsDir   bool `json:"is_dir"`
}



func createStatAsByteArray(o os.FileInfo) ([]byte, error) {
	s := Stat{
		Name:    o.Name(),
		Size:    o.Size(),
		ModTime: o.ModTime(),
		IsDir:   o.IsDir(),
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return b, err

}
