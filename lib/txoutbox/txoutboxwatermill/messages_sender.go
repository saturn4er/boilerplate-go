package txoutboxwatermill

import (
	"context"
	"fmt"

	millmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"

	"github.com/saturn4er/boilerplate-go/lib/txoutbox"
)

const OrderingKetMetadataKey = "ordering_key"
const IdempotencyKeyMetadataKey = "idempotency_key"

type MessagesSender struct {
	topicPublishers  map[string]millmessage.Publisher
	defaultPublisher millmessage.Publisher
}

var _ txoutbox.MessageSender = new(MessagesSender)

func NewMessagesSender(
	topicPublishers map[string]millmessage.Publisher,
	defaultPublisher millmessage.Publisher,
) *MessagesSender {
	return &MessagesSender{
		topicPublishers:  topicPublishers,
		defaultPublisher: defaultPublisher,
	}
}

func (m MessagesSender) SendMessage(ctx context.Context, message *txoutbox.Message) error {
	publisher, ok := m.topicPublishers[message.Topic]
	if !ok {
		if m.defaultPublisher == nil {
			return fmt.Errorf("no publisher for topic %s", message.Topic)
		}

		publisher = m.defaultPublisher
	}
	watermillMessage := millmessage.NewMessage(uuid.New().String(), message.Data)
	watermillMessage.Metadata[IdempotencyKeyMetadataKey] = message.IdempotencyKey
	watermillMessage.Metadata[OrderingKetMetadataKey] = message.OrderingKey

	return publisher.Publish(message.Topic, watermillMessage)
}
