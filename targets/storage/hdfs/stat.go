package s3

import (
	"encoding/json"
	"os"
	"time"
)

type Stat struct {
	Name    string    `json:"name"`
	Size    int64     `json:"name"`
	ModTime time.Time `json:"name"`
	IsDir   bool      `json:"name"`
}

func createStat(o os.FileInfo) Stat {
	return Stat{
		Name:    o.Name(),
		Size:    o.Size(),
		ModTime: o.ModTime(),
		IsDir:   o.IsDir(),
	}
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
