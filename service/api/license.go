package api

type LicenseComponent struct {
	Component  string   `json:"Component"`
	LicAuthors []string `json:"LicAuthors,omitempty"`
	LicAddText string   `json:"LicAddText,omitempty"`
}

type License struct {
	Components []LicenseComponent `json:"Components"`
}

type LicenseText struct {
	LicenseText string `json:"LicenseText"`
}

func EmptyLicense() License {
	return License{
		Components: make([]LicenseComponent, 0),
	}
}
