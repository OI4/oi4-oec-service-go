package v1

type DataSetMetaDataType struct {

	// 	Value range: <name of DataSet>
	// Type: String
	// Requirement: Mandatory
	Name string `json:"Name"`

	// 	Value range: <LocalizedText>
	// Type: LocalizedText (see 9.2.7)
	// Requirement: Optional
	Description LocalizedText `json:"Description,omitempty"`

	// 	Value range: <array of FieldMetaData>
	// Type: FieldMetaData object (see 9.2.6)
	// Requirement: Optional
	Fields []FieldMetaData `json:"Fields,omitempty"`

	// 	Value range: <GUID> in accordance to OPC UA Part 6-5.1.3: 16 Byte as JSON string with separator (Part 6-5.4.2.7).
	// Type: String
	// Requirement: Optional
	//
	// Example: "f1875b4a-3209-431b-a38d-2df5758f92c8"
	// The DataSetClassId allows to refer to the DataSetClass describing the structure of the message. The DataSetClassId identifies a well defined DataSet specified by the Alliance or some other standards.
	// For some recurring use cases, such as MasterAssetModel, fixed GUIDs are specified from the Alliance and must be used (A2).
	//
	// NOTE The DataSetClassId must be present for all resources defined by the Open Industry 4.0 Alliance, which are having a DataSetClassId. This guarantees the availability of a fully schema verifiable messaging.
	DataSetClassId string `json:"DataSetClassId,omitempty"`

	// 	Value range: <ConfigurationVersionDataType>
	// Type: ConfigurationVersionDataType object (see 9.2.5)
	// Requirement: Optional
	ConfigurationVersion ConfigurationVersionDataType `json:"ConfigurationVersion,omitempty"`

	// 	Value range: <array of Namespaces names>
	// Type: String
	// Requirement: Mandatory
	//
	// For details see OPC UA specification Annex A: Part 14-A1.1
	Namespaces []string `json:"Namespaces"`

	StructureDataTypes []StructureDescription `json:"StructureDataTypes"`
}
