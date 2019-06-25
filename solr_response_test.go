package solr

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ParseAdds(t *testing.T) {
	ids, err := UpdateResponse{Adds: []interface{}{"id1", 123912368721, "id2", 1238717123}}.AddedIDs()
	require.Nil(t, err)
	require.Equal(t, []string{"id1", "id2"}, ids)
}
