package topic

import (
	"errors"
	"fmt"
	oi4 "github.com/OI4/oi4-oec-service-go/api/pkg/types"
	"strings"
)

const Oi4Namespace = "Oi4"

type Error struct {
	Message string
	Err     error
}

func (w *Error) Error() string {
	return fmt.Sprintf("%s: %v", w.Message, w.Err)
}

type Topic struct {
	ServiceType   oi4.ServiceType
	Oi4Identifier oi4.Oi4Identifier
	Method        oi4.MethodType
	Resource      oi4.ResourceType
	Source        *oi4.Oi4Identifier
	Category      *string
	Filter        *string
}

func NewTopic(serviceType oi4.ServiceType, oi4Identifier oi4.Oi4Identifier, method oi4.MethodType, resource oi4.ResourceType, source *oi4.Oi4Identifier, category *string, filter *string) Topic {
	return Topic{
		ServiceType:   serviceType,
		Oi4Identifier: oi4Identifier,
		Method:        method,
		Resource:      resource,
		Source:        source,
		Category:      category,
		Filter:        filter,
	}
}

func ParseTopic(topic string) (*Topic, error) {
	parts := strings.Split(topic, "/")
	if len(parts) < 8 {
		return nil, errors.New("invalid topic, to few parts")
	}
	if parts[0] != Oi4Namespace {
		return nil, errors.New("invalid topic, wrong namespace")
	}
	serviceType, err := oi4.ParseServiceType(parts[1])
	if err != nil {
		return nil, &Error{
			Message: "invalid service type",
			Err:     err,
		}
	}

	oi4Identifier, err := oi4.ParseOi4IdentifierFromArray(parts[2:6], true)
	if err != nil {
		return nil, &Error{
			Message: "invalid oi4 identifier",
			Err:     err,
		}
	}

	method, err := oi4.ParseMethodType(parts[6])
	if err != nil {
		return nil, &Error{
			Message: "invalid method type",
			Err:     err,
		}
	}

	resource, err := oi4.ParseResourceType(parts[7])
	if err != nil {
		return nil, &Error{
			Message: "invalid resource type",
			Err:     err,
		}
	}

	var source *oi4.Oi4Identifier
	if len(parts) >= 12 {
		source, err = oi4.ParseOi4IdentifierFromArray(parts[8:12], true)
		if err != nil {
			return nil, &Error{
				Message: "invalid source",
				Err:     err,
			}
		}
	}

	var category *string
	if len(parts) >= 13 {
		category = &parts[12]
	}

	var filter *string
	if len(parts) >= 14 {
		filter = &parts[13]
	}

	result := NewTopic(*serviceType, *oi4Identifier, *method, *resource, source, category, filter)
	return &result, nil
}

func (t *Topic) ToString() string {
	topic := fmt.Sprintf("%s/%s/%s/%s/%s", Oi4Namespace, t.ServiceType, t.Oi4Identifier.ToString(), t.Method, t.Resource)
	if t.Source != nil {
		topic = fmt.Sprintf("%s/%s", topic, t.Source.ToString())
	}
	if t.Category != nil {
		topic = fmt.Sprintf("%s/%s", topic, *t.Category)
	}
	if t.Filter != nil {
		topic = fmt.Sprintf("%s/%s", topic, *t.Filter)
	}

	return topic
}
