package types

import (
	"github.com/OI4/oi4-oec-service-go/dnp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOi4Identifier(t *testing.T) {
	identifier := NewOi4Identifier("manufacturerUri", "model", "productCode", "serialNumber")

	assert.Equal(t, "manufacturerUri", identifier.ManufacturerUri)
	assert.Equal(t, "model", identifier.Model)
	assert.Equal(t, "productCode", identifier.ProductCode)
	assert.Equal(t, "serialNumber", identifier.SerialNumber)
}

func TestNewOi4IdentifierFromString_HappyPath(t *testing.T) {
	identifier, err := ParseOi4Identifier("manufacturerUri/model/productCode/serialNumber", false)

	assert.NoError(t, err)
	assert.Equal(t, "manufacturerUri", identifier.ManufacturerUri)
	assert.Equal(t, "model", identifier.Model)
	assert.Equal(t, "productCode", identifier.ProductCode)
	assert.Equal(t, "serialNumber", identifier.SerialNumber)
}

func TestParseOi4Identifier_InvalidIdentifier(t *testing.T) {
	_, err := ParseOi4Identifier("invalidIdentifier", false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid identifier")

	_, err = ParseOi4Identifier("1/,,/3/4,", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid model")

	_, err = ParseOi4Identifier("1/2/,,/4,", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid product code")

	_, err = ParseOi4Identifier("1/2/3/,,,", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid serial number")
}

func TestOi4Identifier_ToString(t *testing.T) {
	identifier := NewOi4Identifier("manufacturerUri", "model", "productCode", "serialNumber")

	assert.Equal(t, "manufacturerUri/model/productCode/serialNumber", identifier.ToString())
}

func TestOi4Identifier_ToDnpString(t *testing.T) {
	identifier := NewOi4Identifier("manufacturerUri", "model", "productCode", "serialNumber")

	expected := "manufacturerUri/" + dnp.Encode("model") + "/" + dnp.Encode("productCode") + "/" + dnp.Encode("serialNumber")
	assert.Equal(t, expected, identifier.ToDnpString())
}

func TestParseOi4IdentifierFromArray_HappyPath(t *testing.T) {
	identifier, err := ParseOi4IdentifierFromArray([]string{"manufacturerUri", "model", "productCode", "serialNumber"}, false)

	assert.NoError(t, err)
	assert.Equal(t, "manufacturerUri", identifier.ManufacturerUri)
	assert.Equal(t, "model", identifier.Model)
	assert.Equal(t, "productCode", identifier.ProductCode)
	assert.Equal(t, "serialNumber", identifier.SerialNumber)
}

func TestParseOi4IdentifierFromArray_InvalidIdentifier(t *testing.T) {
	_, err := ParseOi4IdentifierFromArray([]string{"invalidIdentifier"}, false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid identifier")
}

func TestParseOi4IdentifierFromArray_EmptyArray(t *testing.T) {
	_, err := ParseOi4IdentifierFromArray([]string{}, false)

	assert.Error(t, err)
}

func TestParseOi4IdentifierFromArray_DecodeTrue(t *testing.T) {
	identifier, err := ParseOi4IdentifierFromArray([]string{"manufacturerUri", dnp.Encode("model"), dnp.Encode("productCode"), dnp.Encode("serialNumber")}, true)

	assert.NoError(t, err)
	assert.Equal(t, "manufacturerUri", identifier.ManufacturerUri)
	assert.Equal(t, "model", identifier.Model)
	assert.Equal(t, "productCode", identifier.ProductCode)
	assert.Equal(t, "serialNumber", identifier.SerialNumber)
}

func TestGetPartFails(t *testing.T) {
	_, err := getPart([]string{}, 0, false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid index: 0")
}
