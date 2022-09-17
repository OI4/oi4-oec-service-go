package v1

type StatusCode uint32

const (
	Status_Good                                                            StatusCode = 0x00000000 //The operation succeeded.
	Status_Uncertain                                                       StatusCode = 0x40000000 //The operation was uncertain.
	Status_Bad                                                             StatusCode = 0x80000000 //The operation failed.
	Status_BadUnexpectedError                                              StatusCode = 0x80010000 //An unexpected error occurred.
	Status_BadInternalError                                                StatusCode = 0x80020000 //An internal error occurred as a result of a programming or configuration error.
	Status_BadOutOfMemory                                                  StatusCode = 0x80030000 //Not enough memory to complete the operation.
	Status_BadResourceUnavailable                                          StatusCode = 0x80040000 //An operating system resource is not available.
	Status_BadCommunicationError                                           StatusCode = 0x80050000 //A low level communication error occurred.
	Status_BadEncodingError                                                StatusCode = 0x80060000 //Encoding halted because of invalid data in the objects being serialized.
	Status_BadDecodingError                                                StatusCode = 0x80070000 //Decoding halted because of invalid data in the stream.
	Status_BadEncodingLimitsExceeded                                       StatusCode = 0x80080000 //The message encoding/decoding limits imposed by the stack have been exceeded.
	Status_BadRequestTooLarge                                              StatusCode = 0x80B80000 //The request message size exceeds limits set by the server.
	Status_BadResponseTooLarge                                             StatusCode = 0x80B90000 //The response message size exceeds limits set by the client.
	Status_BadUnknownResponse                                              StatusCode = 0x80090000 //An unrecognized response was received from the server.
	Status_BadTimeout                                                      StatusCode = 0x800A0000 //The operation timed out.
	Status_BadServiceUnsupported                                           StatusCode = 0x800B0000 //The server does not support the requested service.
	Status_BadShutdown                                                     StatusCode = 0x800C0000 //The operation was cancelled because the application is shutting down.
	Status_BadServerNotConnected                                           StatusCode = 0x800D0000 //The operation could not complete because the client is not connected to the server.
	Status_BadServerHalted                                                 StatusCode = 0x800E0000 //The server has stopped and cannot process any requests.
	Status_BadNothingToDo                                                  StatusCode = 0x800F0000 //There was nothing to do because the client passed a list of operations with no elements.
	Status_BadTooManyOperations                                            StatusCode = 0x80100000 //The request could not be processed because it specified too many operations.
	Status_BadTooManyMonitoredItems                                        StatusCode = 0x80DB0000 //The request could not be processed because there are too many monitored items in the subscription.
	Status_BadDataTypeIdUnknown                                            StatusCode = 0x80110000 //The extension object cannot be (de)serialized because the data type id is not recognized.
	Status_BadCertificateInvalid                                           StatusCode = 0x80120000 //The certificate provided as a parameter is not valid.
	Status_BadSecurityChecksFailed                                         StatusCode = 0x80130000 //An error occurred verifying security.
	Status_BadCertificatePolicyCheckFailed                                 StatusCode = 0x81140000 //The certificate does not meet the requirements of the security policy.
	Status_BadCertificateTimeInvalid                                       StatusCode = 0x80140000 //The certificate has expired or is not yet valid.
	Status_BadCertificateIssuerTimeInvalid                                 StatusCode = 0x80150000 //An issuer certificate has expired or is not yet valid.
	Status_BadCertificateHostNameInvalid                                   StatusCode = 0x80160000 //The HostName used to connect to a server does not match a HostName in the certificate.
	Status_BadCertificateUriInvalid                                        StatusCode = 0x80170000 //The URI specified in the ApplicationDescription does not match the URI in the certificate.
	Status_BadCertificateUseNotAllowed                                     StatusCode = 0x80180000 //The certificate may not be used for the requested operation.
	Status_BadCertificateIssuerUseNotAllowed                               StatusCode = 0x80190000 //The issuer certificate may not be used for the requested operation.
	Status_BadCertificateUntrusted                                         StatusCode = 0x801A0000 //The certificate is not trusted.
	Status_BadCertificateRevocationUnknown                                 StatusCode = 0x801B0000 //It was not possible to determine if the certificate has been revoked.
	Status_BadCertificateIssuerRevocationUnknown                           StatusCode = 0x801C0000 //It was not possible to determine if the issuer certificate has been revoked.
	Status_BadCertificateRevoked                                           StatusCode = 0x801D0000 //The certificate has been revoked.
	Status_BadCertificateIssuerRevoked                                     StatusCode = 0x801E0000 //The issuer certificate has been revoked.
	Status_BadCertificateChainIncomplete                                   StatusCode = 0x810D0000 //The certificate chain is incomplete.
	Status_BadUserAccessDenied                                             StatusCode = 0x801F0000 //User does not have permission to perform the requested operation.
	Status_BadIdentityTokenInvalid                                         StatusCode = 0x80200000 //The user identity token is not valid.
	Status_BadIdentityTokenRejected                                        StatusCode = 0x80210000 //The user identity token is valid but the server has rejected it.
	Status_BadSecureChannelIdInvalid                                       StatusCode = 0x80220000 //The specified secure channel is no longer valid.
	Status_BadInvalidTimestamp                                             StatusCode = 0x80230000 //The timestamp is outside the range allowed by the server.
	Status_BadNonceInvalid                                                 StatusCode = 0x80240000 //The nonce does appear to be not a random value or it is not the correct length.
	Status_BadSessionIdInvalid                                             StatusCode = 0x80250000 //The session id is not valid.
	Status_BadSessionClosed                                                StatusCode = 0x80260000 //The session was closed by the client.
	Status_BadSessionNotActivated                                          StatusCode = 0x80270000 //The session cannot be used because ActivateSession has not been called.
	Status_BadSubscriptionIdInvalid                                        StatusCode = 0x80280000 //The subscription id is not valid.
	Status_BadRequestHeaderInvalid                                         StatusCode = 0x802A0000 //The header for the request is missing or invalid.
	Status_BadTimestampsToReturnInvalid                                    StatusCode = 0x802B0000 //The timestamps to return parameter is invalid.
	Status_BadRequestCancelledByClient                                     StatusCode = 0x802C0000 //The request was cancelled by the client.
	Status_BadTooManyArguments                                             StatusCode = 0x80E50000 //Too many arguments were provided.
	Status_BadLicenseExpired                                               StatusCode = 0x810E0000 //The server requires a license to operate in general or to perform a service or operation, but existing license is expired.
	Status_BadLicenseLimitsExceeded                                        StatusCode = 0x810F0000 //The server has limits on number of allowed operations / objects, based on installed licenses, and these limits where exceeded.
	Status_BadLicenseNotAvailable                                          StatusCode = 0x81100000 //The server does not have a license which is required to operate in general or to perform a service or operation.
	Status_GoodSubscriptionTransferred                                     StatusCode = 0x002D0000 //The subscription was transferred to another session.
	Status_GoodCompletesAsynchronously                                     StatusCode = 0x002E0000 //The processing will complete asynchronously.
	Status_GoodOverload                                                    StatusCode = 0x002F0000 //Sampling has slowed down due to resource limitations.
	Status_GoodClamped                                                     StatusCode = 0x00300000 //The value written was accepted but was clamped.
	Status_BadNoCommunication                                              StatusCode = 0x80310000 //Communication with the data source is defined, but not established, and there is no last known value available.
	Status_BadWaitingForInitialData                                        StatusCode = 0x80320000 //Waiting for the server to obtain values from the underlying data source.
	Status_BadNodeIdInvalid                                                StatusCode = 0x80330000 //The syntax of the node id is not valid.
	Status_BadNodeIdUnknown                                                StatusCode = 0x80340000 //The node id refers to a node that does not exist in the server address space.
	Status_BadAttributeIdInvalid                                           StatusCode = 0x80350000 //The attribute is not supported for the specified Node.
	Status_BadIndexRangeInvalid                                            StatusCode = 0x80360000 //The syntax of the index range parameter is invalid.
	Status_BadIndexRangeNoData                                             StatusCode = 0x80370000 //No data exists within the range of indexes specified.
	Status_BadDataEncodingInvalid                                          StatusCode = 0x80380000 //The data encoding is invalid.
	Status_BadDataEncodingUnsupported                                      StatusCode = 0x80390000 //The server does not support the requested data encoding for the node.
	Status_BadNotReadable                                                  StatusCode = 0x803A0000 //The access level does not allow reading or subscribing to the Node.
	Status_BadNotWritable                                                  StatusCode = 0x803B0000 //The access level does not allow writing to the Node.
	Status_BadOutOfRange                                                   StatusCode = 0x803C0000 //The value was out of range.
	Status_BadNotSupported                                                 StatusCode = 0x803D0000 //The requested operation is not supported.
	Status_BadNotFound                                                     StatusCode = 0x803E0000 //A requested item was not found or a search operation ended without success.
	Status_BadObjectDeleted                                                StatusCode = 0x803F0000 //The object cannot be used because it has been deleted.
	Status_BadNotImplemented                                               StatusCode = 0x80400000 //Requested operation is not implemented.
	Status_BadMonitoringModeInvalid                                        StatusCode = 0x80410000 //The monitoring mode is invalid.
	Status_BadMonitoredItemIdInvalid                                       StatusCode = 0x80420000 //The monitoring item id does not refer to a valid monitored item.
	Status_BadMonitoredItemFilterInvalid                                   StatusCode = 0x80430000 //The monitored item filter parameter is not valid.
	Status_BadMonitoredItemFilterUnsupported                               StatusCode = 0x80440000 //The server does not support the requested monitored item filter.
	Status_BadFilterNotAllowed                                             StatusCode = 0x80450000 //A monitoring filter cannot be used in combination with the attribute specified.
	Status_BadStructureMissing                                             StatusCode = 0x80460000 //A mandatory structured parameter was missing or null.
	Status_BadEventFilterInvalid                                           StatusCode = 0x80470000 //The event filter is not valid.
	Status_BadContentFilterInvalid                                         StatusCode = 0x80480000 //The content filter is not valid.
	Status_BadFilterOperatorInvalid                                        StatusCode = 0x80C10000 //An unrecognized operator was provided in a filter.
	Status_BadFilterOperatorUnsupported                                    StatusCode = 0x80C20000 //A valid operator was provided, but the server does not provide support for this filter operator.
	Status_BadFilterOperandCountMismatch                                   StatusCode = 0x80C30000 //The number of operands provided for the filter operator was less then expected for the operand provided.
	Status_BadFilterOperandInvalid                                         StatusCode = 0x80490000 //The operand used in a content filter is not valid.
	Status_BadFilterElementInvalid                                         StatusCode = 0x80C40000 //The referenced element is not a valid element in the content filter.
	Status_BadFilterLiteralInvalid                                         StatusCode = 0x80C50000 //The referenced literal is not a valid value.
	Status_BadContinuationPointInvalid                                     StatusCode = 0x804A0000 //The continuation point provide is longer valid.
	Status_BadNoContinuationPoints                                         StatusCode = 0x804B0000 //The operation could not be processed because all continuation points have been allocated.
	Status_BadReferenceTypeIdInvalid                                       StatusCode = 0x804C0000 //The reference type id does not refer to a valid reference type node.
	Status_BadBrowseDirectionInvalid                                       StatusCode = 0x804D0000 //The browse direction is not valid.
	Status_BadNodeNotInView                                                StatusCode = 0x804E0000 //The node is not part of the view.
	Status_BadNumericOverflow                                              StatusCode = 0x81120000 //The number was not accepted because of a numeric overflow.
	Status_BadServerUriInvalid                                             StatusCode = 0x804F0000 //The ServerUri is not a valid URI.
	Status_BadServerNameMissing                                            StatusCode = 0x80500000 //No ServerName was specified.
	Status_BadDiscoveryUrlMissing                                          StatusCode = 0x80510000 //No DiscoveryUrl was specified.
	Status_BadSempahoreFileMissing                                         StatusCode = 0x80520000 //The semaphore file specified by the client is not valid.
	Status_BadRequestTypeInvalid                                           StatusCode = 0x80530000 //The security token request type is not valid.
	Status_BadSecurityModeRejected                                         StatusCode = 0x80540000 //The security mode does not meet the requirements set by the server.
	Status_BadSecurityPolicyRejected                                       StatusCode = 0x80550000 //The security policy does not meet the requirements set by the server.
	Status_BadTooManySessions                                              StatusCode = 0x80560000 //The server has reached its maximum number of sessions.
	Status_BadUserSignatureInvalid                                         StatusCode = 0x80570000 //The user token signature is missing or invalid.
	Status_BadApplicationSignatureInvalid                                  StatusCode = 0x80580000 //The signature generated with the client certificate is missing or invalid.
	Status_BadNoValidCertificates                                          StatusCode = 0x80590000 //The client did not provide at least one software certificate that is valid and meets the profile requirements for the server.
	Status_BadIdentityChangeNotSupported                                   StatusCode = 0x80C60000 //The server does not support changing the user identity assigned to the session.
	Status_BadRequestCancelledByRequest                                    StatusCode = 0x805A0000 //The request was cancelled by the client with the Cancel service.
	Status_BadParentNodeIdInvalid                                          StatusCode = 0x805B0000 //The parent node id does not to refer to a valid node.
	Status_BadReferenceNotAllowed                                          StatusCode = 0x805C0000 //The reference could not be created because it violates constraints imposed by the data model.
	Status_BadNodeIdRejected                                               StatusCode = 0x805D0000 //The requested node id was reject because it was either invalid or server does not allow node ids to be specified by the client.
	Status_BadNodeIdExists                                                 StatusCode = 0x805E0000 //The requested node id is already used by another node.
	Status_BadNodeClassInvalid                                             StatusCode = 0x805F0000 //The node class is not valid.
	Status_BadBrowseNameInvalid                                            StatusCode = 0x80600000 //The browse name is invalid.
	Status_BadBrowseNameDuplicated                                         StatusCode = 0x80610000 //The browse name is not unique among nodes that share the same relationship with the parent.
	Status_BadNodeAttributesInvalid                                        StatusCode = 0x80620000 //The node attributes are not valid for the node class.
	Status_BadTypeDefinitionInvalid                                        StatusCode = 0x80630000 //The type definition node id does not reference an appropriate type node.
	Status_BadSourceNodeIdInvalid                                          StatusCode = 0x80640000 //The source node id does not reference a valid node.
	Status_BadTargetNodeIdInvalid                                          StatusCode = 0x80650000 //The target node id does not reference a valid node.
	Status_BadDuplicateReferenceNotAllowed                                 StatusCode = 0x80660000 //The reference type between the nodes is already defined.
	Status_BadInvalidSelfReference                                         StatusCode = 0x80670000 //The server does not allow this type of self reference on this node.
	Status_BadReferenceLocalOnly                                           StatusCode = 0x80680000 //The reference type is not valid for a reference to a remote server.
	Status_BadNoDeleteRights                                               StatusCode = 0x80690000 //The server will not allow the node to be deleted.
	Status_UncertainReferenceNotDeleted                                    StatusCode = 0x40BC0000 //The server was not able to delete all target references.
	Status_BadServerIndexInvalid                                           StatusCode = 0x806A0000 //The server index is not valid.
	Status_BadViewIdUnknown                                                StatusCode = 0x806B0000 //The view id does not refer to a valid view node.
	Status_BadViewTimestampInvalid                                         StatusCode = 0x80C90000 //The view timestamp is not available or not supported.
	Status_BadViewParameterMismatch                                        StatusCode = 0x80CA0000 //The view parameters are not consistent with each other.
	Status_BadViewVersionInvalid                                           StatusCode = 0x80CB0000 //The view version is not available or not supported.
	Status_UncertainNotAllNodesAvailable                                   StatusCode = 0x40C00000 //The list of references may not be complete because the underlying system is not available.
	Status_GoodResultsMayBeIncomplete                                      StatusCode = 0x00BA0000 //The server should have followed a reference to a node in a remote server but did not. The result set may be incomplete.
	Status_BadNotTypeDefinition                                            StatusCode = 0x80C80000 //The provided Nodeid was not a type definition nodeid.
	Status_UncertainReferenceOutOfServer                                   StatusCode = 0x406C0000 //One of the references to follow in the relative path references to a node in the address space in another server.
	Status_BadTooManyMatches                                               StatusCode = 0x806D0000 //The requested operation has too many matches to return.
	Status_BadQueryTooComplex                                              StatusCode = 0x806E0000 //The requested operation requires too many resources in the server.
	Status_BadNoMatch                                                      StatusCode = 0x806F0000 //The requested operation has no match to return.
	Status_BadMaxAgeInvalid                                                StatusCode = 0x80700000 //The max age parameter is invalid.
	Status_BadSecurityModeInsufficient                                     StatusCode = 0x80E60000 //The operation is not permitted over the current secure channel.
	Status_BadHistoryOperationInvalid                                      StatusCode = 0x80710000 //The history details parameter is not valid.
	Status_BadHistoryOperationUnsupported                                  StatusCode = 0x80720000 //The server does not support the requested operation.
	Status_BadInvalidTimestampArgument                                     StatusCode = 0x80BD0000 //The defined timestamp to return was invalid.
	Status_BadWriteNotSupported                                            StatusCode = 0x80730000 //The server does not support writing the combination of value, status and timestamps provided.
	Status_BadTypeMismatch                                                 StatusCode = 0x80740000 //The value supplied for the attribute is not of the same type as the attribute's value.
	Status_BadMethodInvalid                                                StatusCode = 0x80750000 //The method id does not refer to a method for the specified object.
	Status_BadArgumentsMissing                                             StatusCode = 0x80760000 //The client did not specify all of the input arguments for the method.
	Status_BadNotExecutable                                                StatusCode = 0x81110000 //The executable attribute does not allow the execution of the method.
	Status_BadTooManySubscriptions                                         StatusCode = 0x80770000 //The server has reached its maximum number of subscriptions.
	Status_BadTooManyPublishRequests                                       StatusCode = 0x80780000 //The server has reached the maximum number of queued publish requests.
	Status_BadNoSubscription                                               StatusCode = 0x80790000 //There is no subscription available for this session.
	Status_BadSequenceNumberUnknown                                        StatusCode = 0x807A0000 //The sequence number is unknown to the server.
	Status_GoodRetransmissionQueueNotSupported                             StatusCode = 0x00DF0000 //The Server does not support retransmission queue and acknowledgement of sequence numbers is not available.
	Status_BadMessageNotAvailable                                          StatusCode = 0x807B0000 //The requested notification message is no longer available.
	Status_BadInsufficientClientProfile                                    StatusCode = 0x807C0000 //The client of the current session does not support one or more Profiles that are necessary for the subscription.
	Status_BadStateNotActive                                               StatusCode = 0x80BF0000 //The sub-state machine is not currently active.
	Status_BadAlreadyExists                                                StatusCode = 0x81150000 //An equivalent rule already exists.
	Status_BadTcpServerTooBusy                                             StatusCode = 0x807D0000 //The server cannot process the request because it is too busy.
	Status_BadTcpMessageTypeInvalid                                        StatusCode = 0x807E0000 //The type of the message specified in the header invalid.
	Status_BadTcpSecureChannelUnknown                                      StatusCode = 0x807F0000 //The SecureChannelId and/or TokenId are not currently in use.
	Status_BadTcpMessageTooLarge                                           StatusCode = 0x80800000 //The size of the message chunk specified in the header is too large.
	Status_BadTcpNotEnoughResources                                        StatusCode = 0x80810000 //There are not enough resources to process the request.
	Status_BadTcpInternalError                                             StatusCode = 0x80820000 //An internal error occurred.
	Status_BadTcpEndpointUrlInvalid                                        StatusCode = 0x80830000 //The server does not recognize the QueryString specified.
	Status_BadRequestInterrupted                                           StatusCode = 0x80840000 //The request could not be sent because of a network interruption.
	Status_BadRequestTimeout                                               StatusCode = 0x80850000 //Timeout occurred while processing the request.
	Status_BadSecureChannelClosed                                          StatusCode = 0x80860000 //The secure channel has been closed.
	Status_BadSecureChannelTokenUnknown                                    StatusCode = 0x80870000 //The token has expired or is not recognized.
	Status_BadSequenceNumberInvalid                                        StatusCode = 0x80880000 //The sequence number is not valid.
	Status_BadProtocolVersionUnsupported                                   StatusCode = 0x80BE0000 //The applications do not have compatible protocol versions.
	Status_BadConfigurationError                                           StatusCode = 0x80890000 //There is a problem with the configuration that affects the usefulness of the value.
	Status_BadNotConnected                                                 StatusCode = 0x808A0000 //The variable should receive its value from another variable, but has never been configured to do so.
	Status_BadDeviceFailure                                                StatusCode = 0x808B0000 //There has been a failure in the device/data source that generates the value that has affected the value.
	Status_BadSensorFailure                                                StatusCode = 0x808C0000 //There has been a failure in the sensor from which the value is derived by the device/data source.
	Status_BadOutOfService                                                 StatusCode = 0x808D0000 //The source of the data is not operational.
	Status_BadDeadbandFilterInvalid                                        StatusCode = 0x808E0000 //The deadband filter is not valid.
	Status_UncertainNoCommunicationLastUsableValue                         StatusCode = 0x408F0000 //Communication to the data source has failed. The variable value is the last value that had a good quality.
	Status_UncertainLastUsableValue                                        StatusCode = 0x40900000 //Whatever was updating this value has stopped doing so.
	Status_UncertainSubstituteValue                                        StatusCode = 0x40910000 //The value is an operational value that was manually overwritten.
	Status_UncertainInitialValue                                           StatusCode = 0x40920000 //The value is an initial value for a variable that normally receives its value from another variable.
	Status_UncertainSensorNotAccurate                                      StatusCode = 0x40930000 //The value is at one of the sensor limits.
	Status_UncertainEngineeringUnitsExceeded                               StatusCode = 0x40940000 //The value is outside of the range of values defined for this parameter.
	Status_UncertainSubNormal                                              StatusCode = 0x40950000 //The value is derived from multiple sources and has less than the required number of Good sources.
	Status_GoodLocalOverride                                               StatusCode = 0x00960000 //The value has been overridden.
	Status_BadRefreshInProgress                                            StatusCode = 0x80970000 //This Condition refresh failed, a Condition refresh operation is already in progress.
	Status_BadConditionAlreadyDisabled                                     StatusCode = 0x80980000 //This condition has already been disabled.
	Status_BadConditionAlreadyEnabled                                      StatusCode = 0x80CC0000 //This condition has already been enabled.
	Status_BadConditionDisabled                                            StatusCode = 0x80990000 //Property not available, this condition is disabled.
	Status_BadEventIdUnknown                                               StatusCode = 0x809A0000 //The specified event id is not recognized.
	Status_BadEventNotAcknowledgeable                                      StatusCode = 0x80BB0000 //The event cannot be acknowledged.
	Status_BadDialogNotActive                                              StatusCode = 0x80CD0000 //The dialog condition is not active.
	Status_BadDialogResponseInvalid                                        StatusCode = 0x80CE0000 //The response is not valid for the dialog.
	Status_BadConditionBranchAlreadyAcked                                  StatusCode = 0x80CF0000 //The condition branch has already been acknowledged.
	Status_BadConditionBranchAlreadyConfirmed                              StatusCode = 0x80D00000 //The condition branch has already been confirmed.
	Status_BadConditionAlreadyShelved                                      StatusCode = 0x80D10000 //The condition has already been shelved.
	Status_BadConditionNotShelved                                          StatusCode = 0x80D20000 //The condition is not currently shelved.
	Status_BadShelvingTimeOutOfRange                                       StatusCode = 0x80D30000 //The shelving time not within an acceptable range.
	Status_BadNoData                                                       StatusCode = 0x809B0000 //No data exists for the requested time range or event filter.
	Status_BadBoundNotFound                                                StatusCode = 0x80D70000 //No data found to provide upper or lower bound value.
	Status_BadBoundNotSupported                                            StatusCode = 0x80D80000 //The server cannot retrieve a bound for the variable.
	Status_BadDataLost                                                     StatusCode = 0x809D0000 //Data is missing due to collection started/stopped/lost.
	Status_BadDataUnavailable                                              StatusCode = 0x809E0000 //Expected data is unavailable for the requested time range due to an un-mounted volume, an off-line archive or tape, or similar reason for temporary unavailability.
	Status_BadEntryExists                                                  StatusCode = 0x809F0000 //The data or event was not successfully inserted because a matching entry exists.
	Status_BadNoEntryExists                                                StatusCode = 0x80A00000 //The data or event was not successfully updated because no matching entry exists.
	Status_BadTimestampNotSupported                                        StatusCode = 0x80A10000 //The client requested history using a timestamp format the server does not support (i.e requested ServerTimestamp when server only supports SourceTimestamp).
	Status_GoodEntryInserted                                               StatusCode = 0x00A20000 //The data or event was successfully inserted into the historical database.
	Status_GoodEntryReplaced                                               StatusCode = 0x00A30000 //The data or event field was successfully replaced in the historical database.
	Status_UncertainDataSubNormal                                          StatusCode = 0x40A40000 //The value is derived from multiple values and has less than the required number of Good values.
	Status_GoodNoData                                                      StatusCode = 0x00A50000 //No data exists for the requested time range or event filter.
	Status_GoodMoreData                                                    StatusCode = 0x00A60000 //The data or event field was successfully replaced in the historical database.
	Status_BadAggregateListMismatch                                        StatusCode = 0x80D40000 //The requested number of Aggregates does not match the requested number of NodeIds.
	Status_BadAggregateNotSupported                                        StatusCode = 0x80D50000 //The requested Aggregate is not support by the server.
	Status_BadAggregateInvalidInputs                                       StatusCode = 0x80D60000 //The aggregate value could not be derived due to invalid data inputs.
	Status_BadAggregateConfigurationRejected                               StatusCode = 0x80DA0000 //The aggregate configuration is not valid for specified node.
	Status_GoodDataIgnored                                                 StatusCode = 0x00D90000 //The request specifies fields which are not valid for the EventType or cannot be saved by the historian.
	Status_BadRequestNotAllowed                                            StatusCode = 0x80E40000 //The request was rejected by the server because it did not meet the criteria set by the server.
	Status_BadRequestNotComplete                                           StatusCode = 0x81130000 //The request has not been processed by the server yet.
	Status_BadTicketRequired                                               StatusCode = 0x811F0000 //The device identity needs a ticket before it can be accepted.
	Status_BadTicketInvalid                                                StatusCode = 0x81200000 //The device identity needs a ticket before it can be accepted.
	Status_GoodEdited                                                      StatusCode = 0x00DC0000 //The value does not come from the real source and has been edited by the server.
	Status_GoodPostActionFailed                                            StatusCode = 0x00DD0000 //There was an error in execution of these post-actions.
	Status_UncertainDominantValueChanged                                   StatusCode = 0x40DE0000 //The related EngineeringUnit has been changed but the Variable Value is still provided based on the previous unit.
	Status_GoodDependentValueChanged                                       StatusCode = 0x00E00000 //A dependent value has been changed but the change has not been applied to the device.
	Status_BadDominantValueChanged                                         StatusCode = 0x80E10000 //The related EngineeringUnit has been changed but this change has not been applied to the device. The Variable Value is still dependent on the previous unit but its status is currently Bad.
	Status_UncertainDependentValueChanged                                  StatusCode = 0x40E20000 //A dependent value has been changed but the change has not been applied to the device. The quality of the dominant variable is uncertain.
	Status_BadDependentValueChanged                                        StatusCode = 0x80E30000 //A dependent value has been changed but the change has not been applied to the device. The quality of the dominant variable is Bad.
	Status_GoodEdited_DependentValueChanged                                StatusCode = 0x01160000 //It is delivered with a dominant Variable value when a dependent Variable has changed but the change has not been applied.
	Status_GoodEdited_DominantValueChanged                                 StatusCode = 0x01170000 //It is delivered with a dependent Variable value when a dominant Variable has changed but the change has not been applied.
	Status_GoodEdited_DominantValueChanged_DependentValueChanged           StatusCode = 0x01180000 //It is delivered with a dependent Variable value when a dominant or dependent Variable has changed but change has not been applied.
	Status_BadEdited_OutOfRange                                            StatusCode = 0x81190000 //It is delivered with a Variable value when Variable has changed but the value is not legal.
	Status_BadInitialValue_OutOfRange                                      StatusCode = 0x811A0000 //It is delivered with a Variable value when a source Variable has changed but the value is not legal.
	Status_BadOutOfRange_DominantValueChanged                              StatusCode = 0x811B0000 //It is delivered with a dependent Variable value when a dominant Variable has changed and the value is not legal.
	Status_BadEdited_OutOfRange_DominantValueChanged                       StatusCode = 0x811C0000 //It is delivered with a dependent Variable value when a dominant Variable has changed, the value is not legal and the change has not been applied.
	Status_BadOutOfRange_DominantValueChanged_DependentValueChanged        StatusCode = 0x811D0000 //It is delivered with a dependent Variable value when a dominant or dependent Variable has changed and the value is not legal.
	Status_BadEdited_OutOfRange_DominantValueChanged_DependentValueChanged StatusCode = 0x811E0000 //It is delivered with a dependent Variable value when a dominant or dependent Variable has changed, the value is not legal and the change has not been applied.
	Status_GoodCommunicationEvent                                          StatusCode = 0x00A70000 //The communication layer has raised an event.
	Status_GoodShutdownEvent                                               StatusCode = 0x00A80000 //The system is shutting down.
	Status_GoodCallAgain                                                   StatusCode = 0x00A90000 //The operation is not finished and needs to be called again.
	Status_GoodNonCriticalTimeout                                          StatusCode = 0x00AA0000 //A non-critical timeout occurred.
	Status_BadInvalidArgument                                              StatusCode = 0x80AB0000 //One or more arguments are invalid.
	Status_BadConnectionRejected                                           StatusCode = 0x80AC0000 //Could not establish a network connection to remote server.
	Status_BadDisconnect                                                   StatusCode = 0x80AD0000 //The server has disconnected from the client.
	Status_BadConnectionClosed                                             StatusCode = 0x80AE0000 //The network connection has been closed.
	Status_BadInvalidState                                                 StatusCode = 0x80AF0000 //The operation cannot be completed because the object is closed, uninitialized or in some other invalid state.
	Status_BadEndOfStream                                                  StatusCode = 0x80B00000 //Cannot move beyond end of the stream.
	Status_BadNoDataAvailable                                              StatusCode = 0x80B10000 //No data is currently available for reading from a non-blocking stream.
	Status_BadWaitingForResponse                                           StatusCode = 0x80B20000 //The asynchronous operation is waiting for a response.
	Status_BadOperationAbandoned                                           StatusCode = 0x80B30000 //The asynchronous operation was abandoned by the caller.
	Status_BadExpectedStreamToBlock                                        StatusCode = 0x80B40000 //The stream did not return all data requested (possibly because it is a non-blocking stream).
	Status_BadWouldBlock                                                   StatusCode = 0x80B50000 //Non blocking behaviour is required and the operation would block.
	Status_BadSyntaxError                                                  StatusCode = 0x80B60000 //A value had an invalid syntax.
	Status_BadMaxConnectionsReached                                        StatusCode = 0x80B70000 //The operation could not be finished because all available connections are in use.
	Status_UncertainTransducerInManual                                     StatusCode = 0x42080000 //The value may not be accurate because the transducer is in manual mode.
	Status_UncertainSimulatedValue                                         StatusCode = 0x42090000 //The value is simulated.
	Status_UncertainSensorCalibration                                      StatusCode = 0x420A0000 //The value may not be accurate due to a sensor calibration fault.
	Status_UncertainConfigurationError                                     StatusCode = 0x420F0000 //The value may not be accurate due to a configuration issue.
	Status_GoodCascadeInitializationAcknowledged                           StatusCode = 0x04010000 //The value source supports cascade handshaking and the value has been Initialized based on an initialization request from a cascade secondary.
	Status_GoodCascadeInitializationRequest                                StatusCode = 0x04020000 //The value source supports cascade handshaking and is requesting initialization of a cascade primary.
	Status_GoodCascadeNotInvited                                           StatusCode = 0x04030000 //The value source supports cascade handshaking, however, the source’s current state does not allow for cascade.
	Status_GoodCascadeNotSelected                                          StatusCode = 0x04040000 //The value source supports cascade handshaking, however, the source has not selected the corresponding cascade primary for use.
	Status_GoodFaultStateActive                                            StatusCode = 0x04070000 //There is a fault state condition active in the value source.
	Status_GoodInitiateFaultState                                          StatusCode = 0x04080000 //A fault state condition is being requested of the destination.
	Status_GoodCascade                                                     StatusCode = 0x04090000 //The value is accurate, and the signal source supports cascade handshaking.
)

func (status StatusCode) ToSymbolicId() string {
	return statusSymbolicIdMapping[status]
}

var statusSymbolicIdMapping = map[StatusCode]string{
	0x00000000: "Good",
	0x40000000: "Uncertain",
	0x80000000: "Bad",
	0x80010000: "BadUnexpectedError",
	0x80020000: "BadInternalError",
	0x80030000: "BadOutOfMemory",
	0x80040000: "BadResourceUnavailable",
	0x80050000: "BadCommunicationError",
	0x80060000: "BadEncodingError",
	0x80070000: "BadDecodingError",
	0x80080000: "BadEncodingLimitsExceeded",
	0x80B80000: "BadRequestTooLarge",
	0x80B90000: "BadResponseTooLarge",
	0x80090000: "BadUnknownResponse",
	0x800A0000: "BadTimeout",
	0x800B0000: "BadServiceUnsupported",
	0x800C0000: "BadShutdown",
	0x800D0000: "BadServerNotConnected",
	0x800E0000: "BadServerHalted",
	0x800F0000: "BadNothingToDo",
	0x80100000: "BadTooManyOperations",
	0x80DB0000: "BadTooManyMonitoredItems",
	0x80110000: "BadDataTypeIdUnknown",
	0x80120000: "BadCertificateInvalid",
	0x80130000: "BadSecurityChecksFailed",
	0x81140000: "BadCertificatePolicyCheckFailed",
	0x80140000: "BadCertificateTimeInvalid",
	0x80150000: "BadCertificateIssuerTimeInvalid",
	0x80160000: "BadCertificateHostNameInvalid",
	0x80170000: "BadCertificateUriInvalid",
	0x80180000: "BadCertificateUseNotAllowed",
	0x80190000: "BadCertificateIssuerUseNotAllowed",
	0x801A0000: "BadCertificateUntrusted",
	0x801B0000: "BadCertificateRevocationUnknown",
	0x801C0000: "BadCertificateIssuerRevocationUnknown",
	0x801D0000: "BadCertificateRevoked",
	0x801E0000: "BadCertificateIssuerRevoked",
	0x810D0000: "BadCertificateChainIncomplete",
	0x801F0000: "BadUserAccessDenied",
	0x80200000: "BadIdentityTokenInvalid",
	0x80210000: "BadIdentityTokenRejected",
	0x80220000: "BadSecureChannelIdInvalid",
	0x80230000: "BadInvalidTimestamp",
	0x80240000: "BadNonceInvalid",
	0x80250000: "BadSessionIdInvalid",
	0x80260000: "BadSessionClosed",
	0x80270000: "BadSessionNotActivated",
	0x80280000: "BadSubscriptionIdInvalid",
	0x802A0000: "BadRequestHeaderInvalid",
	0x802B0000: "BadTimestampsToReturnInvalid",
	0x802C0000: "BadRequestCancelledByClient",
	0x80E50000: "BadTooManyArguments",
	0x810E0000: "BadLicenseExpired",
	0x810F0000: "BadLicenseLimitsExceeded",
	0x81100000: "BadLicenseNotAvailable",
	0x002D0000: "GoodSubscriptionTransferred",
	0x002E0000: "GoodCompletesAsynchronously",
	0x002F0000: "GoodOverload",
	0x00300000: "GoodClamped",
	0x80310000: "BadNoCommunication",
	0x80320000: "BadWaitingForInitialData",
	0x80330000: "BadNodeIdInvalid",
	0x80340000: "BadNodeIdUnknown",
	0x80350000: "BadAttributeIdInvalid",
	0x80360000: "BadIndexRangeInvalid",
	0x80370000: "BadIndexRangeNoData",
	0x80380000: "BadDataEncodingInvalid",
	0x80390000: "BadDataEncodingUnsupported",
	0x803A0000: "BadNotReadable",
	0x803B0000: "BadNotWritable",
	0x803C0000: "BadOutOfRange",
	0x803D0000: "BadNotSupported",
	0x803E0000: "BadNotFound",
	0x803F0000: "BadObjectDeleted",
	0x80400000: "BadNotImplemented",
	0x80410000: "BadMonitoringModeInvalid",
	0x80420000: "BadMonitoredItemIdInvalid",
	0x80430000: "BadMonitoredItemFilterInvalid",
	0x80440000: "BadMonitoredItemFilterUnsupported",
	0x80450000: "BadFilterNotAllowed",
	0x80460000: "BadStructureMissing",
	0x80470000: "BadEventFilterInvalid",
	0x80480000: "BadContentFilterInvalid",
	0x80C10000: "BadFilterOperatorInvalid",
	0x80C20000: "BadFilterOperatorUnsupported",
	0x80C30000: "BadFilterOperandCountMismatch",
	0x80490000: "BadFilterOperandInvalid",
	0x80C40000: "BadFilterElementInvalid",
	0x80C50000: "BadFilterLiteralInvalid",
	0x804A0000: "BadContinuationPointInvalid",
	0x804B0000: "BadNoContinuationPoints",
	0x804C0000: "BadReferenceTypeIdInvalid",
	0x804D0000: "BadBrowseDirectionInvalid",
	0x804E0000: "BadNodeNotInView",
	0x81120000: "BadNumericOverflow",
	0x804F0000: "BadServerUriInvalid",
	0x80500000: "BadServerNameMissing",
	0x80510000: "BadDiscoveryUrlMissing",
	0x80520000: "BadSempahoreFileMissing",
	0x80530000: "BadRequestTypeInvalid",
	0x80540000: "BadSecurityModeRejected",
	0x80550000: "BadSecurityPolicyRejected",
	0x80560000: "BadTooManySessions",
	0x80570000: "BadUserSignatureInvalid",
	0x80580000: "BadApplicationSignatureInvalid",
	0x80590000: "BadNoValidCertificates",
	0x80C60000: "BadIdentityChangeNotSupported",
	0x805A0000: "BadRequestCancelledByRequest",
	0x805B0000: "BadParentNodeIdInvalid",
	0x805C0000: "BadReferenceNotAllowed",
	0x805D0000: "BadNodeIdRejected",
	0x805E0000: "BadNodeIdExists",
	0x805F0000: "BadNodeClassInvalid",
	0x80600000: "BadBrowseNameInvalid",
	0x80610000: "BadBrowseNameDuplicated",
	0x80620000: "BadNodeAttributesInvalid",
	0x80630000: "BadTypeDefinitionInvalid",
	0x80640000: "BadSourceNodeIdInvalid",
	0x80650000: "BadTargetNodeIdInvalid",
	0x80660000: "BadDuplicateReferenceNotAllowed",
	0x80670000: "BadInvalidSelfReference",
	0x80680000: "BadReferenceLocalOnly",
	0x80690000: "BadNoDeleteRights",
	0x40BC0000: "UncertainReferenceNotDeleted",
	0x806A0000: "BadServerIndexInvalid",
	0x806B0000: "BadViewIdUnknown",
	0x80C90000: "BadViewTimestampInvalid",
	0x80CA0000: "BadViewParameterMismatch",
	0x80CB0000: "BadViewVersionInvalid",
	0x40C00000: "UncertainNotAllNodesAvailable",
	0x00BA0000: "GoodResultsMayBeIncomplete",
	0x80C80000: "BadNotTypeDefinition",
	0x406C0000: "UncertainReferenceOutOfServer",
	0x806D0000: "BadTooManyMatches",
	0x806E0000: "BadQueryTooComplex",
	0x806F0000: "BadNoMatch",
	0x80700000: "BadMaxAgeInvalid",
	0x80E60000: "BadSecurityModeInsufficient",
	0x80710000: "BadHistoryOperationInvalid",
	0x80720000: "BadHistoryOperationUnsupported",
	0x80BD0000: "BadInvalidTimestampArgument",
	0x80730000: "BadWriteNotSupported",
	0x80740000: "BadTypeMismatch",
	0x80750000: "BadMethodInvalid",
	0x80760000: "BadArgumentsMissing",
	0x81110000: "BadNotExecutable",
	0x80770000: "BadTooManySubscriptions",
	0x80780000: "BadTooManyPublishRequests",
	0x80790000: "BadNoSubscription",
	0x807A0000: "BadSequenceNumberUnknown",
	0x00DF0000: "GoodRetransmissionQueueNotSupported",
	0x807B0000: "BadMessageNotAvailable",
	0x807C0000: "BadInsufficientClientProfile",
	0x80BF0000: "BadStateNotActive",
	0x81150000: "BadAlreadyExists",
	0x807D0000: "BadTcpServerTooBusy",
	0x807E0000: "BadTcpMessageTypeInvalid",
	0x807F0000: "BadTcpSecureChannelUnknown",
	0x80800000: "BadTcpMessageTooLarge",
	0x80810000: "BadTcpNotEnoughResources",
	0x80820000: "BadTcpInternalError",
	0x80830000: "BadTcpEndpointUrlInvalid",
	0x80840000: "BadRequestInterrupted",
	0x80850000: "BadRequestTimeout",
	0x80860000: "BadSecureChannelClosed",
	0x80870000: "BadSecureChannelTokenUnknown",
	0x80880000: "BadSequenceNumberInvalid",
	0x80BE0000: "BadProtocolVersionUnsupported",
	0x80890000: "BadConfigurationError",
	0x808A0000: "BadNotConnected",
	0x808B0000: "BadDeviceFailure",
	0x808C0000: "BadSensorFailure",
	0x808D0000: "BadOutOfService",
	0x808E0000: "BadDeadbandFilterInvalid",
	0x408F0000: "UncertainNoCommunicationLastUsableValue",
	0x40900000: "UncertainLastUsableValue",
	0x40910000: "UncertainSubstituteValue",
	0x40920000: "UncertainInitialValue",
	0x40930000: "UncertainSensorNotAccurate",
	0x40940000: "UncertainEngineeringUnitsExceeded",
	0x40950000: "UncertainSubNormal",
	0x00960000: "GoodLocalOverride",
	0x80970000: "BadRefreshInProgress",
	0x80980000: "BadConditionAlreadyDisabled",
	0x80CC0000: "BadConditionAlreadyEnabled",
	0x80990000: "BadConditionDisabled",
	0x809A0000: "BadEventIdUnknown",
	0x80BB0000: "BadEventNotAcknowledgeable",
	0x80CD0000: "BadDialogNotActive",
	0x80CE0000: "BadDialogResponseInvalid",
	0x80CF0000: "BadConditionBranchAlreadyAcked",
	0x80D00000: "BadConditionBranchAlreadyConfirmed",
	0x80D10000: "BadConditionAlreadyShelved",
	0x80D20000: "BadConditionNotShelved",
	0x80D30000: "BadShelvingTimeOutOfRange",
	0x809B0000: "BadNoData",
	0x80D70000: "BadBoundNotFound",
	0x80D80000: "BadBoundNotSupported",
	0x809D0000: "BadDataLost",
	0x809E0000: "BadDataUnavailable",
	0x809F0000: "BadEntryExists",
	0x80A00000: "BadNoEntryExists",
	0x80A10000: "BadTimestampNotSupported",
	0x00A20000: "GoodEntryInserted",
	0x00A30000: "GoodEntryReplaced",
	0x40A40000: "UncertainDataSubNormal",
	0x00A50000: "GoodNoData",
	0x00A60000: "GoodMoreData",
	0x80D40000: "BadAggregateListMismatch",
	0x80D50000: "BadAggregateNotSupported",
	0x80D60000: "BadAggregateInvalidInputs",
	0x80DA0000: "BadAggregateConfigurationRejected",
	0x00D90000: "GoodDataIgnored",
	0x80E40000: "BadRequestNotAllowed",
	0x81130000: "BadRequestNotComplete",
	0x811F0000: "BadTicketRequired",
	0x81200000: "BadTicketInvalid",
	0x00DC0000: "GoodEdited",
	0x00DD0000: "GoodPostActionFailed",
	0x40DE0000: "UncertainDominantValueChanged",
	0x00E00000: "GoodDependentValueChanged",
	0x80E10000: "BadDominantValueChanged",
	0x40E20000: "UncertainDependentValueChanged",
	0x80E30000: "BadDependentValueChanged",
	0x01160000: "GoodEdited_DependentValueChanged",
	0x01170000: "GoodEdited_DominantValueChanged",
	0x01180000: "GoodEdited_DominantValueChanged_DependentValueChanged",
	0x81190000: "BadEdited_OutOfRange",
	0x811A0000: "BadInitialValue_OutOfRange",
	0x811B0000: "BadOutOfRange_DominantValueChanged",
	0x811C0000: "BadEdited_OutOfRange_DominantValueChanged",
	0x811D0000: "BadOutOfRange_DominantValueChanged_DependentValueChanged",
	0x811E0000: "BadEdited_OutOfRange_DominantValueChanged_DependentValueChanged",
	0x00A70000: "GoodCommunicationEvent",
	0x00A80000: "GoodShutdownEvent",
	0x00A90000: "GoodCallAgain",
	0x00AA0000: "GoodNonCriticalTimeout",
	0x80AB0000: "BadInvalidArgument",
	0x80AC0000: "BadConnectionRejected",
	0x80AD0000: "BadDisconnect",
	0x80AE0000: "BadConnectionClosed",
	0x80AF0000: "BadInvalidState",
	0x80B00000: "BadEndOfStream",
	0x80B10000: "BadNoDataAvailable",
	0x80B20000: "BadWaitingForResponse",
	0x80B30000: "BadOperationAbandoned",
	0x80B40000: "BadExpectedStreamToBlock",
	0x80B50000: "BadWouldBlock",
	0x80B60000: "BadSyntaxError",
	0x80B70000: "BadMaxConnectionsReached",
	0x42080000: "UncertainTransducerInManual",
	0x42090000: "UncertainSimulatedValue",
	0x420A0000: "UncertainSensorCalibration",
	0x420F0000: "UncertainConfigurationError",
	0x04010000: "GoodCascadeInitializationAcknowledged",
	0x04020000: "GoodCascadeInitializationRequest",
	0x04030000: "GoodCascadeNotInvited",
	0x04040000: "GoodCascadeNotSelected",
	0x04070000: "GoodFaultStateActive",
	0x04080000: "GoodInitiateFaultState",
	0x04090000: "GoodCascade",
}
