package v1

type LicenseComponent struct {
	Component      string   `json:"Component"`
	LicAuthors     []string `json:"LicAuthors,omitempty"`
	LicAddText     string   `json:"LicAddText,omitempty"`
	SPDXIdentifier string   `json:"SPDXIdentifier,omitempty"`
}

type License struct {
	Components []interface{} `json:"Components"`
}

type LicenseText struct {
	LicenseText string `json:"LicenseText"`
}
