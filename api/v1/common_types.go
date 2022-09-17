package v1

type LocalizedText map[string]string
type KeyValuePair struct {
	Key   string      `json:"Key"`
	Value interface{} `json:"Value"`
}
type BuiltInType byte

const (
	BuiltIn_Enumeration     BuiltInType = 0x6
	BuiltIn_ExtensionObject BuiltInType = 0x22
	BuiltIn_Uinteger        BuiltInType = 0x23
)
