package solr

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ParseAdds(t *testing.T) {
	ids, err := UpdateResponse{Adds: []interface{}{"id1", json.Number("123912368721"), "id2", json.Number("1238717123")}}.AddedIDs()
	require.Nil(t, err)
	require.Equal(t, map[string]int64{"id1": 123912368721, "id2": 1238717123}, ids)
}
