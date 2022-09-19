package types

type NetworkMessageType string

const (
	UA_DATA NetworkMessageType = "ua-data"
)

type NetworkMessage struct {

	// Value range: <unixTimestampInMs-PublisherId>
	// Example: "1567062381000-OTConnector/company.com/model/productCode/4711"
	// note The MessageId must be unique. The Open Industry 4.0 Alliance defines the type as a combination of the current timestamp in ms precision with the PublisherId. In some rare cases, the system timestamp might not be precise enough to avoid sending packages with an unique MessageId. In these cases, the application must guarantee the uniqueness by providing an additional  parameter (example. MessageId = <unixTimestampInMs><counter>-<PublisherId>).
	MessageId string `json:"MessageId"`

	// Only ua-data is valid here.
	MessageType NetworkMessageType `json:"MessageType"`

	// Value range: <ServiceType>/<AppId>
	// Type: String consisting of ServiceType (see 8.1.2) and Oi4Identifier (see 3.1 and 8.1.3) - separated by a /.
	// Requirement: Mandatory in Open Industry 4.0 Alliance context, but not in OPC UA context.
	PublisherId string `json:"PublisherId"`

	// Value range: <GUID> in accordance to OPC UA Part 6-5.1.3: 16 Byte as JSON string with separator (Part 6-5.4.2.7).
	// Type: String
	// Requirement: Optional
	// Example: "f1875b4a-3209-431b-a38d-2df5758f92c8"
	// The DataSetClassId allows to refer to the DataSetClass describing the structure of the message. The DataSetClassId identifies a well defined DataSet specified by the Alliance or some other standards and refers to related metadata.
	// For some recurring use cases, such as MasterAssetModel, fixed GUIDs are specified from the Alliance and must be used (A2).
	// Note The DataSetClassId shall be present for all resources defined by the Open Industry 4.0 Alliance, which are having a DataSetClassId. This guarantees the availability of a fully schema verifiable messaging.
	DataSetClassId string `json:"DataSetClassId,omitempty"`

	// Value range: <empty/omitted> or <MessageId>
	// Type: String
	// Requirement: Conditional
	// Shows the flow between the causal event and its consequences. The CorrelationId does not belong to OPC UA DataSetMessage according to Part 14-7.2.3.3.
	// note The CorrelationId is filled in by the first consumer with the MessageId of the original message and then passed on from service to service until the message is no longer processed.
	CorrelationId string `json:"CorrelationId,omitempty"`

	// Value range: <array of DataSetMessage>
	// Type: DataSetMessage object (see 9.2.3)
	// Requirement: Mandatory
	Messages []*DataSetMessage `json:"Messages"`
}
