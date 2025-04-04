package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"testing"
)

func TestUnmarshalMessage(t *testing.T) {
	content, err := os.ReadFile("./testdata/networkMessageWithFilter.json")
	require.NoError(t, err)

	var nm NetworkMessage
	err = json.Unmarshal(content, &nm)
	require.NoError(t, err)

	nmJson, err := json.Marshal(nm)
	require.NoError(t, err)

	var nmWrapped NetworkMessage
	err = json.Unmarshal(nmJson, &nmWrapped)
	require.NoError(t, err)

	assert.True(t, !reflect.DeepEqual(nm, nmWrapped))
}
