package segmentationservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/saturn4er/boilerplate-go/lib/filter"
)

type Service struct {
	Storage Storage
}

func (s *Service) HandleSetUserTagCommand(ctx context.Context, cmd *SetUserTagCommand) error {
	_, err := s.Storage.UserTags().FirstOrCreate(ctx,
		&UserTagFilter{
			UserID: filter.Equals(cmd.Data.UserID),
			Key:    filter.Equals(cmd.Data.Key),
		},
		&UserTag{
			ID:     uuid.New(),
			UserID: cmd.Data.UserID,
			Key:    cmd.Data.Key,
			Value:  cmd.Data.Value,
		},
	)
	if err != nil {
		return fmt.Errorf("upsert user tag: %w", err)
	}

	return nil
}

func (s *Service) GetUserTags(ctx context.Context, userID uuid.UUID) ([]*UserTag, error) {
	tags, err := s.Storage.UserTags().Find(ctx, &UserTagFilter{
		UserID: filter.Equals(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("find user tags: %w", err)
	}

	return tags, nil
}
