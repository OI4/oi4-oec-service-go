package api

type ServiceMessageType string

const (
	MSG ServiceMessageType = "MSG"
)

type ServiceNetworkMessage struct {

	// Value range: <unixTimestampInMs-PublisherId>
	// Type: String
	// Requirement: Mandatory
	//
	// Example: "1567062381000-OTConnector/company.com/type/order/4711"
	// Must be unique for any single package of this PublisherId.
	//
	// note The MessageId must be unique. The Open Industry 4.0 Alliance defiens the type as a combination of the current timestamp in ms precision with the PublisherId. In some rare cases, the system timestamp might not be precise enough to avoid sending packages with an unique MessageId. In these cases, the application must guarantee the uniqueness by providing an additional parameter (example. MessageId = <unixTimestampInMs><counter>-<PublisherId>).
	MessageId string `json:"MessageId"`

	//Value range: "MSG"
	// Type: String
	// Requirement: Mandatory
	//
	// Only MSG is valid here
	MessageType ServiceMessageType `json:"MessageType"`

	// Value range: <serviceType>/<appId>
	// Type: String consisting of ServiceType (see 8.1.2) and Oi4Identifier (see 3.1) - separated by a /.
	// Requirement: Mandatory in Open Industry 4.0 Alliance context, but not in OPC UA context.
	PublisherId string `json:"PublisherId"`

	// Value range: <GUID> in accordance to OPC UA Part 6-5.1.3: 16 Byte as JSON string with separator (Part 6-5.4.2.7).
	// Type: String
	// Requirement: Optional
	// Example: "f1875b4a-3209-431b-a38d-2df5758f92c8"
	// The DataSetClassId allows to refer to the DataSetClass describing the structure of the message. The DataSetClassId identifies a well defined DataSet specified by the Alliance or some other standards.
	// For some recurring use cases, such as NewDataSetWriterId, fixed GUIDs are specified from the Alliance and must be used (A2).
	//
	// Note The DataSetClassId shall be present for all resources defined by the Open Industry 4.0 Alliance, which are having a DataSetClassId. This guarantees the availability of a fully schema verifiable messaging.
	DataSetClassId string `json:"DataSetClassId,omitempty"`

	// Value range: <empty/omitted> or <MessageId>
	// Type: String
	// Requirement: Conditional
	// Shows the flow between the causal event and its consequences.
	// note The CorrelationId is filled in by the first consumer with the MessageId of the original message and then passed on from service to service until the message is no longer processed.
	CorrelationId string `json:"CorrelationId,omitempty"`

	// Value range: <ServiceParametersRequest> or <ServiceParametersResponse>
	// Type: ServiceParametersRequest object (see 9.2.17) or ServiceParametersResponse object (see 9.2.19)
	// Requirement: Mandatory
	Message any `json:"Message"`
}
