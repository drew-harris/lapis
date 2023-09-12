package maps

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func ToMap(value datatypes.JSON) (map[string]interface{}, error) {
	raw, err := value.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var attributes map[string]interface{}
	err = json.Unmarshal(raw, &attributes)
	if err != nil {
		return nil, err
	}
	return attributes, nil
}

func FromMap(value map[string]interface{}) (datatypes.JSON, error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return datatypes.JSON{}, err
	}
	return datatypes.JSON(jsonBytes), nil
}
