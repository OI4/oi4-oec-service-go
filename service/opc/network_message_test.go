package opc

import (
	"encoding/json"
	"github.com/OI4/oi4-oec-service-go/service/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestCreateNetworkMessage(t *testing.T) {
	file, err := os.ReadFile("./testdata/networkMessage.json")
	require.NoError(t, err)

	var expected api.NetworkMessage
	err = json.Unmarshal(file, &expected)
	require.NoError(t, err)

	var serviceType *api.ServiceType
	var applicationOi4Identifier *api.Oi4Identifier

	parts := strings.SplitN(expected.PublisherId, "/", 2)
	assert.Len(t, parts, 2)
	serviceType, err = api.ParseServiceType(parts[0])
	require.NoError(t, err)

	applicationOi4Identifier, err = api.ParseOi4Identifier(parts[1], true)
	require.NoError(t, err)

	content := make([]api.PublicationContent, len(expected.Messages))
	for i, msg := range expected.Messages {
		content[i] = api.PublicationContent{
			Data: msg,
		}
	}

	publication := api.PublicationMessage{
		Resource:      api.ResourceMam,
		Source:        applicationOi4Identifier,
		CorrelationId: expected.CorrelationId,
		Content:       content,
	}

	actual := CreateNetworkMessage(applicationOi4Identifier, *serviceType, publication)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.PublisherId, actual.PublisherId)
	assert.Equal(t, expected.CorrelationId, actual.CorrelationId)
	assert.Equal(t, len(expected.Messages), len(actual.Messages))

	getTail := func(s string) string {
		strings.SplitN(s, "s", 2)
		assert.Len(t, parts, 2)
		return parts[1]
	}

	assert.Equal(t, getTail(expected.MessageId), getTail(actual.MessageId))

	// Check with empty content
	publication.Content = nil
	actual = CreateNetworkMessage(applicationOi4Identifier, *serviceType, publication)
	assert.Nil(t, actual)

	publication.Content = []api.PublicationContent{}
	actual = CreateNetworkMessage(applicationOi4Identifier, *serviceType, publication)
	assert.Nil(t, actual)
}

func TestMessageIDIsUnique(t *testing.T) {
	var wg sync.WaitGroup
	messageIDs := make(map[string]bool)
	mu := sync.Mutex{}

	publisherID := "publisher1"

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			messageID := GetMessageID(publisherID)
			mu.Lock()
			defer mu.Unlock()
			if messageIDs[messageID] {
				t.Errorf("%d - Duplicate message ID generated: %s", i, messageID)
			}
			messageIDs[messageID] = true
		}()
	}
	wg.Wait()
}

func TestMessageIDChangesWithPublisherID(t *testing.T) {
	id1 := GetMessageID("publisher1")
	id2 := GetMessageID("publisher2")
	if id1 == id2 {
		t.Errorf("Message ID should differ for different publisher IDs")
	}
}

func TestMessageIDDoesNotRepeat(t *testing.T) {
	id1 := GetMessageID("publisher1")
	id2 := GetMessageID("publisher1")
	if id1 == id2 {
		t.Errorf("Message ID should not repeat for the same publisher ID")
	}
}
