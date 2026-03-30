package segmentationstorage

import (
	uuid "github.com/google/uuid"
	errors "github.com/pkg/errors"

	segmentationservice "github.com/saturn4er/boilerplate-go/example/segmentation/segmentationservice"
	// user code 'imports'
	// end user code 'imports'
)

type dbUserTag struct {
	ID     uuid.UUID `gorm:"column:id;primaryKey"`
	UserID uuid.UUID `gorm:"column:user_id;"`
	Key    string    `gorm:"column:key;type:text;"`
	Value  string    `gorm:"column:value;type:text;"`
}

func convertUserTagToDB(src *segmentationservice.UserTag) (*dbUserTag, error) {
	result := &dbUserTag{}
	result.ID = src.ID
	result.UserID = src.UserID
	result.Key = src.Key
	result.Value = src.Value
	return result, nil
}

func convertUserTagFromDB(src *dbUserTag) (*segmentationservice.UserTag, error) {
	result := &segmentationservice.UserTag{}
	result.ID = src.ID
	result.UserID = src.UserID
	result.Key = src.Key
	result.Value = src.Value
	return result, nil
}
func (a dbUserTag) TableName() string {
	return "user_tags"
}

type dbSetUserTagCommand struct {
	ID             uuid.UUID                 `gorm:"column:id;primaryKey"`
	Data           jsonSetUserTagCommandData `gorm:"column:data;"`
	IdempotencyKey string                    `gorm:"column:idempotency_key;type:text;"`
}

func convertSetUserTagCommandToDB(src *segmentationservice.SetUserTagCommand) (*dbSetUserTagCommand, error) {
	result := &dbSetUserTagCommand{}
	result.ID = src.ID
	tmp1, err := convertSetUserTagCommandDataToJsonModel(toPtr(src.Data))
	if err != nil {
		return nil, errors.Wrap(err, "convert SetUserTagCommandData to db")
	}
	result.Data = *tmp1
	result.IdempotencyKey = src.IdempotencyKey
	return result, nil
}

func convertSetUserTagCommandFromDB(src *dbSetUserTagCommand) (*segmentationservice.SetUserTagCommand, error) {
	result := &segmentationservice.SetUserTagCommand{}
	result.ID = src.ID
	tmp4, err := convertSetUserTagCommandDataFromJsonModel(toPtr(src.Data))
	if err != nil {
		return nil, err
	}

	result.Data = fromPtr(tmp4)
	result.IdempotencyKey = src.IdempotencyKey
	return result, nil
}
