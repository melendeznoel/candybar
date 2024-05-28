package food

import (
	"encoding/json"
	"errors"
)

func (is *Ingredients) Scan(value interface{}) error {
	if value == nil {
		is = new(Ingredients)
		return nil
	}

	source, ok := value.([]byte)

	if !ok {
		is = new(Ingredients)
		return errors.New("Type assertion failed on Ingredients")
	}

	err := json.Unmarshal(source, &is)

	if err != nil {
		is = new(Ingredients)
		return err
	}

	return nil
}
