package solr

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_CommitWithin(t *testing.T) {

	urlValues := map[string][]string{}
	CommitWithin(650 * time.Microsecond)(urlValues)
	require.Equal(t, map[string][]string{"commitWithin": {"1"}}, urlValues)

	urlValues = map[string][]string{}
	CommitWithin(5 * time.Millisecond)(urlValues)
	require.Equal(t, map[string][]string{"commitWithin": {"5"}}, urlValues)

}
