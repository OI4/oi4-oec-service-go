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

type BasicType byte

const (
	Type_Boolean         BasicType = 1
	Type_SByte           BasicType = 2
	Type_Byte            BasicType = 3
	Type_Int16           BasicType = 4
	Type_UInt16          BasicType = 5
	Type_Int32           BasicType = 6
	Type_UInt32          BasicType = 7
	Type_Int64           BasicType = 8
	Type_UInt64          BasicType = 9
	Type_Float           BasicType = 10
	Type_Double          BasicType = 11
	Type_String          BasicType = 12
	Type_DateTime        BasicType = 13
	Type_Guid            BasicType = 14
	Type_ByteString      BasicType = 15
	Type_XmlElement      BasicType = 16
	Type_NodeId          BasicType = 17
	Type_ExpandedNodeId  BasicType = 18
	Type_StatusCode      BasicType = 19
	Type_QualifiedName   BasicType = 20
	Type_LocalizedText   BasicType = 21
	Type_ExtensionObject BasicType = 22
	Type_DataValue       BasicType = 23
	Type_Variant         BasicType = 24
	Type_DiagnosticInfo  BasicType = 25
)
