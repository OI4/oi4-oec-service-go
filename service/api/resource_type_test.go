package api

import (
	"testing"
)

func TestParseResourceType_ValidResourceTypes(t *testing.T) {
	validResourceTypes := []string{"MAM", "Health", "Config", "License", "LicenseText", "RtLicense", "Data", "Metadata", "Event", "Profile", "PublicationList", "SubscriptionList", "Interfaces", "ReferenceDesignation"}

	for _, resourceType := range validResourceTypes {
		_, err := ParseResourceType(resourceType)
		if err != nil {
			t.Errorf("Expected no error for valid resource type %s, got %v", resourceType, err)
		}
	}
}

func TestParseResourceType_InvalidResourceType(t *testing.T) {
	_, err := ParseResourceType("InvalidResourceType")
	if err == nil {
		t.Errorf("Expected error for invalid resource type, got nil")
	}
}

func TestParseResourceType_EmptyString(t *testing.T) {
	_, err := ParseResourceType("")
	if err == nil {
		t.Errorf("Expected error for empty string, got nil")
	}
}

func TestResourceType_ToDataSetClassId(t *testing.T) {
	rType := ResourceMam
	expected := "360ca8f3-5e66-42a2-8f10-9cdf45f4bf58"
	result := rType.ToDataSetClassId()
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestResourceType_ToDataSetClassId_InvalidResourceType(t *testing.T) {
	rType := ResourceType("InvalidResourceType")
	expected := ""
	result := rType.ToDataSetClassId()
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
