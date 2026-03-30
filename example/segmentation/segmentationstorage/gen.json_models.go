package segmentationstorage

import (
	driver "database/sql/driver"
	json "encoding/json"

	uuid "github.com/google/uuid"
	errors "github.com/pkg/errors"

	segmentationservice "github.com/saturn4er/boilerplate-go/example/segmentation/segmentationservice"
	// user code 'imports'
	// end user code 'imports'
)

type jsonUserTag struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Key      string    `json:"key"`
	ValueVal string    `json:"value"`
}

func (u *jsonUserTag) Scan(value any) error {
	return json.Unmarshal(value.([]byte), u)
}

func (u jsonUserTag) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func convertUserTagToJsonModel(src *segmentationservice.UserTag) (*jsonUserTag, error) {
	result := &jsonUserTag{}
	result.ID = src.ID
	result.UserID = src.UserID
	result.Key = src.Key
	result.ValueVal = src.Value
	return result, nil
}

func convertUserTagFromJsonModel(src *jsonUserTag) (*segmentationservice.UserTag, error) {
	result := &segmentationservice.UserTag{}
	result.ID = src.ID
	result.UserID = src.UserID
	result.Key = src.Key
	result.Value = src.ValueVal
	return result, nil
}

type jsonSetUserTagCommand struct {
	ID             uuid.UUID                 `json:"id"`
	Data           jsonSetUserTagCommandData `json:"data"`
	IdempotencyKey string                    `json:"idempotency_key"`
}

func (s *jsonSetUserTagCommand) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s jsonSetUserTagCommand) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func convertSetUserTagCommandToJsonModel(src *segmentationservice.SetUserTagCommand) (*jsonSetUserTagCommand, error) {
	result := &jsonSetUserTagCommand{}
	result.ID = src.ID
	tmp1, err := convertSetUserTagCommandDataToJsonModel(toPtr(src.Data))
	if err != nil {
		return nil, errors.Wrap(err, "convert SetUserTagCommandData to db")
	}
	result.Data = *tmp1
	result.IdempotencyKey = src.IdempotencyKey
	return result, nil
}

func convertSetUserTagCommandFromJsonModel(src *jsonSetUserTagCommand) (*segmentationservice.SetUserTagCommand, error) {
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

type jsonSetUserTagCommandData struct {
	UserID   uuid.UUID `json:"user_id"`
	Key      string    `json:"key"`
	ValueVal string    `json:"value"`
}

func (s *jsonSetUserTagCommandData) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s jsonSetUserTagCommandData) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func convertSetUserTagCommandDataToJsonModel(src *segmentationservice.SetUserTagCommandData) (*jsonSetUserTagCommandData, error) {
	result := &jsonSetUserTagCommandData{}
	result.UserID = src.UserID
	result.Key = src.Key
	result.ValueVal = src.Value
	return result, nil
}

func convertSetUserTagCommandDataFromJsonModel(src *jsonSetUserTagCommandData) (*segmentationservice.SetUserTagCommandData, error) {
	result := &segmentationservice.SetUserTagCommandData{}
	result.UserID = src.UserID
	result.Key = src.Key
	result.Value = src.ValueVal
	return result, nil
}
