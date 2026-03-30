package authservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	segmentationsvc "github.com/saturn4er/boilerplate-go/example/segmentation/segmentationservice"
	"github.com/saturn4er/boilerplate-go/lib/filter"
)

type Service struct {
	Storage Storage
}

func (s *Service) Register(ctx context.Context, email, name string, role Role) (*User, error) {
	var result *User

	err := s.Storage.ExecuteInTransaction(ctx, func(ctx context.Context, tx Storage) error {
		user, err := tx.Users().Create(ctx, &User{
			ID:    uuid.New(),
			Email: email,
			Name:  name,
			Role:  role,
		})
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		err = tx.SetUserTagCommands().Send(ctx, &segmentationsvc.SetUserTagCommand{
			ID: uuid.New(),
			Data: segmentationsvc.SetUserTagCommandData{
				UserID: user.ID,
				Key:    "registered",
				Value:  "true",
			},
			IdempotencyKey: user.ID.String(),
		})
		if err != nil {
			return fmt.Errorf("send set user tag command: %w", err)
		}

		result = user

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.Storage.Users().First(ctx, &UserFilter{
		ID: filter.Equals(id),
	})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
