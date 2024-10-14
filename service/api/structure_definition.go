package api

type StructureType string

const (
	Structure_0                   StructureType = "Structure_0"
	StructureWithOptionalFields_1 StructureType = "StructureWithOptionalFields_1"
	Union_2                       StructureType = "Union_2"
)

type StructureDefinition struct {
	DefaultEncodingId NodeId `json:"DefaultEncodingId"`
	BaseDataType      NodeId `json:"BaseDataType"`
	StructureType     `json:"StructureType"`
	Fields            []StructureField `json:"Fields"`
}
