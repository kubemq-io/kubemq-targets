package types

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func loadFile(filename string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return data
}

func getRequest(data []byte) *Request {
	return NewRequest().
		SetMetadata(NewMetadata().Set("key", "value")).
		SetData(data)
}

func getTransportRequest(data interface{}) []byte {
	return NewTransportRequest().
		SetMetadata(NewMetadata().Set("key", "value")).
		SetData(data).MarshalBinary()
}

func loadJson() []byte {
	type t struct {
		Param1 string
		Param2 int64
		Param3 bool
		Param4 struct {
			Param1 string
			Param2 int64
			Param3 bool
		}
	}
	test := &t{
		Param1: "test",
		Param2: 1,
		Param3: true,
		Param4: struct {
			Param1 string
			Param2 int64
			Param3 bool
		}{Param1: "test", Param2: 2, Param3: false},
	}
	data, _ := json.Marshal(test)
	return data
}

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name    string
		body    []byte
		want    *Request
		wantErr bool
	}{
		{
			name:    "json",
			body:    getTransportRequest(loadJson()),
			want:    getRequest(loadJson()),
			wantErr: false,
		},
		{
			name:    "string-json",
			body:    getTransportRequest(string(loadJson())),
			want:    getRequest(loadJson()),
			wantErr: false,
		},
		{
			name:    "string",
			body:    getTransportRequest("test-string"),
			want:    getRequest([]byte("test-string")),
			wantErr: false,
		},
		{
			name:    "yaml-file",
			body:    getTransportRequest(loadFile("./testdata/yaml-file.yaml")),
			want:    getRequest(loadFile("./testdata/yaml-file.yaml")),
			wantErr: false,
		},
		{
			name:    "svg-file",
			body:    getTransportRequest(loadFile("./testdata/svg-file.svg")),
			want:    getRequest(loadFile("./testdata/svg-file.svg")),
			wantErr: false,
		},
		{
			name:    "png-file",
			body:    getTransportRequest(loadFile("./testdata/png-file.png")),
			want:    getRequest(loadFile("./testdata/png-file.png")),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRequest(tt.body)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.want, got)
		})
	}
}
