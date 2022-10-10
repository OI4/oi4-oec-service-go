package types

type NodeIdType int

const (
	NodeIdUint32 NodeIdType = iota
	NodeIdString
	NodeIdGUID
	NodeIdByteString
)

type NodeId struct {
	IdType    NodeIdType  `json:"IdType"`
	Id        interface{} `json:"Id"`
	Namespace int32       `json:"Namespace,omitempty"`
}

type BuiltInDataType byte

const (
	BuiltInType_Boolean         BuiltInDataType = 1
	BuiltInType_SByte           BuiltInDataType = 2
	BuiltInType_Byte            BuiltInDataType = 3
	BuiltInType_Int16           BuiltInDataType = 4
	BuiltInType_UInt16          BuiltInDataType = 5
	BuiltInType_Int32           BuiltInDataType = 6
	BuiltInType_UInt32          BuiltInDataType = 7
	BuiltInType_Int64           BuiltInDataType = 8
	BuiltInType_UInt64          BuiltInDataType = 9
	BuiltInType_Float           BuiltInDataType = 10
	BuiltInType_Double          BuiltInDataType = 11
	BuiltInType_String          BuiltInDataType = 12
	BuiltInType_DateTime        BuiltInDataType = 13
	BuiltInType_Guid            BuiltInDataType = 14
	BuiltInType_ByteString      BuiltInDataType = 15
	BuiltInType_XmlElement      BuiltInDataType = 16
	BuiltInType_NodeId          BuiltInDataType = 17
	BuiltInType_ExpandedNodeId  BuiltInDataType = 18
	BuiltInType_StatusCode      BuiltInDataType = 19
	BuiltInType_QualifiedName   BuiltInDataType = 20
	BuiltInType_LocalizedText   BuiltInDataType = 21
	BuiltInType_ExtensionObject BuiltInDataType = 22
	BuiltInType_DataValue       BuiltInDataType = 23
	BuiltInType_Variant         BuiltInDataType = 24
	BuiltInType_DiagnosticInfo  BuiltInDataType = 25
)
