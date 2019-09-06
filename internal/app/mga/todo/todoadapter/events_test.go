package todoadapter

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message/subscriber"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
)

func TestEventDispatcher_MarkedAsDone(t *testing.T) {
	publisher := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})
	const topic = "todo"
	eventBus, _ := cqrs.NewEventBus(publisher, func(_ string) string { return topic }, &cqrs.JSONMarshaler{})

	messages, err := publisher.Subscribe(context.Background(), topic)
	require.NoError(t, err)

	events := NewEventDispatcher(eventBus)

	event := todo.MarkedAsDone{
		ID: "id",
	}

	err = events.MarkedAsDone(context.Background(), event)
	require.NoError(t, err)

	received, all := subscriber.BulkRead(messages, 1, time.Second)
	if !all {
		t.Fatal("no message received")
	}

	var receivedEvent todo.MarkedAsDone

	err = json.Unmarshal(received[0].Payload, &receivedEvent)
	require.NoError(t, err)

	assert.Equal(t, event, receivedEvent)
}
