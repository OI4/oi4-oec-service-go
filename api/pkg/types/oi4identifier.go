package types

import (
	"fmt"

	"github.com/mzeiher/oi4/api/pkg/utils"
)

type Oi4IdentifierPath string

type Oi4Identifier struct {
	ManufacturerUri string `json:"ManufacturerUri"`
	Model           string `json:"Mode"`
	ProductCode     string `json:"ProductCode"`
	SerialNumber    string `json:"SerialNumber"`
}

func (ident *Oi4Identifier) ToString() string {
	return fmt.Sprintf("%s/%s/%s/%s", ident.ManufacturerUri, utils.DNPEncode(ident.Model), utils.DNPEncode(ident.ProductCode), utils.DNPEncode(ident.SerialNumber))
}
