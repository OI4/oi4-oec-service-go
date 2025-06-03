package opc

import (
	"sync"
	"testing"
)

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
