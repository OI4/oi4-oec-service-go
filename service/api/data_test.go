package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseOi4Data_ValidJson(t *testing.T) {
	jsonData := `{"Pv": "primary", "Sv1": "secondary1", "Sv2": "secondary2"}`
	data, err := ParseOi4Data(jsonData)
	require.NoError(t, err)

	assert.Equal(t, "primary", data.PrimaryValue.(string), "expected primary value to be 'primary', got %v", data.PrimaryValue)
	assert.Equal(t, 2, len(data.values), "expected 2 secondary values, got %d", len(data.values))

}

func TestParseOi4Data_InvalidJson(t *testing.T) {
	jsonData := `{"Pv": "primary", "Sv1": "secondary1", "Sv2": "secondary2"`
	_, err := ParseOi4Data(jsonData)
	require.Error(t, err)
}

func TestParseOi4Data_NoPrimaryValue(t *testing.T) {
	jsonData := `{"Sv1": "secondary1", "Sv2": "secondary2"}`
	_, err := ParseOi4Data(jsonData)
	require.Error(t, err)
	assert.Equal(t, "primary value not found: <nil>", err.Error())
}

func TestParseOi4Data_InvalidSecondaryValue(t *testing.T) {
	jsonData := `{"Pv": "primary", "Sv1_abc": "secondary1"}`
	_, err := ParseOi4Data(jsonData)
	require.Error(t, err)
	assert.Equal(t, "secondary value key Sv1_abc is not valid: Tag must be in format Sv[0-9]+: <nil>", err.Error())
}

func TestAddSecondaryData_ValidTag(t *testing.T) {
	data := &Oi4Data{values: make(map[string]any)}
	err := data.AddSecondaryData("Sv1", new(any))
	require.NoError(t, err)
}

func TestAddSecondaryData_InvalidTag(t *testing.T) {
	data := &Oi4Data{values: make(map[string]any)}
	err := data.AddSecondaryData("InvalidTag", new(any))
	require.Error(t, err)

	err = data.AddSecondaryData("Sv", new(any))
	require.Error(t, err)

	err = data.AddSecondaryData("Sv1a", new(any))
	require.Error(t, err)
}

func TestAddSecondaryData_NilValue(t *testing.T) {
	data := &Oi4Data{values: make(map[string]any)}
	err := data.AddSecondaryData("Sv1", new(any))
	require.NoError(t, err)

	err = data.AddSecondaryData("Sv1", nil)
	require.NoError(t, err)

	if _, exists := data.values["Sv1"]; exists {
		t.Errorf("expected Sv1 to be deleted, but it exists")
	}
}

func TestGetData(t *testing.T) {
	data := &Oi4Data{
		PrimaryValue: "primary",
		values:       map[string]any{"Sv1": "secondary1"},
	}
	result := data.GetData()
	if result.(map[string]any)["Pv"] != "primary" {
		t.Errorf("expected primary value to be 'primary', got %v", result.(map[string]any)["Pv"])
	}
}

func TestClear(t *testing.T) {
	data := &Oi4Data{
		PrimaryValue: "primary",
		values:       map[string]any{"Sv1": "secondary1"},
	}
	data.Clear()
	assert.Nil(t, data.PrimaryValue, "expected primary value to be nil, got %v", data.PrimaryValue)
	assert.Equal(t, 0, len(data.values), "expected values to be empty, got %d", len(data.values))
}
