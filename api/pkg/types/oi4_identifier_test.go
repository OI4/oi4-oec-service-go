package types

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOi4Identifier(t *testing.T) {
	const (
		manufacturerUri = "manufacturerUri"
		model           = "Model"
		productCode     = "product!Code"
		serialNumber    = "FBC#123"
	)

	identifier := NewOi4Identifier(manufacturerUri, model, productCode, serialNumber)

	assert.Equal(t, manufacturerUri, identifier.ManufacturerUri)
	assert.Equal(t, model, identifier.Model)
	assert.Equal(t, productCode, identifier.ProductCode)
	assert.Equal(t, serialNumber, identifier.SerialNumber)
}

func TestParseOi4Identifier_HappyPath(t *testing.T) {
	const (
		manufacturerUri = "manufacturerUri"
		model           = "Model"
		productCode     = "product,2ACode"
		serialNumber    = "FBC,23123"
	)
	const (
		expectedProductCode  = "product*Code"
		expectedSerialNumber = "FBC#123"
	)

	identifier, err := ParseOi4Identifier(fmt.Sprintf("%s/%s/%s/%s", manufacturerUri, model, productCode, serialNumber), false)

	assert.NoError(t, err)
	assert.Equal(t, manufacturerUri, identifier.ManufacturerUri)
	assert.Equal(t, model, identifier.Model)
	assert.Equal(t, productCode, identifier.ProductCode)
	assert.Equal(t, serialNumber, identifier.SerialNumber)

	identifier, err = ParseOi4Identifier(fmt.Sprintf("%s/%s/%s/%s", manufacturerUri, model, productCode, serialNumber), true)

	assert.NoError(t, err)
	assert.Equal(t, manufacturerUri, identifier.ManufacturerUri)
	assert.Equal(t, model, identifier.Model)
	assert.Equal(t, expectedProductCode, identifier.ProductCode)
	assert.Equal(t, expectedSerialNumber, identifier.SerialNumber)
}

func TestParseOi4Identifier_InvalidIdentifier(t *testing.T) {
	_, err := ParseOi4Identifier("invalidIdentifier", false)

	assert.Error(t, err)
}

func TestParseOi4Identifier_InvalidParts(t *testing.T) {
	_, err := ParseOi4Identifier("1/,,/3/4,", true)

	assert.Error(t, err)
}

func TestOi4Identifier_ToString(t *testing.T) {
	identifier := NewOi4Identifier("manufacturerUri", "model", "product*Code", "FBC#123")

	assert.Equal(t, "manufactureruri/model/product,2ACode/FBC,23123", identifier.ToString())
}

func TestToPlainString(t *testing.T) {
	identifier := NewOi4Identifier("manufacturerUri", "model", "productCode", "serialNumber")

	assert.Equal(t, "manufacturerUri/model/productCode/serialNumber", identifier.ToPlainString())
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
}

func TestParseOi4IdentifierFromArray_EmptyArray(t *testing.T) {
	_, err := ParseOi4IdentifierFromArray([]string{}, false)

	assert.Error(t, err)
}

func TestGetPartFails(t *testing.T) {
	_, err := getPart([]string{}, 0, false)

	assert.Error(t, err)
}

func TestParseOi4Identifier_InvalidProductCode(t *testing.T) {
	_, err := ParseOi4Identifier("a/b/,,c/d", true)

	assert.Error(t, err)
	assert.Equal(t, "invalid product code: invalid DNP escape \",,c\"", err.Error())
}

func TestParseOi4Identifier_InvalidSerialNumber(t *testing.T) {
	_, err := ParseOi4Identifier("a/b/c/,,d", true)

	assert.Error(t, err)
	assert.Equal(t, "invalid serial number: invalid DNP escape \",,d\"", err.Error())
}

func TestEqualsWithIdenticalIdentifiers(t *testing.T) {
	identifier1 := NewOi4Identifier("manufacturerUri", "Model", "product!Code", "FBC#123")
	identifier2 := NewOi4Identifier("manufacturerUri", "Model", "product!Code", "FBC#123")

	assert.True(t, identifier1.Equals(identifier2))
}

func TestEqualsWithDifferentIdentifiers(t *testing.T) {
	identifier1 := NewOi4Identifier("manufacturerUri", "Model", "product!Code", "FBC#123")
	identifier2 := NewOi4Identifier("manufacturerUri", "Model", "product!Code", "FBC#124")

	assert.False(t, identifier1.Equals(identifier2))
}
