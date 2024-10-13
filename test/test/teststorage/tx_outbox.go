package teststorage

import (
	"github.com/saturn4er/boilerplate-go/lib/txoutbox"
	testsvc "github.com/saturn4er/boilerplate-go/test/test/testservice"
)

func buildPasswordRecoveryEventMessage(t *testsvc.PasswordRecoveryEvent) (*txoutbox.Message, error) {
	return &txoutbox.Message{
		ID:             1,
		Topic:          "some_topic",
		OrderingKey:    "1",
		IdempotencyKey: "1",
		Data:           nil,
	}, nil
}
