package api

type Profile struct {
	Resources []ResourceType `json:"Resources"`
}

func ProfileApplication() Profile {
	return Profile{
		Resources: []ResourceType{
			ResourceMam,
			ResourceHealth,
			ResourceLicense,
			ResourceLicenseText,
			ResourcePublicationList,
			ResourceProfile,
		},
	}
}

func ProfileDevice() Profile {
	return Profile{
		Resources: []ResourceType{
			ResourceMam,
			ResourceHealth,
			ResourcePublicationList,
			ResourceReferenceDesignation,
			ResourceProfile,
		},
	}
}
