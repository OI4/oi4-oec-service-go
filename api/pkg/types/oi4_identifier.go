package types

import (
	"fmt"
	"github.com/OI4/oi4-oec-service-go/dnp"
	"strings"
)

type Oi4Identifier struct {
	ManufacturerUri string `json:"ManufacturerUri"`
	Model           string `json:"Mode"`
	ProductCode     string `json:"ProductCode"`
	SerialNumber    string `json:"SerialNumber"`
}

func NewOi4Identifier(manufacturerUri, model, productCode, serialNumber string) *Oi4Identifier {
	return &Oi4Identifier{
		ManufacturerUri: manufacturerUri,
		Model:           model,
		ProductCode:     productCode,
		SerialNumber:    serialNumber,
	}
}

func ParseOi4Identifier(id string, decode bool) (*Oi4Identifier, error) {
	parts := strings.Split(id, "/")

	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid identifier: %s", id)
	}

	manufacturerUri, _ := getPart(parts, 0, false)

	model, err := getPart(parts, 1, decode)
	if err != nil {
		return nil, &Error{
			Message: "invalid model",
			Err:     err,
		}
	}

	productCode, err := getPart(parts, 2, decode)
	if err != nil {
		return nil, &Error{
			Message: "invalid product code",
			Err:     err,
		}
	}

	serialNumber, err := getPart(parts, 3, decode)
	if err != nil {
		return nil, &Error{
			Message: "invalid serial number",
			Err:     err,
		}
	}

	return NewOi4Identifier(manufacturerUri, model, productCode, serialNumber), nil
}

func ParseOi4IdentifierFromArray(parts []string, decode bool) (*Oi4Identifier, error) {
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid identifier: %s", parts)
	}
	return ParseOi4Identifier(strings.Join(parts, "/"), decode)
}

func getPart(parts []string, index int, decode bool) (string, error) {
	if index >= len(parts) {
		return "", fmt.Errorf("invalid index: %d", index)
	}
	part := parts[index]
	var err error
	if decode {
		part, err = dnp.Decode(part)

	}
	return part, err
}

func (ident *Oi4Identifier) ToPlainString() string {
	return fmt.Sprintf("%s/%s/%s/%s", ident.ManufacturerUri, ident.Model, ident.ProductCode, ident.SerialNumber)
}

func (ident *Oi4Identifier) ToString() string {
	return fmt.Sprintf("%s/%s/%s/%s", strings.ToLower(ident.ManufacturerUri), dnp.Encode(ident.Model), dnp.Encode(ident.ProductCode), dnp.Encode(ident.SerialNumber))
}

func (ident *Oi4Identifier) Equals(other *Oi4Identifier) bool {
	return ident.ManufacturerUri == other.ManufacturerUri && ident.Model == other.Model && ident.ProductCode == other.ProductCode && ident.SerialNumber == other.SerialNumber
}
