package v1

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
