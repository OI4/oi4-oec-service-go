package topic

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

var appId = api.NewOi4Identifier("acme.com", "FBC", "fbc%183z", "FBC#123")
var source = api.NewOi4Identifier("acme.com", "matches", "m/42-A", "F234#862")
var serviceType = api.ServiceTypeOTConnector
var method = api.MethodGet
var resource = api.ResourceMam

func TestTopicToStringWithAllFieldsSet(t *testing.T) {
	category := "Category"
	filter := api.NewFilter("Filter")

	topic := NewTopic(serviceType, *appId, method, resource, source, &category, filter)
	expected := "Oi4/OTConnector/acme.com/FBC/fbc,25183z/FBC,23123/Get/MAM/acme.com/matches/m,2F42-A/F234,23862/Category/Filter"
	assert.Equal(t, expected, topic.ToString())
}

func TestTopicToStringWithOnlyRequiredFieldsSet(t *testing.T) {
	topic := NewTopic(serviceType, *appId, method, resource, nil, nil, nil)
	expected := "Oi4/OTConnector/acme.com/FBC/fbc,25183z/FBC,23123/Get/MAM"
	assert.Equal(t, expected, topic.ToString())
}

func TestTopicToStringWithSomeOptionalFieldsSet(t *testing.T) {
	topic := NewTopic(serviceType, *appId, method, resource, source, nil, nil)
	expected := "Oi4/OTConnector/acme.com/FBC/fbc,25183z/FBC,23123/Get/MAM/acme.com/matches/m,2F42-A/F234,23862"
	assert.Equal(t, expected, topic.ToString())

}

func TestParseTopicWithValidTopic(t *testing.T) {
	topic := "Oi4/OTConnector/acme.com/FBC/fbc,25183z/FBC,23123/Get/MAM/acme.com/matches/m,2F42-A/F234,23862/Category/Filter"
	result, err := ParseTopic(topic)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, serviceType, result.ServiceType)
	assert.Equal(t, *appId, result.Oi4Identifier)
	assert.Equal(t, method, result.Method)
	assert.Equal(t, resource, result.Resource)
	assert.Equal(t, source, result.Source)
	assert.Equal(t, "Category", *result.Category)
	assert.Equal(t, "Filter", result.Filter.String())

	result, err = ParseTopic("Oi4/OTConnector/acme.com/SampleApplication/FC,23156/F87263976,234/Get/MAM/acme.com/SampleApplication/FC,23156/F87263976,234")
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func TestParseTopicWithEmptyTopic(t *testing.T) {
	topic := ""
	_, err := ParseTopic(topic)
	assert.NotNil(t, err)
	assert.Equal(t, "topic is empty", err.Error())
}

func TestParseTopicWithInvalidNamespace(t *testing.T) {
	topic := "Invalid/OTConnector/acme.com/FBC/fbc,25183z/FBC,23123/Get/MAM/acme.com/matches/m,2F42-A/F234,23862/Category/Filter"
	_, err := ParseTopic(topic)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid topic, wrong namespace", err.Error())
}

func TestParseTopicWithFewParts(t *testing.T) {
	topic := "Oi4/OTConnector/acme.com/FBC/fbc%183z/FBC#123"
	_, err := ParseTopic(topic)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid topic, to few parts", err.Error())
}

func TestParseTopicWithInvalidServiceType(t *testing.T) {
	topic := "Oi4/InvalidServiceType/acme.com/FBC/fbc,25183z/FBC,23123/Get/MAM/acme.com/matches/m,2F42-A/F234,23862/Category/Filter"
	_, err := ParseTopic(topic)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid service type: cannot parse:[InvalidServiceType] as ServiceType", err.Error())
}

func TestParseTopicWithInvalidOi4Identifier(t *testing.T) {
	topic := "Oi4/OTConnector/invalid/identifier/here/,,/Get/MAM/acme.com/matches/m,2F42-A/F234,23862/Category/Filter"
	_, err := ParseTopic(topic)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid oi4 identifier: invalid serial number: invalid DNP escape \",,\"", err.Error())
}

func TestParseTopicWithInvalidMethodType(t *testing.T) {
	topic := "Oi4/OTConnector/acme.com/FBC/fbc,25183z/FBC,23123/InvalidMethod/MAM/acme.com/matches/m,2F42-A/F234,23862/Category/Filter"
	_, err := ParseTopic(topic)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid method type: cannot parse:[InvalidMethod] as MethodType", err.Error())
}

func TestParseTopicWithInvalidResourceType(t *testing.T) {
	topic := "Oi4/OTConnector/acme.com/FBC/fbc,25183z/FBC,23123/Get/InvalidResource/acme.com/matches/m,2F42-A/F234,23862/Category/Filter"
	_, err := ParseTopic(topic)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid resource type: cannot parse:[InvalidResource] as ResourceType", err.Error())
}

func TestParseTopicWithInvalidSource(t *testing.T) {
	topic := "Oi4/OTConnector/acme.com/FBC/fbc,25183z/FBC,23123/Get/MAM/invalid/source/here/,,/Category/Filter"
	_, err := ParseTopic(topic)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid source: invalid serial number: invalid DNP escape \",,\"", err.Error())
}
