package types

type FieldFlag uint16

const (
	FieldFlagPromotedField FieldFlag = 0x1
)

type FieldMetaData struct {
	// Value range: <Unique String representing the name>
	// Type: String
	// Requirement: Mandatory
	//
	// Name of the field. The Name shall be unique in the DataSet.
	Name string `json:"Name"`

	// Value range: <LocalizedText>
	// Type: LocalizedText (see 9.2.7)
	// Requirement: Mandatory
	//
	// Description of the field. The default value shall be a null or empty LocalizedText.
	Description LocalizedText `json:"Description"`

	// Value range: <DataSetFieldFlags>
	// Type: Subtype of UInt16
	// Requirement: Mandatory
	//
	// Flags for the field; see definition in OPC UA Part 14-6.2.3.2.4.
	FieldFlags uint16 `json:"FieldFlags"`

	// Value range: <Byte>
	// Type: Byte
	// Requirement: Mandatory
	//
	// builtInType values are defined in OPC UA Part 6-5.1.2.
	//
	// note The JSON representation of each BuildInType is defined in OPC UA PART 6-5.4.2.
	BuiltInType byte `json:"BuiltInType"`

	// Value range: <NodeId>
	// Type: NodeId object
	// Requirement: Mandatory
	//
	// JSON representation of NodeId is defined in OPC UA Part 6-5.4.2.10.
	//
	// note First pitfall is to wonder about missing Namespace in DataType object. If Namespace is equal to 0, it is not present in most implementations.
	DataType NodeId `json:"DataType"`

	// Value range: <INT32>
	// Type: Int32
	// Requirement: Mandatory
	//
	// Defines if DataType is an array and how many dimensions it has. See OPC UA Part 14-6.2.3.2.3-Table 7 for details.
	ValueRank int32 `json:"ValueRank"`

	// Value range: <array of UINT32>
	// Type: UInt32
	// Requirement: Mandatory
	//
	// This field specifies the maximum length of each dimension. See OPC UA Part 14-6.2.3.2.3-Table 7 for details.
	ArrayDimensions []uint32 `json:"ArrayDimensions"`

	// Value range: <UINT32>
	// Type: UInt32
	// Requirement: Mandatory
	//
	// If the DataType field is a String or ByteString, this field specifies the maximum length of the String or array.
	MaxStringLength uint32 `json:"MaxStringLength"`

	// Value range: <GUID>
	// Type: Guid
	// Requirement: Mandatory
	//
	// The unique ID for the field in the DataSet. GUID-Definition: OPC UA Part 6-5.1.3: 16 Byte as JSON string with separator (Part 6-5.4.2.7)
	DataSetFieldId string `json:"DataSetField"`

	// Value range: <array of KeyValuePair>
	// Type: KeyValuePair
	// Requirement: Mandatory
	//
	// List of Property values providing additional semantics for the field.
	Properties []KeyValuePair `json:"Properties"`
}
