package bigtable

import (
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"testing"
)


func TestParseMetaData(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)

	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
		Request *types.Request
	}{
		{
			name: "valid method write",
			cfg: config.Spec{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			Request: types.NewRequest().
				SetMetadataKeyValue("method", "write").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetMetadataKeyValue("column_family", dat.columnFamily),
			wantErr: false,
		},
		{
			name: "invalid method write - missing column_family",
			cfg: config.Spec{
				Name: "target-gcp-bigtable",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			Request: types.NewRequest().
				SetMetadataKeyValue("method", "write").
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: true,
		},
		{
			name: "valid method write_batch",
			cfg: config.Spec{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			Request: types.NewRequest().
				SetMetadataKeyValue("method", "write_batch").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetMetadataKeyValue("column_family", dat.columnFamily),
			wantErr: false,
		},
		{
			name: "invalid method write_batch - missing column_family",
			cfg: config.Spec{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			Request: types.NewRequest().
				SetMetadataKeyValue("method", "write_batch").
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: true,
		},
		{
			name: "valid method delete_rows",
			cfg: config.Spec{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			Request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_row").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetMetadataKeyValue("row_key_prefix", dat.rowKeyPrefix),
			wantErr: false,
		},
		{
			name: "invalid method delete_row - missing row_key_prefix",
			cfg: config.Spec{
				Name: "target-gcp-bigtable",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			Request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_row").
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: true,
		},
		{
			name: "valid method create_table",
			cfg: config.Spec{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			Request: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetMetadataKeyValue("table_name", dat.tempTable),
			wantErr: false,
		},
		{
			name: "valid method delete_table",
			cfg: config.Spec{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			Request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_table").
				SetMetadataKeyValue("table_name", dat.tempTable),
			wantErr: false,
		},
		{
			name: "valid method get_tables",
			cfg: config.Spec{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			Request: types.NewRequest().
				SetMetadataKeyValue("method", "get_tables"),
			wantErr: false,
		}, {
			name: "invalid method type",
			cfg: config.Spec{
				Name: "target-gcp-bigtable",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			Request: types.NewRequest().
				SetMetadataKeyValue("method", "non_existing_type").
				SetMetadataKeyValue("table_name", dat.tempTable),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := parseMetadata(tt.Request.Metadata)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, m)
		})
	}
}
