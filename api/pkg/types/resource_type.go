package types

import "fmt"

type ResourceType string

const (
	ResourceMam                  ResourceType = "MAM"
	ResourceHealth               ResourceType = "Health"
	ResourceConfig               ResourceType = "Config"
	ResourceLicense              ResourceType = "License"
	ResourceLicenseText          ResourceType = "LicenseText"
	ResourceRtLicense            ResourceType = "RtLicense"
	ResourceData                 ResourceType = "Data"
	ResourceMetadata             ResourceType = "Metadata"
	ResourceEvent                ResourceType = "Event"
	ResourceProfile              ResourceType = "Profile"
	ResourcePublicationList      ResourceType = "PublicationList"
	ResourceSubscriptionList     ResourceType = "SubscriptionList"
	ResourceInterfaces           ResourceType = "Interfaces"
	ResourceReferenceDesignation ResourceType = "ReferenceDesignation"
)

var resourceTypes = map[ResourceType]struct{}{
	ResourceMam:                  {},
	ResourceHealth:               {},
	ResourceConfig:               {},
	ResourceLicense:              {},
	ResourceLicenseText:          {},
	ResourceRtLicense:            {},
	ResourceData:                 {},
	ResourceMetadata:             {},
	ResourceEvent:                {},
	ResourceProfile:              {},
	ResourcePublicationList:      {},
	ResourceSubscriptionList:     {},
	ResourceInterfaces:           {},
	ResourceReferenceDesignation: {},
}

func ParseResourceType(s string) (*ResourceType, error) {
	rType := ResourceType(s)
	_, ok := resourceTypes[rType]
	if !ok {
		return nil, fmt.Errorf(`cannot parse:[%s] as ResourceType`, s)
	}
	return &rType, nil
}

func (r ResourceType) ToDataSetClassId() string {
	return dataSetClassIdMapping[r]
}

type DataSetClassId string

const (
	DataSetClassIdMAM                  DataSetClassId = "360ca8f3-5e66-42a2-8f10-9cdf45f4bf58"
	DataSetClassIdHealth               DataSetClassId = "d8e7b6df-42ba-448a-975a-199f59e8ffeb"
	DataSetClassIdConfig               DataSetClassId = "9d5983db-440d-4474-9fd7-1cd7a6c8b6c2"
	DataSetClassIdLicense              DataSetClassId = "2ae0505e-2830-4980-b65e-0bbdf08e2d45"
	DataSetClassIdLicenseText          DataSetClassId = "a6e6c727-4057-419f-b2ea-3fe9173e71cf"
	DataSetClassIdRtLicense            DataSetClassId = "ebd12d4b-da1c-4671-ab86-db102fecc603"
	DataSetClassIdEvent                DataSetClassId = "543ae05e-b6d9-4161-a0a3-350a0fac5976"
	DataSetClassIdProfile              DataSetClassId = "48017c6a-05c8-48d7-9d85-4b08bbb707f3"
	DataSetClassIdPublicationList      DataSetClassId = "217434d6-6e1e-4230-b907-f52bc9ffe152"
	DataSetClassIdSubscriptionList     DataSetClassId = "e5d68c47-c276-4929-8ab9-4c1090cac785"
	DataSetClassIdInterfaces           DataSetClassId = "96d22d73-bce6-42d3-9949-45e0d04e4d54"
	DataSetClassIdReferenceDesignation DataSetClassId = "27a75019-164a-496d-a38b-90e8a55c2cfa"
	DataSetClassIdFileUpload           DataSetClassId = "3b4a62ba-026f-4ee8-bc99-3a5f85fc9f3b"
	DataSetClassIdFileDownload         DataSetClassId = "760abda2-ba40-4e6e-863a-eea8c002b4e4"
	DataSetClassIdFirmwareUpdate       DataSetClassId = "414e26f6-341b-43b7-90fc-bb9e0b1b0866"
	DataSetClassIdBlink                DataSetClassId = "3b423a40-a676-4ba0-8017-f0b2cd65bc26"
	DataSetClassIdNewDataSetWriterId   DataSetClassId = "2aca55bd-0d6f-41b1-a1c2-2d61afcc21f0"
)

var dataSetClassIdMapping = map[ResourceType]string{
	ResourceMam:                  "360ca8f3-5e66-42a2-8f10-9cdf45f4bf58",
	ResourceHealth:               "d8e7b6df-42ba-448a-975a-199f59e8ffeb",
	ResourceConfig:               "9d5983db-440d-4474-9fd7-1cd7a6c8b6c2",
	ResourceLicense:              "2ae0505e-2830-4980-b65e-0bbdf08e2d45",
	ResourceLicenseText:          "a6e6c727-4057-419f-b2ea-3fe9173e71cf",
	ResourceRtLicense:            "ebd12d4b-da1c-4671-ab86-db102fecc603",
	ResourceEvent:                "543ae05e-b6d9-4161-a0a3-350a0fac5976",
	ResourceProfile:              "48017c6a-05c8-48d7-9d85-4b08bbb707f3",
	ResourcePublicationList:      "217434d6-6e1e-4230-b907-f52bc9ffe152",
	ResourceSubscriptionList:     "e5d68c47-c276-4929-8ab9-4c1090cac785",
	ResourceInterfaces:           "96d22d73-bce6-42d3-9949-45e0d04e4d54",
	ResourceReferenceDesignation: "27a75019-164a-496d-a38b-90e8a55c2cfa",
}
