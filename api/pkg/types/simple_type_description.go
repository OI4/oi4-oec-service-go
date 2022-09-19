package types

type SimpleTypeDescription struct {
	BaseDataType NodeId `json:"BaseDataType"`
	BuiltInType  `json:"BuiltInType"`
	DataTypeId   NodeId `json:"DataTypeId"`
	Name         string `json:"Name"`
}