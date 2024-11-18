package api

type ReferenceDesignationParent struct {
	Value         string        `json:"Value"`
	Local         string        `json:"Local,omitempty"`
	Oi4Identifier Oi4Identifier `json:"Oi4Identifier,omitempty"`
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
	Function *ReferenceDesignationFunction `json:"Function,omitempty"`
	Product  *ReferenceDesignationFunction `json:"Product,omitempty"`
	Location *ReferenceDesignationFunction `json:"Location,omitempty"`
}
