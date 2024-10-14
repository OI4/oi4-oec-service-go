package api

type MasterAssetModel struct {
	Manufacturer       LocalizedText `json:"Manufacturer"`
	ManufacturerUri    string        `json:"ManufacturerUri"`
	Model              LocalizedText `json:"Model"`
	ProductCode        string        `json:"ProductCode"`
	HardwareRevision   string        `json:"HardwareRevision"`
	SoftwareRevision   string        `json:"SoftwareRevision"`
	DeviceRevision     string        `json:"DeviceRevision"`
	DeviceManual       string        `json:"DeviceManual"`
	DeviceClass        string        `json:"DeviceClass"`
	SerialNumber       string        `json:"SerialNumber"`
	ProductInstanceUri string        `json:"ProductInstanceUri"`
	RevisionCounter    int32         `json:"RevisionCounter"`
	Description        LocalizedText `json:"Description"`
}

func (mam *MasterAssetModel) ToOi4Identifier() *Oi4Identifier {
	return &Oi4Identifier{
		ManufacturerUri: mam.ManufacturerUri,
		Model:           mam.Model.Text,
		ProductCode:     mam.ProductCode,
		SerialNumber:    mam.SerialNumber,
	}
}

func (mam *MasterAssetModel) Payload() any {
	return mam
}
