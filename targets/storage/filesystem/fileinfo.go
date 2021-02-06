package filesystem

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"is_dir"`
}

func newFromOSFileInfo(f os.FileInfo, path string) *FileInfo {
	fi := &FileInfo{
		Name:     f.Name(),
		FullPath: "",
		Size:     f.Size(),
		IsDir:    f.IsDir(),
	}
	fi.FullPath, _ = filepath.Abs(path)
	return fi
}

type FileInfoList []*FileInfo

func (l FileInfoList) Marshal() []byte {
	data, _ := json.Marshal(l)
	return data
}
