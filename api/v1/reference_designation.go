package v1

type ReferenceDesignationParent struct {
	Value         string            `json:"Value"`
	Local         string            `json:"Local,omitempty"`
	Oi4Identifier Oi4IdentifierPath `json:"Oi4Identifier,omitempty"`
}

type ReferenceDesignationFunctionProductOrLocation struct {
	Value  string                      `json:"Value"`
	Local  string                      `json:"Local,omitempty"`
	Parent *ReferenceDesignationParent `json:"Parent,omitempty"`
}

type ReferenceDesignationFunction struct {
	Value    string                                         `json:"Value"`
	Local    string                                         `json:"Local,omitempty"`
	Parent   *ReferenceDesignationParent                    `json:"Parent,omitempty"`
	Product  *ReferenceDesignationFunctionProductOrLocation `json:"Product,omitempty"`
	Location *ReferenceDesignationFunctionProductOrLocation `json:"Location,omitempty"`
}

type ReferenceDesignation struct {
	Function *ReferenceDesignationFunction `json:"Function"`
}
