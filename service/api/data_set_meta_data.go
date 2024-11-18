package api

type DataSetMetaDataMessageType string

const (
	UA_METADATA = "ua-metadata"
)

type DataSetMetaData struct {

	// Value range: <unixTimestampInMs-PublisherId>
	// Type: String
	// Requirement: Mandatory
	//
	// Example: "1567062381000-http://company.com/type/order/4711"
	// Must be unique for any single package of this PublisherId.
	//
	// note The MessageId must be unique. The Open Industry 4.0 Alliance defines the type as a combination of the current timestamp in ms precision with the PublisherId. In some rare cases, the system timestamp might not be precise enough to avoid sending packages with an unique MessageId. In these cases, the application must guarantee the uniqueness by providing an additional parameter (example. MessageId = <unixTimestampInMs><counter>-<PublisherId>).
	MessageId string `json:"MessageId"`

	// Value range: "ua-metadata"
	// Type: String
	// Requirement: Mandatory
	//
	// Only ua-metadata is valid here.
	MessageType DataSetMetaDataMessageType `json:"MessageType"`

	// Value range: <ServiceType>/<AppId>
	// Type: String consisting of ServiceType (see 8.1.2) and Oi4Identifier (see 3.1 and 8.1.3) - separated by a /.
	// Requirement Mandatory
	PublisherId string `json:"PublisherId"`

	// Value range: <UINT16>
	// Type: UInt16
	// Requirement: Mandatory
	//
	// An identifier for DataSetWriter which published the DataSetMetaData. It is unique within the scope of a Publisher. The related DataSetMessage (9.2.3) to this DataSetMetaData contains the same DataSetWriterId.
	//
	// note The DataSetWriterId is not persistent and can change on every power cycle of a DSWID
	DataSetWriterId uint16 `json:"DataSetWriterId"`

	// Value range: <Filter>
	// Type: String
	// Requirement: Mandatory
	// The Filter is mandatory, but does not belong to OPC UA DataSetMetaData according to Part 14-7.2.3.4.2. In combination with the used resource in the topic, the Filter, together with the Source, contains the readable reference to the DataSetWriterId and is identical to the filter in the topic (8.1.7).
	// note The Filter helps to combine the MetaData with the related source. In OPC UA context this is done via DataSetWriterId, but this is not very intuitive and might need additional actions to get missing information via PublicationList defined in 9.3.11.
	Filter interface{} `json:"Filter"`

	// Value range: <Oi4Identifier>
	// Type: String
	// Requirement: Mandatory
	//
	// The Source is mandatory, but does not belong to OPC UA DataSetMessage according to Part 14-7.2.3.3. In combination with the used resource in the topic, the Source, together with the Filter, contains the readable reference to the DataSetWriterId and is identical to the Source in the topic (8.1.6) if present.
	//
	// NOTE The Source helps to combine the MetaData with the related source. In OPC UA context this is done via DataSetWriterId, but this is not very intuitive and might need additional actions to get missing information via PublicationList defined in 9.3.11.
	Source string `json:"Source"`

	// Value range: <empty/omitted> or <MessageId>
	// Type: String
	// Requirement: Conditional
	//
	// Shows the flow between the causal event and its consequences. The CorrelationId does not belong to OPC UA DataSetMessage according to Part 14-7.2.3.3.
	//
	// NOTE The CorrelationId is filled in by the first consumer with the MessageId of the original message and then passed on from service to service until the message is no longer processed.
	CorrelationId string `json:"CorrelationId"`

	// Value range: <DataSetMetaDataType>
	// Type: DataSetMetaDataType object (9.2.4)
	// Requirement: Mandatory
	MetaData DataSetMetaDataType `json:"MetaData"`
}
