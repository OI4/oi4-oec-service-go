package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"testing"
)

func TestParseMessage(t *testing.T) {
	content, err := os.ReadFile("./testdata/dataSetMessage.json")
	require.NoError(t, err)

	var dsm DataSetMessage
	err = json.Unmarshal(content, &dsm)
	require.NoError(t, err)

	dsmJson, err := json.Marshal(dsm)
	require.NoError(t, err)
	var dsmWrapped DataSetMessage
	err = json.Unmarshal(dsmJson, &dsmWrapped)
	require.NoError(t, err)

	assert.True(t, reflect.DeepEqual(dsm, dsmWrapped))
}
