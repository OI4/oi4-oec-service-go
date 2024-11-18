package topic

import (
	"errors"
	"fmt"
	"github.com/OI4/oi4-oec-service-go/service/api"
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
	ServiceType   api.ServiceType
	Oi4Identifier api.Oi4Identifier
	Method        api.MethodType
	Resource      api.ResourceType
	Source        *api.Oi4Identifier
	Category      *string
	Filter        api.Filter
}

func NewTopic(serviceType api.ServiceType, oi4Identifier api.Oi4Identifier, method api.MethodType, resource api.ResourceType, source *api.Oi4Identifier, category *string, filter api.Filter) Topic {
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
	serviceType, err := api.ParseServiceType(parts[1])
	if err != nil {
		return nil, &Error{
			Message: "invalid service type",
			Err:     err,
		}
	}

	oi4Identifier, err := api.ParseOi4IdentifierFromArray(parts[2:6], true)
	if err != nil {
		return nil, &Error{
			Message: "invalid oi4 identifier",
			Err:     err,
		}
	}

	method, err := api.ParseMethodType(parts[6])
	if err != nil {
		return nil, &Error{
			Message: "invalid method type",
			Err:     err,
		}
	}

	resource, err := api.ParseResourceType(parts[7])
	if err != nil {
		return nil, &Error{
			Message: "invalid resource type",
			Err:     err,
		}
	}

	var source *api.Oi4Identifier
	if len(parts) >= 12 {
		source, err = api.ParseOi4IdentifierFromArray(parts[8:12], true)
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

	var filter api.Filter
	if len(parts) >= 14 {
		filter = api.NewStringFilter(parts[13])
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
		topic = fmt.Sprintf("%s/%s", topic, t.Filter)
	}

	return topic
}
