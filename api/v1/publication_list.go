package v1

type PublicationMode string

const (
	PublicationMode_OFF_0                       PublicationMode = "OFF_0"
	PublicationMode_ON_REQUEST_1                PublicationMode = "ON_REQUEST_1"
	PublicationMode_APPLICATION_2               PublicationMode = "APPLICATION_2"
	PublicationMode_SOURCE_3                    PublicationMode = "SOURCE_3"
	PublicationMode_FILTER_4                    PublicationMode = "FILTER_4"
	PublicationMode_APPLICATION_SOURCE_5        PublicationMode = "APPLICATION_SOURCE_5"
	PublicationMode_APPLICATION_FILTER_6        PublicationMode = "APPLICATION_FILTER_6"
	PublicationMode_SOURCE_FILTER_7             PublicationMode = "SOURCE_FILTER_7"
	PublicationMode_APPLICATION_SOURCE_FILTER_8 PublicationMode = "APPLICATION_SOURCE_FILTER_8"
)

type PublicationConfig string

const (
	PublicationConfig_NONE_0              PublicationConfig = "NONE_0"
	PublicationConfig_MODE_1              PublicationConfig = "MODE_1"
	PublicationConfig_INTERVAL_2          PublicationConfig = "INTERVAL_2"
	PublicationConfig_MODE_AND_INTERVAL_3 PublicationConfig = "MODE_AND_INTERVAL_3"
)

type PublicationList struct {
	Resource        `json:"Resource"`
	Source          Oi4IdentifierPath  `json:"Source"`
	Filter          *interface{}       `json:"Filter,omitempty"`
	DataSetWriterId uint16             `json:"DataSetWriterId"`
	Mode            PublicationMode    `json:"Mode"`
	Interval        uint32             `json:"Interval"`
	Precisions      map[string]float32 `json:"Precisions"`
	Config          PublicationConfig  `json:"PublicationConfig"`
}
