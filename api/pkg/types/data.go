package types

import (
	"encoding/json"
	"fmt"
	"regexp"
)

const Pv = "Pv"
const svRegex = "^Sv\\d+$"

type Data interface {
	GetData() any
}

type Oi4Data struct {
	PrimaryValue any
	values       map[string]any
}

func NewOi4Data(PrimaryValue any) *Oi4Data {
	return &Oi4Data{
		PrimaryValue: PrimaryValue,
		values:       make(map[string]any),
	}
}

func ParseOi4Data(jsonData string) (*Oi4Data, error) {
	dataMap := map[string]any{}

	err := json.Unmarshal([]byte(jsonData), &dataMap)
	if err != nil {
		return nil, &Error{
			Message: "invalid data JSON",
			Err:     err,
		}
	}

	primaryValue, ok := dataMap[Pv]
	if !ok {
		return nil, &Error{Message: "primary value not found"}
	}
	delete(dataMap, Pv)

	data := NewOi4Data(primaryValue)

	for key, value := range dataMap {
		err = data.AddSecondaryData(key, &value)
		if err != nil {
			return nil, &Error{
				Message: fmt.Sprintf("secondary value key %s is not valid", key),
				Err:     err,
			}
		}
	}

	return data, nil
}

func (data *Oi4Data) GetData() any {
	data.values["Pv"] = data.PrimaryValue
	return data.values
}

func (data *Oi4Data) AddSecondaryData(tag string, value *any) error {
	re := regexp.MustCompile(svRegex)
	if !re.MatchString(tag) {
		return &Error{Message: "Tag must be in format Sv[0-9]+"}
	}
	if value == nil {
		delete(data.values, tag)
		return nil
	}

	data.values[tag] = value
	return nil
}

func (data *Oi4Data) Clear() {
	clear(data.values)
	data.PrimaryValue = nil
}
