package api

type DataSetMessage struct {

	// Value range: <UINT16>
	// Type: UInt16
	// Requirement: Mandatory
	//
	// An identifier for DataSetWriter which published the DataSetMessage. It is unique within the scope of a Publisher. The related DataSetMetaData (9.2.2) to this DataSetMessage contains the same DataSetWriterId.
	// A range of DataSetWriterIds are reserved for special use cases. A PaginationRequest (9.3.15.1) always uses the DataSetWriterId 1, a Pagination (9.3.15.2) uses the 2 and a Locale (9.3.16) uses the 3.
	// The application starts with DataSetWriterIds from 10, because the DataSetWriterIds up to 9 are reserved for special use cases.
	DataSetWriterId uint16 `json:"DataSetWriterId"`

	// Value range: <UINT32>
	// Type: UInt32
	// Requirement: Optional
	//
	// A strictly sequentially increasing sequence number assigned to the DataSetMessage by the DataSetWriter.
	//
	// NOTE SequenceNumber might be of interest for resources with changing content, such as data, metadata, config, … More static like resources such as mam, health might not benefit from it.
	SequenceNumber *uint32 `json:"SequenceNumber,omitempty"`

	// Value range: <ConfigurationVersionDataType>
	// Type: ConfigurationVersionDataType (9.2.5)
	// Requirement: Optional
	//
	// The MetaDataVersion corresponds with the ConfigurationVersion of a DataSetMetaData message (9.2.5).
	//
	// NOTE MetaDataVersion might be of interest for resources with changing parameter sets, such as data. Resources with fixed metadata set do not benefit from it.
	MetaDataVersion *ConfigurationVersionDataType `json:"MetaDataVersion,omitempty"`

	// Value range: <DateTime>
	// Type: String
	// Requirement: Optional
	//
	// Example: "2019-06-26T13:16:00.000+01:00"
	// Timestamp of type DateTime according to ISO 8601-1:2019 and OPC UA Part 6-5.4.2.6. serialized as String.
	// The time of the data acquisition is indicated! Milliseconds might be of interest.
	//
	// NOTE Timestamp might be of interest for resources with changing content, such as health, data, metadata, config. More static like resources such as mam might not benefit from it.
	Timestamp *string `json:"Timestamp,omitempty"`

	// Value range: <StatusCode>
	// Type: UInt32
	// Requirement: Optional
	// Status code to be used as defined in OPC UA Part 4-7.34.2 and CSV-File.
	// NOTE The Status is not mandatory and shall not be send, when the status is OK. When the Status is unequal to OK, the status codes, provided from OPC Foundation, are used.
	Status *StatusCode `json:"Status,omitempty"`

	// Value range: <Filter>
	// Type: String
	// Requirement: Conditional
	//
	// Depending on related use case, the Filter might be mandatory or optional, but does not belong to OPC UA DataSetMessage according to Part 14-7.2.3.3. In combination with the used resource in the topic, the Filter, together with the Source, contains the readable reference to the DataSetWriterId and is identical to the filter in the topic (8.1.7) if present.
	//
	// Note The Filter helps to combine the DataSet in the Payload with the related source. In OPC UA context this is done via DataSetWriterId, but this is not very intuitive and might need additional actions to get missing information via PublicationList  defined in 9.3.11.
	//
	// note Several resources such as MAM or Health and others do not make use of Filter in Message Bus topic and DataSetMessage.
	Filter Filter `json:"Filter,omitempty"`

	// Value range: <Oi4Identifier>
	// Type: String
	// Requirement: Mandatory
	//
	// The Source is mandatory, but does not belong to OPC UA DataSetMessage according to Part 14-7.2.3.3. In combination with the used resource in the topic, the Source, together with the Filter, contains the readable reference to the DataSetWriterId and is identical to the Source in the topic (8.1.6) if present. The Source always describes the asset providing the information. Therefore the Source is the Oi4Identifer of the application or device.
	//
	// NOTE The Source helps to combine the data in the Payload with the related source. In OPC UA context this is done via DataSetWriterId, but this is not very intuitive and might need additional actions to get missing information via PublicationList defined in 9.3.11.
	Source string `json:"Source"`

	// Value range: <Object>
	// Type: Object
	// Requirement: Mandatory
	//
	// This object contains the name-value pairs specified by the PublishedDataSet.
	//
	// note In general, all built-in data types should be possible, but it seems to be problematic to use ExtensionObject, Variant, DataValue, DiagnosticInfo and in some cases NodeId.
	Payload any `json:"Payload"`
}
