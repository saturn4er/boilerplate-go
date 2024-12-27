package txoutboxwatermill

import (
	"context"
	"fmt"

	millmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/saturn4er/boilerplate-go/lib/txoutbox"
)

type MessagesSender struct {
	topicPublishers map[string]millmessage.Publisher
}

var _ txoutbox.MessageSender = new(MessagesSender)

func (m MessagesSender) SendMessage(ctx context.Context, message *txoutbox.Message) error {
	publisher, ok := m.topicPublishers[message.Topic]
	if !ok {
		return fmt.Errorf("no publisher for topic %s", message.Topic)
	}
	watermillMessage := millmessage.NewMessage(uuid.New().String(), message.Data)
	watermillMessage.Metadata["idempotency_key"] = message.IdempotencyKey
	watermillMessage.Metadata["ordering_key"] = message.OrderingKey

	return publisher.Publish(message.Topic)
}
