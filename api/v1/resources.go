package v1

type Resource string

const (
	Resource_MAM                  Resource = "MAM"
	Resource_Health               Resource = "Health"
	Resource_Config               Resource = "Config"
	Resource_License              Resource = "License"
	Resource_LicenseText          Resource = "LicenseText"
	Resource_RtLicense            Resource = "RtLicense"
	Resource_Data                 Resource = "Data"
	Resource_Metadata             Resource = "Metadata"
	Resource_Event                Resource = "Event"
	Resource_Profile              Resource = "Profile"
	Resource_PublicationList      Resource = "PublicationList"
	Resource_SubscriptionList     Resource = "SubscriptionList"
	Resource_Interfaces           Resource = "Interfaces"
	Resource_ReferenceDesignation Resource = "ReferenceDesignation"
)

func (r Resource) ToDataSetClassId() string {
	return dataSetClassIdMapping[r]
}

type DataSetClassId string

const (
	DataSetClassId_MAM                  DataSetClassId = "360ca8f3-5e66-42a2-8f10-9cdf45f4bf58"
	DataSetClassId_Health               DataSetClassId = "d8e7b6df-42ba-448a-975a-199f59e8ffeb"
	DataSetClassId_Config               DataSetClassId = "9d5983db-440d-4474-9fd7-1cd7a6c8b6c2"
	DataSetClassId_License              DataSetClassId = "2ae0505e-2830-4980-b65e-0bbdf08e2d45"
	DataSetClassId_LicenseText          DataSetClassId = "a6e6c727-4057-419f-b2ea-3fe9173e71cf"
	DataSetClassId_RtLicense            DataSetClassId = "ebd12d4b-da1c-4671-ab86-db102fecc603"
	DataSetClassId_Event                DataSetClassId = "543ae05e-b6d9-4161-a0a3-350a0fac5976"
	DataSetClassId_Profile              DataSetClassId = "48017c6a-05c8-48d7-9d85-4b08bbb707f3"
	DataSetClassId_PublicationList      DataSetClassId = "217434d6-6e1e-4230-b907-f52bc9ffe152"
	DataSetClassId_SubscriptionList     DataSetClassId = "e5d68c47-c276-4929-8ab9-4c1090cac785"
	DataSetClassId_Interfaces           DataSetClassId = "96d22d73-bce6-42d3-9949-45e0d04e4d54"
	DataSetClassId_ReferenceDesignation DataSetClassId = "27a75019-164a-496d-a38b-90e8a55c2cfa"
	DataSetClassId_FileUpload           DataSetClassId = "3b4a62ba-026f-4ee8-bc99-3a5f85fc9f3b"
	DataSetClassId_FileDownload         DataSetClassId = "760abda2-ba40-4e6e-863a-eea8c002b4e4"
	DataSetClassId_FirmwareUpdate       DataSetClassId = "414e26f6-341b-43b7-90fc-bb9e0b1b0866"
	DataSetClassId_Blink                DataSetClassId = "3b423a40-a676-4ba0-8017-f0b2cd65bc26"
	DataSetClassId_NewDataSetWriterId   DataSetClassId = "2aca55bd-0d6f-41b1-a1c2-2d61afcc21f0"
	DataSetClassId_Reserved0            DataSetClassId = "32c0a57a-24fd-4378-8eb3-4a42a481bf56"
	DataSetClassId_Reserved1            DataSetClassId = "0721c712-f7b5-4ad0-b2a0-58f0a5744963"
	DataSetClassId_Reserved2            DataSetClassId = "d08ade63-777b-464a-8f58-73649dc3ec1c"
	DataSetClassId_Reserved3            DataSetClassId = "965ab782-4669-480e-bab4-c6a3b14ee4cf"
	DataSetClassId_Reserved4            DataSetClassId = "d2ee7674-cad5-4b59-a37c-3631bb0cd409"
	DataSetClassId_Reserved5            DataSetClassId = "a9031581-4c03-4c6a-9c12-dd7fe0123ef5"
	DataSetClassId_Reserved6            DataSetClassId = "04cac81f-9056-4af3-a8a0-686fc499b5f9"
	DataSetClassId_Reserved7            DataSetClassId = "a04f4f5f-b2cc-40d7-9a6a-0884233ccce7"
	DataSetClassId_Reserved8            DataSetClassId = "2b967128-5b7e-4c0f-a6eb-c59af1118ce6"
	DataSetClassId_Reserved9            DataSetClassId = "54b8a5c2-812d-442e-8a91-257af7304161"
	DataSetClassId_Reserved10           DataSetClassId = "2b4e35a2-437a-4a3d-b007-9b1c27b08127"
	DataSetClassId_Reserved11           DataSetClassId = "97439c90-7d10-47b8-8ada-6913a466107d"
	DataSetClassId_Reserved12           DataSetClassId = "4292bf9f-b9d0-4de1-a33e-f4736356d4c1"
	DataSetClassId_Reserved13           DataSetClassId = "0b3b07af-9cf9-4b40-977f-75b5b02bb0d5"
	DataSetClassId_Reserved14           DataSetClassId = "31c6e1c0-88c7-453b-ae9d-85fb6bc50275"
	DataSetClassId_Reserved15           DataSetClassId = "cfa96d87-e65c-42f4-9d9e-026797c37f4c"
	DataSetClassId_Reserved16           DataSetClassId = "0c959a0f-ec83-43c1-9ccc-20ab20a0e171"
	DataSetClassId_Reserved17           DataSetClassId = "bb725543-9de1-4f22-883f-d32c2de5dea1"
)

var dataSetClassIdMapping = map[Resource]string{
	Resource_MAM:                  "360ca8f3-5e66-42a2-8f10-9cdf45f4bf58",
	Resource_Health:               "d8e7b6df-42ba-448a-975a-199f59e8ffeb",
	Resource_Config:               "9d5983db-440d-4474-9fd7-1cd7a6c8b6c2",
	Resource_License:              "2ae0505e-2830-4980-b65e-0bbdf08e2d45",
	Resource_LicenseText:          "a6e6c727-4057-419f-b2ea-3fe9173e71cf",
	Resource_RtLicense:            "ebd12d4b-da1c-4671-ab86-db102fecc603",
	Resource_Event:                "543ae05e-b6d9-4161-a0a3-350a0fac5976",
	Resource_Profile:              "48017c6a-05c8-48d7-9d85-4b08bbb707f3",
	Resource_PublicationList:      "217434d6-6e1e-4230-b907-f52bc9ffe152",
	Resource_SubscriptionList:     "e5d68c47-c276-4929-8ab9-4c1090cac785",
	Resource_Interfaces:           "96d22d73-bce6-42d3-9949-45e0d04e4d54",
	Resource_ReferenceDesignation: "27a75019-164a-496d-a38b-90e8a55c2cfa",
}
