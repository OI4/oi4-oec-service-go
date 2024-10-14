package api

type StructureField struct {

	// Value range: <unique name for field in StructureDefinition>
	// Type: String
	// Requirement: Mandatory
	Name string `json:"Name"`

	// Value range: <LocalizedText>
	// Type: LocalizedText (see 9.2.7)
	// Requirement: Mandatory
	Description LocalizedText `json:"Description"`

	// Value range: <NodeId>
	// Type: NodeId object
	// Requirement: Mandatory
	//
	// JSON representation of NodeId is defined in OPC UA Part 6-5.4.2.10-Table 23.
	Datatype NodeId `json:"Datatype"`

	// Value range: <INT32>
	// Type: Int32
	// Requirement: Mandatory
	//
	// Scalar (-1) or fixed rank Array (>=1)
	ValueRank int32 `json:"ValueRank"`

	// Value range: <array of UINT32>
	// Type: UInt32
	// Requirement: Mandatory
	ArrayDimensions []uint32 `json:"ArrayDimensions"`

	// Value range: <UINT32>
	// Type: UInt32
	// Requirement: Mandatory
	MaxStringLength uint32 `json:"MaxStringLength"`

	// Value range: <BOOLEAN>
	// Type: Boolean
	// Requirement: Mandatory
	IsOptional bool `json:"IsOptional"`
}
