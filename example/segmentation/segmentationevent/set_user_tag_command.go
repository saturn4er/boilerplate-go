package segmentationevent

import (
	"encoding/json"
	"fmt"

	segmentationsvc "github.com/saturn4er/boilerplate-go/example/segmentation/segmentationservice"
	"github.com/saturn4er/boilerplate-go/lib/txoutbox"
)

func BuildSetUserTagCommandMessage(cmd *segmentationsvc.SetUserTagCommand) (*txoutbox.Message, error) {
	data, err := json.Marshal(cmd.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal command data: %w", err)
	}

	return &txoutbox.Message{
		Topic:          "segmentation.set-user-tag-command",
		OrderingKey:    cmd.Data.UserID.String(),
		IdempotencyKey: cmd.IdempotencyKey,
		Data:           data,
	}, nil
}
