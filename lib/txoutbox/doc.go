// Package txoutbox implements the transactional outbox pattern for reliable
// event delivery in event-driven architectures.
//
// The outbox pattern ensures that domain operations and event publishing are
// atomic: events are persisted to the same database transaction as the domain
// changes, then delivered asynchronously by a background processor.
//
// # Core Types
//
// [Outbox] is the generic interface for sending events:
//
//	type Outbox[Entity any] interface {
//	    Send(ctx context.Context, model *Entity) error
//	}
//
// [GormStorage] implements Outbox by converting domain events into [Message]
// records stored in the tx_outbox.messages table. A BuildMessage function
// transforms the domain type into the wire format.
//
// [MessagesSender] is a background job (implementing jobber.Job) that polls
// for unsent messages, delivers them via a [MessageSender], and deletes them
// on success. It uses SELECT FOR UPDATE to handle concurrent processors safely.
//
// # Message Processing
//
// [MessageProcessor] functions run before a message is persisted, allowing
// enrichment of metadata (e.g., adding trace context or ordering keys).
//
// # Usage with Generated Code
//
// The boilerplate-go generator produces Outbox accessors on Storage interfaces:
//
//	err := tx.SetUserTagCommands().Send(ctx, &segmentationsvc.SetUserTagCommand{
//	    ID:             uuid.New(),
//	    Data:           data,
//	    IdempotencyKey: user.ID.String(),
//	})
//
// Messages are committed with the transaction and delivered asynchronously.
//
// # Watermill Integration
//
// The txoutboxwatermill sub-package bridges outbox messages to Watermill
// publishers for integration with message brokers (Kafka, NATS, Google Pub/Sub, etc.).
package txoutbox
